package server_interfaces

import (
	"net"
	"time"
)

type AbstractConnectionAcceptor interface {
	NewConnectionAcceptor(maxInactivity time.Duration, addr *net.TCPAddr,
		tracker *AbstractSpeedTracker) *AbstractConnectionAcceptor
	Launch()
}
