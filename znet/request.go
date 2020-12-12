package znet

import "zinx/zicafe"

type Request struct {
	conn zicafe.IConnection
	msg zicafe.IMessage
}

func (r Request) GetConnection() zicafe.IConnection {
	return r.conn
}

func (r Request) GetData() []byte {
	return r.msg.GetData()
}

func (r Request) GetMsgId() uint32 {
	return r.msg.GetMsgId()
}
