package server_interfaces

import (
	"errors"
	"time"
)

var (
	TrackerAlreadyLaunched = errors.New("tracker already launched")
)

type AbstractSpeedTracker interface {
	NewSpeedTracker(trackDelay time.Duration) (tracker *AbstractSpeedTracker)
	AddConnection(tag string) (data *ConnectionInfo)
	DeleteConnection(tag string)
	GetSpeedInfo(gat string) (data *ConnectionInfo)
	Launch() (err error)
}

type ConnectionInfo struct {
	expired    bool
	Start      time.Time
	LastUpdate time.Time
	total      uint64
	speed      float64
	avg        float64
}

func (s *ConnectionInfo) AddRecordedQ(recorded uint64) {
	now := time.Now()
	avgDiff := now.Sub(s.LastUpdate).Seconds()
	totalDiff := now.Sub(s.Start).Seconds()

	s.LastUpdate = now
	s.total += recorded
	s.speed = float64(recorded) / avgDiff
	s.avg = float64(s.total) / totalDiff
}

func (s *ConnectionInfo) Speed() float64 {
	return s.speed
}

func (s *ConnectionInfo) Avg() float64 {
	return s.avg
}

func (s *ConnectionInfo) Total() uint64 {
	return s.total
}

func (s *ConnectionInfo) Expired() bool {
	return s.expired
}

func (s *ConnectionInfo) MarkAsExpired() {
	s.expired = true
}
