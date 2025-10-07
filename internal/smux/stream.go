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

// Read reads data from the stream into the provided buffer.
//
// Possible errors:
//   - ErrCtxCancelled: context was cancelled when trying to read data
func (s *Stream) Read(b []byte) (n int, err error) {
	select {
	case <-s.ctx.Done():
		return 0, intErrors.ErrCtxCancelled
	default:
		return s.stream.Read(b)
	}
}

// Write writes data to the connection.
//
// Possible errors:
//   - ErrCtxCancelled: context was cancelled when trying to write data
func (s *Stream) Write(b []byte) (n int, err error) {
	select {
	case <-s.ctx.Done():
		return 0, intErrors.ErrCtxCancelled
	default:
		return s.stream.Write(b)
	}
}

// Send sends raw bytes with timeout.
//
// Possible errors:
//   - ErrWrite: failed to send data through the stream
//   - ErrShortWrite: data sent was shorter than expected
func (s *Stream) Send(b []byte, timeout time.Duration) error {
	return transfer.Send(s.stream, b, timeout)
}

// NewStreamedSender returns an io.WriteCloser that allows
// the caller to directly write data to the stream and set a timeout.
func (s *Stream) NewStreamedSender(timeout time.Duration) io.WriteCloser {
	return transfer.NewStreamedSender(s.stream, timeout)
}

// SendSerialized sends serialized data with length header.
// SendSerialized adds 8 bytes to the data for length prefixing.
//
// Possible errors:
//   - ErrWrite: failed to send data through the stream
//   - ErrShortWrite: data sent was shorter than expected
func (s *Stream) SendSerialized(b []byte, timeout time.Duration) error {
	return transfer.SendSerialized(s.stream, b, timeout)
}

// SendEncrypted sends AES-GCM encrypted data.
// SendEncrypted adds 24 bytes to the data for encryption and length prefixing.
//
// Possible errors:
//   - ErrWrite: failed to send data through the stream
//   - ErrShortWrite: data sent was shorter than expected
//   - ErrAesKey: the aesKey is nil, meaning authentication is not enabled
//   - ErrCipher: the aesKey has an invalid key size
//   - ErrGCM: failed to create GCM
func (s *Stream) SendEncrypted(b []byte, timeout time.Duration) error {
	return transfer.SendEncrypted(s.stream, b, s.aesKey, timeout)
}

// NewStreamedEncryptedSender returns an io.WriteCloser that encrypts data as it is written to the stream.
//
// Possible errors:
//   - ErrAesKey: the aesKey is nil, meaning authentication is not enabled
//   - ErrWrite: failed to send the nonce through the stream
//   - ErrShortWrite: the nonce sent was shorter than expected
//   - ErrStreamCipher: failed to create an AES-CTR stream
//   - ErrCipher: invalid key size
func (s *Stream) NewStreamedEncryptedSender(timeout time.Duration) (io.WriteCloser, error) {
	return transfer.NewStreamedEncryptedSender(s.stream, s.aesKey, timeout)
}

// Receive reads data into buffer with timeout.
//
// Possible errors:
//   - ErrRead: failed to receive data from the stream
func (s *Stream) Receive(b []byte, timeout time.Duration) error {
	return transfer.Receive(s.stream, b, timeout)
}

// NewStreamedReceiver returns an io.ReadCloser that allows
// the caller to directly read data from the stream and set a timeout.
func (s *Stream) NewStreamedReceiver(timeout time.Duration) io.ReadCloser {
	return transfer.NewStreamedReceiver(s.stream, timeout)
}

// ReceiveSerialized reads serialized data with length header.
// The buffer provided should be at least 8 bytes bigger than the data expected to receive.
//
// Possible errors:
//   - ErrSmallBuffer: the buffer provided was too small to receive the incoming data
//   - ErrRead: failed to receive data from the stream
func (s *Stream) ReceiveSerialized(b []byte, timeout time.Duration) (uint64, error) {
	return transfer.ReceiveSerialized(s.stream, b, timeout)
}

// ReceiveEncrypted reads AES-GCM encrypted data.
// The buffer provided should be at least 28 bytes bigger than the data expected to receive.
//
// Possible errors:
//   - ErrAesKey: the aesKey is nil, meaning authentication is not enabled
//   - ErrSmallBuffer: the buffer provided was too small to receive the incoming data
//   - ErrRead: failed to receive data from the stream
//   - ErrCipher: invalid key size
//   - ErrGCM: failed to create GCM
//   - ErrShort: ciphertext received is malformed because it is too short
//   - ErrDecrypt: failed to decrypt the received data
func (s *Stream) ReceiveEncrypted(b []byte, timeout time.Duration) (uint64, error) {
	return transfer.ReceiveEncrypted(s.stream, b, s.aesKey, timeout)
}

// NewStreamedEncryptedReceiver returns an io.ReadCloser that decrypts data as it comes from the stream.
//
// Possible errors:
//   - ErrAesKey: the aesKey is nil, meaning authentication is not enabled
//   - ErrRead: failed to receive the nonce from the stream
//   - ErrStreamCipher: failed to create an AES-CTR stream
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
