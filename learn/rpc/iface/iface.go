package iface

import (
	"context"
)

type RPCServer interface {
	//注册服务实例，rcvr是receiver的意思，它是我们对外暴露的方法的实现者，
	// metaData是注册服务时携带的额外的元数据，它描述了rcvr的其他信息
	Register(rcvr interface{}, metaData map[string]string) error
	//开始对外提供服务
	Server(newwork string, add string) error
}
type Call struct {
	ServiceMethod string      //服务名.方法名
	Args          interface{} //参数
	Reply         interface{} //返回值（指针类型）
	Error         error       // 错误信息
	Done          chan *Call  // 在调用结束时激活
}

type RPCClient interface {
	// Go 表示异步调用
	Go(ctx context.Context, serviceMethod string, arg interface{}, reply interface{}, done chan *Call) *Call

	//Call表示异步调用
	Call(ctx context.Context, serviceMethod string, arg interface{}, reply interface{}) error
	Close() error
}
