# Reflect

## Interface 和 Reflect

Golang 关于类型设计的原则：
```
    变量包括（type, value）两部分
    type 包括 static type 和 concrete type。static type 在编码可以看见的类型（int, string），concrete type 是 runtime 系统看见的类型。
    类型断言能否成功，取决于变量的 concrete type，而不是 static type。因此，一个 reader 变量如果它的 concrete type 实现了 writer 方法的话，也可以被断言为 writer。
```
反射是建立在类型之上的，Golang 的指定类型的变量的类型是静态的（指定 int, string 这些变量，它们的 type 是 static type），在创建变量的时候已经确定。反射主要与 interface 类型有关（它的 type 是 concrete type）。

每个 interface 变量都有一个对应 pair，pair 中记录了实现变量的值和类型：
```
    (value, type)
```
value 是实际变量值，type 是实际变量的类型。一个 interface{} 类型的变量包含了 2 个指针，一个指针指向值得类型（对应得 concrete type），另外一个指向实际得值（对应 value）。
例如：创建类型为 *os.File 的变量，然后将其赋值给一个接口变量 r:
```go
    tty, err := os.OpenFile("/dev/tty", os.O_RDWR, 0)
    var r io.Reader
    r = tty
```
接口变量 r 的 pair 中将记录如下信息：（tty, *os.File），将 r 赋值给另一个接口变量 w ：
```go
    var w io.Writer
    w = r.(io.Writer)
```
接口变量 w 的 pair 和 r 的 pair 相同，都是（tty, *os.File），即使 w 是空类型接口， pair 也是不变的。
反射就是用来检测存储在接口变量内部（value, concrete type）pair 对的一种机制。

# Golang 的发射 reflect
## reflect 的基本功能 TpyeOf 和 ValueOf
反射是用来检测存储在接口变量内部 pair 对的一种机制。那用 reflect 反射包用什么方式可以直接获取到变量内部的信息？它提供了两种类型，分别是 reflect.ValueOf 和 reflect.TypeOf:
```go
// TypeOf returns the reflection Type that represents the dynamic type of i.
// If i is a nil interface value, TypeOf returns nil.
func TypeOf(i interface{}) Type {
	eface := *(*emptyInterface)(unsafe.Pointer(&i))
	return toType(eface.typ)
}

// ValueOf returns a new Value initialized to the concrete value
// stored in the interface i. ValueOf(nil) returns the zero Value.
func ValueOf(i interface{}) Value {
	if i == nil {
		return Value{}
	}
	escapes(i)
	return unpackEface(i)
}
```
reflect.TypeOf() 是获取 pair 中的 type，reflect.ValueOf() 是获取 pair 中的 value。
```go
package main

import(
    "fmt"
    "reflect"
)

func main() {
    var a = "hi"
    fmt.Println("type: "reflect.Typeof(a))
    fmt.Println("value: "reflect.ValueOf(b))
}

running result:
    type: string
    value: hi
```
**说明**
1. reflect.TypeOf: 直接给到我们想要的 type 类型，如 string, int, 各种 pointer, sturct 等等真实的类型
2. reflect.ValueOf: 直接给到我们想要的具体值，如 hi 这个具体值。
3. 也就是说明反射可以将 “接口类型变量” 转换为 “反射类型对象”，反射类型指的是 reflect.Type 和 reflect.Value 这两种

## 从 reflect.Value 中获取接口 interface 的信息
当执行 reflect.ValueOf( interface )之后，就得到了一个类型为 "reflect.Value" 变量，可以通过它本身的 Interface() 方法获得接口变量的真实内容，然后他就可以通过类型进行转换，转换为原有的真实类型。可能是**已知原有类型**，也有可能是**未知原有类型**。
### 已知原有类型（进行强制转换）

已知类型后转换为其对应的类型做法如下，直接通过 Interface 方法然后强制转换：
```
    realValue := value.Interface().(已知类型)
```

[example2](./example2.go)

**说明**

1. 转换的时候，如果转换的类型不完全符合，则直接 painc，类型要求非常严格！
2. 转换的时候，需要区分指针还是值
3. 反射可以将"反射类型对象"再重新转换为"接口类型变量"

### 未知原有类型（遍历探测其 Filed）

[example3](./example3.go)

**说明**

获取未知类型的 interface 的具体变量及类型的步骤为：  
1.  先获取 interface 的 reflect.Type, 然后通过 NumField 进行遍历
2.  再通过 reflect.Type 的 Field 获取其 Field
3.  最后通过 Field 的 interface() 得到对应的 value

获取未知类型的 interface 的所属方法的步骤为：  
1. 先获取interface的reflect.Type，然后通过NumMethod进行遍历
2. 再分别通过reflect.Type的Method获取对应的真实的方法
3. 最后对结果取其Name和Type得知具体的方法名

反射可以将“反射类型对象”再重新转换为“接口类型变量”，struct 或者 struct 的嵌套都是一样的判断处理方式

### 通过 reflect.Value 设置实际变量的值

reflect.Value 可以通过 reflect.ValueOf(X) 获得，只有当 X 是指针的时候，才可以通过 reflect.Value 修改实际变量 X 的值，即：要修改反射类型的对象就一定保证其值是 "addressable" 的

[example4](./example4.go)

**说明**  
1. 传入的参数是指针，然后可以通过 pointer.Elem() 去获取所指向的 Value
2. 如果传入的参数不是指针，而是变量，那么  
*  通过 Elem 获取原始值对应的对象则直接 panic
*  通过 CanSet 方法查询是否可以设置返回 false
3. newValue.CantSet()表示是否可以重新设置其值，如果输出的是true则可修改，否则不能修改
4. reflect.Value.Elem() 表示获取原始值对应的反射对象，只有原始对象才能修改，当前反射对象是不能修改的
5. 也就是说如果要修改反射类型对象，其值必须是“addressable”【对应的要传入的是指针，同时要通过Elem方法获取原始值对应的反射对象】
6. struct 或者 struct 的嵌套都是一样的判断处理方式

### 通过 reflect.ValueOf 来进行方法的调用

[example5](./example5.go)

**说明**   
1. 要通过反射来调用起对应的方法，必须要先通过reflect.ValueOf(interface)来获取到reflect.Value
2. reflect.Value.MethodByName这.MethodByName，需要指定准确真实的方法名字，如果错误将直接panic，MethodByName返回一个函数值对应的reflect.Value方法的名字
3. []reflect.Value，这个是最终需要调用的方法的参数，可以没有或者一个或者多个，根据实际参数来定。
4. reflect.Value的 Call 这个方法，这个方法将最终调用真实的方法，参数务必保持一致，如果reflect.Value'Kind不是一个方法，那么将直接panic。


## Golang 的 反射 reflect 性能

1. 涉及到内存分配以及后续的GC；
2. reflect实现里面有大量的枚举，也就是for循环，比如类型之类的

所以反射很慢

