# HTTP 源码分析

## Server 端

以一个简单例子入手：

```go
package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/ping", func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprintf(writer, "pong")
	})
	http.ListenAndServe(":8088", nil)
}
```

运行代码，在浏览器打开 **localhost:8088/ping**  就可以看到 **pong**，或者在终端 **curl localhost:8088/ping**，然后 http.ListenAndServe 开启监听。当请求过来，则根据路由执行对应的 handler 函数。

另外一种实现方式：

```go
package main

import (
	"fmt"
	"net/http"
)

type DemoHandler struct{
}

func (d *DemoHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	fmt.Fprintf(writer, "hello world")
}

func main() {
	http.Handle("/ping", &DemoHandler{})	
	http.ListenAndServe(":8088", nil)
}
```

### 注册路由

通过例子发现 **http.HandleFunc** 和 **http.Handle** 都是路由注册。

- http.HandleFunc 第二参数是一个具有 func(writer http.ResponseWriter, request *http.Request) 签名的函数；
- http.Handle 第二个参数是一个结构体，该结构体实现了 http.Handler 接口。  

```go
func HandleFunc(pattern string, handler func(ResponseWriter, *Request)) {
	DefaultServeMux.HandleFunc(pattern, handler)
}

// HandleFunc registers the handler function for the given pattern.
func (mux *ServeMux) HandleFunc(pattern string, handler func(ResponseWriter, *Request)) {
	if handler == nil {
		panic("http: nil handler")
	}
	mux.Handle(pattern, HandlerFunc(handler))
}
```

```go
func Handle(pattern string, handler Handler) { DefaultServeMux.Handle(pattern, handler) }
```

通过源码可以看出：两个注册函数最后都是通过 **DefaultServeMux** 和 **Handle** 方法来调用的

#### Handler

handler 是一个接口：

```go
type Handler interface {
	ServeHTTP(ResponseWriter, *Request)
}
```

该接口声明了一个 ServeHTTP 的函数签名，任何结构体实现了这个接口，那么结构体就是一个 Handler 对象。go 的 http 服务都是基于 Handler 处理的，而 Handler 对象的 ServeHTTP 方法也是处理 request 构建 response 的核心。

回到上面的函数 **HandleFunc**

```go
mux.Handle(pattern, HandlerFunc(handler))
```

