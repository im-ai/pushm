package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"
)

var (
	path   string
	profix string
)

func main() {
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
			readfileAndRename(pathname, profix, filePath)

		}
	}
	return err
}

func readfileAndRename(pathname, profix, filePath string) {

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
	newpath := pathname + "\\" + profix + "-" + sfm + ".txt"
	fmt.Println(newpath)
	err := os.Rename(filePath, newpath)
	if err != nil {
		fmt.Println(err)
	}
	time.Sleep(1 * time.Second)
}
