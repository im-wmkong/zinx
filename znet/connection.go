package znet

import (
	"errors"
	"fmt"
	"io"
	"net"
	"zinx/zicafe"
)

type Connection struct {
	Conn     *net.TCPConn
	ConnID   uint32
	isClosed bool
	Router   zicafe.IRouter
	ExitChan chan bool
}

func NewConnection(conn *net.TCPConn, connID uint32, router zicafe.IRouter) *Connection {
	return &Connection{
		Conn:     conn,
		ConnID:   connID,
		isClosed: false,
		Router:   router,
		ExitChan: make(chan bool, 1),
	}
}

func (c *Connection) StartReader() {
	fmt.Println("Connection Reader Goroutine is running")
	defer fmt.Printf("Connection Reader is exit, ConnID: %d, RemoteAddr: %s\n", c.ConnID, c.Conn.RemoteAddr().String())
	defer c.Stop()

	for {
		dp := NewDataPack()

		headData := make([]byte, dp.GetHeadLen())
		if _, err := io.ReadFull(c.GetTCPConnection(), headData); err != nil {
			fmt.Printf("Connection read error, ConnID: %d, error: %s\n", c.ConnID, err)
			break
		}

		msg, err := dp.Unpack(headData)
		if err != nil {
			fmt.Printf("Unpack error: %s\n", err)
		}

		var data []byte
		if msg.GetMsgLen() > 0 {
			data = make([]byte, msg.GetMsgLen())
			if _, err := io.ReadFull(c.GetTCPConnection(), data); err != nil {
				fmt.Printf("Connection read error, ConnID: %d, error: %s\n", c.ConnID, err)
				break
			}
		}
		msg.SetData(data)

		req := Request{
			conn: c,
			msg:  msg,
		}
		go func(request zicafe.IRequest) {
			c.Router.PreHandle(req)
			c.Router.Handle(req)
			c.Router.PostHandle(req)
		}(&req)

	}
}

func (c *Connection) Start() {
	fmt.Printf("Connection Start, ConnID: %d\n", c.ConnID)
	go c.StartReader()
}
func (c *Connection) Stop() {
	fmt.Printf("Connection Stop, ConnID: %d\n", c.ConnID)

	if c.isClosed == true {
		return
	}
	c.isClosed = true

	c.Conn.Close()

	close(c.ExitChan)
}
func (c *Connection) GetTCPConnection() *net.TCPConn {
	return c.Conn
}
func (c *Connection) GetConnID() uint32 {
	return c.ConnID
}
func (c *Connection) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}
func (c *Connection) SendMsg(msgId uint32, data []byte) error {
	if c.isClosed {
		fmt.Println("Connection is closed when send msg")
		return errors.New("connection is closed when send msg")
	}

	dp := NewDataPack()
	binaryMsg, err := dp.Pack(NewMessage(msgId, data))
	if err != nil {
		fmt.Printf("Pack error, msg id: %d\n", msgId)
		return err
	}

	if _, err := c.Conn.Write(binaryMsg); err != nil {
		fmt.Printf("Connection wirte error, msg id: %d, error: %s\n", msgId, err)
		return err
	}

	return nil
}
