package server

import (
	"errors"
	"flag"
	"net"
	"time"
)

var (
	PortNotSpecified      = errors.New("port not specified")
	IncorrectPort         = errors.New("incorrect port")
	DirectoryNotSpecified = errors.New("directory not specified")
)

func ParseFlags(defaultInactivity time.Duration,
	defaultSpeedTrackDelay time.Duration, defaultDownloadDir string) (
	host string, port int,
	downloadDir string,
	maxInactivity time.Duration, speedTrackDelay time.Duration, err error) {
	flag.StringVar(&host, "address", "", "Start address")
	flag.IntVar(&port, "port", -1, "Start port")
	flag.StringVar(
		&downloadDir, "downloaddir", defaultDownloadDir, "Download directory path",
	)
	flag.DurationVar(
		&maxInactivity, "inactivitytimeout", defaultInactivity,
		"Max connection inactivity time",
	)
	flag.DurationVar(
		&speedTrackDelay, "speedtrackdelay",
		defaultSpeedTrackDelay, "Speed track delay",
	)
	flag.Parse()

	err = IsCorrectPort(port)
	if err != nil {
		return "", 0, "", 0, 0, err
	}

	err = IsCorrectAddress(host)
	if err != nil {
		return "", 0, "", 0, 0, err
	}

	err = IsCorrectDir(downloadDir)
	if err != nil {
		return "", 0, "", 0, 0, err
	}
	return
}

func IsCorrectAddress(addr string) (err error) {
	if addr == "" {
		return nil
	}

	_, err = net.ResolveIPAddr("ip", addr)
	if err != nil {
		return err
	}
	return nil
}

func IsCorrectPort(port int) (err error) {
	if port == -1 {
		return PortNotSpecified
	}

	if !(0 <= port && port <= 65536) {
		return IncorrectPort
	}
	return nil
}

func IsCorrectDir(dir string) (err error) {
	if dir == "" {
		return DirectoryNotSpecified
	}
	return nil
}
