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

	err := request.GetConnection().SendMsg(0, []byte("PingRouter..."))
	if err != nil {
		fmt.Println(err)
	}
}

type HelloRouter struct {
	znet.BaseRouter
}

func (r *HelloRouter) Handle(request ziface.IRequest) {
	fmt.Println("Call HelloRouter Handle")
	fmt.Printf("receive from client: msgid = %d ,data = %s", request.GetMsgID(), request.GetData())

	err := request.GetConnection().SendMsg(1, []byte("HelloRouter..."))
	if err != nil {
		fmt.Println(err)
	}
}

func DoOnStart(c ziface.IConnection) {
	fmt.Println("-->DoOnStart is start")
	err := c.SendMsg(2, []byte("-->DoOnStart begin..."))
	if err != nil {
		fmt.Println(err)
		return
	}
}
func DoOnStop(c ziface.IConnection) {
	fmt.Printf("-->DoOnStop is start,%d", c.GetConnID())
}

func main() {
	s := znet.NewServer()

	s.SetOnConnStart(DoOnStart)
	s.SetOnConnStop(DoOnStop)

	s.AddRouter(0, &PingRouter{})
	s.AddRouter(1, &HelloRouter{})

	s.Serve()
}
