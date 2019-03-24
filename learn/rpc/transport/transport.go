package transport

import (
	"io"
	"net"
)

//传输层的定义，用于读取数据
type Transport interface {
	Dial(network, addr string) error
	//这里直接内嵌了ReadWriteCloser接口，包含Read、Write和Close方法
	io.ReadWriteCloser
	RemoteAddr() net.Addr
	LocalAddr() net.Addr
}

//服务端监听器定义，用于监听端口和建立连接
type Listener interface {
	Listen(network, addr string) error
	Accept() (Transport, error)
	//这里直接内嵌了Closer接口，包含Close方法
	io.Closer
}
