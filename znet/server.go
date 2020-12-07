package znet

import (
	"fmt"
	"net"
	"zinx/zicafe"
)

type Server struct {
	Name      string
	IPVersion string
	IP        string
	Port      int
}

func CallBackToClient(conn *net.TCPConn, data []byte, cnt int) error {
	fmt.Println("Connection handle CallBackToClient")
	if _, err := conn.Write(data[:cnt]); err != nil {
		fmt.Printf("Connection write error: %s\n", err)
		return err
	}
	return nil
}

func (s Server) Start() {
	fmt.Printf("[Start] Server Listener at IP: %s, Port: %d, is starting\n", s.IP, s.Port)

	go func() {
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
			dealConn := NewConnection(conn, cid, CallBackToClient)
			cid++
			go dealConn.Start()
		}
	}()
}

func (s Server) Stop() {

}

func (s Server) Serve() {
	s.Start()

	select {}
}

func NewServer(name string) zicafe.IServer {
	return Server{
		Name:      name,
		IPVersion: "tcp4",
		IP:        "0.0.0.0",
		Port:      8999,
	}
}
