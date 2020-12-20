package main

import (
	"fmt"
	"zinx/zicafe"
	"zinx/znet"
)

func ConnectionStart(conn zicafe.IConnection) {
	fmt.Println("Connection start")
}

func ConnectionStop(conn zicafe.IConnection) {
	fmt.Println("Connection stop")
}

type PingRouter struct {
	znet.BaseRouter
}

func (pr *PingRouter) Handle(request zicafe.IRequest) {
	fmt.Printf("Recv from client msg, msgId: %d, data: %s\n", request.GetMsgId(), string(request.GetData()))

	if err := request.GetConnection().SendMsg(1, []byte("Pong\n")); err != nil {
		fmt.Printf("Sned error: %s\n", err)
	}
}

type HelloRouter struct {
	znet.BaseRouter
}

func (hr *HelloRouter) Handle(request zicafe.IRequest) {
	fmt.Printf("Recv from client msg, msgId: %d, data: %s\n", request.GetMsgId(), string(request.GetData()))

	if err := request.GetConnection().SendMsg(1, []byte("Hello Zinx\n")); err != nil {
		fmt.Printf("Sned error: %s\n", err)
	}
}

func main() {
	s := znet.NewServer()

	s.SetOnConnStart(ConnectionStart)
	s.SetOnConnStop(ConnectionStop)

	s.AddRouter(0, &PingRouter{})
	s.AddRouter(1, &HelloRouter{})
	s.Serve()
}
