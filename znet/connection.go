package znet

import (
	"fmt"
	"net"
	"zinx/utils"
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
		buf := make([]byte, utils.GlobalObject.MaxPackageSize)
		_, err := c.Conn.Read(buf)
		if err != nil {
			fmt.Printf("Connection read error, ConnID: %d, error: %s\n", c.ConnID, err)
			continue
		}
		req := Request{
			conn: c,
			data: buf,
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
func (c *Connection) Send(data []byte) error {
	return nil
}
