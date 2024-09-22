package server_interfaces

import (
	"net"
	"time"
)

type AbstractSpeedTracker interface {
	NewSpeedTracker(trackDelay time.Duration) *AbstractSpeedTracker
	AddConnection(conn *net.Conn) (data *ConnectionInfo)
	DeleteConnection(conn *net.Conn)
	GetSpeedInfo(conn *net.Conn) *ConnectionInfo
	Launch()
}

type ConnectionInfo struct {
	Expired    bool
	Start      time.Time
	LastUpdate time.Time
	total      uint64
	Speed      float64
	Avg        float64
}

func (s *ConnectionInfo) UpdateSpeed(recorded uint64) (avg float64) {
	now := time.Now()
	avgDiff := now.Sub(s.LastUpdate).Seconds()
	totalDiff := now.Sub(s.Start).Seconds()

	s.LastUpdate = now
	s.total += recorded
	s.Speed = float64(recorded) / avgDiff
	s.Avg = float64(s.total) / totalDiff

	return s.Speed
}
