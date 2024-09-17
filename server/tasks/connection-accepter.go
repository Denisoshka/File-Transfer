package tasks

import (
	lab2 "lab2"
	"net"
	"os"
	"sync"
)

type ConnectionAcceptor struct {
	addr    *net.TCPAddr
	dirPath string
}

func NewConnectionAcceptor(addr *net.TCPAddr, dirPath string) *ConnectionAcceptor {
	return &ConnectionAcceptor{
		addr:    addr,
		dirPath: dirPath,
	}
}

func (c *ConnectionAcceptor) StartNetWorker(g *sync.WaitGroup) {
	err := os.MkdirAll(c.dirPath, 0666)
	if err != nil {
		panic(err)
	}

	addr := c.addr
	conn, err := net.ListenTCP("tcp", addr)
	if err != nil {
		panic(err)
	}

	defer func(conn net.Listener) {
		_ = conn.Close()
	}(conn)

	for {
		conn, err := conn.AcceptTCP()
		if err != nil {
			lab2.Log.Infoln("new connection accept failed " + err.Error())
			continue
		}
		worker := NewConnection(conn)
		go worker.StartNetWorker()
	}
}
