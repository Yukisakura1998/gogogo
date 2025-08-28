package utils

import (
	"encoding/json"
	"os"
	"zinx/ziface"
)

type GlobalObj struct {
	TcpServer ziface.IServer
	Host      string `json:"Host"`
	TcpPort   int    `json:"TcpPort"`
	Name      string `json:"Name"`
	Version   string

	MaxPacketSize    uint32
	MaxConn          int    `json:"MaxConn"`
	WorkerPoolSize   uint32 `json:"WorkerPoolSize"`
	MaxWorkerTaskLen uint32
	ConfigFilePath   string
	MaxMsgChanLen    uint32
}

func (g *GlobalObj) Reload() {
	data, err := os.ReadFile("conf/zinx.json") //文件路径问题
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(data, &GlobalObject)
	if err != nil {
		panic(err)
	}
}

var GlobalObject *GlobalObj

func init() {
	GlobalObject = &GlobalObj{
		Host:             "0.0.0.0",
		TcpPort:          7777,
		Name:             "ZinxServerApp",
		Version:          "V0.4",
		MaxPacketSize:    4096,
		MaxConn:          12000,
		WorkerPoolSize:   10,
		MaxWorkerTaskLen: 1024,
		MaxMsgChanLen:    1024,
	}
	GlobalObject.Reload()
}
