package server

import (
	"errors"
	serverinterfaces "lab2/server/server-interfaces"
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
	err := createUploadsDir(a.saveDir)
	if err != nil {
		panic(err)
	}
	listener, err := net.ListenTCP("tcp", a.addr)
	if err != nil {
		panic(err)
	}

	defer func(listener *net.TCPListener) { _ = listener.Close() }(listener)
	go func() {
		err = a.speedTracker.Launch()
		if err != nil {
			panic(err)
		}
	}()

	for {
		conn, err := listener.AcceptTCP()
		if err != nil {
			if errors.Is(err, syscall.EINVAL) {
				continue
			}
		}

		LOG.Infoln("new connection from ", conn.RemoteAddr())
		d := NewTCPDownloader(a.saveDir, a.maxInactivity, conn, a.speedTracker)
		go func() { _ = d.Launch() }()
	}
}
