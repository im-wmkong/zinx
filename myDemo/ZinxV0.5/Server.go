package main

import (
	"fmt"
	"zinx/zicafe"
	"zinx/znet"
)

type PingRouter struct {
	znet.BaseRouter
}

func (pr *PingRouter) Handle(request zicafe.IRequest) {
	fmt.Printf("Recv from client msg, msgId: %d, data: %s\n", request.GetMsgId(), string(request.GetData()))

	if err := request.GetConnection().SendMsg(1, []byte("Ping\n")); err != nil {
		fmt.Printf("Sned error: %s\n", err)
	}
}

func main() {
	s := znet.NewServer("[zinx V0.5]")
	s.AddRouter(&PingRouter{})
	s.Serve()
}
