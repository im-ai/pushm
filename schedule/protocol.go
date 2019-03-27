package main

import (
	"fmt"
	"hash/crc32"
)

//数据包类型
const (
	HEART_BEAT_PACKET = 0x00
	REPORT_PACKET     = 0x01
)

//数据包
type Packet struct {
	PacketType    byte
	PacketContent []byte
}

//心跳包
type HeartPacket struct {
	Version         string  `json:"version"`
	Timestamp       int64   `json:"timestamp"`
	Gonumber        int     `json:"gonumber"`
	Responsetime    float64 `json:"responsetime"`
	Responsemaxtime float64 `json:"responsemaxtime"` //  最大响应时间
}

//数据包
type ReportPacket struct {
	Content   string `json:"content"`
	Rand      int    `json:"rand"`
	Timestamp int64  `json:"timestamp"`
}

//使用的协议与服务器端保持一致
func EnPackSendData(sendBytes []byte) []byte {
	packetLength := len(sendBytes) + 8
	result := make([]byte, packetLength)
	result[0] = 0xFF
	result[1] = 0xFF
	result[2] = byte(uint16(len(sendBytes)) >> 8)
	result[3] = byte(uint16(len(sendBytes)) & 0xFF)
	copy(result[4:], sendBytes)
	sendCrc := crc32.ChecksumIEEE(sendBytes)
	result[packetLength-4] = byte(sendCrc >> 24)
	result[packetLength-3] = byte(sendCrc >> 16 & 0xFF)
	result[packetLength-2] = 0xFF
	result[packetLength-1] = 0xFE
	fmt.Println(result)
	return result
}

type PressureBody struct {
	TypeId int    // 1: http get 2: http post  3: ws
	Url    string // 请求 url
	Json   string // post参数
	Number int    // 每秒开启 gorouting 次数
}
