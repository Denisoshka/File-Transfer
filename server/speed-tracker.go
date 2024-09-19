package server

import (
	"fmt"
	serverinterfaces "lab2/server-interfaces"

	"net"
	"sync"
	"time"
)

/*
	type SpeedInfo struct {
		lastUpdate time.Time
		recordedQ  int
	}
*/

type SpeedTracker struct {
	serverinterfaces.AbstractSpeedTracker
	trackDelay time.Duration
	data       map[string]*serverinterfaces.SpeedInfo
	mux        sync.Mutex
}

func NewSpeedTracker(trackDelay time.Duration) *SpeedTracker {
	return &SpeedTracker{
		trackDelay: trackDelay,
		data:       make(map[string]*serverinterfaces.SpeedInfo),
		mux:        sync.Mutex{},
	}
}

func (s *SpeedTracker) AddConnection(conn *net.Conn) (data *serverinterfaces.SpeedInfo) {
	s.mux.Lock()
	defer s.mux.Unlock()
	addr := (*conn).RemoteAddr().String()
	data = &serverinterfaces.SpeedInfo{LastUpdate: time.Now()}
	s.data[addr] = data
	return data
}

func (s *SpeedTracker) DeleteConnection(conn *net.Conn) {
	s.mux.Lock()
	defer s.mux.Unlock()
	addr := (*conn).RemoteAddr().String()
	delete(s.data, addr)
}

func (s *SpeedTracker) GetSpeedInfo(conn *net.Conn) *serverinterfaces.
	SpeedInfo {
	s.mux.Lock()
	defer s.mux.Unlock()
	addr := (*conn).RemoteAddr().String()
	return s.data[addr]
}

func (s *SpeedTracker) Launch() {
	for {
		s.mux.Lock()
		fmt.Print("\033[H\033[2J")
		for addr, speedInfo := range s.data {
			lastUpdate := speedInfo.LastUpdate
			writed := speedInfo.RecordedQ
			fmt.Println(
				addr, " bytes in second: ",
				float64(writed)/time.Since(lastUpdate).Seconds(),
			)
		}
		s.mux.Unlock()
		time.Sleep(s.trackDelay)
	}
}
