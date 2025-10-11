package transfer

import (
	"net"
	"time"
)

type timedReadWriteCloser struct {
	conn    net.Conn
	timeout time.Duration
}

func (t *timedReadWriteCloser) Read(p []byte) (n int, err error) {
	if t.timeout > 0 {
		t.conn.SetReadDeadline(time.Now().Add(t.timeout))
	}
	defer t.conn.SetReadDeadline(time.Time{})
	return t.conn.Read(p)
}

func (t *timedReadWriteCloser) Write(p []byte) (n int, err error) {
	if t.timeout > 0 {
		t.conn.SetWriteDeadline(time.Now().Add(t.timeout))
	}
	defer t.conn.SetWriteDeadline(time.Time{})
	return t.conn.Write(p)
}

func (t *timedReadWriteCloser) Close() error {
	return t.conn.Close()
}
