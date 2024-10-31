package utils

import (
	"encoding/binary"
	"io"
	"lab2/utils/requests"
	"net"
	"time"
)

func ConnReadN(conn net.Conn, data []byte, n int,
	maxSocketInactivity time.Duration) (total int, err error) {
	total = 0
	for start := 0; total < n; {

		if err = conn.SetReadDeadline(time.Now().Add(maxSocketInactivity)); err != nil {
			break
		}
		if start, err = conn.Read(data[total:n]); err != nil {
			break
		}
		total += start
	}
	return total, err
}

/*func ConnWriteN(conn net.Conn, data []byte, n int) (total int, err error) {
	total = 0
	for start := 0; total < n; total += start {
		start, err = conn.Write(data[total:n])
		if err != nil {
			break
		}
	}
	return total, err
}*/

/*func FileWriteN(file *os.File, data []byte, n int) (total int, err error) {
	total = 0
	for start := 0; total < n; total += start {
		start, err = file.Write(data[total:n])
		if err != nil {
			break
		}
	}
	return total, err
}*/

/*func FileReadN(file *os.File, data []byte, n int) (total int, err error) {
	total = 0
	for start := 0; total < n; total += start {
		start, err = file.Read(data[total:n])
		if err != nil {
			break
		}
	}
	return total, err
}*/

func ReadRequest(conn net.Conn, maxConnInactivityDelay time.Duration,
	req requests.AbstractRequest, maxReqSize int64, buf []byte) (err error) {
	reqSizeBuf := make([]byte, 4)
	err = conn.SetReadDeadline(time.Now().Add(maxConnInactivityDelay))
	if err != nil {
		return err
	}
	_, err = io.ReadFull(conn, reqSizeBuf[:4])
	//_, err = ConnReadN(conn, reqSizeBuf, 4, maxConnInactivityDelay)
	if err != nil {
		return err
	}
	initialReqSize := int32(binary.BigEndian.Uint32(reqSizeBuf))
	if int64(initialReqSize) > maxReqSize {
		return requests.IncorrectRequestSize
	}

	if buf == nil || len(buf) < int(initialReqSize) {
		buf = make([]byte, initialReqSize)
	}

	if err = conn.SetReadDeadline(time.Now().Add(maxConnInactivityDelay)); err != nil {
		return err
	}
	if _, err = io.ReadFull(conn, buf[4:int(initialReqSize)-4]); err != nil {
		return err
	}
	/*_, err = ConnReadN(
		conn, buf[4:], int(initialReqSize)-4, maxConnInactivityDelay,
	)*/
	binary.BigEndian.PutUint32(buf[0:4], uint32(initialReqSize))
	err = req.DecodeFrom(buf)
	if err != nil {
		return err
	}
	return nil
}
