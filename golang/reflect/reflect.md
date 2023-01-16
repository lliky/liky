# Reflect

## Interface 和 Reflect

Golang 关于类型设计的原则：
```
    变量包括（type, value）两部分
    type 包括 static type 和 concrete type。static type 在编码可以看见的类型（int, string），concrete type 是 runtime 系统看见的类型。
    类型断言能否成功，取决于变量的 concrete type，而不是 static type。因此，一个 reader 变量如果它的 concrete type 实现了 writer 方法的话，也可以被断言为 writer。
```
