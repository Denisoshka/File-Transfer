package server

import (
	"flag"
	"net"
	"strconv"
	"time"
)

const DefaultInactivity = time.Second * 100
const DefaultSpeedTrackDelay = time.Second * 3

func Start() {
	host, port, downloadDir, inactivityTimeout, speedTrackDelay, err := ParseFlags(
		DefaultInactivity, DefaultSpeedTrackDelay,
	)
	if err != nil {
		flag.PrintDefaults()

		LOG.Errorln(err)
		return
	}
	addr, err := net.ResolveTCPAddr(
		"tcp", net.JoinHostPort(host, strconv.Itoa(port)),
	)
	if err != nil {
		LOG.Errorln(err)
		return
	}
	spdtrcr := NewSpeedTracker(speedTrackDelay, speedTrackDelay*10)
	LOG.Infoln(
		"start server with args\n\t "+
			"downloadDir ", downloadDir, "\n\t"+
			"inactivityTimeout: ", inactivityTimeout, "\n\t"+
			"addr", addr,
	)
	acc := NewConnectionAcceptor(
		downloadDir, inactivityTimeout, addr, spdtrcr,
	)
	acc.Launch()
}
