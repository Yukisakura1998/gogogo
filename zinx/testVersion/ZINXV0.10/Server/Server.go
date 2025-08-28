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
	fmt.Printf("receive from client: msgid = %d ,data = %s\n", request.GetMsgID(), request.GetData())

	err := request.GetConnection().SendMsg(0, []byte("PingRouter..."))
	if err != nil {
		fmt.Println(err)
	}
}

func DoOnStart(c ziface.IConnection) {
	fmt.Println("-->DoOnStart is start")
	c.SetProperty("name", "yuki")
	fmt.Println("-->set name yuki")
	err := c.SendMsg(2, []byte("-->DoOnStart begin..."))
	if err != nil {
		fmt.Println(err)
		return
	}
}
func DoOnStop(c ziface.IConnection) {
	fmt.Printf("-->DoOnStop is start,%s\n", c.GetConnID())

	if name, err := c.GetProperty("name"); err != nil {
		fmt.Println("-->get property name :", name)
	}
}

func main() {
	s := znet.NewServer()

	s.SetOnConnStart(DoOnStart)
	s.SetOnConnStop(DoOnStop)

	s.AddRouter(0, &PingRouter{})

	s.Serve()
}
