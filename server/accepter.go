package server

import (
	"errors"
	"lab2"
	serverinterfaces "lab2/server-interfaces"
	"lab2/utils"
	"net"
	"syscall"
	"time"
)

type ConnectionAcceptor struct {
	serverinterfaces.AbstractConnectionAcceptor
	saveDir       string
	addr          *net.TCPAddr
	maxInactivity time.Duration
	speedTracker  *SpeedTracker
}

func NewConnectionAcceptor(saveDir string,
	maxInactivity time.Duration,
	addr *net.TCPAddr,
	tracker *SpeedTracker) *ConnectionAcceptor {
	return &ConnectionAcceptor{
		saveDir:       saveDir,
		addr:          addr,
		maxInactivity: maxInactivity,
		speedTracker:  tracker,
	}
}

func (a *ConnectionAcceptor) Launch() {
	err := utils.CreateUploadsDir(a.saveDir)
	if err != nil {
		panic(err)
	}
	listener, err := net.ListenTCP("tcp", a.addr)
	if err != nil {
		panic(err)
	}
	defer func(listener *net.TCPListener) { _ = listener.Close() }(listener)

	for {
		conn, err := listener.AcceptTCP()
		if err != nil {
			if errors.Is(err, syscall.EINVAL) {
				continue
			}
		}

		lab2.Log.Infoln("new connection from ", conn.RemoteAddr())
		d := NewTCPDownloader(a.saveDir, a.maxInactivity, conn)
		go func() { _ = d.Launch() }()
	}
}
