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

## 3. 延迟队列

### 3.1 DelayingInterface

在通用队列的基础上了一个 AddAfter 接口，实现了过一会儿加入 item 的功能，使得在一个 item 处理失败之后，能够在指定延时之后重新入队。

```go
type DelayingInterface interface {
	Interface
	// AddAfter adds an item to the workqueue after the indicated duration has passed
	AddAfter(item interface{}, duration time.Duration)
}
```

### 3.2 delayingType struct

结构体定义，各个字段作用如注释。

```go
// client-go/util/workqueue/delaying_queue.go
type delayingType struct {
	// 通用队列的基本功能
	Interface

	// 计时器
	// clock tracks time for delayed firing
	clock clock.Clock

	// 队列关闭信号
	// stopCh lets us signal a shutdown to the waiting loop
	stopCh chan struct{}

	// 保证 shutdown 只执行一次
	// stopOnce guarantees we only signal shutdown a single time
	stopOnce sync.Once

	// 心跳 10s
	// heartbeat ensures we wait no more than maxWait before firing
	heartbeat clock.Ticker

	// 传递 waitfor 的 channel，默认大小 1000
	// waitingForAddCh is a buffered channel that feeds waitingForAdd
	waitingForAddCh chan *waitFor

	// metrics counts the number of retries
	metrics retryMetrics
}
```

### 3.3 waitFor

waitFor 结构定义如下
```go
// client-go/util/workqueue/delaying_queue.go
type waitFor struct {
	data    t	// 准备添加到队列的元素
	readyAt time.Time	// 加入到队列的时间
	index int	// 堆中的索引
}

```
waitForPriorityQueue 小顶堆实现了定时器功能，每次弹出最先该加入队列的元素。
```go
//client-go/util/workqueue/delaying_queue.go
type waitForPriorityQueue []*waitFor

func (pq waitForPriorityQueue) Len() int {
	return len(pq)
}
func (pq waitForPriorityQueue) Less(i, j int) bool {
	return pq[i].readyAt.Before(pq[j].readyAt)
}
func (pq waitForPriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

// Push adds an item to the queue. Push should not be called directly; instead,
// use `heap.Push`.
func (pq *waitForPriorityQueue) Push(x interface{}) {
	n := len(*pq)
	item := x.(*waitFor)
	item.index = n
	*pq = append(*pq, item)
}

// Pop removes an item from the queue. Pop should not be called directly;
// instead, use `heap.Pop`.
func (pq *waitForPriorityQueue) Pop() interface{} {
	n := len(*pq)
	item := (*pq)[n-1]
	item.index = -1
	*pq = (*pq)[0:(n - 1)]
	return item
}
```

### 3.4 NewDelayingQueue

代码中有两个点需要关注：

*	NewNamed  
	>	NewNamed 用于创建通用队列的对应类型 Type 对象，然后值赋值给了 delayingType 的 Interface 字段。
*	go ret.waitingLoop()
	> 3.5 小节解释

```go
//client-go/util/workqueue/delaying_queue.go
func NewDelayingQueue() DelayingInterface {
	return NewDelayingQueueWithCustomClock(clock.RealClock{}, "")
}

func NewDelayingQueueWithCustomClock(clock clock.WithTicker, name string) DelayingInterface {
	return newDelayingQueue(clock, NewNamed(name), name)
}

func newDelayingQueue(clock clock.WithTicker, q Interface, name string) *delayingType {
	ret := &delayingType{
		Interface:       q,
		clock:           clock,
		heartbeat:       clock.NewTicker(maxWait),
		stopCh:          make(chan struct{}),
		waitingForAddCh: make(chan *waitFor, 1000),
		metrics:         newRetryMetrics(name),
	}

	go ret.waitingLoop()
	return ret
}
```


### 3.5 waitingLoop

具体的逻辑如下面的注释

```go
//client-go/util/workqueue/delaying_queue.go
func (q *delayingType) waitingLoop() {
	defer utilruntime.HandleCrash()

	// 占位用的，队列中没有数据用的
	never := make(<-chan time.Time)

	// Make a timer that expires when the item at the head of the waiting queue is ready
	var nextReadyAtTimer clock.Timer

	// 初始化小顶堆
	waitingForQueue := &waitForPriorityQueue{}
	heap.Init(waitingForQueue)

	waitingEntryByData := map[t]*waitFor{}

	for {
		// 如果 queue 关闭，则退出 loop 
		if q.Interface.ShuttingDown() {
			return
		}

		now := q.clock.Now()

		// Add ready entries
		for waitingForQueue.Len() > 0 {
			// 如果堆不为空，则获取堆顶元素
			entry := waitingForQueue.Peek().(*waitFor)
			// 如果堆顶元素大于当前时间，说明里面的元素都还没到加入队列的时间，直接跳出
			if entry.readyAt.After(now) {
				break
			}
			// 如果堆顶元素小于当前时间，则 pop 出堆顶元素，然后加入队列中。
			entry = heap.Pop(waitingForQueue).(*waitFor)
			q.Add(entry.data)
			delete(waitingEntryByData, entry.data)
		}

		// 如果堆为空，则使用 never 做无限时长定时器
		nextReadyAt := never
		// 如果堆不为空，设置最近元素的时间为定时器的时间	
		if waitingForQueue.Len() > 0 {
			if nextReadyAtTimer != nil {
				nextReadyAtTimer.Stop()
			}
			// 获取堆顶元素
			entry := waitingForQueue.Peek().(*waitFor)
			// 实例化 timer 定时器	
			nextReadyAtTimer = q.clock.NewTimer(entry.readyAt.Sub(now))
			nextReadyAt = nextReadyAtTimer.C()
		}

		// 阻塞这里选择一个
		select {
			// 队列关闭，直接跳出 loop
		case <-q.stopCh:
			return

			// 心跳， 10s 一次，重新进行选择最近的定时任务
		case <-q.heartbeat.C():
			// continue the loop, which will add ready items

			// 上次计算的最近元素定时器已到期，进行下次循环，然后处理该元素
		case <-nextReadyAt:
			// continue the loop, which will add ready items

			// 收到新添加的定时器任务
		case waitEntry := <-q.waitingForAddCh:

			// 如果新对象还没有到期，把定时器任务放入到小顶堆中
			if waitEntry.readyAt.After(q.clock.Now()) {
				insert(waitingForQueue, waitingEntryByData, waitEntry)
			} else {
				// 如果新对象到期，则直接放入到队列中
				q.Add(waitEntry.data)
			}

			// 优化点，通过一个循环，将 waitingForAddCh 中所有的元素都消费掉，根据情况要么插入小顶堆要么放入队列
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

func insert(q *waitForPriorityQueue, knownEntries map[t]*waitFor, entry *waitFor) {
	existing, exists := knownEntries[entry.data]
	// 如果存在，就直接更新 readyAt 时间
	if exists {
		if existing.readyAt.After(entry.readyAt) {
			existing.readyAt = entry.readyAt
			heap.Fix(q, existing.index)
		}

		return
	}
	// 如果不存在，直接插入到小顶堆中
	heap.Push(q, entry)
	knownEntries[entry.data] = entry
}
```

### 3.6 AddAfter

和通用队列相比，延迟队列就是多了 AddAfter 方法;  
该方法作用就是，在指定的延时时间到达之后，将元素加入到队列中

```go
//client-go/util/workqueue/delaying_queue.go
func (q *delayingType) AddAfter(item interface{}, duration time.Duration) {
	// don't add if we're already shutting down
	// 队列关闭，直接返回
	if q.ShuttingDown() {
		return
	}

	q.metrics.retry()

	// immediately add things with no delay
	// 如果时间到了，直接入队
	if duration <= 0 {
		q.Add(item)
		return
	}

	select {
	case <-q.stopCh:
		// unblock if ShutDown() is called
		// 构造 waitFor，直接放入到 channel 中
	case q.waitingForAddCh <- &waitFor{data: item, readyAt: q.clock.Now().Add(duration)}:
	}
}

```

## 4. 限速队列

限速队列是在延迟队列的基础上实现的。

元素是通过 AddRateLimited 加入队列的，加入的时候需要先经过 rateLimiter 计算是否触发限速。如果需要则计算该元素的 delay 时长，把该对象放到延迟队列中。

### 4.1 RateLimitingInterface 

```go
// client-go/util/workqueue/rate_limiting_queue.go
type RateLimitingInterface interface {
	DelayingInterface

	// AddRateLimited adds an item to the workqueue after the rate limiter says it's ok
	AddRateLimited(item interface{})

	// Forget indicates that an item is finished being retried.  Doesn't matter whether it's for perm failing
	// or for success, we'll stop the rate limiter from tracking it.  This only clears the `rateLimiter`, you
	// still have to call `Done` on the queue.
	Forget(item interface{})

	// NumRequeues returns back how many times the item was requeued
	NumRequeues(item interface{}) int
}
```

### 4.2 rateLimitingType struct

```go
// client-go/util/workqueue/rate_limiting_queue.go
type rateLimitingType struct {
	// 延迟队列功能
	DelayingInterface
	// 限速组件
	rateLimiter RateLimiter
}
```

### 4.3 rateLimitingType 的实现

在初始化 rateLimitingType，创建了 DelayingInterface 和 RateLimiter

AddRateLimited 方法是通过限速计算出需要等待的时长，然后调用 DelayingInterface.AddAfter 方法决定把对象扔到延迟队列里还是队列里
NumRequeues 方法和 Forget 方法要根据不同的限速器具体分析

```go
// client-go/util/workqueue/rate_limiting_queue.go
func NewRateLimitingQueue(rateLimiter RateLimiter) RateLimitingInterface {
	return &rateLimitingType{
		DelayingInterface: NewDelayingQueue(),
		rateLimiter:       rateLimiter,
	}
}

func (q *rateLimitingType) AddRateLimited(item interface{}) {
	q.DelayingInterface.AddAfter(item, q.rateLimiter.When(item))
}

func (q *rateLimitingType) NumRequeues(item interface{}) int {
	return q.rateLimiter.NumRequeues(item)
}

func (q *rateLimitingType) Forget(item interface{}) {
	q.rateLimiter.Forget(item)
}

```

### 4.3 RateLimiter

RateLimiter 也是一个接口, 在 workqueue 内置了几个 RateLimiter 限速器的实现。也可以自己去实现限速器

*	BucketRateLimiter
*	ItemExponentialFailureRateLimiter
*	ItemFastSlowRateLimiter
*	WithMaxWaitRateLimiter
*	MaxOfRateLimiter

```go
// client-go/util/workqueue/default_rate_limiters.go
type RateLimiter interface {
	// 获取元素入队时间
	When(item interface{}) time.Duration
	// 删除该元素的记录
	Forget(item interface{})
	// 记录该对象次数
	NumRequeues(item interface{}) int
}
```

#### 4.3.1 BucketRateLimiter 限速器

令牌桶限速器 when 方法直接返回了该元素入队时间

令牌桶原理可参考：[令牌桶解析](./../../golang/rate-limiter/limiter.md)

```go
// client-go/util/workqueue/default_rate_limiters.go
type BucketRateLimiter struct {
	*rate.Limiter
}

func (r *BucketRateLimiter) When(item interface{}) time.Duration {
	return r.Limiter.Reserve().Delay()
}

func (r *BucketRateLimiter) NumRequeues(item interface{}) int {
	return 0
}

func (r *BucketRateLimiter) Forget(item interface{}) {
}
```

#### 4.3.2 ItemExponentialFailureRateLimiter

该限速器使用了一个 map 记录了元素的次数，后通过 backoff 算法可以求出当前需要等待的时长，只要 Forget 不擦除，下次就是上次的 2 倍，有一个最大时间。

看默认的函数：DefaultItemBasedRateLimiter

*	baseDelay: 1 ms
*	maxDelay: 1000s

延时时间就是： 1ms，2ms，4ms，8ms，16ms ... 上限时 1000s
```go
// client-go/util/workqueue/default_rate_limiters.go
type ItemExponentialFailureRateLimiter struct {
	failuresLock sync.Mutex
	failures     map[interface{}]int

	baseDelay time.Duration
	maxDelay  time.Duration
}

func NewItemExponentialFailureRateLimiter(baseDelay time.Duration, maxDelay time.Duration) RateLimiter {
	return &ItemExponentialFailureRateLimiter{
		failures:  map[interface{}]int{},
		baseDelay: baseDelay,
		maxDelay:  maxDelay,
	}
}

func DefaultItemBasedRateLimiter() RateLimiter {
	return NewItemExponentialFailureRateLimiter(time.Millisecond, 1000*time.Second)
}

func (r *ItemExponentialFailureRateLimiter) When(item interface{}) time.Duration {
	r.failuresLock.Lock()
	defer r.failuresLock.Unlock()

	// 获取上次计数，然后把失败次数加一
	exp := r.failures[item]
	r.failures[item] = r.failures[item] + 1

	// 计算 backoff 时长，是上次的两倍
	backoff := float64(r.baseDelay.Nanoseconds()) * math.Pow(2, float64(exp))
	if backoff > math.MaxInt64 {
		// 不能超过 最大延迟时间
		return r.maxDelay
	}

	calculated := time.Duration(backoff)
	if calculated > r.maxDelay {
		return r.maxDelay
	}

	return calculated
}

func (r *ItemExponentialFailureRateLimiter) NumRequeues(item interface{}) int {
	r.failuresLock.Lock()
	defer r.failuresLock.Unlock()

	return r.failures[item]
}

// 不记录该对象的次数，那么下次等待时间就从 baseDelay 算起
func (r *ItemExponentialFailureRateLimiter) Forget(item interface{}) {
	r.failuresLock.Lock()
	defer r.failuresLock.Unlock()

	delete(r.failures, item)
}

```

#### 4.3.3 ItemFastSlowRateLimiter

该限速器使用了一个 map 记录了元素的次数，然后配置两个快慢时间，以及一个快时间的最大尝试次数。

```go
// client-go/util/workqueue/default_rate_limiters.go
type ItemFastSlowRateLimiter struct {
	failuresLock sync.Mutex
	// 记录次数
	failures     map[interface{}]int

	// 快记录最大尝试次数
	maxFastAttempts int
	// 快延迟
	fastDelay       time.Duration
	// 慢延迟
	slowDelay       time.Duration
}

var _ RateLimiter = &ItemFastSlowRateLimiter{}

func NewItemFastSlowRateLimiter(fastDelay, slowDelay time.Duration, maxFastAttempts int) RateLimiter {
	return &ItemFastSlowRateLimiter{
		failures:        map[interface{}]int{},
		fastDelay:       fastDelay,
		slowDelay:       slowDelay,
		maxFastAttempts: maxFastAttempts,
	}
}

func (r *ItemFastSlowRateLimiter) When(item interface{}) time.Duration {
	r.failuresLock.Lock()
	defer r.failuresLock.Unlock()

	r.failures[item] = r.failures[item] + 1

	if r.failures[item] <= r.maxFastAttempts {
		return r.fastDelay
	}

	return r.slowDelay
}

func (r *ItemFastSlowRateLimiter) NumRequeues(item interface{}) int {
	r.failuresLock.Lock()
	defer r.failuresLock.Unlock()

	return r.failures[item]
}

func (r *ItemFastSlowRateLimiter) Forget(item interface{}) {
	r.failuresLock.Lock()
	defer r.failuresLock.Unlock()

	delete(r.failures, item)
}
```

#### 4.3.4 WithMaxWaitRateLimiter

和其他限速器一起使用，配置了需要等待最长时间值，如果从其他限速器获取的等待时间值大于配置的值，就使用配置的，否则就用限速器获取的值。

```go
// client-go/util/workqueue/default_rate_limiters.go
type WithMaxWaitRateLimiter struct {
	limiter  RateLimiter
	maxDelay time.Duration
}

func NewWithMaxWaitRateLimiter(limiter RateLimiter, maxDelay time.Duration) RateLimiter {
	return &WithMaxWaitRateLimiter{limiter: limiter, maxDelay: maxDelay}
}

func (w WithMaxWaitRateLimiter) When(item interface{}) time.Duration {
	delay := w.limiter.When(item)
	if delay > w.maxDelay {
		return w.maxDelay
	}

	return delay
}

func (w WithMaxWaitRateLimiter) Forget(item interface{}) {
	w.limiter.Forget(item)
}

func (w WithMaxWaitRateLimiter) NumRequeues(item interface{}) int {
	return w.limiter.NumRequeues(item)
}
```

#### 4.3.5 MaxOfRateLimiter

可以传入多个 RateLimiter 限速器实例, 使用 When() 求等待间隔时间, 遍历计算所有的 RateLimiter 实例, 求最大的时长. Forget 同理, 也是对所有的 RateLimiter 集合遍历调用.

```go
// client-go/util/workqueue/default_rate_limiters.go
type MaxOfRateLimiter struct {
	limiters []RateLimiter
}

func (r *MaxOfRateLimiter) When(item interface{}) time.Duration {
	ret := time.Duration(0)
	for _, limiter := range r.limiters {
		curr := limiter.When(item)
		if curr > ret {
			ret = curr
		}
	}
	return ret
}

func NewMaxOfRateLimiter(limiters ...RateLimiter) RateLimiter {
	return &MaxOfRateLimiter{limiters: limiters}
}

func (r *MaxOfRateLimiter) NumRequeues(item interface{}) int {
	ret := 0
	for _, limiter := range r.limiters {
		curr := limiter.NumRequeues(item)
		if curr > ret {
			ret = curr
		}
	}

	return ret
}

func (r *MaxOfRateLimiter) Forget(item interface{}) {
	for _, limiter := range r.limiters {
		limiter.Forget(item)
	}
}
```

## 5. 总结

Interface 是通用队列；
DelayingInterface 是在 Interface 基础上实现的延迟队列；
RateLimitingInterface 是在 DelayingInterface 基础上实现的限速队列；

限速队列中的限速器我们可以使用它们提供的，也可以使用我们自定义的

在我们自定义代码中，我们可以直接实例化一个限速队列

然后在 ResourceEventHandlerFuncs 中把元素加入队列，我们控制器编写代码

queue.Get() -> queue.Done -> queue.Forget() -> queue.AddRateLimited()
