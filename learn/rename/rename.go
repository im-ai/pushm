package main

import (
	"bufio"
	"fmt"
	"github.com/henrylee2cn/pholcus/common/util"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

var (
	path string
)

func main() {
	zhibiaomap := initmap()
	fmt.Printf("请把目录拖进来: ")
	fmt.Scanln(&path)
	GetAllFile(path, zhibiaomap)

}

func GetAllFile(pathname string, zhibiaomap map[int]string) error {
	rd, err := ioutil.ReadDir(pathname)
	for _, fi := range rd {
		if fi.IsDir() {
			GetAllFile(pathname+"\\"+fi.Name()+"\\", zhibiaomap)
		} else {
			filePath := pathname + "\\" + fi.Name()
			readfileAndRename(pathname, filePath)
			//readfileAndRename2(pathname, filePath, zhibiaomap)

		}
	}
	return err
}

func readfileAndRename2(pathname, filePath string, zhibiaomap map[int]string) {

	f, err := os.Open(filePath)
	if err != nil {
		fmt.Println(err)
		return
	}
	buf := bufio.NewReader(f)
	bzname := ""
	flag := 0
	lendd := "0"
	preline := ""
	for {
		line, err := buf.ReadString('\n')
		line = strings.TrimSpace(line)
		flag++

		if line == "32767" {
			//bzname = Substr(bzname, 0, len(bzname)-1)
			lasidx := strings.LastIndex(bzname, "_")
			bzname = Substr(bzname, 0, lasidx)
			lendd = preline
			break
		}

		if flag > 2 {
			s := zhibiaomap[util.Atoi(line)]
			if s == "" {
				s = line
			}
			bzname = bzname + s + "_"
		}

		preline = line

		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println(err)
			break
		}
	}
	lasidx := strings.LastIndex(bzname, "_")
	bzname = Substr(bzname, 0, lasidx)

	fmt.Println(bzname)
	fmt.Println(lendd)
	now := time.Now()

	y := fmt.Sprintf("%d", now.Year())
	m := fmt.Sprintf("%d", now.Month())
	d := fmt.Sprintf("%d", now.Day())
	h := fmt.Sprintf("%d", now.Hour())
	mm := fmt.Sprintf("%d", now.Minute())
	s := fmt.Sprintf("%d", now.Second())

	if len(m) == 1 {
		m = "0" + m
	}
	if len(d) == 1 {
		d = "0" + d
	}
	if len(h) == 1 {
		h = "0" + h
	}
	if len(mm) == 1 {
		mm = "0" + mm
	}
	if len(s) == 1 {
		s = "0" + s
	}

	f.Close()

	sfm := y + m + d + h + mm + s
	newpath := pathname + "\\" + bzname + "-I-" + string(lendd) + "-0001-" + sfm + ".txt"
	fmt.Println(newpath)
	err2 := os.Rename(filePath, newpath)
	if err2 != nil {
		fmt.Println(err2)
	}
	time.Sleep(1 * time.Second)
}

func readfileAndRename(pathname, filePath string) {

	now := time.Now()
	y := fmt.Sprintf("%d", now.Year())
	m := fmt.Sprintf("%d", now.Month())
	d := fmt.Sprintf("%d", now.Day())
	h := fmt.Sprintf("%d", now.Hour())
	mm := fmt.Sprintf("%d", now.Minute())
	s := fmt.Sprintf("%d", now.Second())

	if len(m) == 1 {
		m = "0" + m
	}
	if len(d) == 1 {
		d = "0" + d
	}
	if len(h) == 1 {
		h = "0" + h
	}
	if len(mm) == 1 {
		mm = "0" + mm
	}
	if len(s) == 1 {
		s = "0" + s
	}

	sfm := y + m + d + h + mm + s
	newpath := pathname + "\\" + "窦性心动过缓_ST段压低(显著)_T波改变(小于0.05mV)-I-1-0001-" + sfm + ".txt"
	fmt.Println(newpath)
	err := os.Rename(filePath, newpath)
	if err != nil {
		fmt.Println(err)
	}
	time.Sleep(1 * time.Second)
}
func Substr(str string, start int, end int) string {
	rs := str[:]
	length := len(rs)

	if start < 0 || start > length {
		panic("start is wrong")
	}

	if end < 0 || end > length {
		panic("end is wrong")
	}

	return string(rs[start:end])
}

func initmap() map[int]string {
	zhibiaomap := map[int]string{}
	zhibiaomap[1] = "窦性心律"
	zhibiaomap[2] = "窦性心律,心电图未见异常"
	zhibiaomap[3] = "窦性心动过速"
	zhibiaomap[4] = "窦性心动过缓"
	zhibiaomap[5] = "窦性停搏"
	zhibiaomap[6] = "心房颤动"
	zhibiaomap[7] = "房性早搏"
	zhibiaomap[8] = "偶发房性早搏"
	zhibiaomap[9] = "频发房性早搏"
	zhibiaomap[10] = "房性早搏二联律"
	zhibiaomap[11] = "房性早搏三联律"
	zhibiaomap[12] = "成对性房性早搏"
	zhibiaomap[13] = "短阵房性心动过速"
	zhibiaomap[14] = "室性早搏"
	zhibiaomap[15] = "偶发室性早搏"
	zhibiaomap[16] = "频发室性早搏"
	zhibiaomap[17] = "室性早搏二联律"
	zhibiaomap[18] = "室性早搏三联律"
	zhibiaomap[19] = "成对室性早搏"
	zhibiaomap[20] = "短阵室性心动过速"
	zhibiaomap[21] = "室上性心动过速"
	zhibiaomap[22] = "一度房室阻滞"
	zhibiaomap[23] = "ST段抬高(显著)"
	zhibiaomap[24] = "ST段压低(显著)"
	zhibiaomap[25] = "QT/QTc间期延长"
	zhibiaomap[26] = "RR长间歇"
	zhibiaomap[27] = "心室内差异传导"
	zhibiaomap[28] = "干扰波"
	zhibiaomap[29] = "导联脱落"
	zhibiaomap[30] = "心房扑动"
	zhibiaomap[31] = "短PR间期"
	zhibiaomap[32] = "二度房室阻滞"
	zhibiaomap[33] = "P波增高"
	zhibiaomap[34] = "P波增宽"
	zhibiaomap[35] = "QRS波群呈XX型"
	zhibiaomap[36] = "R波高电压"
	zhibiaomap[37] = "室内阻滞"
	zhibiaomap[38] = "T波改变"
	zhibiaomap[39] = "短QT/QTc间期"
	zhibiaomap[40] = "心电图未见明显异常"
	return zhibiaomap
}
