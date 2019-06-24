package main

import (
	"fmt"
	"github.com/smallnest/rpcx"
)



func main() {
	s := &rpcx.DirectClientSelector{Network: "tcp", Address: "127.0.0.1:8972"}
	client := rpcx.NewClient(s)

	args := &Args{7, 8}
	var reply Reply
	err := client.Call("Arith.Mul", args, &reply)
	if err != nil {
		fmt.Printf("error for Arith: %d*%d, %v \n", args.A, args.B, err)
	} else {
		fmt.Printf("Arith: %d*%d=%d \n", args.A, args.B, reply.C)
	}

	client.Close()
}