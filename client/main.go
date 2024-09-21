package client

import (
	"fmt"
	"net"
	"strconv"
	"time"
)

var (
	maxConnectionDelay = time.Hour * 356
)

func main() {
	host, port, filePath, fileName, err := ParseFlags()
	if err != nil {
		fmt.Println(err)
		return
	}
	addr, err := net.ResolveTCPAddr(
		"tcp", net.JoinHostPort(host, strconv.Itoa(port)),
	)
	if err != nil {
		fmt.Println(err)
		return
	}

	uploader := NewTCPUploader(addr, maxConnectionDelay)
	err = uploader.Launch(filePath, fileName)
	if err != nil {
		fmt.Println(err)
	}
}
