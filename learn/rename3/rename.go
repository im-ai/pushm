package main

import (
	"crypto/md5"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"sync"
	"time"
)

var (
	lock    sync.Mutex
	nt      time.Time
	path    string
	tarpath string
)

func main() {
	nt = time.Now()
	nt.AddDate(-1, 0, 0)
	fmt.Printf("请把目录拖进来: ")
	fmt.Scanln(&path)
	fmt.Printf("请把新目录拖进来: ")
	fmt.Scanln(&tarpath)
	GetAllFile(path, tarpath)

}

func GetAllFile(pathname, tarpath string) error {
	rd, err := ioutil.ReadDir(pathname)
	for _, fi := range rd {
		if fi.IsDir() {
			GetAllFile(pathname+"\\"+fi.Name()+"\\", tarpath)
		} else {
			filePath := pathname + "\\" + fi.Name()
			readfileAndRename(filePath, tarpath)
		}
	}
	return err
}

func readfileAndRename(filePath, tarpath string) string {
	//substr := str.Substr(filePath, len(filePath)-18, len(filePath))
	substr := getMd5ByFile(filePath) // 获取文件MD5
	fmt.Println(substr)
	newpath := tarpath + "\\" + substr + ".txt"
	fmt.Println("newpath:" + newpath)

	src, err := os.OpenFile(filePath, os.O_RDWR, 0644)
	if err != nil {
		fmt.Println(err)
		return newpath
	}
	defer src.Close()

	dst, err := os.OpenFile(newpath, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		fmt.Println(err)
		return newpath
	}
	defer dst.Close()

	payload, err := ioutil.ReadAll(src)
	n, err := dst.Write(payload)
	if err != nil {
		fmt.Println(err)
	}
	n = n

	return newpath
}

func getMd5ByFile(path string) string {
	f, err := os.OpenFile(path, os.O_RDWR, 0644)
	if err != nil {
		fmt.Println("Open", err)
		return ""
	}

	defer f.Close()

	md5hash := md5.New()
	if _, err := io.Copy(md5hash, f); err != nil {
		fmt.Println("Copy", err)
		return ""
	}
	sprintf := fmt.Sprintf("%x", md5hash.Sum(nil))
	return sprintf
	//fmt.Printf("%x\n", md5hash.Sum(nil))
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
