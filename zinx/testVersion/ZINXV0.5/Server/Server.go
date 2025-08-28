package main

import (
	"fmt"
	"zinx/ziface"
	"zinx/znet"
)

type PingRouter struct {
	znet.BaseRouter
}

func (r *PingRouter) Handle(request ziface.IRequest) {
	fmt.Println("Call PingRouter Handle")
	fmt.Printf("receive from client: msgid = %d ,data = %s", request.GetMsgID(), request.GetData())

	err := request.GetConnection().SendMsg(1, []byte("ping..."))
	if err != nil {
		fmt.Println(err)
	}
}
func main() {
	s := znet.NewServer()

	s.AddRouter(&PingRouter{})

	s.Serve()
}
