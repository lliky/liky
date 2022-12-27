# 条件变量（Cond）

## 使用场景
```
 sync.Cond 想访问共享资源的 goroutine ，当共享状态发生变化时候，可以用来通知被互斥锁阻塞的 goroutine。经常用于多个 goroutine 等待，一个 goroutine 通知的场景。
```
sync.Cond 是基于互斥锁或读写锁的

## sync.Cond 一个函数和三个方法

sync.Cond 定义如下
```go
// Each Cond has an associated Locker L (often a *Mutex or *RWMutex),
// which must be held when changing the condition and calling the Waie method.

// A Cond must not be copied after first use.
type Cond struct {
    // L is held while observing or changing the condition
    L Locker
    // contains filtered or unexported fields
}
```
每个 Cond 实例都会关联一个锁 L （互斥锁或者读写锁），当改变条件或调用 Wait 方法时，必须加锁。

### NewCond 创建实例
```go
func NewCond(l Locker) *Cond
```
NewCond 创建实例，需要关联一个锁。

### Wait()
```go
// Wait atomically unlocks c.L and suspends execution of the calling goroutine.
// After later resuming execution, Wait locks c.L before returning. Unlike in 
// other systems, Wait cannot return unless awoken by Broadcast and Signal.

// Because c.L is not locked when Wait first resumes, the caller typically cannot
// assume that the condition is true when Wait returns.Instead, the caller should 
//Wait in a loop:

//    c.L.Lock()
//    for !condition() {
//        c.Wait()    
//    }
//    ... make use of condition
//    c.L.Unlock()

func (c *Cond)Wait()
```
调用 Wait 会自动释放锁 c.L，并挂起调用者所在的 goroutine ，因此当前携程会阻塞在 Wait 方法调用的地方。如果其他协程调用了 Signal 或 Broadcast 唤醒了该协程，那么 Wait 方法在结束阻塞时，会重新给 c.L 加锁，并继续执行 Wait 后面的代码。

对条件的检查，使用 for 而非 if，是因为当前协程被唤醒时，条件不一定符合要求，需要再次 Wait 等待下次被唤醒，为了保险，使用 for 能够确保条件符合要求后，再执行后续的代码

**调用 Wait()，如果想使用 condition ,就需要加锁**

### Signal()
```go
// Signal wakes one goroutine waiting on c, if there is any.

// It is allowed but not required for the caller to hold c.L during the call.
func (c *Cond)Signal()
```
Signal 只唤醒任意一个等待 c 的 goroutine。
调用 Signal 的时候，可以加锁，也可以不加锁。

### Broadcast()
```go
// Broadcast wakes all goroutine waiting on c.

// It is allowed but not required for the caller to hold c.L during the call.
func (c *Cond)Broadcast()
```
Broadcast 唤醒所有等待 c 的 goroutine。
调用 Broadcast 的时候，可以加锁，也可以不加锁。

## 使用示例

![here](./cond.go)