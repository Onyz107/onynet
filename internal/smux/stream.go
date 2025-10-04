package smux

import (
	"context"
	"io"
	"net"
	"time"

	intErrors "github.com/Onyz107/onynet/errors"
	"github.com/Onyz107/onynet/internal/transfer"
	"github.com/xtaci/smux"
)

type Stream struct {
	stream *smux.Stream
	aesKey []byte
	ctx    context.Context
}

func (s *Stream) Read(b []byte) (n int, err error) {
	select {
	case <-s.ctx.Done():
		return 0, intErrors.ErrCtxCancelled
	default:
		return s.stream.Read(b)
	}
}

func (s *Stream) Write(b []byte) (n int, err error) {
	select {
	case <-s.ctx.Done():
		return 0, intErrors.ErrCtxCancelled
	default:
		return s.stream.Write(b)
	}
}

// Send sends raw bytes with timeout.
func (s *Stream) Send(b []byte, timeout time.Duration) error {
	return transfer.Send(s.stream, b, timeout)
}

// NewStreamedSender returns an io.WriteCloser that allows
// the caller to directly write data to the stream and set a timeout.
func (s *Stream) NewStreamedSender(timeout time.Duration) io.WriteCloser {
	return transfer.NewStreamedSender(s.stream, timeout)
}

// SendSerialized sends serialized data with length header.
func (s *Stream) SendSerialized(b []byte, timeout time.Duration) error {
	return transfer.SendSerialized(s.stream, b, timeout)
}

// SendEncrypted sends AES-GCM encrypted data.
func (s *Stream) SendEncrypted(b []byte, timeout time.Duration) error {
	return transfer.SendEncrypted(s.stream, b, s.aesKey, timeout)
}

// NewStreamedEncryptedSender returns an io.WriteCloser that encrypts data as it is written to the stream.
func (s *Stream) NewStreamedEncryptedSender(timeout time.Duration) (io.WriteCloser, error) {
	return transfer.NewStreamedEncryptedSender(s.stream, s.aesKey, timeout)
}

// Receive reads data into buffer with timeout.
func (s *Stream) Receive(b []byte, timeout time.Duration) error {
	return transfer.Receive(s.stream, b, timeout)
}

// NewStreamedReceiver returns an io.ReadCloser that allows
// the caller to directly read data from the stream and set a timeout.
func (s *Stream) NewStreamedReceiver(timeout time.Duration) io.ReadCloser {
	return transfer.NewStreamedReceiver(s.stream, timeout)
}

// ReceiveSerialized reads serialized data with length header.
func (s *Stream) ReceiveSerialized(b []byte, timeout time.Duration) (uint64, error) {
	return transfer.ReceiveSerialized(s.stream, b, timeout)
}

// ReceiveEncrypted reads AES-GCM encrypted data.
func (s *Stream) ReceiveEncrypted(b []byte, timeout time.Duration) (uint64, error) {
	return transfer.ReceiveEncrypted(s.stream, b, s.aesKey, timeout)
}

// NewStreamedEncryptedReceiver returns an io.ReadCloser that decrypts data as it comes from the stream.
func (s *Stream) NewStreamedEncryptedReceiver(timeout time.Duration) (io.ReadCloser, error) {
	return transfer.NewStreamedEncryptedReceiver(s.stream, s.aesKey, timeout)
}

func (s *Stream) Close() error {
	return s.stream.Close()
}

func (s *Stream) LocalAddr() net.Addr {
	return s.stream.LocalAddr()
}

func (s *Stream) RemoteAddr() net.Addr {
	return s.stream.RemoteAddr()
}

func (s *Stream) SetDeadline(t time.Time) error {
	return s.stream.SetDeadline(t)
}

func (s *Stream) SetReadDeadline(t time.Time) error {
	return s.stream.SetReadDeadline(t)
}

func (s *Stream) SetWriteDeadline(t time.Time) error {
	return s.stream.SetWriteDeadline(t)
}
