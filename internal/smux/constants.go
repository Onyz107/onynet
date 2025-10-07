package smux

import (
	"io"
	"net"
	"sync"
	"time"
)

type Communicator interface {
	net.Conn

	Send(b []byte, timeout time.Duration) error
	NewStreamedSender(timeout time.Duration) io.WriteCloser
	SendSerialized(b []byte, timeout time.Duration) error
	SendEncrypted(b []byte, timeout time.Duration) error
	NewStreamedEncryptedSender(timeout time.Duration) (io.WriteCloser, error)

	Receive(b []byte, timeout time.Duration) error
	NewStreamedReceiver(timeout time.Duration) io.ReadCloser
	ReceiveSerialized(b []byte, timeout time.Duration) (uint64, error)
	ReceiveEncrypted(b []byte, timeout time.Duration) (uint64, error)
	NewStreamedEncryptedReceiver(timeout time.Duration) (io.ReadCloser, error)
}

var headerPool = sync.Pool{
	New: func() any {
		buf := make([]byte, 2)
		return &buf
	},
}
