# GRPC

## protoc

它是Protocol Buffers 的编译器，用于将 .proto 文件编译为不同编程语言的代码。

### 安装

1. 可以前往 [GitHub](https://github.com/protocolbuffers/protobuf/releases) ，选择版本安装；将 文件解压放到可执行路径就行

   ```shell
   $ unzip protoc-25.1-linux-x86_64.zip -d /usr/local/protoc
   $ export PATH="$PATH:/usr/local/protoc/bin"
   ```

2. 命令安装

   ```shell
   $ sudo apt-get update
   $ sudo apt-get install -y protobuf-compiler
   ```

3. 最后查看版本

   ```shell
   $ protoc --version
   ```

   

### 参数

- -I：指定 .proto 文件的搜索目录
- `--cpp_out`：生产 C++ 代码
- `--java_out`：生成 Java 代码
- `--python_out`：生成 Python 代码
- `--csharp_out`：生成 C# 代码
- **`--go_out`：生成 Go 代码**
- **`--go-grpc_out`：生成 gRPC 代码**
- `--plugin`：指定插件
- **`--proto_path`：指定 `.proto` 文件的搜索路径**
- `--descriptor_set_out`：生成描述符集文件
- `--include_imports`：包含导入的 `.proto` 文件

### protoc-gen-go & protoc-gen-go-grpc

生成 go 代码要依赖这两个插件

- protoc-gen-go：用于从 .proto 文件生成 Go 语言代码
- protoc-gen-grpc：用于从 .proto 文件生成 gRPC 服务的 Go 代码。

```shell
$ go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
$ go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2	
$ export PATH="$PATH:$(go env GOPATH)/bin"
```

总结：protoc 是主要的 Protocol Buffers 编译器，protoc--gen-go 和 protoc-gen-grpc 是与 protoc 一起使用的插件，分别从  .proto 文件生成 Go 代码 和 gRPC 代码。

### 用法

```shell
$ protoc --proto_path=proto proto/*.proto  --go_out=pb --go-grpc_out=pb
```

- `--proto_path=proto`：指定 .proto 文件搜索路径，可选的
- `proto/*.proto`：要编译的 .proto 文件文件，通配符说明编译所有文件，可以指定单个文件
- `--go_out=pb`：生成的 Go 语言代码的输出路径 path1，需要和 .proto 文件里面的 `option go_package="path2;pacakge";`搭配使用，最终代码路径是 `path1 + path2`
- `--go-grpc_out=pb`: 生成 gRPC 代码，路径和go_out 一样



## protobuf 语法

Protocol buffers 是 Google 的语言中立、平台中立、可扩展的结构化数据序列化。我们以 [proto3](https://protobuf.dev/programming-guides/proto3/)  为例。

### 基本规范

文件以 **.proto** 作为文件后缀。

除结构定义外，其他语句以分号结尾

rpc 方法定义结尾的分号可有可无

结构定义：message、service、enum

message 命名采用驼峰命名方式，字段采用小写字母 + 下划线

enum 命名采用驼峰命名方式，字段采用大写字母 + 下划线

service 和 rpc 方法名统一采用驼峰命名

```protobuf
syntax="proto3";

enum Unit {
	UNKNOWN = 0;
	BIT = 1;
	BYTE = 2;
}
message RateLaptopRequest {
	string laptop_id = 1;
}
message RateLaptopResponse {
	string rated_count = 1;
	Unit unit = 2;
}

service LaptopService {
	rpc RateLaptop(RateLaptopRequest)returns(RateLaptopResponse){};  // 可有可无
}
```

### 字段规则

字段格式：**限定修饰符|数据类型|字段名称| = |字段编码值| **

限定修饰符：

- optional: 
  - 没有设置：返回默认值，不会被序列化
  - 设置：显式设置，会被序列化，也可解析出来
- repeated: 可以重复 0 次或者多次，会保留重复的次序。我认为就是数组
- map: 键值对

数据类型：

| .proto tpye | Notes                                                  | Go type |
| ----------- | ------------------------------------------------------ | ------- |
| double      |                                                        | float64 |
| float       |                                                        | float32 |
| int32       | 使用可变长度编码。负数效率低，可以用 sint32            | int32   |
| int64       | 使用可变长度编码。负数效率低，可以用 sint64            | int64   |
| uint32      | 使用可变长度编码。                                     | uint32  |
| uint64      | 使用可变长度编码。                                     | uint64  |
| sint32      | 使用可变长度编码。可以有效的编码负数                   | int32   |
| sint64      | 使用可变长度编码。可以有效的编码负数                   | int64   |
| fixed32     | 始终是 4 字节，值如果大于 2 的 28 次方，比 uint32 有效 | uint32  |
| fixed64     | 始终是 8 字节，值如果大于 2 的 56 次方，比 uint64 有效 | uint64  |
| sfixed32    | 始终是 4 字节                                          | int32   |
| sfixed64    | 始终是 8 字节                                          | int64   |
| bool        |                                                        | bool    |
| string      | 字符串包含 utf-8 编码，长度不超过 2 的 32 次方         | string  |
| bytes       | 任意直接序，长度不超过 2 的 32 次方                    | []byte  |

字段名称：建议采用下划线分割

字段编码值：

* 每个编号值介于 1 - 536,870,911。

* **给定的编号在该消息中必须唯一**。

* 值 19000 - 19999 保留字段号，不可用。
* **一旦字段被使用，就不要更改**
* 建议使用 1 - 15 编号，需要一个字节去编码



### 枚举

编号必须从 0 开始，因为可以作为默认值



### 导入 .proto

```protobuf
import "myproject/other_protos.proto";
```

### service 定义

```protobuf
service LaptopService {
  // unary RPC
  rpc CreateLaptop(CreateLaptopRequest) returns (CreateLaptopResponse) {};
  // server streaming PRC
  rpc SearchLaptop( SearchLaptopRequest) returns (stream SearchLaptopResponse) {};
  // client streaming RPC
  rpc UploadImage(stream UploadImageRequest) returns (UploadImageResponse) {};
  //  bidirectional streaming RPC
  rpc RateLaptop(stream RateLaptopRequest) returns (stream RateLaptopResponse) {};
}
```

可以定义 4 中 RPC

### message 定义

```protobuf
message SearchResponse {
  repeated Result results = 1; // Get result
}

message Result {
  string url = 1;  
  string title = 2;
  repeated string snippets = 3;
}

message Outer {                // Level 0
        message MiddleAA {        // Level 1
            message Inner {        // Level 2
                int64 ival = 1;
                bool  booly = 2;
            }
        }
        message MiddleBB {         // Level 1
            message Inner {         // Level 2
                int32 ival = 1;
                bool  booly = 2;
            }
        }
    }
```

可以通过 // 去注释， 嵌套定义

### Map 类型

```
map<key_type, value_type> map_field = N;
```

- 键、值类型可以是内置的类型，也可以是自定义message类型
- 字段不支持repeated属性



### Oneof

如果消息包含多个字段，但是最多设置一个字段，这样就可以省内存

```protobuf
message UploadImageRequest {
  oneof data {
    ImageInfo info = 1;
    bytes chunk_data = 2;
  }
}
```

会覆盖，按照定义顺序，得到最后一个。



## gRPC 反射 evans CLI

#### 安装

可以前往 [GitHub](https://github.com/ktr0731/evans/releases) ，选择版本安装；将 文件解压放到可执行路径就行

```shell
$ tar -xzvf evans_linux_amd64.tar.gz -C /usr/local/evans/
$ export PATH="$PATH:/usr/local/evans"
```

#### 登录

```shell
$ evans -r repl -p 8080
# 检查 gRPC 服务器有哪些服务
$ evans -p 8080 -r cli list
```

#### 关键字

show、service、message、desc、call

##### show

```shell
127.0.0.1:8080> show service
+---------------+--------------+---------------------+----------------------+
|    SERVICE    |     RPC      |    REQUEST TYPE     |    RESPONSE TYPE     |
+---------------+--------------+---------------------+----------------------+
| AuthService   | Login        | LoginRequest        | LoginResponse        |
| LaptopService | CreateLaptop | CreateLaptopRequest | CreateLaptopResponse |
| LaptopService | SearchLaptop | SearchLaptopRequest | SearchLaptopResponse |
| LaptopService | UploadImage  | UploadImageRequest  | UploadImageResponse  |
| LaptopService | RateLaptop   | RateLaptopRequest   | RateLaptopResponse   |
+---------------+--------------+---------------------+----------------------+
```

##### service

```shell
127.0.0.1:8080> service AuthService

AuthService@127.0.0.1:8080> 
```

##### desc

```shell
AuthService@127.0.0.1:8080> show message
+----------------------+
|       MESSAGE        |
+----------------------+
| CreateLaptopRequest  |
| CreateLaptopResponse |
| LoginRequest         |
| LoginResponse        |
| RateLaptopRequest    |
| RateLaptopResponse   |
| SearchLaptopRequest  |
| SearchLaptopResponse |
| UploadImageRequest   |
| UploadImageResponse  |
+----------------------+

AuthService@127.0.0.1:8080> desc LoginRequest
+----------+-------------+----------+
|  FIELD   |    TYPE     | REPEATED |
+----------+-------------+----------+
| password | TYPE_STRING | false    |
| username | TYPE_STRING | false    |
+----------+-------------+----------+
```



##### call a RPC

```shell
127.0.0.1:8080> service AuthService

AuthService@127.0.0.1:8080> call Login
username (TYPE_STRING) => admin1
password (TYPE_STRING) => secret
{
  "accessToken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MDU2NDc4MDgsInVzZXJuYW1lIjoiYWRtaW4xIiwicm9sZSI6ImFkbWluIn0.9fC-ThbLHbAOYl1qikXPmnFsgiLNn8FM9R2IIA4oMjI"
}
```



## 4 种模式

### Simple RPC( Unary RPC)

#### 如何定义

```protobuf
syntax = "proto3";
option go_package = "./;pb";

message CreateLaptopRequest {
  Laptop laptop = 1;
}

message CreateLaptopResponse {
  string id = 1;
}

service LaptopService {
  rpc CreateLaptop(CreateLaptopRequest) returns (CreateLaptopResponse) {};
}
```

#### 代码实现

server side 参考 [CreateLaptop](./service/laptop_server.go)

client side 参考 [ CreateLaptop](./client/laptop_client.go)

### Server Streaming RPC

#### 如何定义

需要在前面加关键字 **stream**

```protobuf
syntax = "proto3";
option go_package = "./;pb";

message SearchLaptopRequest {
  Filter filter = 1;
}

message SearchLaptopResponse {
 Laptop laptop = 1;
}
service LaptopService {
  rpc SearchLaptop( SearchLaptopRequest) returns (stream SearchLaptopResponse) {};
}
```

#### 代码实现

server side 参考 [SearchLaptop](./service/laptop_server.go)

client side 参考 [ SearchLaptop](./client/laptop_client.go)

### Client Streaming RPC

#### 如何定义

```protobuf
syntax = "proto3";
option go_package = "./;pb";

message UploadImageRequest {
  oneof data {
    ImageInfo info = 1;
    bytes chunk_data = 2;
  }
}

message ImageInfo {
  string laptop_id = 1;
  string image_type = 2;
}

message UploadImageResponse {
  string id = 1;
  uint32 size = 2;
}
service LaptopService {
  rpc UploadImage(stream UploadImageRequest) returns (UploadImageResponse) {};
}
```

#### 代码实现

server side 参考 [UploadImage](./service/laptop_server.go)

client side 参考 [ UploadImage](./client/laptop_client.go)

### Bidirectional Streaming RPC

#### 如何定义

```protobuf
syntax = "proto3";
option go_package = "./;pb";

message RateLaptopRequest {
  string laptop_id = 1;
  double score = 2;
}

message RateLaptopResponse {
  string laptop_id = 1;
  uint32 rated_count = 2;
  double average_score = 3;
}
service LaptopService {
  rpc RateLaptop(stream RateLaptopRequest) returns (stream RateLaptopResponse) {};
}
```

#### 代码实现

server side 参考 [RateLaptop](./service/laptop_server.go)

client side 参考 [ RateLaptop](./client/laptop_client.go)

## GRPC 的其他特性

### Interceptor

### Deadlines

Deadlines 和 timeout 在分布式系统中，是常用的两种模式，它们都是用于控制程序和操作的操作时间，有不同的含义和应用场景

#### deadlines

指某个操作的最后期限或者截至时间，即操作必须在这个时间之前完成，否则会被取消或进行相应的处理。

#### Timeout

是执行某个操作的最长允许时间，超过这个时间就被会取消或者超时处理。

#### 应用场景

timeout 通常用于控制网络请求、IO读取或者其他超时处理的场景。比如 HTTP 客户端请求超时；deadline 用于需要明确指定操作截止时间的场景，比如分布式 RPC 调用，防止运行时间过长或停滞操作影响整个系统。

#### 代码

```go
// client
ctx, cancel := context.WithDeadline(context.Background(), 5*time.Second)
	defer cancel()
rsp, err := laptopClient.service.CreateLaptop(ctx, req)

// server
func (server *LaptopServer) CreateLaptop(ctx context.Context, req *pb.CreateLaptopRequest) (*pb.CreateLaptopResponse, error) {
	...
	//some heavy processsing
	//time.Sleep(6 * time.Second)
    if context.DeadlineExceeded == ctx.Err() {
		return nil, err
	}
	...
	return res, nil
}
```

在 client 和 server 都可以独立的判断 RPC 是否成功，那么可能会导致他们之间的结论一致。当 client 已经 DeadlineExceeded 的时候，server 可能还在继续响应，所以需要在 server  加判断。

### Cancellation

### Error Handling

当调用 gRPC 时，客户端会收到成功状态或者相对应的错误代码；gRPC 定义了一组 status Codes

```go
// A Code is an unsigned 32-bit error code as defined in the gRPC spec.
type Code uint32

const (
	// OK is returned on success.
	OK Code = 0

	// Canceled indicates the operation was canceled (typically by the caller).
	Canceled Code = 1

	// Unknown error.
	Unknown Code = 2

	// InvalidArgument indicates client specified an invalid argument.
	InvalidArgument Code = 3

	// The deadline expired before the opteration could complete.
	DeadlineExceeded Code = 4

	// NotFound means some requested entity (e.g., file or directory) was not found.
	NotFound Code = 5

	// AlreadyExists means an attempt to create an entity failed because one already exists.
	AlreadyExists Code = 6

	// PermissionDenied indicates the caller does not have permission to execute the specified operation.
	PermissionDenied Code = 7

	// ResourceExhausted indicates some resource has been exhauste.
	ResourceExhausted Code = 8

	// FailedPrecondition indicates operation was rejected because the
	// system is not in a state required for the operation's execution.
	FailedPrecondition Code = 9

	// Aborted indicates the operation was aborted, typically due to a
	// concurrency issue like sequencer check failures, transaction aborts,
	// etc.
	Aborted Code = 10

	// OutOfRange means operation was attempted past the valid range.
	OutOfRange Code = 11

	// Unimplemented indicates operation is not implemented or not supported/enabled in this service.
	Unimplemented Code = 12

	// Internal errors.
	Internal Code = 13

	// Unavailable indicates the service is currently unavailable.
	Unavailable Code = 14

	// DataLoss indicates unrecoverable data loss or corruption.
	DataLoss Code = 15

	// Unauthenticated indicates the request does not have valid
	// authentication credentials for the operation.
	Unauthenticated Code = 16

	_maxCode = 17
)
```

最主要就是 New()

```go
// New returns a Status representing c and msg.
func New(c codes.Code, msg string) *Status {
	return status.New(c, msg)
}

// Error returns an error representing c and msg.  If c is OK, returns nil.
func Error(c codes.Code, msg string) error {
	return New(c, msg).Err()
}
```



### Multiplexing

### Metadata

### Load Balancing

### TLS
