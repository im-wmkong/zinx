package znet

import "zinx/zicafe"

type Request struct {
	conn zicafe.IConnection
	data []byte
}

func (r Request) GetConnection() zicafe.IConnection {
	return r.conn
}

func (r Request) GetData() []byte {
	return r.data
}
