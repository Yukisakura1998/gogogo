package znet

import (
	"errors"
	"fmt"
	"io"
	"net"
	"sync"
	"zinx/utils"
	"zinx/ziface"
)

type Connection struct {
	TcpServer ziface.IServer
	Conn      *net.TCPConn
	ConnID    string
	isClosed  bool

	//handlerAPI ziface.HandFunc
	//注册的方法
	//Router ziface.IRouter

	MsgHandler   ziface.IMsgHandler
	ExitBuffChan chan bool
	msgChan      chan []byte
	msgBuffChan  chan []byte
	property     map[string]interface{}
	propertyLock sync.RWMutex
}

func (c *Connection) SetProperty(key string, value interface{}) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()

	c.property[key] = value
}

func (c *Connection) GetProperty(key string) (value interface{}, err error) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()

	if value, ok := c.property[key]; ok {
		return value, nil
	} else {
		return nil, errors.New("get property error,not found " + key)
	}
}

func (c *Connection) RemoveProperty(key string) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()

	delete(c.property, key)
}

func (c *Connection) SendMsg(msgId uint32, data []byte) error {
	if c.isClosed == true {
		return errors.New("在发送消息时链接被关闭")
	}
	dp := NewDataPack()
	msg, err := dp.Pack(NewMessage(msgId, data))
	if err != nil {
		fmt.Printf("Pack error : %s, msgid: %d", err, msgId)
		return errors.New("消息打包失败")
	}
	/*if _, err := c.Conn.Write(msg); err != nil {
		fmt.Printf("Write error : %s, msg id: %d", err, msgId)
		c.ExitBuffChan <- true
		return errors.New("conn Write error")
	}*/
	c.msgChan <- msg
	return nil
}
func (c *Connection) SendBuffMsg(msgId uint32, data []byte) error {
	if c.isClosed == true {
		return errors.New("在发送消息时链接被关闭")
	}
	dp := NewDataPack()
	msg, err := dp.Pack(NewMessage(msgId, data))
	if err != nil {
		fmt.Printf("Pack error : %s, msgid: %d", err, msgId)
		return errors.New("消息打包失败")
	}
	c.msgChan <- msg
	return nil
}

func NewConnection(server ziface.IServer, conn *net.TCPConn, connID string, msgHandler ziface.IMsgHandler) *Connection {
	c := &Connection{
		TcpServer: server,
		Conn:      conn,
		ConnID:    connID,
		isClosed:  false,
		//handlerAPI:   callbackApi,
		MsgHandler:   msgHandler,
		ExitBuffChan: make(chan bool, 1),
		msgChan:      make(chan []byte),
		msgBuffChan:  make(chan []byte, utils.GlobalObject.MaxMsgChanLen),
		property:     make(map[string]interface{}),
	}
	c.TcpServer.GetConnMgr().Add(c)
	return c
}

func (c *Connection) StarReader() {
	fmt.Println("Reader Goroutine is Running")
	defer fmt.Println(c.RemoteAddr().String(), " conn reader exit!")
	defer c.Stop()
	for {
		dp := NewDataPack()
		headData := make([]byte, dp.GetHeadLen())
		if _, err := io.ReadFull(c.GetTCPConnection(), headData); err != nil {
			fmt.Println("read head error ", err)
			c.ExitBuffChan <- true
			continue
		}
		/*cnt, err := c.Conn.Read(buf)
		if err != nil {
			fmt.Println("read error ", err)
			c.ExitBuffChan <- true
			return
		}
		if err := c.handlerAPI(c.Conn, buf[:cnt], cnt); err != nil {
			fmt.Println("connID: ", c.ConnID, " handle is error :", err)
			c.ExitBuffChan <- true
			return
		}*/
		msg, err := dp.Unpack(headData)
		if err != nil {
			fmt.Println("Unpack error ", err)
			c.ExitBuffChan <- true
			continue
		}
		var data []byte
		if msg.GetDataLength() > 0 {
			data = make([]byte, msg.GetDataLength())
			if _, err := io.ReadFull(c.GetTCPConnection(), data); err != nil {
				fmt.Println("read msg error ", err)
				c.ExitBuffChan <- true
				continue
			}

		}

		msg.SetData(data)

		req := &Request{
			conn: c,
			msg:  msg,
			//data: append([]byte(nil), buf[:cnt]...), // 防御性拷贝
		}

		/*go func(request ziface.IRequest) {
			defer func() {
				if r := recover(); r != nil {
					fmt.Printf("handler panic: %v", r)
				}
			}()

			c.Router.PreHandle(request)
			c.Router.Handle(request)
			c.Router.PostHandle(request)
		}(req)*/
		if utils.GlobalObject.WorkerPoolSize > 0 {
			c.MsgHandler.SendMsgToTaskQueue(req)
		} else {
			go c.MsgHandler.DoMsgHandler(req)
		}
	}
}
func (c *Connection) StarWriter() {
	fmt.Println("Writer Goroutine is Running")
	defer fmt.Println(c.RemoteAddr().String(), " conn writer exit!")
	defer c.Stop()
	for {
		select {
		case data := <-c.msgChan:
			if _, err := c.Conn.Write(data); err != nil {
				fmt.Println("send msg error ", err)
				return
			}
		case data, ok := <-c.msgBuffChan:
			if ok {
				if _, err := c.Conn.Write(data); err != nil {
					fmt.Println("send msg error ", err)
					return
				}
			} else {
				fmt.Println("msgBuffChan is closed")
				break
			}
		case <-c.ExitBuffChan:
			return
		}
	}
}

func (c *Connection) Start() {
	go c.StarReader()
	go c.StarWriter()
	c.TcpServer.CallOnConnStart(c)
	for {
		select {
		case <-c.ExitBuffChan:
			return
		}
	}
}

func (c *Connection) Stop() {
	if c.isClosed == true {
		return
	}
	c.isClosed = true
	//用于调用用户注册的关闭时的业务
	c.TcpServer.CallOnConnStop(c)

	if err := c.Conn.Close(); err != nil {
		return
	}

	c.ExitBuffChan <- true

	c.TcpServer.GetConnMgr().Remove(c)

	close(c.ExitBuffChan)
	close(c.msgBuffChan)
}

func (c *Connection) GetTCPConnection() *net.TCPConn {
	return c.Conn
}

func (c *Connection) GetConnID() string {
	return c.ConnID
}

func (c *Connection) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}
