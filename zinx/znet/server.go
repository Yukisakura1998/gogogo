package znet

import (
	"fmt"
	"github.com/google/uuid"
	"net"
	"zinx/utils"
	"zinx/ziface"
)

type Server struct {
	//服务器的名称
	Name string
	//tcp4 or other
	IPVersion string
	//服务绑定的IP地址
	IP string
	//服务绑定的端口
	Port int
	//注册的路由方法
	//Router ziface.IRouter
	//绑定消息
	msgHandler ziface.IMsgHandler
	//连接管理
	ConnMgr ziface.IConnManager

	OnConnStart func(conn ziface.IConnection)
	OnConnStop  func(conn ziface.IConnection)
}

/*
	func CallBackToClient(conn *net.TCPConn, data []byte, cnt int) error {
		fmt.Println("[Conn Handler] CallBackToClient...")
		if _, err := conn.Write(data[:cnt]); err != nil {
			fmt.Println("write error ", err)
			return errors.New("write error")
		}
		return nil
	}
*/
func (s *Server) Start() {
	fmt.Printf("[START] Server name: %s, listenner at IP : %s,Port %d, is starting\n", s.Name, s.IP, s.Port)
	fmt.Printf("[Zinx] Version: %s ,MaxConn : %d , MaxPacketSize: %d\n", utils.GlobalObject.Version, utils.GlobalObject.MaxConn, utils.GlobalObject.MaxPacketSize)
	//开启一个GO去做服务器段的监听业务
	go func() {
		//0.启动工作池
		s.msgHandler.StartWorkerPool()
		//1.获取一个addr
		addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))
		if err != nil {
			fmt.Println("resolve tcp addr err: ", err)
			return
		}
		//2.监听服务器连接
		listener, err := net.ListenTCP(s.IPVersion, addr)
		if err != nil {
			fmt.Println("listen ", s.IPVersion, " err: ", err)
			return
		}
		//监听OK
		fmt.Println("start Zinx server ", s.Name, "success, now listening...")

		//3启动server
		for {
			//自动生成ID uuid
			newUuid := uuid.New()
			cid := newUuid.String()

			//3.1 阻塞等待客户端请求连接
			conn, err := listener.AcceptTCP()
			if err != nil {
				fmt.Println("Accept err ", err)
				continue
			}
			//3.2 设置服务器的最大连接
			if s.ConnMgr.Len() >= utils.GlobalObject.MaxConn {
				err := conn.Close()
				if err != nil {
					continue
				}
				continue
			}

			//3.3 TODO Server.Start() 处理新连接请求的业务方法

			/*//暂时做回显方法
			go func() {
				//不断循环接受数据
				for {
					buf := make([]byte, 512)
					cnt, err := conn.Read(buf)
					if err != nil {
						fmt.Println("read error ", err)
						continue
					}
					if _, err := conn.Write(buf[:cnt]); err != nil {
						fmt.Println("write error ", err)
						continue
					}
					fmt.Println(string(buf[:cnt]))
				}
			}()*/
			dealConn := NewConnection(s, conn, cid, s.msgHandler)

			//启动业务
			go dealConn.Start()
		}
	}()
}

func (s *Server) Stop() {
	//TODO Server.Stop()清理连接和其他信息
	fmt.Println("[STOP] Zinx server ,name :", s.Name)
	s.ConnMgr.ClearConn()
}

func (s *Server) Serve() {
	s.Start()
	//TODO Server.serve() 可以添加其他操作
	select {}
}
func (s *Server) AddRouter(msgId uint32, router ziface.IRouter) {
	s.msgHandler.AddRouter(msgId, router)
	fmt.Println("Add router success...")
}
func (s *Server) GetConnMgr() ziface.IConnManager {
	return s.ConnMgr
}
func (s *Server) SetOnConnStart(f func(c ziface.IConnection)) {
	s.OnConnStart = f
}

func (s *Server) SetOnConnStop(f func(c ziface.IConnection)) {
	s.OnConnStop = f
}

func (s *Server) CallOnConnStart(c ziface.IConnection) {
	if s.OnConnStart != nil {
		fmt.Println("-->call OnConnStart")
		s.OnConnStart(c)
	}
}

func (s *Server) CallOnConnStop(c ziface.IConnection) {
	if s.OnConnStop != nil {
		fmt.Println("-->call OnConnStop")
		s.OnConnStop(c)
	}
}
func NewServer() ziface.IServer {
	utils.GlobalObject.Reload()
	s := &Server{
		Name:       utils.GlobalObject.Name,
		IPVersion:  "tcp4",
		IP:         utils.GlobalObject.Host,
		Port:       utils.GlobalObject.TcpPort,
		msgHandler: NewMsgHandler(),
		ConnMgr:    NewConnManager(),
	}
	return s
}
