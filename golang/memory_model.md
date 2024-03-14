# 1. 内存模型

golang 内存模型有几个核心要点：

* 以空间换时间，一次缓存，多次复用

  由于每次向操作系统申请内存操作很重，所以就一次多申请一些，以备后用。

  Golang 中的堆 mheap 正是基于这个思想，产生的数据结构，我们可以从两个视角来看 Golang 运行时的堆：

  1. 对操作系统，是用户进程中缓存的内存
  2. 对于 go 进程内部，堆是所有对象的内存起源

* 多级缓存，实现无/细锁化

  ![memory_1](../image/go/memory/memory_1.png)

  堆是 Go 运行时最大的临界共享资源，意味着每次存取都要加锁。

  为了解决这个问题，Golang 在堆 mheap 上，依次细粒度化，建立了 mcentral、mcache 的模型，对三者说明如下：

  * mheap：全局的内存起源，访问要加全局锁

  * mcentral ：每种对象大小规格对应的缓存，锁的粒度仅限于同一种规格以内
  * mcache : 每个 P 持有一份的内存缓存，访问时无锁
  
* 多级规格，提高利用率

  ![](../image/go/memory/memory_2.png)  

    

  Page 和 mspan 概念：

  1. page：最小存储单元

     golang 借鉴操作系统分页管理的思想，每个最小存储单元称之为 page，大小为 8K。

  2. mspan：最小管理单元

     mspan 大小为 page 的整数倍，且从 8B  到 32KB 被划分为 67 种不同的规格，分配对象时，会根据大小映射到不同规格的 mspan，从中获取空间。

  多规格 mspan 特点：

  1. 根据规格大小，产生了等级的制度
  2. 消除了外部碎片，但不可避免产生内部碎片
  3. 宏观上能提高整体空间利用率
  4. 有个规格等级的概念，才支持 mcentral 实现细锁化

* 全局总览

  ![TCmalloc](../image/go/memory/memory_3.png)  

  上图就是 golang 的整体架构图，它是借鉴 Thread Caching Malloc 的内存模型

# 2. 核心概念

## 2.1 内存单元 mspan

![mspan](../image/go/memory/memory_4.png) 

mspan 的特质：

* mspan 是 golang 内存管理最小单元

* mspan 大小是 page 的整数倍（Go 中的 page 是 8KB），且内部的 page 是连续的。

* 每个 mspan 根据空间大小以及面向分配对象的大小，会被划分为不同等级

* 同等级的 mspan 会从属同一个 mcentral ，最终会被组织成链表，因此带有前后指针

* 由于同等级的 mspan 内聚于同一个 mcentral，所以会基于同一把互斥锁管理

* mspan 会基于 bitMap 辅助快速找到空闲内存块（object）

  ![](../image/go/memory/memory_5.png)  

mspan 的源码位于 runtime/mheap.go 文件中：

```go
type mspan struct {
	next *mspan     // next span in list, or nil if none
	prev *mspan     // previous span in list, or nil if none
	// 起始地址
	startAddr uintptr // address of first byte of span aka s.base()
  // 该 mspan 包含几页
	npages    uintptr // number of pages in span

	// 标记此前位置的 object 都已经被占用
	freeindex uintptr
	// 可以存放多少个 object
	nelems uintptr // number of object in the span.

	// bitMap 每个 bit 对应一个 object，标识该块是否被占用 和 freeindex 一起使用
	allocCache uint64

	// 标识 mspan 等级，包含 class 和 noscan 两部分信息
	spanclass             spanClass     // size class and noscan (uint8)
	// ...
}
```

