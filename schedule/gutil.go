package main

import (
	"bytes"
	"fmt"
	"os"
)

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

//处理错误，根据实际情况选择这样处理，还是在函数调之后不同的地方不同处理
func CheckErr(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}
