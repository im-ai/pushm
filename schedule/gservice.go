package main

import (
	"encoding/json"
	"fmt"
	"net"
)

//在这里处理收到的包，就和一般的逻辑一样了，根据类型进行不同的处理，因人而异
//我这里处理了心跳和一个上报数据包
//服务器往客户端的数据包很简单地以\n换行结束了，偷了一个懒:)，正常情况下也可根据自己的协议来封装好
//然后在客户端写一个状态来处理
func ProcessRecvData(packet *Packet, conn net.Conn) {
	switch packet.PacketType {
	case HEART_BEAT_PACKET:
		var beatPacket HeartPacket
		json.Unmarshal(packet.PacketContent, &beatPacket)

		config := GetConfigChange()
		bytesa, e := json.Marshal(config)
		if e != nil {
			fmt.Println(e)
			return
		}
		if config.Number == 0 {
			bytesCombineInit = BytesCombine(bytesa, []byte("\n"))
			conn.Write(bytesCombineInit)

			if beatPacket.Gonumber > 0 {
				config.Number = -beatPacket.Gonumber
				beatPacket.Gonumber = 0
			}

			bytesa, e := json.Marshal(config)
			if e != nil {
				fmt.Println(e)
				return
			}

			bytesCombineInit = BytesCombine(bytesa, []byte("\n"))
			conn.Write(bytesCombineInit)

			beatPacket.Responsetime = 0
			beatPacket.Responsemaxtime = 0
		}

		if beatPacket.Gonumber < 0 {
			beatPacket.Gonumber = 0
		}
		if beatPacket.Responsetime < 0 {
			beatPacket.Responsetime = 0
		}
		if beatPacket.Responsemaxtime < 0 {
			beatPacket.Responsemaxtime = 0
		}
		goroutinemap <- beatPacket.Gonumber
		goresptimemap <-  beatPacket.Responsetime
		gorespmaxtimemap <- beatPacket.Responsemaxtime

		if goroutineflag > goclientnumber {
			goroutineflag = 0
			goroutinenumber = 0
		}

		goroutineflag++
		goroutinenumber = goroutinenumber+beatPacket.Gonumber

		if goroutinenumber > gonumber {
			return
		}
		fmt.Println("send message:", string(bytesCombine))
		conn.Write(bytesCombine)
		return
	case REPORT_PACKET:
		var reportPacket ReportPacket
		json.Unmarshal(packet.PacketContent, &reportPacket)
		conn.Write([]byte("Report data has recive\n"))
		return
	}
}
