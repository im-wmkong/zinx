package main

import (
	"fmt"
	"zinx/zicafe"
	"zinx/znet"
)

type PingRouter struct {
	znet.BaseRouter
}

func (pr *PingRouter) PreHandle(request zicafe.IRequest) {
	fmt.Println("Call Router PreHandle")
	_, err := request.GetConnection().GetTCPConnection().Write([]byte("Before ping\n"))
	if err != nil {
		fmt.Printf("Before ping error: %s\n", err)
	}
}

func (pr *PingRouter) Handle(request zicafe.IRequest) {
	fmt.Println("Call Router Handle")
	_, err := request.GetConnection().GetTCPConnection().Write([]byte("Ping\n"))
	if err != nil {
		fmt.Printf("Ping error: %s\n", err)
	}
}

func (pr *PingRouter) PostHandle(request zicafe.IRequest) {
	fmt.Println("Call Router Handle")
	_, err := request.GetConnection().GetTCPConnection().Write([]byte("After ping\n"))
	if err != nil {
		fmt.Printf("After ping error: %s\n", err)
	}
}

func main() {
	s := znet.NewServer("[zinx V0.3]")
	s.AddRouter(&PingRouter{})
	s.Serve()
}
