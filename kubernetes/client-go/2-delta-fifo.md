## 1. DeltaFIFO 概述
 
 DeltaFIFO 是一个先进先出的队列，delta 代表变化的资源对象，其包含资源对象数据本身及其变化类型。

 Delta 结构
 ```go
// tools/cache/delta_fifo.go
type Delta struct {
	Type   DeltaType
	Object interface{}
}
 ```

 DeltaFIFO 结构

 ```go
// tools/cache/delta_fifo.go
type DeltaFIFO struct {
    // ...

	// `items` maps a key to a Deltas.
	// Each such Deltas has at least one Delta.
	items map[string]Deltas
	// `queue` maintains FIFO order of keys for consumption in Pop().
	// There are no duplicates in `queue`.
	// A key is in `queue` if and only if it is in `items`.
	queue []string
    
    // ...
}
 ```  

 DeltaFIFO 存储着 map[object key]Deltas 以及 object key 的 queue，Delta 装有对象数据及对象的变化类型。输入输出方面，Reflector 负责 DeltaFIFO 的输入， controller 负责处理 DeltaFIFO 的输出。

 一个对象能算出一个唯一的 object key，其实对应一个 Deltas，所有一个对象对应一个 Deltas。  

 DeltaType 有四种：
 ```go
    Added   DeltaType = "Added"
	Updated DeltaType = "Updated"
	Deleted DeltaType = "Deleted"
	Replaced DeltaType = "Replaced"
	Sync DeltaType = "Sync"
 ```

 针对同一个对象，可能有多个不同 Type 的 Delta 元素在 Deltas 中，表示对该对象做了不同的操作，多个相同的 Type 的 Delta 元素在 Deltas 中（除 Deleted 外，Deleted  类型会被去重），比如短时间内，多次对某个对象进行了更新操作，那么就会有多个 Updated 类型的 Delta 放入 Deltas 中。

## 2. DeltaFIFO 的定义与初始化分析

### 2.1 DeltaFIFO struct

DeltaFIFO struct 定义了 DeltaFIFO 的一些属性  

1.  lock：读写锁，操作 DeltaFIFO 中的 items 与 queue 之前都要先加锁；
2.  items：是 map，key 根据对象算出， value 为 Deltas 类型；
3.  queue：存储对象 key 的队列；
4.  keyFunc：计算对象 key 的函数；

```go
// tools/cache/delta_fifo.go
type DeltaFIFO struct {
	// lock/cond protects access to 'items' and 'queue'.
	lock sync.RWMutex
	cond sync.Cond

	// `items` maps a key to a Deltas.
	// Each such Deltas has at least one Delta.
	items map[string]Deltas

	// `queue` maintains FIFO order of keys for consumption in Pop().
	// There are no duplicates in `queue`.
	// A key is in `queue` if and only if it is in `items`.
	queue []string

	// populated is true if the first batch of items inserted by Replace() has been populated
	// or Delete/Add/Update/AddIfNotPresent was called first.
	populated bool
	// initialPopulationCount is the number of items inserted by the first call of Replace()
	initialPopulationCount int

	// keyFunc is used to make the key used for queued item
	// insertion and retrieval, and should be deterministic.
	keyFunc KeyFunc

	// knownObjects list keys that are "known" --- affecting Delete(),
	// Replace(), and Resync()
	knownObjects KeyListerGetter

	// Used to indicate a queue is closed so a control loop can exit when a queue is empty.
	// Currently, not used to gate any of CRUD operations.
	closed bool

	// emitDeltaTypeReplaced is whether to emit the Replaced or Sync
	// DeltaType when Replace() is called (to preserve backwards compat).
	emitDeltaTypeReplaced bool
}
```  
**Type Deltas**

Deltas 类型，是 Delta 的切片类型

```go
// DeltaType is the type of a change (addition, deletion, etc)
type DeltaType string
```

**Type Delta**

Delta 类型，有两个属性：

*   Type：代表的是 Delta 的类型，有 Add、Updated、Deleted、Replaced、Sync 五种类型；
*   Object：存储的资源对象，如 pod 等资源对象；

```go
type Delta struct {
	Type   DeltaType
	Object interface{}
}
```

### 2.2 DeltaFIFO 的初始化

NewDeltaFIFO 初始化了一个 items 和 queue 都为空的 DeltaFIFO 并返回。入参可以传入三个参数。
```go
// client-go/tools/cache/delta_fifo.go

// NewDeltaFIFOWithOptions returns a Queue which can be used to process changes to
// items. See also the comment on DeltaFIFO.
func NewDeltaFIFOWithOptions(opts DeltaFIFOOptions) *DeltaFIFO {
	if opts.KeyFunction == nil {
		opts.KeyFunction = MetaNamespaceKeyFunc
	}

	f := &DeltaFIFO{
		items:        map[string]Deltas{},
		queue:        []string{},
		keyFunc:      opts.KeyFunction,
		knownObjects: opts.KnownObjects,

		emitDeltaTypeReplaced: opts.EmitDeltaTypeReplaced,
	}
	f.cond.L = &f.lock
	return f
}
```

## DeltaFIFO 核心处理方法

在 sharedIndexInformer.Run 方法中调用 NewDeltaFIFOWithOptions 初始化 DeltaFIFO，然后将 DeltaFIFO 作为参数赋值给初始化的 Config。

```go
// client-go/tools/cache/shared_informer.go

func (s *sharedIndexInformer) Run(stopCh <-chan struct{}) {
	
	// ...

	fifo := NewDeltaFIFOWithOptions(DeltaFIFOOptions{
		KnownObjects:          s.indexer,
		EmitDeltaTypeReplaced: true,
	})

	cfg := &Config{
		Queue:            fifo,

		// ...

	}

	func() {
		s.startedLock.Lock()
		defer s.startedLock.Unlock()

		s.controller = New(cfg)
		s.controller.(*controller).clock = s.clock
		s.started = true
	}()

	// ...

	s.controller.Run(stopCh)
}

// New makes a new Controller from the given Config.
func New(c *Config) Controller {
	ctlr := &controller{
		config: *c,
		clock:  &clock.RealClock{},
	}
	return ctlr
}
```
在 controller.Run 方法中，调用 NewReflector 初始化 Reflector，将之前传入 DeltaFIFO 赋值给 Reflector 的 store，所以 r.store 就是 DeltaFIFO，而调用 r.store.Add、r.store.Update、r.store.Delete、r.store.Replace 方法就是 DeleteFIFO 的 Add、Update、Delte、Replace 方法。

```go
// client-go/tools/cache/controller.go
func (c *controller) Run(stopCh <-chan struct{}) {
	
	// ...

	r := NewReflector(
		c.config.ListerWatcher,
		c.config.ObjectType,
		c.config.Queue,
		c.config.FullResyncPeriod,
	)

	// ...

	wg.StartWithChannel(stopCh, r.Run)

	// ...
}
```

```go
// client-go/tools/cache/reflector.go
func NewReflector(lw ListerWatcher, expectedType interface{}, store Store, resyncPeriod time.Duration) *Reflector {
	return NewNamedReflector(naming.GetNameFromCallsite(internalPackages...), lw, expectedType, store, resyncPeriod)
}

func NewNamedReflector(name string, lw ListerWatcher, expectedType interface{}, store Store, resyncPeriod time.Duration) *Reflector {
	r := &Reflector{

		// ... 	

		store:         store,

		// ...		
	}
	// ...
	return r
}
```

### 3.1 DeltaFIFO.Add

DeltaFIFO 的 Add 操作，主要逻辑：  
（1）、加锁；  
（2）、调用 f.queueActionLocked，操作 DeltaFIFO 中的 queue 与 Deltas，根据对象 key 构造 Added 类型的新 delta 追加到相应的 Deltas 中；  
（3）、释放锁。

```go
// client-go/tools/cache/delta_fifo.go
func (f *DeltaFIFO) Add(obj interface{}) error {
	f.lock.Lock()
	defer f.lock.Unlock()
	f.populated = true
	return f.queueActionLocked(Added, obj)
}
```
### 3.2 DeltaFIFO.queueActionLocked

queueActionLocked 负责操作 DeltaFIFO 中的 queue 和 Deltas，根据对象 key 构造新的 Delta 追加到对应的 Deltas 中，主要逻辑：  
（1）计算对象的 key；  
（2）构造新的 Delta，将新的 Delta 追加到 Deltas 末尾；  
（3）调用 dedupDeltas 将 Delta 去重（只是将 Deltas 最末尾的两个 deleted 类型的 Delta 去重）；  
（4）判断对象的 key 是否在 queue 中，不在则加入到 queue 中；  
（5）根据对象 key 更新 items 中的 Deltas；  
（6）通知所有消费者接触阻塞。

```go
// client-go/tools/cache/delta_fifo.go
func (f *DeltaFIFO) queueActionLocked(actionType DeltaType, obj interface{}) error {
	id, err := f.KeyOf(obj)
	if err != nil {
		return KeyError{obj, err}
	}
	oldDeltas := f.items[id]
	newDeltas := append(oldDeltas, Delta{actionType, obj})
	newDeltas = dedupDeltas(newDeltas)

	if len(newDeltas) > 0 {
		if _, exists := f.items[id]; !exists {
			f.queue = append(f.queue, id)
		}
		f.items[id] = newDeltas
		f.cond.Broadcast()
	} else {
		// This never happens, because dedupDeltas never returns an empty list
		// when given a non-empty list (as it is here).
		// If somehow it happens anyway, deal with it but complain.
		if oldDeltas == nil {
			klog.Errorf("Impossible dedupDeltas for id=%q: oldDeltas=%#+v, obj=%#+v; ignoring", id, oldDeltas, obj)
			return nil
		}
		klog.Errorf("Impossible dedupDeltas for id=%q: oldDeltas=%#+v, obj=%#+v; breaking invariant by storing empty Deltas", id, oldDeltas, obj)
		f.items[id] = newDeltas
		return fmt.Errorf("Impossible dedupDeltas for id=%q: oldDeltas=%#+v, obj=%#+v; broke DeltaFIFO invariant by storing empty Deltas", id, oldDeltas, obj)
	}
	return nil
}
```

### 3.3 DeltaFIFO.Update

DeltaFIFO 的 Update 操作，主要逻辑：  
（1）加锁；  
（2）、调用 f.queueActionLocked，操作 DeltaFIFO 中的 queue 与 Deltas，根据对象 key 构造 Updated 类型的新 delta 追加到相应的 Deltas 中；  
（3）、释放锁。

```go
// client-go/tools/cache/delta_fifo.go
func (f *DeltaFIFO) Update(obj interface{}) error {
	f.lock.Lock()
	defer f.lock.Unlock()
	f.populated = true
	return f.queueActionLocked(Updated, obj)
}
```

### 3.4 DeltaFIFO.Delete

（1）计算对象的 key;  
（2）加锁；  
（3）items 中不存在对象 key，这直接 return，跳过处理；  
（4）调用f.queueActionLocked，操作DeltaFIFO中的queue与Deltas，根据对象key构造Deleted类型的新Delta追加到相应的Deltas中；  
（5）释放锁；

```go
// client-go/tools/cache/delta_fifo.go
func (f *DeltaFIFO) Delete(obj interface{}) error {
	id, err := f.KeyOf(obj)
	if err != nil {
		return KeyError{obj, err}
	}
	f.lock.Lock()
	defer f.lock.Unlock()
	f.populated = true
	// informer 中，knownObjects 不为 nil

	// todo knownObjects ?
	if f.knownObjects == nil {
		if _, exists := f.items[id]; !exists {
			// Presumably, this was deleted when a relist happened.
			// Don't provide a second report of the same deletion.
			return nil
		}
	} else {
		// We only want to skip the "deletion" action if the object doesn't
		// exist in knownObjects and it doesn't have corresponding item in items.
		// Note that even if there is a "deletion" action in items, we can ignore it,
		// because it will be deduped automatically in "queueActionLocked"
		_, exists, err := f.knownObjects.GetByKey(id)
		_, itemsExist := f.items[id]
		if err == nil && !exists && !itemsExist {
			// Presumably, this was deleted when a relist happened.
			// Don't provide a second report of the same deletion.
			return nil
		}
	}

	// exist in items and/or KnownObjects
	return f.queueActionLocked(Deleted, obj)
}
```
### 3.5 DeleteFIFO.Replace

1.	加锁；
2.	遍历list，计算对象的key，循环调用f.queueActionLocked，操作DeltaFIFO中的queue与Deltas，根据对象key构造Sync类型的新Delta追加到相应的Deltas中；
3.	对比DeltaFIFO中的items与Replace方法的list，如果DeltaFIFO中的items有，但传进来Replace方法的list中没有某个key，则调用f.queueActionLocked，操作DeltaFIFO中的queue与Deltas，根据对象key构造Deleted类型的新Delta追加到相应的Deltas中（避免重复，使用DeletedFinalStateUnknown包装对象）；
4.	释放锁；


```go
// client-go/tools/cache/delta_fifo.go
func (f *DeltaFIFO) Replace(list []interface{}, _ string) error {
	f.lock.Lock()
	defer f.lock.Unlock()
	keys := make(sets.String, len(list))

	// keep backwards compat for old clients
	action := Sync
	if f.emitDeltaTypeReplaced {
		action = Replaced
	}

	// Add Sync/Replaced action for each new item.
	for _, item := range list {
		key, err := f.KeyOf(item)
		if err != nil {
			return KeyError{item, err}
		}
		keys.Insert(key)
		if err := f.queueActionLocked(action, item); err != nil {
			return fmt.Errorf("couldn't enqueue object: %v", err)
		}
	}

	if f.knownObjects == nil {
		// Do deletion detection against our own list.
		queuedDeletions := 0
		for k, oldItem := range f.items {
			if keys.Has(k) {
				continue
			}
			// Delete pre-existing items not in the new list.
			// This could happen if watch deletion event was missed while
			// disconnected from apiserver.
			var deletedObj interface{}
			if n := oldItem.Newest(); n != nil {
				deletedObj = n.Object
			}
			queuedDeletions++
			if err := f.queueActionLocked(Deleted, DeletedFinalStateUnknown{k, deletedObj}); err != nil {
				return err
			}
		}

		if !f.populated {
			f.populated = true
			// While there shouldn't be any queued deletions in the initial
			// population of the queue, it's better to be on the safe side.
			f.initialPopulationCount = keys.Len() + queuedDeletions
		}

		return nil
	}

	// Detect deletions not already in the queue.
	knownKeys := f.knownObjects.ListKeys()
	queuedDeletions := 0
	for _, k := range knownKeys {
		if keys.Has(k) {
			continue
		}

		deletedObj, exists, err := f.knownObjects.GetByKey(k)
		if err != nil {
			deletedObj = nil
			klog.Errorf("Unexpected error %v during lookup of key %v, placing DeleteFinalStateUnknown marker without object", err, k)
		} else if !exists {
			deletedObj = nil
			klog.Infof("Key %v does not exist in known objects store, placing DeleteFinalStateUnknown marker without object", k)
		}
		queuedDeletions++
		if err := f.queueActionLocked(Deleted, DeletedFinalStateUnknown{k, deletedObj}); err != nil {
			return err
		}
	}

	if !f.populated {
		f.populated = true
		f.initialPopulationCount = keys.Len() + queuedDeletions
	}

	return nil
}
```