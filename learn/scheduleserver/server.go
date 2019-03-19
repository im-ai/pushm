package main

//数据包的类型
const (
	HEART_BEAT_PACKET = 0x00
	REPORT_PACKET     = 0x01
)

var (
	server = ":8001"
)

//这里是包的结构体，其实是可以不需要的
type Packet struct {
	PacketType    byte
	PacketContent []byte
}

//心跳包，这里用了json来序列化，也可以用github上的gogo/protobuf包
//具体见(https://github.com/gogo/protobuf)
type HeartPacket struct {
	Version   string `json:"version"`
	Timestamp int64  `json:"timestamp"`
}
