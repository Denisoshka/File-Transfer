package server

import (
	"lab2"
	"net"
	"strconv"
	"time"
)

const DefaultInactivity = time.Second * 10
const DefaultSpeedTrackDelay = time.Second * 3

func main() {
	host, port, downloadDir, inactivtyTimeout, speedTrackDelay, err := ParseFlags(
		DefaultInactivity, DefaultSpeedTrackDelay,
	)
	if err != nil {
		lab2.Log.Errorln(err)
		return
	}
	addr, err := net.ResolveTCPAddr(
		"tcp", net.JoinHostPort(host, strconv.Itoa(port)),
	)
	if err != nil {
		lab2.Log.Errorln(err)
		return
	}
	spdtrcr := NewSpeedTracker(speedTrackDelay)
	acc := NewConnectionAcceptor(downloadDir, inactivtyTimeout, addr, spdtrcr)
	acc.Launch()
}
