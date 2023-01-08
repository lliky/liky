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

### 延迟队列
```go
// DelayingInterface is an Interface that can Add an item at a later time. This makes it easier to
// requeue items after failures without ending up in a hot-loop.

// DelayingInterface 是一个延时队列，可以在以后的时间来添加元素的接口
// 这使得它更容易在处理失败后重新入队列，而不致于陷入 hot-loop
type DelayingInterface interface {
	// 扩展通用队列
	Interface
	// AddAfter adds an item to the workqueue after the indicated duration has passed
	// 在指定的时间之后将元素添加到工作队列中
	AddAfter(item interface{}, duration time.Duration)
a
```
延时队列就是增加一个函数来实现元素延迟的添加

```go
// delayingType wraps an Interface and provides delayed re-enquing
// delayingType 包装了 Interface 通用接口并提供了延迟重入队列
type delayingType struct {
	Interface		// 一个通用队列

	// clock tracks time for delayed firing
	// 时钟用于跟踪延迟触发的时间
	clock clock.Clock

	// stopCh lets us signal a shutdown to the waiting loop
	// 关闭信号
	stopCh chan struct{}
	// stopOnce guarantees we only signal shutdown a single time
	// 用来保证只发出一次关闭信号
	stopOnce sync.Once

	// heartbeat ensures we wait no more than maxWait before firing
	// 在触发之前确保我们等待的时间不超过 maxWait
	heartbeat clock.Ticker

	// waitingForAddCh is a buffered channel that feeds waitingForAdd
	// 延迟添加的元素封装成 waitFor 放到 channel 中
	waitingForAddCh chan *waitFor

	// metrics counts the number of retries
	// 记录重试次数
	metrics retryMetrics
}

// waitFor holds the data to add and the time it should be added
// waitFor 持有要添加的数据和应该添加的时间
type waitFor struct {
	data    t				// 添加的数据元素
	readyAt time.Time		// 在什么时候添加到队列中
	// index in the priority queue (heap)
	index int				// 优先级队列中的索引
}
```
在延迟队列的实现 delayingType 结构体包含一个通用队列 Interface 的实现，最重要的属性就是 **waitingForAddCh** ，这是一个 buffered channel，将延迟队列添加的元素封装成 **waitFor** 放到 channel 中，意思就是当到了指定时间后就将元素添加到通用队列中进行处理，还没到时间的话就放到 buffered channel。

```go
/ waitForPriorityQueue implements a priority queue for waitFor items.
//
// waitForPriorityQueue implements heap.Interface. The item occurring next in
// time (i.e., the item with the smallest readyAt) is at the root (index 0).
// Peek returns this minimum item at index 0. Pop returns the minimum item after
// it has been removed from the queue and placed at index Len()-1 by
// container/heap. Push adds an item at index Len(), and container/heap
// percolates it into the correct location.
// waitForPriorityQueue 为 waifFor 的元素实现了一个优先级队列
// 把需要延迟的元素放到一个队列中，然后在队列中按照元素的延时添加时间（readyAt）从小到大排序，最小堆
type waitForPriorityQueue []*waitFor // 优先级队列
```
#### 延时队列的实现

waitForPriorityQueue 的实现
```go
// 获取队列长度
func (pq waitForPriorityQueue) Len() int {
	return len(pq)
}
// 判断索引 i 和 j 上元素大小
func (pq waitForPriorityQueue) Less(i, j int) bool {
	// 根据时间先后顺序来决定先后顺序
	// i 位置的元素时间在 j 之前，则证明索引 i 的元素小于索引 j 的元素
	return pq[i].readyAt.Before(pq[j].readyAt)
}
// 交换索引 i 和 j 的元素
func (pq waitForPriorityQueue) Swap(i, j int) {
	// 交换元素
	pq[i], pq[j] = pq[j], pq[i]
	// 更新元素里面的索引信息
	pq[i].index = i
	pq[j].index = j
}

// Push adds an item to the queue. Push should not be called directly; instead,
// use `heap.Push`.
// 添加元素到队列中
// 不要直接调用 push 函数，应该使用 heap.Push
func (pq *waitForPriorityQueue) Push(x interface{}) {
	n := len(*pq)
	item := x.(*waitFor)
	item.index = n
	*pq = append(*pq, item)
}

// Pop removes an item from the queue. Pop should not be called directly;
// instead, use `heap.Pop`.
// 从队列弹出最后一个元素
func (pq *waitForPriorityQueue) Pop() interface{} {
	n := len(*pq)
	item := (*pq)[n-1]
	item.index = -1
	*pq = (*pq)[0:(n - 1)]
	return item
}

// Peek returns the item at the beginning of the queue, without removing the
// item or otherwise mutating the queue. It is safe to call directly.
// 直接获取队列开头的元素，不会删除元素或改变队列
func (pq waitForPriorityQueue) Peek() interface{} {
	return pq[0]
}
```

#### AddAfter() 方法
```go
// AddAfter adds the given item to the work queue after the given delay
// 在指定的延时时间之后将元素 item 添加到队列中
func (q *delayingType) AddAfter(item interface{}, duration time.Duration) {
	// don't add if we're already shutting down
	// 如果队列关闭直接退出
	if q.ShuttingDown() {
		return
	}

	q.metrics.retry()

	// immediately add things with no delay
	// 如果延迟的时间 <= 0，相当于通用队列一样直接添加元素
	if duration <= 0 {
		q.Add(item)
		return
	}

	// select 没有 default ，所以可能阻塞
	select {
		// 如果调用 ShutDown() 则解除阻塞
	case <-q.stopCh:
		// unblock if ShutDown() is called
		// 把元素封装成 waitFor 传给 waitingForAddCh
	case q.waitingForAddCh <- &waitFor{data: item, readyAt: q.clock.Now().Add(duration)}:
	}
}
```

#### waitingLoop()

**AddAfter** 往 **waitingForAddCh** 里面添加数据，如何从这 channel 消费数据的？就是 **waitingLoop** 函数，这个函数在实例话 **DelayInterface** 后启动一个协程 **newDelayingQueue**
```go
// waitingLoop runs until the workqueue is shutdown and keeps a check on the list of items to be added.
// waitingLoop 一直运行直到队列关闭并对要添加的元素列表进行检查
func (q *delayingType) waitingLoop() {
	defer utilruntime.HandleCrash()

	// Make a placeholder channel to use when there are no items in our list
	// 创建一个占位符通道，当列表中没有元素的时候，利用这个变量实现长时间等待
	never := make(<-chan time.Time)

	// Make a timer that expires when the item at the head of the waiting queue is ready
	// 构造一个定时器，当等待队列头部的元素准备好时，该定时器就会失效
	var nextReadyAtTimer clock.Timer
	// 构造一个优先级队列
	waitingForQueue := &waitForPriorityQueue{}
	// 构造最小堆
	heap.Init(waitingForQueue)
	// 用来避免元素重复添加，如果重复添加了就只更新时间
	waitingEntryByData := map[t]*waitFor{}
	// 死循环
	for {
		// 如果队列关闭，则直接退出
		if q.Interface.ShuttingDown() {
			return
		}
		// 获取当前时间
		now := q.clock.Now()

		// Add ready entries
		// 如果优先队列中有元素的话
		for waitingForQueue.Len() > 0 {
			// 获取第一个元素
			entry := waitingForQueue.Peek().(*waitFor)
			// 如果第一个元素指定的时间还没到，跳出循环，第一个元素时最小的
			if entry.readyAt.After(now) {
				break
			}
			// 时间过了，从优先队列拿出来放到通用队列里面
			// 同时要把元素从上面提到的 map 删除，因为不再判断重复添加了
			entry = heap.Pop(waitingForQueue).(*waitFor)
			q.Add(entry.data)
			delete(waitingEntryByData, entry.data)
		}

		// Set up a wait for the first item's readyAt (if one exists)
		// 如果优先队列中还有元素，那就用第一个元素指定的时间剪去当前时间作为等待时间
		// 因为优先队列是用时间排序的，后面需要更长的时间，先处理前面的元素
		nextReadyAt := never
		if waitingForQueue.Len() > 0 {
			if nextReadyAtTimer != nil {
				nextReadyAtTimer.Stop()
			}
			// 获取第一个元素
			entry := waitingForQueue.Peek().(*waitFor)
			// 第一个元素的时间减去当前时间作为等待时间
			nextReadyAtTimer = q.clock.NewTimer(entry.readyAt.Sub(now))
			nextReadyAt = nextReadyAtTimer.C()
		}

		select {
			// 退出信号
		case <-q.stopCh:
			return
			// 定时器，每过一段时间没有任何数据，就执行一次大循环
		case <-q.heartbeat.C():
			// continue the loop, which will add ready items
			// 上面的等待时间信号，继续循环，添加准备好的元素
		case <-nextReadyAt:
			// continue the loop, which will add ready items
			// AddAfter 函数中放入到 channel 的元素，从 channel 获取
		case waitEntry := <-q.waitingForAddCh:
			// 如果时间过了，直接加入通用队列里面，没过加入优先队列里面
			if waitEntry.readyAt.After(q.clock.Now()) {
				insert(waitingForQueue, waitingEntryByData, waitEntry)
			} else {
				// 放入通用队列
				q.Add(waitEntry.data)
			}
			// 下面就是把 channel 里面的数据全部取出来，如果没有数据就退出，继续大循环
			drained := false
			for !drained {
				select {
				case waitEntry := <-q.waitingForAddCh:
					if waitEntry.readyAt.After(q.clock.Now()) {
						insert(waitingForQueue, waitingEntryByData, waitEntry)
					} else {
						q.Add(waitEntry.data)
					}
				default:
					drained = true
				}
			}
		}
	}
}

// insert adds the entry to the priority queue, or updates the readyAt if it already exists in the queue
// 插入元素到优先队列，如果已经存在则更新时间
func insert(q *waitForPriorityQueue, knownEntries map[t]*waitFor, entry *waitFor) {
	// if the entry already exists, update the time only if it would cause the item to be queued sooner
	existing, exists := knownEntries[entry.data]
	if exists {
		// 元素存在，谁的时间新就用谁的
		if existing.readyAt.After(entry.readyAt) {
			existing.readyAt = entry.readyAt
			// 调整优先级队列
			heap.Fix(q, existing.index)
		}

		return
	}
	// 放入元素到优先队列
	heap.Push(q, entry)
	// 更新 map
	knownEntries[entry.data] = entry
}

```

### 限速队列
原理：利用延迟队列的特性，延迟某个元素的插入时间来达到限速的目的，限速队列是延迟队列的扩展，增加了 **AddRateLimited** 、 **Forget** 、 **NumRequeues** 3个方法。
#### 接口
```go
// RateLimitingInterface is an interface that rate limits items being added to the queue.
// RateLimitingInterface 是对加入队列的元素进行速率限制的接口
type RateLimitingInterface interface {
	// 延时队列
	DelayingInterface

	// AddRateLimited adds an item to the workqueue after the rate limiter says it's ok
	// 在限速器说 ok 后，将元素加入到工作队列中
	AddRateLimited(item interface{})

	// Forget indicates that an item is finished being retried.  Doesn't matter whether it's for perm failing
	// or for success, we'll stop the rate limiter from tracking it.  This only clears the `rateLimiter`, you
	// still have to call `Done` on the queue.
	// 丢弃指定元素
	Forget(item interface{})

	// NumRequeues returns back how many times the item was requeued
	// 查询元素放入队列的次数
	NumRequeues(item interface{}) int
}
```
#### 实现
```go
// rateLimitingType wraps an Interface and provides rateLimited re-enquing
// 限速队列的实现
type rateLimitingType struct {
	// 集成了延迟队列
	DelayingInterface
	// 限速器
	rateLimiter RateLimiter
}

// AddRateLimited AddAfter's the item based on the time when the rate limiter says it's ok
// 通过限速器获取延迟时间，然后加入到延时队列
func (q *rateLimitingType) AddRateLimited(item interface{}) {
	q.DelayingInterface.AddAfter(item, q.rateLimiter.When(item))
}
// 直接通过限速器获取元素放入队列的次数
func (q *rateLimitingType) NumRequeues(item interface{}) int {
	return q.rateLimiter.NumRequeues(item)
}
// 直接通过限速器丢弃指定的元素
func (q *rateLimitingType) Forget(item interface{}) {
	q.rateLimiter.Forget(item)
}
```

#### 限速器
接口定义
```go
type RateLimiter interface {
	// When gets an item and gets to decide how long that item should wait
	// 获取指定的元素需要等待多久
	When(item interface{}) time.Duration
	// Forget indicates that an item is finished being retried.  Doesn't matter whether it's for failing
	// or for success, we'll stop tracking it
	// 释放指定元素，表示该元素已经处理
	Forget(item interface{})
	// NumRequeues returns back how many failures the item has had
	// 返回某个对象被重新入队多少次，监控用
	NumRequeues(item interface{}) int
}
```

##### BucketRateLimiter
```go
// BucketRateLimiter adapts a standard bucket to the workqueue ratelimiter API
// 令牌桶限速器，固定速率 qps
type BucketRateLimiter struct {
	// go 自带的
	*rate.Limiter
}
// 判断是否实现 RateLimiter 接口
var _ RateLimiter = &BucketRateLimiter{}

func (r *BucketRateLimiter) When(item interface{}) time.Duration {
	// 获取需要等待的时间（延迟），而且这个延迟是一个相对固定的周期
	return r.Limiter.Reserve().Delay()
}

func (r *BucketRateLimiter) NumRequeues(item interface{}) int {
	// 固定频率，不需要重试
	return 0
}

func (r *BucketRateLimiter) Forget(item interface{}) {
	// 不需要重试，也不需要丢弃
}
```
令牌桶限速器里面直接包装一个令牌桶 Limiter 对象。

#### ItemExponentialFailureRateLimiter
指数增长限速器，元素错误次数指数递增限速器，根据元素错误次数逐渐累加等待时间
```go
// ItemExponentialFailureRateLimiter does a simple baseDelay*2^<num-failures> limit
// dealing with max failures and expiration are up to the caller
// 当处理对象失败的时候，其再次入队的等待时间 *2，到 MaxDelay 为止，直到超过最大失败次数
type ItemExponentialFailureRateLimiter struct {
	// 修改失败次数用到的锁
	failuresLock sync.Mutex
	// 记录每个元素错误次数
	failures     map[interface{}]int
	// 元素延迟基数
	baseDelay time.Duration
	// 元素最大的延迟时间
	maxDelay  time.Duration
}

func (r *ItemExponentialFailureRateLimiter) When(item interface{}) time.Duration {
	r.failuresLock.Lock()
	defer r.failuresLock.Unlock()
	// 累加错误计数
	exp := r.failures[item]
	r.failures[item] = r.failures[item] + 1

	// The backoff is capped such that 'calculated' value never overflows.
	// 通过错误次数计算延迟时间： 2^i * baseDelay
	backoff := float64(r.baseDelay.Nanoseconds()) * math.Pow(2, float64(exp))
	if backoff > math.MaxInt64 {
		// 最大延迟时间
		return r.maxDelay
	}
	// 取计算的延迟值和最大延迟值的最小值
	calculated := time.Duration(backoff)
	if calculated > r.maxDelay {
		return r.maxDelay
	}

	return calculated
}
// 元素错误次数，直接从 failures 中取
func (r *ItemExponentialFailureRateLimiter) NumRequeues(item interface{}) int {
	r.failuresLock.Lock()
	defer r.failuresLock.Unlock()

	return r.failures[item]
}
// 直接从 failures 删除指定元素
func (r *ItemExponentialFailureRateLimiter) Forget(item interface{}) {
	r.failuresLock.Lock()
	defer r.failuresLock.Unlock()

	delete(r.failures, item)
}
```
#### ItemFastSlowRateLimiter
快慢限速器，先用快尝试超过阈值之后，用慢尝试，很少使用该限速器
```go
// ItemFastSlowRateLimiter does a quick retry for a certain number of attempts, then a slow retry after that
// 快慢限速器，先以 fastDelay 为周期进行尝试，超过 maxFastAttempts 次数后，以 slowDelay 为周期进行尝试
type ItemFastSlowRateLimiter struct {
	failuresLock sync.Mutex
	// 错误次数计数
	failures     map[interface{}]int
	// 错误尝试阈值
	maxFastAttempts int
	// 短延迟时间
	fastDelay       time.Duration
	// 长延迟时间
	slowDelay       time.Duration
}
```

#### MaxOfRateLimiter 
混合限速器，内部有多个限速器，选择所有限速器中**速度最慢**的一种方案，比如内部有三个限速器，When 接口返回延迟最大的那个。
```go
// MaxOfRateLimiter calls every RateLimiter and returns the worst case response
// When used with a token bucket limiter, the burst could be apparently exceeded in cases where particular items
// were separately delayed a longer time.
// 选择所有限速器中速度最慢的一种方案
type MaxOfRateLimiter struct {
	// 限速器数组
	limiters []RateLimiter
}

func (r *MaxOfRateLimiter) When(item interface{}) time.Duration {
	ret := time.Duration(0)
	// 获取所有限速器里面时间最大的延迟时间
	for _, limiter := range r.limiters {
		curr := limiter.When(item)
		if curr > ret {
			ret = curr
		}
	}

	return ret
}

func (r *MaxOfRateLimiter) NumRequeues(item interface{}) int {
	ret := 0
	// 次数也是取最大值
	for _, limiter := range r.limiters {
		curr := limiter.NumRequeues(item)
		if curr > ret {
			ret = curr
		}
	}

	return ret
}

func (r *MaxOfRateLimiter) Forget(item interface{}) {
	// 循环遍历 Forget
	for _, limiter := range r.limiters {
		limiter.Forget(item)
	}
}
```

Kubernetes 中默认的控制器限速器初始化就是使用的混合限速器

```go
// DefaultControllerRateLimiter is a no-arg constructor for a default rate limiter for a workqueue.  It has
// both overall and per-item rate limiting.  The overall is a token bucket and the per-item is exponential
// 实例化默认的限速器，由 ItemExponentialFailureRateLimiter 和 BucketRateLimiter 组成的混合限速器
func DefaultControllerRateLimiter() RateLimiter {
	return NewMaxOfRateLimiter(
		NewItemExponentialFailureRateLimiter(5*time.Millisecond, 1000*time.Second),
		// 10 qps, 100 bucket size.  This is only for retry speed and its only the overall factor (not per item)
		// 10 qps, 100 buckent 容量
		&BucketRateLimiter{Limiter: rate.NewLimiter(rate.Limit(10), 100)},
	)
}

func NewItemExponentialFailureRateLimiter(baseDelay time.Duration, maxDelay time.Duration) RateLimiter {
	return &ItemExponentialFailureRateLimiter{
		failures:  map[interface{}]int{},
		baseDelay: baseDelay,
		maxDelay:  maxDelay,
	}
}
```