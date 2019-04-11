package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

func InitCfg() error {
	numberstr := GetCfg("number", "sfig.ini")
	//fmt.Println("number:", numberstr)
	nubmers, _ := strconv.Atoi(numberstr)
	nubmer = nubmers
	typeIdstr := GetCfg("typeId", "sfig.ini")
	//fmt.Println("typeId:", typeIdstr)
	typeId, _ := strconv.Atoi(typeIdstr)
	urlstr := GetCfg("url", "sfig.ini")
	//fmt.Println("url:", urlstr)
	jsonstr := GetCfg("json", "sfig.ini")
	//fmt.Println("json:", jsonstr)
	gonumber = GetConfigByKey("gonumber")
	config := &PressureBody{
		TypeId: typeId,
		Url:    urlstr,
		Json:   jsonstr,
		Number: nubmer,
	}
	bytesa, e := json.Marshal(config)
	if e != nil {
		fmt.Println(e)
		return e
	}
	bytesCombine = BytesCombine(bytesa, []byte("\n"))

	return nil
}

func GetConfigByKey(key string) int {
	numberstr := GetCfg(key, "sfig.ini")
	//fmt.Println("number:", numberstr)
	nubmers, _ := strconv.Atoi(numberstr)
	return nubmers
}

func GetCfg(tag string, filepath string) string {
	dat, err := ioutil.ReadFile(filepath) //读取文件
	CheckErr(err)                         //检查是否有错误
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

func GetConfig() *PressureBody {
	numberstr := GetCfg("number", "sfig.ini")
	//fmt.Println("number:", numberstr)
	nubmers, _ := strconv.Atoi(numberstr)
	nubmer = nubmers
	typeIdstr := GetCfg("typeId", "sfig.ini")
	//fmt.Println("typeId:", typeIdstr)
	typeId, _ := strconv.Atoi(typeIdstr)
	urlstr := GetCfg("url", "sfig.ini")
	//fmt.Println("url:", urlstr)
	jsonstr := GetCfg("json", "sfig.ini")
	//fmt.Println("json:", jsonstr)

	gonumber = GetConfigByKey("gonumber")

	config := &PressureBody{
		TypeId: typeId,
		Url:    urlstr,
		Json:   jsonstr,
		Number: nubmer,
	}
	return config
}

var(
	TypeId int    // 1: http get 2: http post  3: ws
	Urlstr    string // 请求 url
	Jsonstr   string // post参数
	Number int    // 每秒开启 gorouting 次数
)
func GetConfigChange() *PressureBody {

	config := &PressureBody{
		TypeId: TypeId,
		Url:    Urlstr,
		Json:   Jsonstr,
		Number: nubmer,
	}

	return config
}


func Changeconf(w http.ResponseWriter, r *http.Request) {
	numbers, ok := r.URL.Query()["number"]
	if !ok || len(numbers) < 1 {
		log.Println("Url Param 'number' is missing")
		return
	}

	gonumbers, ok := r.URL.Query()["gonumber"]
	typeIds, ok := r.URL.Query()["typeId"]
	urls, ok := r.URL.Query()["url"]
	jsons, ok := r.URL.Query()["json"]

	// Query()["key"] will return an array of items,
	// we only want the single item.
	numbert, _ := strconv.Atoi(string(numbers[0]))

	//log.Println("Url Param 'number' 1 is: ", nubmer)

	configt := GetConfig()
	configt.Number = numbert
	nubmer = numbert

	if len(gonumbers) > 0 {
		gonumbern, _ := strconv.Atoi(string(gonumbers[0]))
		gonumber = gonumbern
	}
	if len(typeIds) > 0 {
		typeId, _ := strconv.Atoi(string(typeIds[0]))
		configt.TypeId = typeId
		TypeId = typeId
	}
	if len(urls) > 0 {
		configt.Url = urls[0]
		Urlstr = urls[0]
	}
	if len(jsons) > 0 {
		configt.Json = jsons[0]
		Jsonstr = jsons[0]
	}

	bytesa, e := json.Marshal(configt)
	if e != nil {
		fmt.Println(e)
		return
	}
	bytesCombine = BytesCombine(bytesa, []byte("\n"))

	//log.Println("Url Param 'number' 2 is: ", nubmer)
}
