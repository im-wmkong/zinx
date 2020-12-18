package zicafe

type IConnManager interface {
	Add(connection IConnection)
	Remove(connection IConnection)
	Get(connID uint32) (IConnection, error)
	Len() int
	Clear()
}
