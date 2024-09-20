package utils

import (
	"net"
	"os"
	"time"
)

const (
	DirPerm  = 0644
	FilePerm = 0644
)

func CreateUploadsDir(dirPath string) error {
	return os.MkdirAll(dirPath, DirPerm)
}

func PrepareFile(filePath string, fileSize int64) (file *os.File,
	fileExists bool, err error) {
	file, err = os.OpenFile(filePath, os.O_CREATE|os.O_EXCL|os.O_RDWR, FilePerm)
	if err != nil {
		if os.IsExist(err) {
			return nil, true, err
		}
		return nil, false, err
	}
	err = file.Truncate(0)
	if err != nil {
		_ = file.Close()
		return nil, false, err
	}
	err = file.Truncate(fileSize)
	if err != nil {
		_ = file.Close()
		return nil, false, err
	}
	return file, false, nil
}

func ConnReadN(conn net.Conn, data []byte, n int,
	maxSocketInactivity time.Duration) (total int, err error) {
	total = 0
	for i := 0; i < n; i++ {
		err := conn.SetReadDeadline(time.Now().Add(maxSocketInactivity))
		if err != nil {
			return 0, err
		}
		start, err := conn.Read(data[total:n])
		if err != nil {
			break
		}
		total += start
	}
	return total, err
}

func ConnWriteN(conn net.Conn, data []byte, n int) (total int, err error) {
	total = 0
	for i := 0; i < n; i++ {
		start, err := conn.Write(data[total:n])
		if err != nil {
			break
		}
		total += start
	}
	return total, err
}

func FileWriteN(file *os.File, data []byte, n int) (total int, err error) {
	total = 0
	for i := 0; i < n; i++ {
		start, err := file.Write(data[total:n])
		if err != nil {
			break
		}
		total += start
	}

}
