# 限流器

限制某个服务每秒的调用本服务的频率客户端请求太多，超出服务端的服务能力，导致服务不可用。

## 令牌桶
```
令牌桶就是一个固定大小的桶，系统会以恒定的速率向桶中放 Token ，桶满则暂时不放。在请求比较少的时候桶可以先“攒” Token ，应对突发流量，如果桶中有剩余 Token 就可以一直取。如果没有剩余 Token ，则需要等到桶中被放置了 Token 才行。
```

Golang 官方扩展库自带了限流算法实现，**golang.org/x/time/rate 。是基于令牌桶实现的

### Limiter
```go
// Limiter has three main methods, Allow, Reserve, and Wait.
// Most callers should use Wait.
//
// Each of the three methods consumes a single token.
// They differ in their behavior when no token is available.
// If no token is available, Allow returns false.
// If no token is available, Reserve returns a reservation for a future token
// and the amount of time the caller must wait before using it.
// If no token is available, Wait blocks until one can be obtained
// or its associated context.Context is canceled.
//
// The methods AllowN, ReserveN, and WaitN consume n tokens.
type Limiter struct {
	mu     sync.Mutex
	limit  Limit            // 往桶里放 Token 速率
	burst  int              // 令牌桶的大小
	tokens float64          // 桶中的令牌
	// last is the last time the limiter's tokens field was updated
	last time.Time          // 上次往桶中放 Token 的时间
	// lastEvent is the latest time of a rate-limited event (past or future)
	lastEvent time.Time     // 上次发生限速器事件的事件
}
```

Limiter 提供了三种方法 Allow, Reserve, Wait 去消费 Token，每次可以消费一个，也可以消费多个。当 Token 不足的时候，每种方法都有不同的处理。

### 三种方法

#### Allow/AllowN
```go
// Allow reports whether an event may happen now.
func (lim *Limiter) Allow() bool {
	return lim.AllowN(time.Now(), 1)
}

// AllowN reports whether n events may happen at time t.
// Use this method if you intend to drop / skip events that exceed the rate limit.
// Otherwise use Reserve or Wait.
func (lim *Limiter) AllowN(t time.Time, n int) bool {
}
```
**Allow** 实际就是对 **AllowN(time.Now(),1)**进行简化的函数
**AllowN** 表示 截至到某一时刻 t ，桶中的 Token 数量是否至少为 n 个，满足返回 true，同时消费桶中 n 个 Token，反之不消费返回 false
应用场景：如果请求速率超过限制，直接丢弃或跳过请求用它，否则用 **Reserve** 或 **Wait**

#### Reserve/ReserveN

```go
// Reserve is shorthand for ReserveN(time.Now(), 1).
func (lim *Limiter) Reserve() *Reservation {
	return lim.ReserveN(time.Now(), 1)
}

// ReserveN returns a Reservation that indicates how long the caller must wait before n events happen.
// The Limiter takes this Reservation into account when allowing future events.
// The returned Reservation’s OK() method returns false if n exceeds the Limiter's burst size.
// Use this method if you wish to wait and slow down in accordance with the rate limit without dropping events.
// If you need to respect a deadline or cancel the delay, use Wait instead.
// To drop or skip events exceeding rate limit, use Allow instead.
func (lim *Limiter) ReserveN(t time.Time, n int) *Reservation {
	r := lim.reserveN(t, n, InfDuration)
	return &r
}
```
当调用完成，无论 Token 是否充足，都会返回一个 ***Reservation** 对象。可以调用该对象的 **Delay()** 方法，该方法返回的参数类型为 **time.Duration**，反映了需要等待的时间，必须等到等待时间之后，才能进行接下来的工作。如果不想等待可以调用 **Cancel()** 方法，该方法会将 Token 归还。
应用场景：如果希望根据速率等待或减速而不丢弃用它

#### Wait/WaitN
```go
// Wait is shorthand for WaitN(ctx, 1).
func (lim *Limiter) Wait(ctx context.Context) (err error) {
	return lim.WaitN(ctx, 1)
}

// WaitN blocks until lim permits n events to happen.
// It returns an error if n exceeds the Limiter's burst size, the Context is
// canceled, or the expected wait time exceeds the Context's Deadline.
// The burst limit is ignored if the rate limit is Inf.
func (lim *Limiter) WaitN(ctx context.Context, n int) (err error) {
}
```
使用 **WaitN()** 方法消费 Token 时，如果此时桶内 Token 数量不足（小于 N），那么 **Wait** 方法将会阻塞，直到 Token 满足条件。也可以设置 context 的 Deadline 或者 Timeout 来决定此次 Wait 的最长时间。

### 动态调整桶大小和速率
```go
// SetBurst is shorthand for SetBurstAt(time.Now(), newBurst).
func (lim *Limiter) SetBurst(newBurst int) {
	lim.SetBurstAt(time.Now(), newBurst)
}
// SetLimit is shorthand for SetLimitAt(time.Now(), newLimit).
func (lim *Limiter) SetLimit(newLimit Limit) {
	lim.SetLimitAt(time.Now(), newLimit)
}
```

### time/rate 源码分析
#### 两个转换
##### durationFromTokens
```go
// durationFromTokens is a unit conversion function from the number of tokens to the duration
// of time it takes to accumulate them at a rate of limit tokens per second.
func (limit Limit) durationFromTokens(tokens float64) time.Duration {
	if limit <= 0 {
		return InfDuration
	}
	seconds := tokens / float64(limit)
	return time.Duration(float64(time.Second) * seconds)
}
```
生成 N 个 Token 所需要的时间

##### tokensFromDuration
```go
// tokensFromDuration is a unit conversion function from a time duration to the number of tokens
// which could be accumulated during that duration at a rate of limit tokens per second.
func (limit Limit) tokensFromDuration(d time.Duration) float64 {
	if limit <= 0 {
		return 0
	}
	return d.Seconds() * float64(limit)
}
```
给定一段时间 d ，这段时间可以生成多少个 Token

整体流程：
1.  计算上次取 Token 的时间到当前时刻，期间一共产生了多少个 Token
    ```
    当前 Token = 新产生 Token + 之前剩余的 Token - 要消费的 Token
    ```
2.  如果消费后剩余的 Token 大于 0 ，说明此时 Token 桶内不为空，此时 Token 充足，无需调用方等待。如果 Token 小于 0，则需要等待
3.  将需要等待的时间等相关结果返回给调用方

#### Token 归还
没看明白

>	在 reverveN 中实现的，有 update state

[源码注释](./rate.go)