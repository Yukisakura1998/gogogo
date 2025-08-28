package main

import (
	"fmt"
	"net"
	"time"
)

func mainV01() {
	fmt.Println("Client0 Test... Start...")
	//3s发送延时
	time.Sleep(3 * time.Second)
	conn, err := net.Dial("tcp", "127.0.0.1:7777")
	if err != nil {
		fmt.Println("client connect err : ", err)
		return
	}
	for {
		_, err := conn.Write([]byte("hello"))
		if err != nil {
			fmt.Println("write error : ", err)
			return
		}
		buf := make([]byte, 512)
		cnt, err := conn.Read(buf)
		if err != nil {
			fmt.Println("read error : ", err)
			return
		}

		fmt.Printf("server call back : %s , cnt = %d\n", buf, cnt)

		time.Sleep(1 * time.Second)
	}
}
