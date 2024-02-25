date: 2024-02-26

goversion:  go1.21.6



# 1. 核心数据结构

```go
type hchan struct {
	qcount   uint           // total data in the queue
	dataqsiz uint           // size of the circular queue
	buf      unsafe.Pointer // points to an array of dataqsiz elements
	elemsize uint16
	closed   uint32
	elemtype *_type // element type
	sendx    uint   // send index
	recvx    uint   // receive index
	recvq    waitq  // list of recv waiters
	sendq    waitq  // list of send waiters

	lock mutex
}
```

hcan：channel 数据结构

* qcount：当前 channel 存在多少个元素
* dataqsiz： 当前 channel 能够存放元素数量
* buf：存放元素的环形队列
* elemsize：元素类型大小
* closed：标识 channel 是否关闭
* elemtype：channel 元素类型
* sendx：发送元素环形队列的 index
* recvx：接收元素环形队列的 index
* sendq：因发送而陷入阻塞的协程队列
* recvq：因接受而陷入阻塞的协程队列
