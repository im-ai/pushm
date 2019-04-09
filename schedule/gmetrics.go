package main

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/shirou/gopsutil/mem"
	"log"
	"net/http"
	"os"
	"time"
)

func InitMetrics() {
	//初始化日志服务
	logger := log.New(os.Stdout, "[Memory]", log.Lshortfile|log.Ldate|log.Ltime)

	//初始一个http handler
	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/conf", Changeconf)
	//初始化一个容器
	diskPercent := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "go_memeory_percent",
		Help: "memeory use percent",
	},
		[]string{"percent"},
	)
	goroutineCount := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "go_open_goroutine_count",
		Help: " go open goroutine number",
	},
		[]string{"number"},
	)
	responseTime := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "go_average_response_time",
		Help: " go average response time",
	},
		[]string{"number"},
	)
	responseMaxTime := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "go_max_response_time",
		Help: " go max response time",
	},
		[]string{"number"},
	)
	prometheus.MustRegister(diskPercent)
	prometheus.MustRegister(goroutineCount)
	prometheus.MustRegister(responseTime)
	prometheus.MustRegister(responseMaxTime)

	// 启动web服务，监听1010端口
	go func() {
		logger.Println("ListenAndServe at:192.168.1.70:1010")
		err := http.ListenAndServe(":1010", nil)
		if err != nil {
			logger.Fatal("ListenAndServe: ", err)
		}
	}()

	//收集内存使用的百分比
	for {
		//logger.Println("start collect memory used percent!")
		v, err := mem.VirtualMemory()
		if err != nil {
			logger.Println("get memeory use percent error:%s", err)
		}
		usedPercent := v.UsedPercent
		fmt.Println("get memeory use percent:", usedPercent)
		diskPercent.WithLabelValues("usedMemory").Set(usedPercent)

		// open goroutine size
		tmp := 0
		for _, vnum := range goroutinemap {
			tmp += vnum
		}
		fmt.Println("get open goroutine number :", tmp)
		goroutineCount.WithLabelValues("openGoroutineNumber").Set(float64(tmp))

		// 平均响应时间
		tmptime := 0.0
		for _, vnum := range goresptimemap {
			tmptime += vnum
		}
		ilen := len(goresptimemap)
		if ilen > 0 {
			tmptime = tmptime / float64(ilen)
		}
		fmt.Println("go average response time :", tmptime)
		responseTime.WithLabelValues("averageResponTime").Set(tmptime)

		// 最大响应时间
		tmpmaxtime := 0.0
		for _, vnum := range gorespmaxtimemap {
			if vnum > tmpmaxtime {
				tmpmaxtime = vnum
			}
		}
		fmt.Println("go average response time :", tmpmaxtime)
		responseMaxTime.WithLabelValues("maxResponTime").Set(tmpmaxtime)

		fmt.Println("")
		time.Sleep(time.Second * 2)
	}

}
