package main

import "github.com/gogf/gf/g/database/gdb"

// 全局变量  gmain 统一定义
var (
	server           = ":9876"
	bytesCombine     []byte
	bytesCombineInit []byte
	goroutinenumber  = 0
	goroutinemap     = make(map[string]int)
	goresptimemap    = make(map[string]float64) //平均响应时间
	gorespmaxtimemap = make(map[string]float64) //最大响应时间
	nubmer           int
	gonumber         int
	db               gdb.DB
)

func main() {

	// 初始化配置
	cfgerr := InitCfg()
	if cfgerr != nil {
		return
	}

	// 启动 指标服务器
	go InitMetrics()

	// 启动控制页面
	go InitControl()
	InitControl()

	// 启动 tcp 分发服务器
	StartTcpServer()
}
