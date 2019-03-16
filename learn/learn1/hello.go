package main

import (
	"fmt"
	"github.com/im-ai/stringutil"
	"net/http"
	"os"
)

// init 先于 main 先执行    _ 这种默认导入，可以先执行 init函数
func init() {
	fmt.Println("init")
}

func Add(a, b int) int {
	return a + b
}

func main() {
	fmt.Println("imput args:", os.Args[:])
	fmt.Println("hello")
	fmt.Println(stringutil.Reverse("hello"))

	http.ListenAndServe(":8080", nil)
}
