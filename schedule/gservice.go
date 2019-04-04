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
		fmt.Printf("recieve heat beat from [%s] ,data is [%v]\n", conn.RemoteAddr().String(), beatPacket)
		_, ok := goroutinemap[conn.RemoteAddr().String()]
		if !ok {
			config := GetConfig()
			fmt.Println("Init remote client gonumber ", conn.RemoteAddr().String(), " goroutineNumber: ", config.Number)
			bytesa, e := json.Marshal(config)
			if e != nil {
				fmt.Println(e)
				return
			}
			bytesCombineInit = BytesCombine(bytesa, []byte("\n"))
			conn.Write(bytesCombineInit)

			if config.Number == 0 {
				if beatPacket.Gonumber > 0 {
					config.Number = -beatPacket.Gonumber
				}
				bytesa, e := json.Marshal(config)
				if e != nil {
					fmt.Println(e)
					return
				}
				bytesCombineInit = BytesCombine(bytesa, []byte("\n"))
				conn.Write(bytesCombineInit)
			}
		}
		fmt.Println("RemoteAddr()", conn.RemoteAddr())
		fmt.Println("Gonumber:", beatPacket.Gonumber)
		goroutinemap[conn.RemoteAddr().String()] = beatPacket.Gonumber
		goresptimemap[conn.RemoteAddr().String()] = beatPacket.Responsetime
		gorespmaxtimemap[conn.RemoteAddr().String()] = beatPacket.Responsemaxtime

		if goroutinenumber > gonumber {
			fmt.Println("The maximum value has been reduced to goroutine  number:", gonumber)
			return
		}
		fmt.Println("send message:", string(bytesCombine))
		conn.Write(bytesCombine)
		return
	case REPORT_PACKET:
		var reportPacket ReportPacket
		json.Unmarshal(packet.PacketContent, &reportPacket)
		fmt.Printf("recieve report data from [%s] ,data is [%v]\n", conn.RemoteAddr().String(), reportPacket)
		conn.Write([]byte("Report data has recive\n"))
		return
	}
}
