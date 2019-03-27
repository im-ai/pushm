package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	lock   sync.Mutex
	nt     time.Time
	path   string
	profix string
)

func main() {
	nt = time.Now()
	nt.AddDate(-1, 0, 0)
	fmt.Printf("请输入前缀: ")
	fmt.Scanln(&profix)
	fmt.Printf("请把目录拖进来: ")
	fmt.Scanln(&path)
	GetAllFile(path, profix)

}

func GetAllFile(pathname, profix string) error {
	rd, err := ioutil.ReadDir(pathname)
	for _, fi := range rd {
		if fi.IsDir() {
			GetAllFile(pathname+"\\"+fi.Name()+"\\", profix)
		} else {
			filePath := pathname + "\\" + fi.Name()
			newFilePath := readfileAndRename(pathname, profix, filePath)
			seek2line(newFilePath)
		}
	}
	return err
}

func seek2line(filePath string) {

	f, err := os.OpenFile(filePath, os.O_RDWR, 0644)
	if err != nil {
		fmt.Println(err)
	}
	defer f.Close()

	rd := bufio.NewReader(f)
	flag := 0
	typenum := 0  //  第二行数字 + 3 必得 250
	flagline := 0 //  0 ： 未处理  1,2 前添加 2,1,2    1: 处理过
	var buffer bytes.Buffer
	line1 := ""
	line2 := ""
	line3 := ""
	var valid = regexp.MustCompile("[0-9]")
	for {

		flag++
		line, err := rd.ReadString('\n') //以'\n'为结束符读入一行
		if err != nil || io.EOF == err {
			break
		}

		if flag == 1 {
			line1 = line
			fmt.Println("line1", line1)
		} else if flag == 2 {
			line2 = line
			linetmp := valid.FindAllString(line2, -1)
			i, err := strconv.Atoi(linetmp[0])
			if err != nil {
				typenum = 0
			} else {
				typenum = i
			}
			fmt.Println("line2", line2)
		} else if flag == 3 {
			line3 = line
			fmt.Println("line3", line3)
			t2 := valid.FindAllString(line2, -1)
			joint2 := strings.Join(t2, "")
			t3 := valid.FindAllString(line3, -1)
			joint3 := strings.Join(t3, "")
			fmt.Println(joint2)
			fmt.Println(joint3)
			if joint2 == "1" && joint3 == "2" {
				flagline = 1
				buffer.WriteString(line1)
				buffer.WriteString("2\r\n")
				buffer.WriteString(line2)
				buffer.WriteString(line3)
			} else {
				buffer.WriteString(line1)
				buffer.WriteString(line2)
				buffer.WriteString(line3)
			}
			fmt.Println(buffer.String())
		} else if flagline == 1 {
			line250 := valid.FindAllString(line, -1)
			join250 := strings.Join(line250, "")
			if join250 != "250" {
				buffer.WriteString("250\r\n")
				buffer.WriteString("60\r\n")
				if join250 != "32767" {
					buffer.WriteString("32767\r\n")
				}
				buffer.WriteString(line)
			} else {
				buffer.WriteString(line)
			}
			flagline = 2
		} else if flag == (typenum+3) && flagline != 2 {
			fmt.Println("line250", line)
			//  应该是 250
			line250 := valid.FindAllString(line, -1)
			join250 := strings.Join(line250, "")
			fmt.Println(join250)
			if join250 != "250" {
				buffer.WriteString("250\r\n")
				buffer.WriteString("60\r\n")
				if join250 != "32767" {
					buffer.WriteString("32767\r\n")
				}
				buffer.WriteString(line)
			} else {
				buffer.WriteString(line)
			}
		} else if flag == (typenum+5) && flagline != 2 {
			line32767 := valid.FindAllString(line, -1)
			join32767 := strings.Join(line32767, "")
			if join32767 != "32767" {
				buffer.WriteString("32767\r\n")
				buffer.WriteString(line)
			} else {
				buffer.WriteString(line)
			}
		} else {
			buffer.WriteString(line)
		}
	}
	f.Seek(0, 0)
	f.WriteString(buffer.String())

}

//
//else if flag == (typenum + 3) {
//fmt.Println("line250",line)
////  应该是 250
//line250 := valid.FindAllString(line, -1)
//if line250[0] != "250" {
//res = res+"250\r\n"
//res = res+"60\r\n"
//res = res+line
//}else{
//res = res+line
//}
//}else if flag == (typenum+5){
//fmt.Println("line32767",line)
////  应该是 32767
//line32767 := valid.FindAllString(line, -1)
//if line32767[0] != "32767" {
//res = res+"32767\r\n"
//res = res+line
//}else{
//res = res+line
//}
//}

func readfileAndRename(pathname, profix, filePath string) string {
	lock.Lock()
	sfm := getTime()
	lock.Unlock()
	newpath := pathname + "\\" + profix + "-" + sfm + ".txt"
	fmt.Println(newpath)
	err := os.Rename(filePath, newpath)
	if err != nil {
		fmt.Println(err)
	}
	return newpath
}

func getTime() string {
	y := fmt.Sprintf("%d", nt.Year())
	m := fmt.Sprintf("%d", nt.Month())
	d := fmt.Sprintf("%d", nt.Day())
	h := fmt.Sprintf("%d", nt.Hour())
	mm := fmt.Sprintf("%d", nt.Minute())
	s := fmt.Sprintf("%d", nt.Second())
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

	duration, _ := time.ParseDuration("1s")
	nt = nt.Add(duration)
	return sfm
}
