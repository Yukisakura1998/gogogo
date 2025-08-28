package main

import (
	"fmt"
	"io"
	"net"
	"time"
	"zinx/znet"
)

func main() {
	fmt.Println("Client0 Test... Start...")
	time.Sleep(3 * time.Second)
	conn, err := net.Dial("tcp", "127.0.0.1:7777")
	if err != nil {
		fmt.Println("client connect err : ", err)
		return
	}
	for {
		dp := znet.NewDataPack()
		msg, _ := dp.Pack(znet.NewMessage(0, []byte("ZinxV0.5 Test Start.")))
		if _, err = conn.Write(msg); err != nil {
			fmt.Println("write error err : ", err)
			return
		}

		headData := make([]byte, dp.GetHeadLen())
		if _, err = io.ReadFull(conn, headData); err != nil {
			fmt.Println("read head error err : ", err)
			return
		}
		msgHead, err := dp.Unpack(headData)

		if err != nil {
			fmt.Println("server unpack error err : ", err)
			return
		}

		if msgHead.GetDataLength() > 0 {
			msg := msgHead.(*znet.Message)
			msg.Data = make([]byte, msg.GetDataLength())
			_, err := io.ReadFull(conn, msg.Data)
			if err != nil {
				fmt.Println("server unpack data error err : ", err)
				return
			}
			fmt.Printf("==> Receive messageId ID = %d , len = %d ,data = %s\n", msg.Id, msg.GetDataLength(), msg.GetData())
		}
		time.Sleep(1 * time.Second)
	}
}
