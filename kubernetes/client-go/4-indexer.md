# Indexer 源码分析

## 1 Indexer 概述

Indexer 中有 informer 维护的指定资源对象的相对于 etcd 数据的一份本地缓存，可通过该缓存获取资源对象，以减少 API Server 、对 etcd 的请求压力。

```go
// client-go/tools/cache/thread_safe_store.go
type threadSafeMap struct {
	lock  sync.RWMutex
	items map[string]interface{}

	// indexers maps a name to an IndexFunc
	indexers Indexers
	// indices maps a name to an Index
	indices Indices
}
```

informer 所维护的缓存依赖于 theadSafeMap 结构体中的 items 属性，其本质上是一个用 map 构建的键值对，资源对象都存在 items 这个 map 中，key 为资源对象的 **namespace/name** 组成，value 为资源对象。

Indexer 除了维护了本地内存缓存外，还有索引功能。索引的目的就是为了快速查找，比如我们要查找某个节点上的所有 pod、查找某个命名空间下的所有 pod ，利用索引，可快速查找。索引依赖 indexers 和 indices 字段。

## 2. Indexer 的结构定义

### 2.1 Indexer Interface
Indexer 接口继承了 Store 接口，以及包含几个 index 索引相关的方法声明。

```go
// client-go/tools/cache/index.go
type Indexer interface {
	Store
	Index(indexName string, obj interface{}) ([]interface{}, error)
	IndexKeys(indexName, indexedValue string) ([]string, error)
	ListIndexFuncValues(indexName string) []string
	ByIndex(indexName, indexedValue string) ([]interface{}, error)
	GetIndexers() Indexers
	AddIndexers(newIndexers Indexers) error
}
```

### 2.2 Store Interface
Store 接口是实现了 Add、Update、Delete、Get 等方法，用于操作 informer 的本地缓存
```go
// client-go/tools/cache/store.go
type Store interface {

	// Add adds the given object to the accumulator associated with the given object's key
	Add(obj interface{}) error

	// Update updates the given object in the accumulator associated with the given object's key
	Update(obj interface{}) error

	// Delete deletes the given object from the accumulator associated with the given object's key
	Delete(obj interface{}) error

	// List returns a list of all the currently non-empty accumulators
	List() []interface{}

	// ListKeys returns a list of all the keys currently associated with non-empty accumulators
	ListKeys() []string

	// Get returns the accumulator associated with the given object's key
	Get(obj interface{}) (item interface{}, exists bool, err error)

	// GetByKey returns the accumulator associated with the given key
	GetByKey(key string) (item interface{}, exists bool, err error)

	// Replace will delete the contents of the store, using instead the
	// given list. Store takes ownership of the list, you should not reference
	// it after calling this function.
	Replace([]interface{}, string) error

	// Resync is meaningless in the terms appearing here but has
	// meaning in some implementations that have non-trivial
	// additional behavior (e.g., DeltaFIFO).
	Resync() error
}
```

### 2.3 cache struct

cache 是对 Indexer interface 的一个实现，也是 Store interface 的实现，其中包含了一个 ThreadSafeStore 接口的实现，以及一个可以计算 object key 的函数 keyFunc。

cache 会根据 keyFunc 生成某个 obj 对象对应的一个唯一 key，然后调用 ThreadSafeStore 接口中的方法来操作本地缓存中的对象。

```go
// client-go/tools/cache/store.go
type cache struct {
	// cacheStorage bears the burden of thread safety for the cache
	cacheStorage ThreadSafeStore
	// keyFunc is used to make the key for objects stored in and retrieved from items, and
	// should be deterministic.
	keyFunc KeyFunc
}
```

### 2.4 ThreadSafeStore interface

ThreadSafeStore 接口包含了操作本地缓存的增删改查方法以及索引相关的方法，名称和 indexer 方法相似，就是该接口每个方法多了一个 key 参数，是由 cache struct 中的 keyFunc 计算 object 所得。

```go
// client-go/tools/cache/thread_safe_store.go
type ThreadSafeStore interface {
	Add(key string, obj interface{})
	Update(key string, obj interface{})
	Delete(key string)
	Get(key string) (item interface{}, exists bool)
	List() []interface{}
	ListKeys() []string
	Replace(map[string]interface{}, string)
	Index(indexName string, obj interface{}) ([]interface{}, error)
	IndexKeys(indexName, indexKey string) ([]string, error)
	ListIndexFuncValues(name string) []string
	ByIndex(indexName, indexKey string) ([]interface{}, error)
	GetIndexers() Indexers

	// AddIndexers adds more indexers to this store.  If you call this after you already have data
	// in the store, the results are undefined.
	AddIndexers(newIndexers Indexers) error
	// Resync is a no-op and is deprecated
	Resync() error
}
```

### 2.5 threadSafeMap struct

threadSafeMap 是 ThreadSafeStore interface 的实现，items 是用 map 构建的键值对，资源对象都存在 items 中，key 是根据资源对象算出来的，value 就是资源对象本身。indexers 和 indices 与索引相关。

```go
// client-go/tools/cache/thread_safe_store.go
type threadSafeMap struct {
	lock  sync.RWMutex
	items map[string]interface{}

	// indexers maps a name to an IndexFunc
	indexers Indexers
	// indices maps a name to an Index
	indices Indices
}
```

### 2.6 总结

1.  Store interface: 定义了Add、Update、Delete、List、Get等一些对象增删改查的方法声明，用于操作informer的本地缓存；
2.  Indexer interface: 继承了一个Store接口（实现本地缓存），以及包含几个index索引相关的方法声明（实现索引功能）；
3.  cache struct: Indexer接口的一个实现，所以自然也是Store接口的一个实现，cache struct包含一个ThreadSafeStore接口的实现，以及一个计算object key的函数KeyFunc；
4.  ThreadSafeStore interface: 包含了操作本地缓存的增删改查方法以及索引功能的相关方法，其方法名称与Indexer接口的类似，最大区别是ThreadSafeStore接口的增删改查方法入参基本都有key，由cache struct中的KeyFunc函数计算得出object key；
5.  threadSafeMap struct: ThreadSafeStore接口的一个实现，其最重要的一个属性便是items了，items是用map构建的键值对，资源对象都存在items这个map中，key根据资源对象来算出，value为资源对象本身，这里的items即为informer的本地缓存了，而indexers与indices属性则与索引功能有关


## 3. Indexer 的索引功能

在 threadSafeMap 中，indexers 和 indices 两个属性是和索引相关的。

```go
// client-go/tools/cache/thread_safe_store.go
type threadSafeMap struct {
	lock  sync.RWMutex
	items map[string]interface{}

	// indexers maps a name to an IndexFunc
	indexers Indexers
	// indices maps a name to an Index
	indices Indices
}
```

```go
// client-go/tools/cache/index.go
// Index maps the indexed value to a set of keys in the store that match on that value
type Index map[string]sets.String

// Indexers maps a name to an IndexFunc
type Indexers map[string]IndexFunc

// Indices maps a name to an Index
type Indices map[string]Index
```

### 3.1 Indexers 和 IndexFunc

```go
// client-go/tools/cache/index.go
// IndexFunc knows how to compute the set of indexed values for an object.
type IndexFunc func(obj interface{}) ([]string, error)
```

Indexers 包含了所有索引器（索引分类）及其索引器函数 IndexFunc，IndexFunc 为计算某个索引键下的所有对象键列表的方法；

```json
Indexers: {
    "索引器1": 索引函数1,
    "索引器2": 索引函数2,
}
```

**示例**：
```json
Indexers: {
    "namespace": MetaNamespaceIndexFunc,
    "nodeName": NodeNameIndexFunc,
}
```
```go
func MetaNamespaceIndexFunc(obj interface{}) ([]string, error) {
	meta, err := meta.Accessor(obj)
	if err != nil {
		return []string{""}, fmt.Errorf("object has no meta: %v", err)
	}
	return []string{meta.GetNamespace()}, nil
}

func NodeNameIndexFunc(obj interface{}) ([]string, error) {

    //....

}
```

### 3.2 Indices 和 Index

```go
// client-go/tools/cache/index.go
type Index map[string]sets.String

type Indices map[string]Index

type String map[string]Empty
```

Indices 包含了所有索引器（索引分类）及其所有的索引数据 index；而 Index 则包含了索引键以及索引键下的所有对象键的列表；

```json
Indices: {
    "索引器1": {
        "索引键1": ["资源对象key1","资源对象key2"],
        "索引键2": ["资源对象key3","资源对象key4"]
    }
    "索引器2": {
        "索引键3": ["资源对象key1","资源对象key3"],
        "索引键4": ["资源对象key2","资源对象key4"]
    }
}
```

**数据示例**：
```go
pod1 := &v1.Pod {
    ObjectMeta: metav1.ObjectMeta {
        Name: "pod-1",
        Namespace: "default",
    },
    Spec: v1.PodSpec{
        NodeName: "node1",
    }
}

pod2 := &v1.Pod {
    ObjectMeta: metav1.ObjectMeta {
        Name: "pod-2",
        Namespace: "default",
    },
    Spec: v1.PodSpec{
        NodeName: "node2",
    }
}

pod3 := &v1.Pod {
    ObjectMeta: metav1.ObjectMeta {
        Name: "pod-3",
        Namespace: "kube-system",
    },
    Spec: v1.PodSpec{
        NodeName: "node2",
    }
}
```

```json
Indexers: {
    "namespace": namespaceIndexFunc,
    "nodeNmae": nodeNameIndexFunc
}

Indices: {
    "namespace": {
        "default": ["default/pod-1", "default/pod-2"],
        "kube-system": ["kube-system/pod-3"]
    }
    "nodeName": {
        "node1": ["default/pod-1"],
        "node2": ["default/pod-2", "kube-system/pod-3"]
    }
}
```

### 3.3 索引小结

明确一点：每个 informer 是对应一个资源类型，比如 pod, deployment 等。 

1.  先有 indexers 的索引分类，以及对应的索引分类函数；
2.  indices ，可以根据索引分类函数计算每个对象的属于哪一类，比如 default，kube-system；
3. 在根据将资源对象key 添加到对应的 Index 中。

### 3.4 索引函数分析

可以看到有几个函数是和索引相关的。

函数的介绍基于下面的数据：

```json
Indexers: {
    "namespace": namespaceIndexFunc,
    "nodeNmae": nodeNameIndexFunc
}

Indices: {
    "namespace": {
        "default": ["default/pod-1", "default/pod-2"],
        "kube-system": ["kube-system/pod-3"]
    }
    "nodeName": {
        "node1": ["default/pod-1"],
        "node2": ["default/pod-2", "kube-system/pod-3"]
    }
}
```


```go
// client-go/tools/cache/index.go
type Indexer interface {
	Store
	Index(indexName string, obj interface{}) ([]interface{}, error)
	IndexKeys(indexName, indexedValue string) ([]string, error)
	ListIndexFuncValues(indexName string) []string
	ByIndex(indexName, indexedValue string) ([]interface{}, error)
	GetIndexers() Indexers
	AddIndexers(newIndexers Indexers) error
}
```

#### 3.4.1 Index(indexName string, obj interface{}) ([]interface{}, error)

返回资源对象列表是和给定资源对象 obj 用索引函数计算的索引值匹配

**示例**：
```go
items, err := indexer.Index("namespace", &metav1.ObjectMeta{Namespace: "default"})
for _, pod := range items {
    fmt.Println(pod.(*v1.Pod).Name)
}
```
输出：
```sh
pod-1
pod-2
```

分析： 如下注释

```go
// client-go/tools/cache/store.go
func (c *cache) Index(indexName string, obj interface{}) ([]interface{}, error) {
	return c.cacheStorage.Index(indexName, obj)
}
```

```go
// client-go/tools/cache/thread_safe_store.go
func (c *threadSafeMap) Index(indexName string, obj interface{}) ([]interface{}, error) {
	c.lock.RLock()
	defer c.lock.RUnlock()

    // indexName = namespace
    // indexFunc = namespaceIndexFunc
	indexFunc := c.indexers[indexName]
	if indexFunc == nil {
		return nil, fmt.Errorf("Index with name %s does not exist", indexName)
	}

    // indexedValues = []{"default"}
	indexedValues, err := indexFunc(obj)
	if err != nil {
		return nil, err
	}
	index := c.indices[indexName]

    /*
     index := {
        "default": ["default/pod-1", "default/pod-2"]
        "kube-system": ["kube-system/pod-3"]
    }
    */

	var storeKeySet sets.String
	if len(indexedValues) == 1 {
		// In majority of cases, there is exactly one value matching.
		// Optimize the most common path - deduping is not needed here.
		storeKeySet = index[indexedValues[0]]
	} else {
		// Need to de-dupe the return list.
		// Since multiple keys are allowed, this can happen.
		storeKeySet = sets.String{}
		for _, indexedValue := range indexedValues {
			for key := range index[indexedValue] {
				storeKeySet.Insert(key)
			}
		}
	}

    // storeKeySet = ["default/pod-1", "default/pod-2"]

	list := make([]interface{}, 0, storeKeySet.Len())
	for storeKey := range storeKeySet {
		list = append(list, c.items[storeKey])
	}
	return list, nil
}
```
