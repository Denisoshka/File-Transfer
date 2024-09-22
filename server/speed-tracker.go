package server

import (
	"fmt"
	serverinterfaces "lab2/server/server-interfaces"
	"sync"
	"time"
)

type SpeedTracker struct {
	serverinterfaces.AbstractSpeedTracker
	expiredCheckDelay time.Duration
	trackDelay        time.Duration
	data              map[string]*serverinterfaces.ConnectionInfo
	mux               sync.RWMutex
}

func NewSpeedTracker(trackDelay time.Duration, expiredCheckDelay time.Duration) *SpeedTracker {
	return &SpeedTracker{
		expiredCheckDelay: expiredCheckDelay,
		trackDelay:        trackDelay,
		data:              make(map[string]*serverinterfaces.ConnectionInfo),
		mux:               sync.RWMutex{},
	}
}

func (s *SpeedTracker) AddConnection(tag string) (data *serverinterfaces.ConnectionInfo) {
	s.mux.Lock()
	data = &serverinterfaces.ConnectionInfo{
		Expired:    false,
		LastUpdate: time.Now(),
		Start:      time.Now(),
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

func (s *SpeedTracker) GetSpeedInfo(tag string) (data *serverinterfaces.ConnectionInfo) {
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
				"connection ", addr, " speed ", speedInfo.Speed, " avg ", speedInfo.Avg,
			)
		}
		s.mux.RUnlock()

		if time.Now().Sub(lastCheck) > s.expiredCheckDelay {
			s.mux.Lock()
			for addr, info := range s.data {
				if info.Expired {
					delete(s.data, addr)
				}
			}
			s.mux.Unlock()
		}
		time.Sleep(s.trackDelay)
	}
}

func (s *SpeedTracker) cleanUp() {}
