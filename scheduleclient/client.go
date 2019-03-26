package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/jander/golog/logger"
	"github.com/kardianos/service"
	"github.com/panjf2000/ants"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
	"log"
	"math/rand"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strconv"
	"time"
)

//默认的服务器地址
var (
	serveradd = "192.168.1.70:9876"
	gonumber  = 0
)

//客户端对象
type TcpClient struct {
	connection *net.TCPConn
	hawkServer *net.TCPAddr
	stopChan   chan struct{}
	netStop    int
}

type program struct{}

func (p *program) Start(s service.Service) error {
	go p.run()
	return nil
}

func (p *program) run() {
	// 代码写在这儿

	//拿到服务器地址信息
INIT:
	hawkServer, err := net.ResolveTCPAddr("tcp", serveradd)
	if err != nil {
		fmt.Printf("hawk server [%s] resolve error: [%s]\n", serveradd, err.Error())
		fmt.Println("重试连接服务器...")
		time.Sleep(10 * time.Second)
		goto INIT
	}
	//连接服务器
	connection, err := net.DialTCP("tcp", nil, hawkServer)
	if err != nil {
		fmt.Printf("connect to hawk server error: [%s]\n", err.Error())
		fmt.Println("重试连接服务器...")
		time.Sleep(10 * time.Second)
		goto INIT
	}
	client := &TcpClient{
		connection: connection,
		hawkServer: hawkServer,
		stopChan:   make(chan struct{}),
		netStop:    0,
	}

	//启动接收
	go client.receivePackets()

	//发送心跳的goroutine
	go func() {
		heartBeatTick := time.Tick(1 * time.Second)
		for {
			select {
			case <-heartBeatTick:
				client.sendHeartPacket()
				if client.netStop == 1 {
					fmt.Println("net stop common !")
					return
				}
			case <-client.stopChan:
				break
			}
		}
	}()
	<-client.stopChan
	fmt.Println("重试连接服务器...")
	time.Sleep(10 * time.Second)
	goto INIT
	//测试用的，开300个goroutine每秒发送一个包
	//for i := 0; i < 300; i++ {
	//	go func() {
	//		sendTimer := time.After(1 * time.Second)
	//		for {
	//			select {
	//			case <-sendTimer:
	//				client.sendReportPacket()
	//				sendTimer = time.After(1 * time.Second)
	//			case <-client.stopChan:
	//				return
	//			}
	//		}
	//	}()
	//}
	//等待退出
}

func (p *program) Stop(s service.Service) error {
	return nil
}

func main() {

	rand.Seed(time.Now().UnixNano())
	itoa := strconv.Itoa(rand.Intn(100))
	svcConfig := &service.Config{
		Name:        "scheduleclient" + itoa, //服务显示名称
		DisplayName: "scheduleclient" + itoa, //服务名称
		Description: "scheduleclient" + itoa, //服务描述

	}
	prg := &program{}
	s, err := service.New(prg, svcConfig)
	if err != nil {
		logger.Fatal(err)
	}
	if err != nil {
		logger.Fatal(err)
	}

	if len(os.Args) > 1 {
		if os.Args[1] == "install" {
			s.Install()
			logger.Println("服务安装成功")
			return
		}
		if os.Args[1] == "remove" {
			s.Uninstall()
			logger.Println("服务卸载成功")
			return
		}
	}

	err = s.Run()

	if err != nil {
		logger.Error(err)
	}

}

// 接收数据包
func (client *TcpClient) receivePackets() {
	pool, _ := ants.NewPool(200000)
	defer pool.Release()

	reader := bufio.NewReader(client.connection)
	for {
		//承接上面说的服务器端的偷懒，我这里读也只是以\n为界限来读区分包
		msg, err := reader.ReadString('\n')
		if err != nil {
			//在这里也请处理如果服务器关闭时的异常
			fmt.Println("server close \n")
			close(client.stopChan)
			client.netStop = 1
			break
		}
		fmt.Print(msg)
		var pbody PressureBody
		_ = json.Unmarshal([]byte(msg), &pbody)
		v, _ := mem.VirtualMemory()
		if v.UsedPercent > 80.0 {
			continue
		}
		cc, _ := cpu.Percent(time.Second, false)
		if cc[0] > 80.0 {
			continue
		}
		gonumber = gonumber + pbody.Number
		for i := 0; i < pbody.Number; i++ {
			if pbody.TypeId == 3 {
				_ = pool.Submit(func() {
					go connectWs(pbody.Url, client)
				})
			} else if pbody.TypeId == 2 {
				_ = pool.Submit(func() {
					connectHttpPost(pbody.Url, pbody.Json)
				})
			} else if pbody.TypeId == 1 {
				_ = pool.Submit(func() {
					connectHttpGet(pbody.Url)
				})
			}
		}
	}
}

//发送数据包
//仔细看代码其实这里做了两次json的序列化，有一次其实是不需要的
func (client *TcpClient) sendReportPacket() {
	reportPacket := ReportPacket{
		Content:   getRandString(),
		Timestamp: time.Now().Unix(),
		Rand:      rand.Int(),
	}
	packetBytes, err := json.Marshal(reportPacket)
	if err != nil {
		fmt.Println(err.Error())
	}
	//这一次其实可以不需要，在封包的地方把类型和数据传进去即可
	packet := Packet{
		PacketType:    REPORT_PACKET,
		PacketContent: packetBytes,
	}
	sendBytes, err := json.Marshal(packet)
	if err != nil {
		fmt.Println(err.Error())
	}
	//发送
	client.connection.Write(EnPackSendData(sendBytes))
	fmt.Println("Send metric data success!")
}

//发送心跳包，与发送数据包一样
func (client *TcpClient) sendHeartPacket() {
	heartPacket := HeartPacket{
		Version:   "1.0",
		Timestamp: time.Now().Unix(),
		Gonumber:  gonumber,
	}
	packetBytes, err := json.Marshal(heartPacket)
	if err != nil {
		fmt.Println(err.Error())
	}
	packet := Packet{
		PacketType:    HEART_BEAT_PACKET,
		PacketContent: packetBytes,
	}
	sendBytes, err := json.Marshal(packet)
	if err != nil {
		fmt.Println(err.Error())
	}
	client.connection.Write(EnPackSendData(sendBytes))
	fmt.Println("Send heartbeat data success!")
}

//拿一串随机字符
func getRandString() string {
	length := rand.Intn(50)
	strBytes := make([]byte, length)
	for i := 0; i < length; i++ {
		strBytes[i] = byte(rand.Intn(26) + 97)
	}
	return string(strBytes)
}

func connectHttpGet(url string) {
	fmt.Println("GET:" + url)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()
	defer func() {
		gonumber--
	}()
}

func connectHttpPost(url, json string) {
	fmt.Println("POST:" + url)
	var jsonStr = []byte(json)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()
	defer func() {
		gonumber--
	}()
}
func connectWs(urls string, client *TcpClient) {
	fmt.Println("ws:" + urls)
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u := url.URL{Scheme: "ws", Host: urls, Path: "/ws"}
	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	done := make(chan struct{})

	defer func() {
		gonumber--
	}()
	// 读取数据
	go func() {
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			if client.netStop == 1 {
				log.Println("stop common")
				return
			}

			v, _ := mem.VirtualMemory()
			if v.UsedPercent > 80.0 {
				return
			}
			cc, _ := cpu.Percent(time.Second, false)
			if cc[0] > 80 {
				return
			}
			log.Printf("recv: %s", message)
		}
	}()

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-done:
			return
		case t := <-ticker.C:
			if client.netStop == 1 {
				log.Println("stop common")
				return
			}

			v, _ := mem.VirtualMemory()
			if v.UsedPercent > 80.0 {
				return
			}
			cc, _ := cpu.Percent(time.Second, false)
			if cc[0] > 80 {
				return
			}

			err := c.WriteMessage(websocket.TextMessage, []byte(t.String()))
			if err != nil {
				log.Println("write:", err)
				return
			}
		case <-interrupt:
			log.Println("interrupt")

			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
				return
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return
		}
	}
}
