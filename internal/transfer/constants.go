package transfer

import (
	"net"
	"sync"
	"time"
)

var headerPool = sync.Pool{
	New: func() any {
		buf := make([]byte, 8)
		return &buf
	},
}

var timedReadWriteCloserCache sync.Map // conn net.Conn -> timedReadWriteCloser

func getTimedReadWriteCloser(conn net.Conn, timeout time.Duration) *timedReadWriteCloser {
	if val, ok := timedReadWriteCloserCache.Load(conn); ok {
		return val.(*timedReadWriteCloser)
	}

	timedConn := &timedReadWriteCloser{conn: conn, timeout: timeout}
	timedReadWriteCloserCache.Store(conn, timedConn)
	return timedConn
}
