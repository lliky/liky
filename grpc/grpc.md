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



