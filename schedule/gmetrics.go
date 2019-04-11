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

	idx := 0
	totalnum := 0
	arrgnum :=make([]int,1000,1000)
	go func() {
		for{
			select {
			case gnum := <-goroutinemap:
				arrgnum[idx] = gnum
				if idx > goclientnumber {
					idx = 0
				}
				idx++
				fmt.Printf("arr %v\n",arrgnum)
				for i := 0; i< goclientnumber	;i++  {
					totalnum = totalnum+arrgnum[i]
				}

				fmt.Println("go open goroutine number ",totalnum)
				goroutineCount.WithLabelValues("openGoroutineNumber").Set(float64(totalnum))
			}
		}
	}()

	avgtimecfg := 0.0
	go func() {
		for{
			select {
			case avgtime := <-goresptimemap:
				avgtimecfg = (avgtimecfg+avgtime)/2
				fmt.Println("go average response time ",avgtimecfg)
				responseTime.WithLabelValues("averageResponTime").Set(avgtimecfg)
			}
		}
	}()

	maxtimecfg := 0.0
	go func() {
		for{
			select {
			case maxtime := <-gorespmaxtimemap:
				if maxtime > maxtimecfg{
					maxtimecfg = maxtime
				}
				fmt.Println("go max response time ", maxtimecfg)
				responseMaxTime.WithLabelValues("maxResponTime").Set(maxtimecfg)
			}
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

		fmt.Println("")
		time.Sleep(time.Second * 2)
	}

}
