package zicafe

type IServer interface {
	Start()
	Stop()
	Serve()
}
