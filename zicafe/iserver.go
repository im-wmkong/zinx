package zicafe

type IServer interface {
	Start()
	Stop()
	Serve()
	AddRouter(msgID uint32, router IRouter)
}
