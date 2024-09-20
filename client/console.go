package client

import (
	"errors"
	"flag"
	"net"
	"os"
)

var (
	PortNotSpecified = errors.New("port not specified")
	IncorrectPort    = errors.New("incorrect port")
	SpecifyFileName  = errors.New("specify file name")
)

func ParseFlags() (host string, port int, filePath string, fileName string, err error) {
	flag.StringVar(&host, "address", "", "Server address")
	flag.IntVar(&port, "port", -1, "Server port")
	flag.StringVar(&filePath, "filepath", "", "file path for upload")
	flag.StringVar(&fileName, "name", "", "file name for upload")
	flag.Parse()

	err = IsCorrectPort(port)
	if err != nil {
		return "", 0, "", "", err
	}

	err = IsCorrectAddress(host)
	if err != nil {
		return "", 0, "", "", err
	}

	err = IsCorrectFilePath(filePath)
	if err != nil {
		return "", 0, "", "", err
	}

	err = IsCorrectFileName(fileName)
	if err != nil {
		return "", 0, "", "", err
	}

	return host, port, filePath, fileName, nil
}

func IsCorrectAddress(addr string) (err error) {
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

func IsCorrectFilePath(filePath string) (err error) {
	_, err = os.Stat(filePath)
	return err
}

func IsCorrectFileName(fileName string) (err error) {
	if fileName == "" {
		return SpecifyFileName
	}
	return nil
}
