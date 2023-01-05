# Client-go

## Indexer
```go
// k8s.io/client-go/tools/cache/indexer.go

// 用于计算一个对象的索引键集合
type IndexFunc func(obj interface{})([]string, error)

// 索引键与对象键集合的映射
type Index map[string]sets.String

// 索引器名称（或者索引分类）与 IndexFunc 的映射，相当于存储索引的各种分类
type Indexer map[string]IndexFunc

// 索引器名称与 Index 索引的映射
type Indices map[string]Index
```
-   IndexFunc: 索引器函数，用于计算一个资源对象的索引键列表。比如 命名空间，Label 标签，Annotation 等属性来生成索引键列表
-   Index: 存储数据。如果要查找某个命名空间下的 Pod，那就让 Pod 按照其命名空间进行索引，对应的 Index 类型就是： **map[namespace]sets.Pod**
-   Indexers: 存储索引器，key 为索引器名称，value 为索引器实现的函数
-   Indics：存储缓存器，key 为索引器名称，value 为缓存的数据

```json
// Indexers 就是包含的所有索引器(分类)以及对应实现
Indexers: {  
  "namespace": NamespaceIndexFunc,
  "nodeName": NodeNameIndexFunc,
}
// Indices 就是包含的所有索引分类中所有的索引数据
Indices: {
 "namespace": {  //namespace 这个索引分类下的所有索引数据
  "default": ["pod-1", "pod-2"],  // Index 就是一个索引键下所有的对象键列表
  "kube-system": ["pod-3"]   // Index
 },
 "nodeName": {  //nodeName 这个索引分类下的所有索引数据(对象键列表)
  "node1": ["pod-1"],  // Index
  "node2": ["pod-2", "pod-3"]  // Index
 }
}
```

## ShareInformer
### 作用
主要负责完成两大类功能：
1. 缓存我们关注的资源对象的最新状态的数据
2. 根据资源对象的变化事件来通知我们注册的事件处理方法


## WorkQueue
### 通用队列

抽象
```go
type Interface interface {
	Add(item interface{})                     // 向队列添加一个元素
	Len() int                                 // 队列长度，元素个数
	Get() (item interface{}, shutdown bool)   // 从队列中获取一个元素，双返回值
	Done(item interface{})                    // 告知队列该元素已经处理完了
	ShutDown()                                // 关闭队列
	ShutDownWithDrain()                       
	ShuttingDown() bool                       // 查询队列时候正在关闭
}
```
实现
```go
type Type struct {
	// queue defines the order in which we will work on items. Every
	// element of queue should be in the dirty set and not in the
	// processing set.
	queue []t             // 元素数组

	// dirty defines all of the items that need to be processed.
	dirty set             // dirty 的元素集合

	// Things that are currently being processed are in the processing set.
	// These things may be simultaneously in the dirty set. When we finish
	// processing something and remove it from this set, we'll check if
	// it's in the dirty set, and if so, add it to the queue.
	processing set        // 正在处理的元素集合

	cond *sync.Cond       // 条件变量，通知队列有数据了

	shuttingDown bool     // 关闭标记
	drain        bool

	metrics queueMetrics  // 上报数据的那种吧，prometheus 差不多

	unfinishedWorkUpdatePeriod time.Duration
	clock                      clock.WithTicker
}

type empty struct{}
type t interface{}
type set map[t]empty
```

#### Add() 方法
```go
// Add marks item as needing processing.
func (q *Type) Add(item interface{}) {
  // 加锁解锁
	q.cond.L.Lock()
	defer q.cond.L.Unlock()
  // 队列正在关闭，直接返回
	if q.shuttingDown {
		return
	}
  // 已经标记为 dirty 数据，也直接返回，因为存储在 dirty 数据的集合中
	if q.dirty.has(item) {
		return
	}
  // 告知 metrics 添加了元素
	q.metrics.add(item)
  // 添加到脏数据集合中
	q.dirty.insert(item)
  // 元素刚被拿走处理，直接返回
	if q.processing.has(item) {
		return
	}
  // 添加到元素数组里面
	q.queue = append(q.queue, item)
  // 通知有新的数据，阻塞的协程会被唤醒
	q.cond.Signal()
}
```
分析队列添加元素的状态：  
1.  队列关闭了，所以不接受任何数据。
2.  队列中没有数据，直接添加在队列中。
3.  队列中已经有了数据，如何判断？map 类型肯定最快，数组顺序遍历效率太低，这就是 dirty 存在的价值之一。
4.  队列曾经存储过该元素，但是已经被拿走还没有调用 Done 时，也就是正在处理中的元素，此时再添加当前的元素应该是最新的，处理中的应该是旧的，也就是脏的。

在正常情况下，元素只会在 dirty 和 processing 存在一份，同时存在就说明该元素在被处理的同时又被添加了一次，那么先前的那次可理解为 dirty，后续添加的还要被处理。

#### Get() 方法
```go
// Get blocks until it can return an item to be processed. If shutdown = true,
// the caller should end their goroutine. You must call Done with item when you
// have finished processing it.
func (q *Type) Get() (item interface{}, shutdown bool) {
  // 加锁解锁
	q.cond.L.Lock()
	defer q.cond.L.Unlock()
  // 没有数据，阻塞
	for len(q.queue) == 0 && !q.shuttingDown {
		q.cond.Wait()
	}
  // 被唤醒，队列为空，说明队列被关闭了
	if len(q.queue) == 0 {
		// We must be shutting down.
		return nil, true
	}

	item = q.queue[0]
	// The underlying array still exists and reference this object, so the object will not be garbage collected.
	q.queue[0] = nil
	q.queue = q.queue[1:]
  // 告知 metrics 元素被取走
	q.metrics.get(item)
  // 从 dirty 集合移走，加入到 processing
	q.processing.insert(item)
	q.dirty.delete(item)

	return item, false
}
```

#### Done() 方法
```go
// Done marks item as done processing, and if it has been marked as dirty again
// while it was being processed, it will be re-added to the queue for
// re-processing.
func (q *Type) Done(item interface{}) {
  // 加锁解锁
	q.cond.L.Lock()
	defer q.cond.L.Unlock()
  // 告知 metrics 数据处理完成
	q.metrics.done(item)
  // 从 processing 集合中移除
	q.processing.delete(item)
  // 判断 dirty 集合，在处理期间是不是新添加进去，如果是，就添加到队列中去
	if q.dirty.has(item) {
		q.queue = append(q.queue, item)
		q.cond.Signal()
	} else if q.processing.len() == 0 {
		q.cond.Signal()
	}
}
```