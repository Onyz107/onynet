package transfer

import (
	"net"
	"time"
)

type timedReader struct {
	r       net.Conn
	timeout time.Duration
}

func (t *timedReader) Read(p []byte) (n int, err error) {
	if t.timeout > 0 {
		t.r.SetReadDeadline(time.Now().Add(t.timeout))
	}
	defer t.r.SetReadDeadline(time.Time{})
	return t.r.Read(p)
}

func (t *timedReader) Close() error {
	return t.r.Close()
}

type timedWriter struct {
	w       net.Conn
	timeout time.Duration
}

func (t *timedWriter) Write(p []byte) (n int, err error) {
	if t.timeout > 0 {
		t.w.SetWriteDeadline(time.Now().Add(t.timeout))
	}
	defer t.w.SetWriteDeadline(time.Time{})
	return t.w.Write(p)
}

func (t *timedWriter) Close() error {
	return t.w.Close()
}
