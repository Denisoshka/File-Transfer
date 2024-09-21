package server

import (
	"fmt"
	serverinterfaces "lab2/server-interfaces"
	"sync"
	"time"
)

type ConnectionInfo struct {
	expired    bool
	start      time.Time
	lastUpdate time.Time
	total      uint64
	speed      float64
	avg        float64
}

func (s *ConnectionInfo) updateSpeed(recorded uint64) (avg float64) {
	now := time.Now()
	avgDiff := now.Sub(s.lastUpdate).Seconds()
	totalDiff := now.Sub(s.start).Seconds()

	s.lastUpdate = now
	s.total += recorded
	s.speed = float64(recorded) / avgDiff
	s.avg = float64(s.total) / totalDiff

	return s.speed
}

type SpeedTracker struct {
	serverinterfaces.AbstractSpeedTracker
	expiredCheckDelay time.Duration
	trackDelay        time.Duration
	data              map[string]*ConnectionInfo
	mux               sync.RWMutex
}

func NewSpeedTracker(trackDelay time.Duration, expiredCheckDelay time.Duration) *SpeedTracker {
	return &SpeedTracker{
		expiredCheckDelay: expiredCheckDelay,
		trackDelay:        trackDelay,
		data:              make(map[string]*ConnectionInfo),
		mux:               sync.RWMutex{},
	}
}

func (s *SpeedTracker) AddConnection(tag string) (data *ConnectionInfo) {
	s.mux.Lock()
	data = &ConnectionInfo{
		expired:    false,
		lastUpdate: time.Now(),
		start:      time.Now(),
	}
	s.data[tag] = data
	s.mux.Unlock()
	return data
}

func (s *SpeedTracker) DeleteConnection(tag string) {
	s.mux.Lock()
	delete(s.data, tag)
	s.mux.Unlock()
}

func (s *SpeedTracker) GetSpeedInfo(tag string) (data *ConnectionInfo) {
	s.mux.RLock()
	data = s.data[tag]
	s.mux.RUnlock()
	return
}

func (s *SpeedTracker) Launch() {
	lastCheck := time.Now()
	for {
		s.mux.RLock()
		fmt.Print("\033[H\033[2J")
		for addr, speedInfo := range s.data {
			fmt.Println(
				"connection ", addr, " speed ", speedInfo.speed, " avg ", speedInfo.avg,
			)
		}
		s.mux.RUnlock()

		if time.Now().Sub(lastCheck) > s.expiredCheckDelay {
			s.mux.Lock()
			for addr, info := range s.data {
				if info.expired {
					delete(s.data, addr)
				}
			}
			s.mux.Unlock()
		}
		time.Sleep(s.trackDelay)
	}
}

func (s *SpeedTracker) cleanUp() {}
