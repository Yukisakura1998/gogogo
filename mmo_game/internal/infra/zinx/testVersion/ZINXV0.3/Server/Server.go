package main

import (
	"fmt"
	"zinx/ziface"
	"zinx/znet"
)

type PingRouterV03 struct {
	znet.BaseRouter
}

func (r *PingRouterV03) HandleV03(request ziface.IRequest) {
	fmt.Println("Call PingRouter Handle")
	_, err := request.GetConnection().GetTCPConnection().Write([]byte("ping router..."))
	if err != nil {
		fmt.Println("Call PingRouter err:", err)
		return
	}
}

func main3() {
	s := znet.NewServer()

	s.AddRouter(&PingRouterV03{})

	s.Serve()
}
