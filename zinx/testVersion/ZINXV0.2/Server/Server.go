package main

import (
	"fmt"
	"zinx/ziface"
	"zinx/znet"
)

type PingrouterV02 struct {
	znet.BaseRouter
}

func (r *PingrouterV02) PreHandle(request ziface.IRequest) {
	fmt.Println("Call PreHandle")
	if _, err := request.GetConnection().GetTCPConnection().Write([]byte("before main handle ")); err != nil {
		fmt.Println("PreHandle error :", err)
		return
	}
}

func (r *PingrouterV02) Handle(request ziface.IRequest) {
	fmt.Println("Call Handle")
	if _, err := request.GetConnection().GetTCPConnection().Write([]byte("main handle ")); err != nil {
		fmt.Println("Handle error :", err)
		return
	}
}

func (r *PingrouterV02) PostHandle(request ziface.IRequest) {
	fmt.Println("Call PostHandle")
	if _, err := request.GetConnection().GetTCPConnection().Write([]byte("after main handle ")); err != nil {
		fmt.Println("PostHandle error :", err)
		return
	}
}
func main2() {
	s := znet.NewServer()
	s.AddRouter(&PingrouterV02{})
	s.Serve()
}
