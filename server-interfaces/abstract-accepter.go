package server_interfaces

import (
	"net"
	"time"
)

type SpeedInfo struct {
	LastUpdate time.Time
	RecordedQ  int
}

type AbstractSpeedTracker interface {
	NewSpeedTracker(trackDelay time.Duration) *AbstractSpeedTracker
	AddConnection(conn *net.Conn) (data *SpeedInfo)
	DeleteConnection(conn *net.Conn)
	GetSpeedInfo(conn *net.Conn) *SpeedInfo
	Launch()
}
