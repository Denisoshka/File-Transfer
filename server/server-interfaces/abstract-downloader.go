package server_interfaces

import (
	"net"
	"time"
)

type AbstractTCPDownloader interface {
	NewTCPDownloader(dirPath string, maxDelay time.Duration,
		conn *net.TCPConn, tracker *AbstractSpeedTracker) *AbstractTCPDownloader
	Launch() (err error)
}
