package main

// 可以不用，主要是为了调用包捏的 init 函数
import _ "fmt"

func main() {

}

//// 别名
//import io "fmt"
//
//func main()  {
//	io.Println("aaa")
//}
//
//// . 点操作
//import . "fmt"
//
//func main() {
//	Println()
//}

//import (
//	"fmt"
//	"os"
//)
//
//func main() {
//	fmt.Println("this is a learn")
//	fmt.Println("os.Args = ",os.Args)
//}
