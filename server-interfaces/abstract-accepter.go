package server_interfaces

import (
	"lab2/server"
	"net"
	"time"
)

type AbstractSpeedTracker interface {
	NewSpeedTracker(trackDelay time.Duration) *AbstractSpeedTracker
	AddConnection(conn *net.Conn) (data *server.ConnectionInfo)
	DeleteConnection(conn *net.Conn)
	GetSpeedInfo(conn *net.Conn) *server.ConnectionInfo
	Launch()
}
