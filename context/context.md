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
		child.cancel(false, err, cause)
	}
	c.children = nil
	c.mu.Unlock()

	if removeFromParent {
		removeChild(c.Context, c)
	}
}
```

