package znet

import (
	"errors"
	"fmt"
	"io"
	"net"
	"sync"
	"zinx/utils"
	"zinx/zicafe"
)

type Connection struct {
	TcpServer    zicafe.IServer
	Conn         *net.TCPConn
	ConnID       uint32
	isClosed     bool
	ExitChan     chan bool
	MsgChan      chan []byte
	MsgHandler   zicafe.IMsgHandler
	property     map[string]interface{}
	propertyLock sync.RWMutex
}

func NewConnection(server zicafe.IServer, conn *net.TCPConn, connID uint32, msgHandler zicafe.IMsgHandler) *Connection {
	c := &Connection{
		TcpServer:  server,
		Conn:       conn,
		ConnID:     connID,
		isClosed:   false,
		MsgHandler: msgHandler,
		ExitChan:   make(chan bool, 1),
		MsgChan:    make(chan []byte),
		property:   make(map[string]interface{}),
	}

	c.TcpServer.GetConnManager().Add(c)

	return c
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
		if utils.GlobalObject.WorkerPoolSize > 0 {
			c.MsgHandler.SendMsgToTaskQueue(req)
		} else {
			go c.MsgHandler.DoMsgHandler(req)
		}
	}
}

func (c *Connection) StartWriter() {
	fmt.Println("Connection Writer Goroutine is running")
	defer fmt.Printf("Connection Writer is exit, ConnID: %d, RemoteAddr: %s\n", c.ConnID, c.Conn.RemoteAddr().String())

	for {
		select {
		case data := <-c.MsgChan:
			if _, err := c.Conn.Write(data); err != nil {
				fmt.Printf("Connection Send error: %s\n", err)
				return
			}
		case <-c.ExitChan:
			return
		}
	}
}

func (c *Connection) Start() {
	fmt.Printf("Connection Start, ConnID: %d\n", c.ConnID)
	go c.StartReader()
	go c.StartWriter()
	c.TcpServer.CallOnConnStart(c)
}
func (c *Connection) Stop() {
	fmt.Printf("Connection Stop, ConnID: %d\n", c.ConnID)

	if c.isClosed == true {
		return
	}
	c.isClosed = true

	c.TcpServer.CallOnConnStop(c)

	c.Conn.Close()

	c.ExitChan <- true

	c.TcpServer.GetConnManager().Remove(c)

	close(c.ExitChan)
	close(c.MsgChan)
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

	c.MsgChan <- binaryMsg

	return nil
}

func (c *Connection) SetProperty(key string, property interface{}) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()

	c.property[key] = property
}

func (c *Connection) GetProperty(key string) (interface{}, error) {
	c.propertyLock.RLock()
	defer c.propertyLock.RUnlock()

	if property, ok := c.property[key]; ok {
		return property, nil
	}

	return nil, errors.New("property not exists")
}

func (c *Connection) RemoveProperty(key string) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()

	delete(c.property, key)
}
