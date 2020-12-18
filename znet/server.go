package znet

import (
	"fmt"
	"net"
	"zinx/utils"
	"zinx/zicafe"
)

type Server struct {
	Name       string
	IPVersion  string
	IP         string
	Port       int
	MsgHandler zicafe.IMsgHandler
	ConnMgr zicafe.IConnManager
}

func (s *Server) Start() {
	fmt.Printf("[Start] Server Listener at IP: %s, Port: %d, is starting\n", s.IP, s.Port)

	go func() {
		s.MsgHandler.StartWorkerPool()
		// 获取一个tcp addr
		addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))
		if err != nil {
			fmt.Println("Resolve tcp address error: ", err)
			return
		}

		// 监听服务器地址
		listener, err := net.ListenTCP(s.IPVersion, addr)
		if err != nil {
			fmt.Printf("Listen %s error: %s\n", s.IP, err)
			return
		}
		fmt.Printf("Start Zinx server success, %s listening...\n", s.Name)

		// 阻塞等待客户端链接
		var cid uint32
		cid = 0
		for {
			conn, err := listener.AcceptTCP()
			if err != nil {
				fmt.Printf("Accept err %s\n", err)
				continue
			}

			if s.ConnMgr.Len() >= utils.GlobalObject.MaxConn {
				conn.Close()
				continue
			}

			dealConn := NewConnection(s, conn, cid, s.MsgHandler)
			cid++
			go dealConn.Start()
		}
	}()
}

func (s *Server) Stop() {
	s.ConnMgr.Clear()
}

func (s *Server) Serve() {
	s.Start()

	select {}
}

func (s *Server) AddRouter(msgID uint32, router zicafe.IRouter) {
	s.MsgHandler.AddRouter(msgID, router)
}

func (s *Server) GetConnManager() zicafe.IConnManager {
	return s.ConnMgr
}

func NewServer() zicafe.IServer {
	return &Server{
		Name:       utils.GlobalObject.Name,
		IPVersion:  "tcp4",
		IP:         utils.GlobalObject.Host,
		Port:       utils.GlobalObject.TcpPort,
		MsgHandler: NewMsgHandle(),
		ConnMgr: NewConnManager(),
	}
}
