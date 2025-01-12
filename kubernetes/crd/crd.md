
## CRD
可参考该链接
[CRD](https://kubernetes.io/zh-cn/docs/tasks/extend-kubernetes/custom-resources/custom-resource-definitions/)

## 命令
```shell
kubectl get crd
kubectl api-resources
```

### finalizers

[finalizers](https://kubernetes.io/zh-cn/docs/tasks/extend-kubernetes/custom-resources/custom-resource-definitions/#finalizers)  
修改对象，将你开始执行删除的时间添加到 metadata.deletionTimestamp 字段。
禁止对象被删除，直到其 metadata.finalizers 字段内的所有项被删除。


### 合法性验证
```yaml
properties:
  name:
    type: string
    pattern: '^test$'
```

### 附加字段

[additionalPrinterColumns](https://kubernetes.io/zh-cn/docs/tasks/extend-kubernetes/custom-resources/custom-resource-definitions/#additional-printer-columns)

### 子资源

[subresources](https://kubernetes.io/zh-cn/docs/tasks/extend-kubernetes/custom-resources/custom-resource-definitions/#subresources)

CRD 支持 statue 和 scale 子资源

### 设置默认值
```yaml
properties:
  name:
    type: string
    default: "demo"
```

### 多版本