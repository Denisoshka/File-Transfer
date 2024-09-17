package tasks

import (
	"net"
	"sync"
	lab2 "xyi"
)

type ConnectionAccept struct {
	addr *net.TCPAddr
}

func (c *ConnectionAccept) StartNetWorker(g *sync.WaitGroup) {
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
