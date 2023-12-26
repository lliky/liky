# Context

**主要在异步场景用于实现协程或者函数间传递取消信号以及共享数据**

## 什么是 context.Context

```go
type Context interface {
    Deadline() (deadline time.Time, ok bool)
    Done() <-chan struct{}
    Err() error
    Value(key any) any
}
```

Context 是一个 interface，定义四个 api:

- Deadline: 返回 context 的过期时间（有些不一定有 cancelCtx）；
- Done: 返回 context 中的 channel；
- Err: 返回错误；两种常见的错误，超时，取消
- Value: 返回 context 中对应的 key 的值。

当协程运行时间达到 deadline 时，就会调用取消函数，关闭 done 通道，这时所有监听 done 通道的子协程都会收到该消息，知道父协程已经关闭，自己也需要结束。



## 2 源码分析

### 2.1 Error

```go
var Canceled = errors.New("context canceled")
var DeadlineExceeded error = deadlineExceededError{}

type deadlineExceededError struct{}

func (deadlineExceededError) Error() string   { return "context deadline exceeded" }
func (deadlineExceededError) Timeout() bool   { return true }
func (deadlineExceededError) Temporary() bool { return true }
```

- Canceled: context 被 cancel 时会报此错误
- DeadlineExceeded: context 超时时会报此错误

### 2.2 emptyCtx

```go
// An emptyCtx is never canceled, has no values, and has no deadline. It is not
// struct{}, since vars of this type must have distinct addresses.
type emptyCtx int

func (*emptyCtx) Deadline() (deadline time.Time, ok bool) {
	return
}

func (*emptyCtx) Done() <-chan struct{} {
	return nil
}

func (*emptyCtx) Err() error {
	return nil
}

func (*emptyCtx) Value(key any) any {
	return nil
}
```

- emptyCtx 是一个空的 context，本质上类型为一个整型
- Deadline 方法返回一个默认时间以及 false 的 flag，标识当前 context 不存在过期时间；
- Done 方法返回一个 nil 值，用户无论往 nil 中写入或者读取数据，都会阻塞；
- Err 方法返回的错误永远为 nil；
- Value 方法返回的 value 同样永远为 nil；

### 2.3 Background and TODO

```go
var (
	background = new(emptyCtx)
	todo       = new(emptyCtx)
)
func Background() Context {
	return background
}
func TODO() Context {
	return todo
}
```

* Background: 通常 main 函数，初始化，test 等
* TODO: 还不确定用在什么地方

### 2.4 cancelCtx

#### 2.4.1 数据结构

```go
type cancelCtx struct {
	Context

	mu       sync.Mutex            // protects following fields
	done     atomic.Value          // of chan struct{}, created lazily, closed by first cancel call
	children map[canceler]struct{} // set to nil by the first cancel call
	err      error                 // set to non-nil by the first cancel call
	cause    error                 // set to non-nil by the first cancel call
}

type canceler interface {
	cancel(removeFromParent bool, err, cause error)
	Done() <-chan struct{}
}
```

- 内嵌了一个 context 作为其父 context。可见，cancelCtx 必然为某个 context 的子 context；
- 内置了一把锁，用以协调并发场景下的资源获取；
- done  实际类型为 chan struct{}，用以反映 cancelCtx 生命周期的通道；
- children：一个 map，指向 cancelCtx 的所有子 context；
- err：记录当前 cancelCtx 的错误，必然为某个 context 的子 context；

#### 2.4.2 Deadline 方法

cancelCtx 为实现该方法，仅是内嵌了一个带有 Deadline 方法的 Context interface，因此倘若直接调用会报错。应该是直接复用 父 context 的 Deadline。

#### 2.4.3 Done 方法

```go
func (c *cancelCtx) Done() <-chan struct{} {
	d := c.done.Load()
	if d != nil {
		return d.(chan struct{})
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	d = c.done.Load()
	if d == nil {
		d = make(chan struct{})
		c.done.Store(d)
	}
	return d.(chan struct{})
}
```

* 基于 atomic 包，读取 cancelCtx 中的 chan；若存在，则直接返回；
* 加锁后，检查 chan 是否存在，若存在则返回；( double check)
* 初始化 chan 存储到 atomic.Value 当中，并返回。（懒加载机制）

#### 2.4.4 Err 方法

```go
func (c *cancelCtx) Err() error {
	c.mu.Lock()
	err := c.err
	c.mu.Unlock()
	return err
}
```

* 加锁；
* 读取 cancelCtx.err；
* 解锁；
* 返回结果；

#### 2.4.5 Value 方法

```go
func (c *cancelCtx) Value(key any) any {
	if key == &cancelCtxKey {
		return c
	}
	return value(c.Context, key)
}
```

* 若 key  特定值 &cancelCtxKey，则返回 cancelCtx 自身的指针；
* 否则遵循 valueCtx 的思路取值返回。

### 2.5 context.WithCancel

#### 2.5.1 context.WithCancel()

```go
// Canceled is the error returned by Context.Err when the context is canceled.
var Canceled = errors.New("context canceled")

func WithCancel(parent Context) (ctx Context, cancel CancelFunc) {
	c := withCancel(parent)
	return c, func() { c.cancel(true, Canceled, nil) }
}

func withCancel(parent Context) *cancelCtx {
	if parent == nil {
		panic("cannot create context from nil parent")
	}
	c := newCancelCtx(parent)
	propagateCancel(parent, c)
	return c
}
```

- 校验父 context 非空；
- 注入父 context 构造好一个新的 cancelCtx；
- 在 propagateCancel 方法内启动一个守护进程，以保证父 context 终止时，该 cancelCtx 也会被终止；
- 将 cancelCtx 返回，连带返回一个用以终止该 cancelCtl 的闭包函数；

#### 2.5.2 newCancelCtx

```go
func newCancelCtx(parent Context) *cancelCtx {
	return &cancelCtx{Context: parent}
}
```

-  注入父 context 后，返回一个新的 cancelCtx；

#### 2.5.3 propagateCancel

父 context 取消，子 context 也要相应取消

```go
// propagateCancel arranges for child to be canceled when parent is.
func propagateCancel(parent Context, child canceler) {
	done := parent.Done()
	if done == nil {
		return // parent is never canceled
	}

	select {
	case <-done:
		// parent is already canceled
		child.cancel(false, parent.Err(), Cause(parent))
		return
	default:
	}
	// 父 context 也是 cancelCtx 
	if p, ok := parentCancelCtx(parent); ok {
		p.mu.Lock()
		if p.err != nil {
			// parent has already been canceled
			child.cancel(false, p.err, p.cause)
		} else {
			if p.children == nil {
				p.children = make(map[canceler]struct{})
			}
			p.children[child] = struct{}{}
		}
		p.mu.Unlock()
	} else {
		goroutines.Add(1)
		go func() {
			select {
			case <-parent.Done():
				child.cancel(false, parent.Err(), Cause(parent))
			case <-child.Done():
			}
		}()
	}
}
```

propagateCancel 方法，用以传递父子 context 之间的 cancel 事件：

- 若 parent 是不会被 cancel 的类型（如 empty），则直接返回；
- 若 parent 已经被 cancel，则直接终止子 context，并以 parent 的 err 作为子 context 的 err；
- 假如 parent 是 cancelCtx 的类型，则加锁，并将子 context 添加到 parent 的 children map 当中；
- 假如 parent 不是 cancelCtx 类型，但又存在 cancel 的能力（比如用户自定义实现的 context），则启动一个协程，通过多路复用的方式监控 parent 状态，若其终止，则同时终止子 context，并透传 parent 的 err。

parentCancelCtx 是如何校验 parent 是否为 cancelCtx 的类型：

```go
func parentCancelCtx(parent Context) (*cancelCtx, bool) {
	done := parent.Done()
	if done == closedchan || done == nil {
		return nil, false
	}
	p, ok := parent.Value(&cancelCtxKey).(*cancelCtx)
	if !ok {
		return nil, false
	}
	pdone, _ := p.done.Load().(chan struct{})
	if pdone != done {
		return nil, false
	}
	return p, true
}
```

- 若 parent 的 chan 已关闭或者是不会被 cancel 的类型，则返回 false；
- 若以特定的 cancelCtxKey 从 parent 中取值，取得的 value 是 parent 本身，则返回 true（基于 cancelCtxKey 为 key 取值时返回 cancelCtx 自身，是 cancelCtx 特有的协议）。

#### 2.5.4 cancelCtl.cancel

```go
func (c *cancelCtx) cancel(removeFromParent bool, err, cause error) {
	if err == nil {
		panic("context: internal error: missing cancel error")
	}
	if cause == nil {
		cause = err
	}
	c.mu.Lock()
	if c.err != nil {
		c.mu.Unlock()
		return // already canceled
	}
	c.err = err
	c.cause = cause
	d, _ := c.done.Load().(chan struct{})
	if d == nil {
		c.done.Store(closedchan)
	} else {
		close(d)
	}
	for child := range c.children {
		// NOTE: acquiring the child's lock while holding parent's lock.
		child.cancel(false, err, cause) //为什么传递 false 因为父 children 会置为nil
	}
	c.children = nil
	c.mu.Unlock()

	if removeFromParent {
		removeChild(c.Context, c)
	}
}
```

- cancelCtx.cancel 方法有三个入参：
  - removeFromParent：表示当前 context 是否需要父 context 的 children map 中删除；
  - err：cancel 后的错误
  - cause: cancel 后的错误原因
- 进入主体，首先校验传入的 err 是否为 nil, 若为空则 panic；
- 检查 cause 是否为 nil，若为空则将 err 赋值到 cause；
- 加锁 ；
- 检查 cancelCtx 自带的 err 是否为空，若非空说明已经被 cancel，则解锁返回；
- 将 err, cause  赋值给 cancelCtx.err, cancelCtx.cause ；
- 处理 cancelCtx 的 chan，若之前未初始化，则直接注入一个 closeChan，否则关闭该 chan；
- 遍历当前 cancelCtx.children map，依次将 children context 都进行取消，置为 nil；
- 解锁；
- 根据传入的 removeFromParent flag 判断是否需要手动把 cancelCtx 从 parent 的 children map 中移除；  

如何将 cancelCtx 从 parent 的 children map 中移除？

```go
func removeChild(parent Context, child canceler) {
	p, ok := parentCancelCtx(parent)
	if !ok {
		return
	}
	p.mu.Lock()
	if p.children != nil {
		delete(p.children, child)
	}
	p.mu.Unlock()
}
```

- 如果 parent 不是 cancelCtx，直接返回（以为只有 cancelCtx 才有 children map）
- 加锁；
- 从 parent 的 children map 中删除对应的 child；
- 解锁；

### 2.6 timerCtx

#### 2.6.1 数据结构

```go
type timerCtx struct {
	*cancelCtx
	timer *time.Timer // Under cancelCtx.mu.

	deadline time.Time
}
```

timerCtx 在 cancelCtx 基础之上又做了一层封装，除了继承 cancelCtx  能力之外新增了两个字段：

- timer：用于定时终止 context；
- deadline：用于字段 timerCtx 的过期时间；

#### 2.6.2 timerCtx.Deadline

```
func (c *timerCtx) Deadline() (deadline time.Time, ok bool) {
	return c.deadline, true
}
```

context.Context interface 下的 Deadline 仅在 timerCtx 中有效，展示过期时间；

#### 2.6.3 timerCtx.cancel

```go
func (c *timerCtx) cancel(removeFromParent bool, err, cause error) {
	c.cancelCtx.cancel(false, err, cause)
	if removeFromParent {
		// Remove this timerCtx from its parent cancelCtx's children.
		removeChild(c.cancelCtx.Context, c)
	}
	c.mu.Lock()
	if c.timer != nil {
		c.timer.Stop()
		c.timer = nil
	}
	c.mu.Unlock()
}
```

- 复用继承的 cancelCtx 的 cancel 能力，进行 cancel 处理；
- 判断是否手动从 parent 的 children map 中移除，若是则进行处理；
- 加锁；
- 停止 timer.Timer，释放资源；
- 解锁；

### 2.7 context.WithTimeout & context.WithDeadline

```go
func WithTimeout(parent Context, timeout time.Duration) (Context, CancelFunc) {
	return WithDeadline(parent, time.Now().Add(timeout))
}
```

context.WithTimeout 方法用于构造一个 timerCtx，本质上会调用 context.WithDeadline方法

```go
func WithDeadline(parent Context, d time.Time) (Context, CancelFunc) {
	if parent == nil {
		panic("cannot create context from nil parent")
	}
	if cur, ok := parent.Deadline(); ok && cur.Before(d) {
		// The current deadline is already sooner than the new one.
		return WithCancel(parent)
	}
	c := &timerCtx{
		cancelCtx: newCancelCtx(parent),
		deadline:  d,
	}
	propagateCancel(parent, c)
	dur := time.Until(d)
	if dur <= 0 {
		c.cancel(true, DeadlineExceeded, nil) // deadline has already passed
		return c, func() { c.cancel(false, Canceled, nil) }
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.err == nil {
		c.timer = time.AfterFunc(dur, func() {
			c.cancel(true, DeadlineExceeded, nil)
		})
	}
	return c, func() { c.cancel(true, Canceled, nil) }
}
```

- 校验 parent context 非空；
- 校验 parent 的过期时间是否早于自己，若是，则构造一个 cancelCtx 返回即可；
- 构造出一个新的 timerCtx；
- 启动守护方法，同步 parent 的 cancel  事件到子 context；
- 判断过期时间是否已到，若是，直接 cancel timerCtx，并返回 DeadlineExceeded 的错误；
- 加锁；
- 启动 time.Timer，设定一个延时时间，即达到过期时间后会终止该 timerCtx, 并返回 DeadlineExceeded 的错误；
- 解锁；
- 返回 timerCtx，已经封装了 cancel 逻辑的闭包 cancel 函数。

### 2.8 valueCtx

#### 2.8.1 数据结构

```go
type valueCtx struct {
	Context
	key, val any
}
```

- valueCtx 继承了一个 parent Context；
- 一个 valueCtx 中仅有一组 kv 对；

#### 2.8.2 valueCtx.Value

```go
func (c *valueCtx) Value(key any) any {
	if c.key == key {
		return c.val
	}
	return value(c.Context, key)
}
```

- 假如当前 valueCtx 的 key 等于用户传入的 key，则直接返回其 value；
- 假如不等，则从 parent context 中依次想上寻找。

```
func value(c Context, key any) any {
	for {
		switch ctx := c.(type) {
		case *valueCtx:
			if key == ctx.key {
				return ctx.val
			}
			c = ctx.Context
		case *cancelCtx:
			if key == &cancelCtxKey {
				return c
			}
			c = ctx.Context
		case *timerCtx:
			if key == &cancelCtxKey {
				return ctx.cancelCtx
			}
			c = ctx.Context
		case *emptyCtx:
			return nil
		default:
			return c.Value(key)
		}
	}
}
```

- 启动一个 for 循环，由上而下，由子及父，依次对 key 进行匹配；
- 其中 cancelCtx，timerCtx，emptyCtx 类型会有特殊的处理方式；
- 找到匹配到的 key，则将该组 value 进行返回。

#### 2.8.3 valueCtx 用法小结

valueCtx 不适合作为存储介质，存放大量的 kv 数据，原因有三：

- 一个 valueCtx 实例只能存一个 kv 对，n 个 kv 对会嵌套 n 个 valueCtx，造成空间浪费；
- 基于 k 寻找 v 的过程是线性的，时间复杂度 O(n)；
- 不支持基于 k 的去重，相同 k 可能重复存在，并基于起点的不同，返回不同的 v。valueCtx 的定位类似于请求头，只适合存放少量作用域较大的全局 meta 数据。

#### 2.8.4 context.WithValue

```go
func WithValue(parent Context, key, val any) Context {
	if parent == nil {
		panic("cannot create context from nil parent")
	}
	if key == nil {
		panic("nil key")
	}
	if !reflectlite.TypeOf(key).Comparable() {
		panic("key is not comparable")
	}
	return &valueCtx{parent, key, val}
}
```

- 若 parent context 为空， panic；
- 若 key 为空，panic；
- 若 key 的类型不能比较，panic；
- 包括 parent context 以及 kv 对，返回一个新的 valueCtx。
