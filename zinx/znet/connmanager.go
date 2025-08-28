package znet

import (
	"errors"
	"fmt"
	"sync"
	"zinx/ziface"
)

type ConnManager struct {
	connections map[string]ziface.IConnection
	connLock    sync.RWMutex
}

func (c *ConnManager) Add(conn ziface.IConnection) {
	c.connLock.Lock()
	defer c.connLock.Unlock()

	c.connections[conn.GetConnID()] = conn

	fmt.Printf("connection add to ConnManager successfully :conn num = %d", c.Len())
}

func (c *ConnManager) Remove(conn ziface.IConnection) {
	c.connLock.Lock()
	defer c.connLock.Unlock()

	//c.connections[conn.GetConnID()] = conn
	delete(c.connections, conn.GetConnID())

	fmt.Printf("connection add to ConnManager successfully :conn num = %d", c.Len())
}

func (c *ConnManager) Get(id string) (ziface.IConnection, error) {
	c.connLock.Lock()
	defer c.connLock.Unlock()

	if conn, ok := c.connections[id]; ok {
		return conn, nil
	} else {
		return nil, errors.New("get conn error ,not found")
	}

}

func (c *ConnManager) Len() int {
	return len(c.connections)
}

func (c *ConnManager) ClearConn() {
	c.connLock.Lock()
	defer c.connLock.Unlock()

	for connID, conn := range c.connections {
		conn.Stop()
		delete(c.connections, connID)
	}
	fmt.Printf("Clear All Connection successfully:conn num = %d", c.Len())

}

func NewConnManager() *ConnManager {
	return &ConnManager{
		connections: make(map[string]ziface.IConnection),
	}
}
