package znet

import (
	"errors"
	"sync"
	"zinx/zicafe"
)

type ConnManager struct {
	connections map[uint32]zicafe.IConnection
	connLock sync.RWMutex
}

func NewConnManager() *ConnManager {
	return &ConnManager{
		connections: make(map[uint32]zicafe.IConnection),
	}
}

func (cm *ConnManager) Add(connection zicafe.IConnection) {
	cm.connLock.Lock()
	defer cm.connLock.Unlock()

	cm.connections[connection.GetConnID()] = connection
}

func (cm *ConnManager) Remove(connection zicafe.IConnection) {
	cm.connLock.Lock()
	defer cm.connLock.Unlock()

	delete(cm.connections, connection.GetConnID())
}

func (cm *ConnManager) Get(connID uint32) (zicafe.IConnection, error) {
	cm.connLock.RLock()
	defer cm.connLock.RUnlock()

	if 	conn, ok := cm.connections[connID]; ok {
		return conn, nil
	}
	return nil, errors.New("connection not fount")
}

func (cm *ConnManager) Len() int {
	return len(cm.connections)
}

func (cm *ConnManager) Clear() {
	cm.connLock.Lock()
	defer cm.connLock.Unlock()

	for id, conn := range cm.connections{
		conn.Stop()
		delete(cm.connections, id)
	}
}



