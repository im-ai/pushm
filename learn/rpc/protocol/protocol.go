package protocol

import "io"

type Header struct {
	Seq           uint64            //序号, 用来唯一标识请求或响应
	MessageType   byte              //消息类型，用来标识一个消息是请求还是响应
	CompressType  byte              //压缩类型，用来标识一个消息的压缩方式
	SerializeType byte              //序列化类型，用来标识消息体采用的编码方式
	StatusCode    byte              //状态类型，用来标识一个请求是正常还是异常
	ServiceName   string            //服务名
	MethodName    string            //方法名
	Error         string            //方法调用发生的异常
	MetaData      map[string]string //其他元数据
}

//Messagge表示一个消息体
type Message struct {
	*Header        //head部分, Header的定义参考上一篇文章
	Data    []byte //body部分
}

//Protocol定义了如何构造和序列化一个完整的消息体
type Protocol interface {
	NewMessage() *Message
	DecodeMessage(r io.Reader) (*Message, error)
	EncodeMessage(message *Message) []byte
}
