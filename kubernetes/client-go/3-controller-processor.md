## 1 Controller 和 Processor 概述

### 1.1 Controller

Controller 从 DeltaFIFO 中 pop 出来的 Deltas 出来处理，根据对象的变化更新 Indexer 本地缓存，并通知 Processor 相关对象有变化事件发生。

### 1.2 Processor

Processor 根据 controller 的通知，即根据对象的变化事件类型，调用相应的 ResourceEventHandler 来处理对象的变化。

## 2 Controller 初始化与启动分析

### 2.1 Controller 初始化

New 用于初始化 Controller，返回的是一个 interface，其结构体 controller 实现了 Controller interface。

```go
// client-go/tools/cache/controller.go
func New(c *Config) Controller {
	ctlr := &controller{
		config: *c,
		clock:  &clock.RealClock{},
	}
	return ctlr
}

// `*controller` implements Controller
type controller struct {
	config         Config
	reflector      *Reflector
	reflectorMutex sync.RWMutex
	clock          clock.Clock
}
```

### 2.2 controller 启动

controller.Run 为 controller 的启动方法。

1.  调用 NewReflector，初始化 Reflector；
2.  调用 r.Run，启动 Reflector；
3.  调用 c.processLoop，循环处理从 DeltaFIFO Pop 出来的数据，就是 controller 的核心处理。

```go

// client-go/tools/cache/controller.go
func (c *controller) Run(stopCh <-chan struct{}) {
	defer utilruntime.HandleCrash()
	go func() {
		<-stopCh
		c.config.Queue.Close()
	}()
	r := NewReflector(
		c.config.ListerWatcher,
		c.config.ObjectType,
		c.config.Queue,
		c.config.FullResyncPeriod,
	)
	r.ShouldResync = c.config.ShouldResync
	r.WatchListPageSize = c.config.WatchListPageSize
	r.clock = c.clock
	if c.config.WatchErrorHandler != nil {
		r.watchErrorHandler = c.config.WatchErrorHandler
	}

	c.reflectorMutex.Lock()
	c.reflector = r
	c.reflectorMutex.Unlock()

	var wg wait.Group

	wg.StartWithChannel(stopCh, r.Run)

	wait.Until(c.processLoop, time.Second, stopCh)
	wg.Wait()
}
```

### 2.3 controller 的核心处理分析

controller 的核心处理方法 processLoop 中，就是循环调用 c.config.Queue.Pop 将 DeltaFIFO 中的队头元素给 Pop 出来，然后调用 c.config.Process 方法来处理，当处理出错时，在调用 c.config.Queue.AddIfNotPresent 将对象重新加入到 DeltaFIFO 中。



```go
func (c *controller) processLoop() {
	for {
		obj, err := c.config.Queue.Pop(PopProcessFunc(c.config.Process))
		if err != nil {
			if err == ErrFIFOClosed {
				return
			}
			if c.config.RetryOnError {
				// This is the safe way to re-enqueue.
				c.config.Queue.AddIfNotPresent(obj)
			}
		}
	}
}
```

在 sharedIndexInformer 的初始化启动分析可以知道，c.config.Process 就是 s.HandleDeltas 方法。
```go
client-go/tools/cache/shared_informer.go
func (s *sharedIndexInformer) Run(stopCh <-chan struct{}) {
	
    // ...

	cfg := &Config{

        // ... 

		Process:           s.HandleDeltas,
	}

    // ...

	s.controller.Run(stopCh)
}
```
### 2.4 config.Process 和 s.HandleDeltas

s.HandleDeltas 主要逻辑：

1.  遍历 Deltas；
2.  判断 Delta 类型；
3.  若类型是 Sync、Replaced、Added、Updated，则从 indexer 获取对象；  
    1.  若该对象不存在，调用 s.indexer.Add 添加对象到 indexer，然后构造 addNotification struct，并且调用 s.processor.distribute；
    2.  若对象存在，调用 s.indexer.Update 更新 indexer 里面的对象；  
        1.  若类型是 Sync，distribute 参数是 true；
        2.  若类型是 Replaced，则会通过 resourceVersion 去判断 distribute 参数；
    然后构造 updateNotification struct，并且调用 s.processor.distribute；
4.  若类似是 Deleted，则调用 s.indexer.Delete 删除 index 中的对象，然构造 deleteNotification，并调用 s.processor.distribute；

todo: distribute 里面的 syncingListeners ？？？

```go
// client-go/tools/cache/shared_informer.go
func (s *sharedIndexInformer) HandleDeltas(obj interface{}) error {
	s.blockDeltas.Lock()
	defer s.blockDeltas.Unlock()

	// from oldest to newest
	for _, d := range obj.(Deltas) {
		switch d.Type {
		case Sync, Replaced, Added, Updated:
			s.cacheMutationDetector.AddObject(d.Object)
			if old, exists, err := s.indexer.Get(d.Object); err == nil && exists {
				if err := s.indexer.Update(d.Object); err != nil {
					return err
				}

				isSync := false
				switch {
				case d.Type == Sync:
					// Sync events are only propagated to listeners that requested resync
					isSync = true
				case d.Type == Replaced:
					if accessor, err := meta.Accessor(d.Object); err == nil {
						if oldAccessor, err := meta.Accessor(old); err == nil {
							// Replaced events that didn't change resourceVersion are treated as resync events
							// and only propagated to listeners that requested resync
							isSync = accessor.GetResourceVersion() == oldAccessor.GetResourceVersion()
						}
					}
				}
				s.processor.distribute(updateNotification{oldObj: old, newObj: d.Object}, isSync)
			} else {
				if err := s.indexer.Add(d.Object); err != nil {
					return err
				}
				s.processor.distribute(addNotification{newObj: d.Object}, false)
			}
		case Deleted:
			if err := s.indexer.Delete(d.Object); err != nil {
				return err
			}
			s.processor.distribute(deleteNotification{oldObj: d.Object}, false)
		}
	}
	return nil
}
```

## 3 processor 核心处理方法

### 3.1 sharedIndexInformer.processor.distribute

p.distribute 构造好 addNotification、updateNotification、deleteNotification 对象写入到 p.addCh 中。
sync 是 true ，对象写入到 p.syncingListeners 中，但是 p.syncingListeners 好像没有启动 ？？？说明 sync 会被忽略？？？

```go
// client-go/tools/cache/shared_informer.go
func (p *sharedProcessor) distribute(obj interface{}, sync bool) {
	p.listenersLock.RLock()
	defer p.listenersLock.RUnlock()

	if sync {
		for _, listener := range p.syncingListeners {
			listener.add(obj)
		}
	} else {
		for _, listener := range p.listeners {
			listener.add(obj)
		}
	}
}

func (p *processorListener) add(notification interface{}) {
	p.addCh <- notification
}
```

### 3.2 sharedIndexInformer.processor.run

s.processor.run 启动 processor，listener 有 run 和 pop 两个方法。
在代码中只启动了 p.listeners，没有启动 p.syncingListeners。

```go
// client-go/tools/cache/shared_informer.go
func (p *sharedProcessor) run(stopCh <-chan struct{}) {
	func() {
		p.listenersLock.RLock()
		defer p.listenersLock.RUnlock()
		for _, listener := range p.listeners {
			p.wg.Start(listener.run)
			p.wg.Start(listener.pop)
		}
		p.listenersStarted = true
	}()
	<-stopCh
	p.listenersLock.RLock()
	defer p.listenersLock.RUnlock()
	for _, listener := range p.listeners {
		close(listener.addCh) // Tell .pop() to stop. .pop() will tell .run() to stop
	}
	p.wg.Wait() // Wait for all .pop() and .run() to stop
}
```

### 3.3 processorListener.run 

在 processorListener.run 方法中，从 p.nextCh 读取数据，判断对象类型，执行对应的操作
```go
// client-go/tools/cache/shared_informer.go
func (p *processorListener) run() {
	stopCh := make(chan struct{})
	wait.Until(func() {
		for next := range p.nextCh {
			switch notification := next.(type) {
			case updateNotification:
				p.handler.OnUpdate(notification.oldObj, notification.newObj)
			case addNotification:
				p.handler.OnAdd(notification.newObj)
			case deleteNotification:
				p.handler.OnDelete(notification.oldObj)
			default:
				utilruntime.HandleError(fmt.Errorf("unrecognized notification: %T", next))
			}
		}
		// the only way to get here is if the p.nextCh is empty and closed
		close(stopCh)
	}, 1*time.Second, stopCh)
}
```

p.handler.OnUpdate、p.handler.OnAdd、p.handler.OnDelete 方法就是对 ResourceEventHandler interface 的实现。比如 ResourceEventHandlerFuncs 结构体，我们是可以自定义的，在客户端自定义 Add、Update、Delete 函数。  

```go
// client-go/tools/cache/shared_informer.go
type ResourceEventHandler interface {
	OnAdd(obj interface{})
	OnUpdate(oldObj, newObj interface{})
	OnDelete(obj interface{})
}

// ResourceEventHandlerFuncs is an adaptor to let you easily specify as many or
// as few of the notification functions as you want while still implementing
// ResourceEventHandler.  This adapter does not remove the prohibition against
// modifying the objects.
type ResourceEventHandlerFuncs struct {
	AddFunc    func(obj interface{})
	UpdateFunc func(oldObj, newObj interface{})
	DeleteFunc func(obj interface{})
}

// OnAdd calls AddFunc if it's not nil.
func (r ResourceEventHandlerFuncs) OnAdd(obj interface{}) {
	if r.AddFunc != nil {
		r.AddFunc(obj)
	}
}

// OnUpdate calls UpdateFunc if it's not nil.
func (r ResourceEventHandlerFuncs) OnUpdate(oldObj, newObj interface{}) {
	if r.UpdateFunc != nil {
		r.UpdateFunc(oldObj, newObj)
	}
}

// OnDelete calls DeleteFunc if it's not nil.
func (r ResourceEventHandlerFuncs) OnDelete(obj interface{}) {
	if r.DeleteFunc != nil {
		r.DeleteFunc(obj)
	}
}
```

### 3.4 processorListener.pop

processorListener.pop 方法就是把数据从 p.addCh 拿出来，然后赋值，把数据放入 p.nextCh 中。

todo 

```go
// client-go/tools/cache/shared_informer.go
func (p *processorListener) pop() {
	defer utilruntime.HandleCrash()
	defer close(p.nextCh) // Tell .run() to stop

	var nextCh chan<- interface{}
	var notification interface{}
	for {
		select {
		case nextCh <- notification:
			// Notification dispatched
			var ok bool
			notification, ok = p.pendingNotifications.ReadOne()
			if !ok { // Nothing to pop
				nextCh = nil // Disable this select case
			}
		case notificationToAdd, ok := <-p.addCh:
			if !ok {
				return
			}
			if notification == nil { // No notification to pop (and pendingNotifications is empty)
				// Optimize the case - skip adding to pendingNotifications
				notification = notificationToAdd
				nextCh = p.nextCh
			} else { // There is already a notification waiting to be dispatched
				p.pendingNotifications.WriteOne(notificationToAdd)
			}
		}
	}
}

```

## 4 总结

### 4.1 Controller
Controller  从 DeltaFIFO 中 pop Deltas 出来处理，根据对象的变化更新 Indexer 本地缓存，并通知 Processor 相关对象有变化事件发生

### 4.2 processor
Processor 根据 Controller 的通知，即根据对象的变化事件类型（ addNotification、  updateNotification、 deleteNotification），调用相应的 ResourceEventHandler （addFunc、 updateFunc、deleteFunc）来处理对象的变化。