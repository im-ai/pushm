package main

import "github.com/smallnest/rpcx"

import "fmt"



type Arith int

func (t *Arith) Mul(args *Args, reply *Reply) error {
	reply.C = args.A * args.B
	fmt.Println(reply.C, " = ", args.A, " * ", args.B)
	return nil
}

func (t *Arith) Error(args *Args, reply *Reply) error {
	panic("ERROR")
}

func main() {
	server := rpcx.NewServer()
	server.RegisterName("Arith", new(Arith))
	server.Serve("tcp", "127.0.0.1:8972")
}