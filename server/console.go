package server

import (
	"errors"
	"flag"
	"net"
)

var (
	PortNotSpecified       = errors.New("port not specified")
	IncorrectPort          = errors.New("incorrect port")
	DirrectoryNotSpecified = errors.New("directory not specified")
)

func ParseFlags() (host string, port int, downloadDir string, err error) {
	flag.StringVar(&host, "address", "", "Server address")
	flag.IntVar(&port, "port", -1, "Server port")
	flag.StringVar(&downloadDir, "downloaddir", "", "Download directory path")
	flag.Parse()

	err = IsCorrectPort(port)
	if err != nil {
		return "", 0, "", err
	}

	err = IsCorrectAddress(host)
	if err != nil {
		return "", 0, "", err
	}

	err = IsCorrectDir(downloadDir)
	if err != nil {
		return "", 0, "", err
	}

	return host, port, downloadDir, nil
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
		return DirrectoryNotSpecified
	}
	return nil
}
