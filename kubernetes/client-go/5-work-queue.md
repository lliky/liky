# Work Queue

## 1. Work Queue 概述
在 controller 中的 processor ，我们会自定义自己的 Handler 处理对象的变化，但是我们的处理速度和事件变化的速度不是一个量级的，所以就引入了 **work queue**。
在 client-go 里面有三种队列：  

*   通用队列
*   延迟队列
*   限速队列

在 **doc.go**  的注释文档中，它们不但有队列的基本功能，它们具有以下特性：

```go
// client-go/util/workqueue/doc.go
// Package workqueue provides a simple queue that supports the following
// features:
//  * Fair: items processed in the order in which they are added.
//  * Stingy: a single item will not be processed multiple times concurrently,
//      and if an item is added multiple times before it can be processed, it
//      will only be processed once.
//  * Multiple consumers and producers. In particular, it is allowed for an
//      item to be reenqueued while it is being processed.
//  * Shutdown notifications.

```

*   公平：元素是按照先进先出的顺序处理的；
*   小气：一个元素不会被并发处理，如果一个元素在处理之前被添加多次，则它会被处理一次；
*   支持多个消费者和生产者，特别的，它允许一个在处理中的元素重新入队；
*   关闭通知。

## 2. 通用队列  

### 2.1 Interface interface

通用队列实现 Interface interfac

```go
// client-go/util/workqueue/queue.go
type Interface interface {
	Add(item interface{})
	Len() int
	Get() (item interface{}, shutdown bool)
	Done(item interface{})
	ShutDown()
	ShutDownWithDrain()
	ShuttingDown() bool
}
```

### 2.2 Type struct

Type stuct 实现了 Interface interface , 里面最主要的三个字段:  

*	queue: queue 定义了队列中 items 的顺序, 是 []interface 切片, 可以保存任意数据;
*	dirty: dirty 类型是一个 map, 说明 queue 中的元素都会存放到 dirty 中, 就是待处理的元素;
*	processing: processing 类型也是一个 map, 保存的是正在被处理的元素;

为了不让 queue 中存在重复的元素, 所以加了 dirty, 判断 map 是否存在某个元素比判断 slice 中是否存在快的多.

```
// client-go/util/workqueue/queue.go

type Type struct {
	// queue defines the order in which we will work on items. Every
	// element of queue should be in the dirty set and not in the
	// processing set.
	queue []t

	// dirty defines all of the items that need to be processed.
	dirty set

	// Things that are currently being processed are in the processing set.
	// These things may be simultaneously in the dirty set. When we finish
	// processing something and remove it from this set, we'll check if
	// it's in the dirty set, and if so, add it to the queue.
	processing set

	cond *sync.Cond

	shuttingDown bool
	drain        bool
    
    // ...
}
```

### 2.3 Add()

主要逻辑:  

1.	加锁;
2.	如果队列关闭,直接返回;
3.	如果 dirty 中存在,直接返回, 否则加入 dirty;  
	> 这样子解决了 queue 中重复放入某个元素的问题;  
	> 就是 doc.go 说的第二条, 不允许同一个元素被并发处理;
4.	如果在 processing 存在, 直接返回;
5.	加入队列;
6.	发送一个信号, 唤醒可能正在等待获取元素的 goroutine;
7. 	解锁;
```go
// client-go/util/workqueue/queue.go

// Add marks item as needing processing.
func (q *Type) Add(item interface{}) {
	q.cond.L.Lock()
	defer q.cond.L.Unlock()
	if q.shuttingDown {
		return
	}
	if q.dirty.has(item) {
		return
	}

	q.metrics.add(item)

	q.dirty.insert(item)
	if q.processing.has(item) {
		return
	}

	q.queue = append(q.queue, item)
	q.cond.Signal()
}
```

### 2.4 Get()

主要逻辑:
1.	加锁;
2.	循环判断队列是否为空,并且没有队列没有关闭使当前 goroutine 进入等待状态;
3.	goroutine 被唤醒, 且队列为空,说明队列关闭;
4.	取出第一个元素;
5.	加入到 processing 中;
7.	删除 dirty 中的元素;
8. 	解锁;

```go
// client-go/util/workqueue/queue.go

func (q *Type) Get() (item interface{}, shutdown bool) {
	q.cond.L.Lock()
	defer q.cond.L.Unlock()
	for len(q.queue) == 0 && !q.shuttingDown {
		q.cond.Wait()
	}
	if len(q.queue) == 0 {
		// We must be shutting down.
		return nil, true
	}

	item = q.queue[0]
	// The underlying array still exists and reference this object, so the object will not be garbage collected.
	q.queue[0] = nil
	q.queue = q.queue[1:]

	q.metrics.get(item)

	q.processing.insert(item)
	q.dirty.delete(item)

	return item, false
}
```

### 2.5 Done()

在自定义的代码中, 我们 Get() 元素之后, 处理完成元素之后, 需要调用 Done 函数;

主要逻辑:
1. 	加锁;
2.	删除 processing 中的元素;
3.	如果处理完成的元素还在脏队列中, 说明在处理的时候加入进去的, 这个时候需要将该元素重新入队, 同时发送信号;
4.	如果 processing 中没有元素了, 那么队列可以优雅的关闭;  
	> q.processing.len() == 0 和 ShutDownWithDrain 有关;
5.	解锁;

```go
// client-go/util/workqueue/queue.go

func (q *Type) Done(item interface{}) {
	q.cond.L.Lock()
	defer q.cond.L.Unlock()

	q.metrics.done(item)

	q.processing.delete(item)
	if q.dirty.has(item) {
		q.queue = append(q.queue, item)
		q.cond.Signal()
	} else if q.processing.len() == 0 {
		q.cond.Signal()
	}
}
```

### 2.6 ShutDown() 

ShutDown 关闭主要就是广播给所有等待新元素的 goroutine ,设置标志位

```go 
// client-go/util/workqueue/queue.go
func (q *Type) ShutDown() {
	q.setDrain(false)
	q.shutdown()
}

func (q *Type) shutdown() {
	q.cond.L.Lock()
	defer q.cond.L.Unlock()
	q.shuttingDown = true
	q.cond.Broadcast()
}
```

### 2.7 ShutDownWithDrain()

优雅的关闭  
主要的逻辑就是 关闭主要就是广播给所有等待新元素的,设置标志位;
如果在还有正在处理的元素, 需要正在处理的元素处理完毕之后, 才会真正的关闭退出;
**waitForProcessing** 中 q.cond.Wait 对应了 **Done()** 中 **q.processing.len() == 0**

```go
// client-go/util/workqueue/queue.go

func (q *Type) ShutDownWithDrain() {
	q.setDrain(true)
	q.shutdown()
	for q.isProcessing() && q.shouldDrain() {
		q.waitForProcessing()
	}
}


func (q *Type) isProcessing() bool {
	q.cond.L.Lock()
	defer q.cond.L.Unlock()
	return q.processing.len() != 0
}

func (q *Type) waitForProcessing() {
	q.cond.L.Lock()
	defer q.cond.L.Unlock()
	// Ensure that we do not wait on a queue which is already empty, as that
	// could result in waiting for Done to be called on items in an empty queue
	// which has already been shut down, which will result in waiting
	// indefinitely.
	if q.processing.len() == 0 {
		return
	}
	q.cond.Wait()
}
```

### 2.8 总结

dirty 是对未被处理元素进行去重;
processing 是对正在处理的元素进行去重和重排队;

如果一个元素正在被处理, 再次添加同一个元素, 是放到 dirty 里面, 而不是直接放入 queue, 是因为放入到 queue 中, 并发的场景, 可能导致同一个元素, 同时在处理. 如果是放入到 dirty 里, 在调用 Done 的时候, 说明该元素已经处理完成, 如果还存在 dirty 里, 让该元素重新入队, 然后可以被其他等待的 goroutine 消费, 保证了一个元素不会被同时处理.