package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"hash/crc32"
	"io"
	"io/ioutil"
	"net"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var (
	server       = ":9876"
	nubmer       = 0
	bytesCombine []byte
)

//与服务器相关的资源都放在这里面
type TcpServer struct {
	listener   *net.TCPListener
	hawkServer *net.TCPAddr
}

func main() {

	numberstr := getCfg("number", "sfig.ini")
	fmt.Println("number:", numberstr)
	nubmers, _ := strconv.Atoi(numberstr)
	nubmer = nubmers

	typeIdstr := getCfg("typeId", "sfig.ini")
	fmt.Println("typeId:", typeIdstr)
	typeId, _ := strconv.Atoi(typeIdstr)

	urlstr := getCfg("url", "sfig.ini")
	fmt.Println("url:", urlstr)

	jsonstr := getCfg("json", "sfig.ini")
	fmt.Println("json:", jsonstr)
	config := &PressureBody{
		TypeId: typeId,
		Url:    urlstr,
		Json:   jsonstr,
		Number: nubmer,
	}
	bytesa, e := json.Marshal(config)
	if e != nil {
		fmt.Println(e)
		return
	}
	bytesCombines := BytesCombine(bytesa, []byte("\n"))
	bytesCombine = bytesCombines

	//类似于初始化套接字，绑定端口
	hawkServer, err := net.ResolveTCPAddr("tcp", server)
	checkErr(err)
	//侦听
	listen, err := net.ListenTCP("tcp", hawkServer)
	checkErr(err)
	//记得关闭
	defer listen.Close()

	tcpServer := &TcpServer{
		listener:   listen,
		hawkServer: hawkServer,
	}
	fmt.Println("start server successful......")

	//开始接收请求
	for {
		conn, err := tcpServer.listener.Accept()
		fmt.Println("accept tcp client %s", conn.RemoteAddr().String())
		checkErr(err)

		go Handle(conn)
	}

}

//处理函数，这是一个状态机
//根据数据包来做解析
//数据包的格式为|0xFF|0xFF|len(高)|len(低)|Data|CRC高16位|0xFF|0xFE
//其中len为data的长度，实际长度为len(高)*256+len(低)
//CRC为32位CRC，取了最高16位共2Bytes
//0xFF|0xFF和0xFF|0xFE类似于前导码
func Handle(conn net.Conn) {
	defer conn.Close()
	//状态机状态
	state := 0x00
	//数据包长度
	length := uint16(0)
	// crc 检验和
	crc16 := uint16(0)
	var recvBuffer []byte
	//游标
	cursor := uint16(0)
	bufferReader := bufio.NewReader(conn)
	//状态机处理数据
	for {
		recvByte, err := bufferReader.ReadByte()
		if err != nil {
			//这里因为做了心跳，所以就没有加deadline时间，如果客户端断开连接
			//这里ReadByte方法返回一个io.EOF的错误，具体可考虑文档
			if err == io.EOF {
				fmt.Printf("client %s is close!\n", conn.RemoteAddr().String())
			}
			//在这里直接退出goroutine，关闭由defer操作完成
			return
		}
		//进入状态机，根据不同的状态来处理
		switch state {
		case 0x00:
			if recvByte == 0xFF {
				state = 0x01
				//初始化状态机
				recvBuffer = nil
				length = 0
				crc16 = 0
			} else {
				state = 0x00
			}
			break
		case 0x01:
			if recvByte == 0xFF {
				state = 0x02
			} else {
				state = 0x00
			}
			break
		case 0x02:
			length += uint16(recvByte) * 256
			state = 0x03
			break
		case 0x03:
			length += uint16(recvByte)
			// 一次申请缓存，初始化游标，准备读数据
			recvBuffer = make([]byte, length)
			cursor = 0
			state = 0x04
			break
		case 0x04:
			//不断地在这个状态下读数据，直到满足长度为止
			recvBuffer[cursor] = recvByte
			cursor++
			if cursor == length {
				state = 0x05
			}
			break
		case 0x05:
			crc16 += uint16(recvByte) * 256
			state = 0x06
			break
		case 0x06:
			crc16 += uint16(recvByte)
			state = 0x07
			break
		case 0x07:
			if recvByte == 0xFF {
				state = 0x08
			} else {
				state = 0x00
			}
		case 0x08:
			if recvByte == 0xFE {
				//执行数据包校验
				if (crc32.ChecksumIEEE(recvBuffer)>>16)&0xFFFF == uint32(crc16) {
					var packet Packet
					//把拿到的数据反序列化出来
					json.Unmarshal(recvBuffer, &packet)
					//新开协程处理数据
					go processRecvData(&packet, conn)
				} else {
					fmt.Println("丢弃数据!")
				}
			}
			//状态机归位,接收下一个包
			state = 0x00
		}
	}
}

//在这里处理收到的包，就和一般的逻辑一样了，根据类型进行不同的处理，因人而异
//我这里处理了心跳和一个上报数据包
//服务器往客户端的数据包很简单地以\n换行结束了，偷了一个懒:)，正常情况下也可根据自己的协议来封装好
//然后在客户端写一个状态来处理
func processRecvData(packet *Packet, conn net.Conn) {
	switch packet.PacketType {
	case HEART_BEAT_PACKET:
		var beatPacket HeartPacket
		json.Unmarshal(packet.PacketContent, &beatPacket)
		fmt.Printf("recieve heat beat from [%s] ,data is [%v]\n", conn.RemoteAddr().String(), beatPacket)
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

//处理错误，根据实际情况选择这样处理，还是在函数调之后不同的地方不同处理
func checkErr(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func getCfg(tag string, filepath string) string {
	dat, err := ioutil.ReadFile(filepath) //读取文件
	checkErr(err)                         //检查是否有错误
	cfg := string(dat)                    //将读取到达配置文件转化为字符串
	var str string
	s1 := fmt.Sprintf("[^;]%s *= *.{1,}\\n", tag)
	s2 := fmt.Sprintf("%s *= *", tag)
	reg, err := regexp.Compile(s1)
	if err == nil {
		tag_str := reg.FindString(cfg) //在配置字符串中搜索
		if len(tag_str) > 0 {
			r, _ := regexp.Compile(s2)
			i := r.FindStringIndex(tag_str) //查找配置字符串的确切起始位置
			var h_str = make([]byte, len(tag_str)-i[1])
			copy(h_str, tag_str[i[1]:])
			str1 := fmt.Sprintln(string(h_str))
			str2 := strings.Replace(str1, "\n", "", -1)
			str = strings.Replace(str2, "\r", "", -1)
		}
	}
	return str
}

//BytesCombine 多个[]byte数组合并成一个[]byte
func BytesCombine(pBytes ...[]byte) []byte {
	len := len(pBytes)
	s := make([][]byte, len)
	for index := 0; index < len; index++ {
		s[index] = pBytes[index]
	}
	sep := []byte("")
	return bytes.Join(s, sep)
}