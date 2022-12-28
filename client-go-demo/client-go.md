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