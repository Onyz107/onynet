package smux

import (
	"context"
	"io"
	"net"
	"sync"
	"time"
)

type Handler interface {
	OpenStream(name string, ctx context.Context, timeout time.Duration) (*Stream, error)
	AcceptStream(name string, ctx context.Context, timeout time.Duration) (*Stream, error)
}

type Communicator interface {
	// Conn is a generic stream-oriented network connection.
	//
	// Multiple goroutines may invoke methods on a Conn simultaneously.
	net.Conn

	// Send sends raw bytes with timeout.
	//
	// Possible errors:
	//   - ErrWrite: failed to send data through the stream
	//   - ErrShortWrite: data sent was shorter than expected
	//   - ErrTimeout: timeout occurred when receiving data from the stream
	Send(b []byte, timeout time.Duration) error

	// NewStreamedSender returns an io.WriteCloser that allows
	// the caller to directly write data to the stream and set a timeout.
	NewStreamedSender(timeout time.Duration) io.WriteCloser

	// SendSerialized sends serialized data with length header.
	//
	// SendSerialized adds 8 bytes to the data for length prefixing,
	// so for example if you want to send a 1024 bytes message the function will send 1032 bytes.
	//
	// Possible errors:
	//   - ErrWrite: failed to send data through the stream
	//   - ErrShortWrite: data sent was shorter than expected
	//   - ErrTimeout: timeout occurred when receiving data from the stream
	SendSerialized(b []byte, timeout time.Duration) error

	// SendEncrypted sends AES-GCM encrypted data.
	//
	// SendEncrypted adds 24 bytes to the data for encryption and length prefixing.
	// so for example if you want to send a 1024 bytes message the function will send 1048 bytes.
	//
	// Possible errors:
	//   - ErrWrite: failed to send data through the stream
	//   - ErrShortWrite: data sent was shorter than expected
	//   - ErrAesKey: the aesKey is nil, meaning authentication is not enabled
	//   - ErrCipher: the aesKey has an invalid key size
	//   - ErrGCM: failed to create GCM
	//   - ErrTimeout: timeout occurred when receiving data from the stream
	SendEncrypted(b []byte, timeout time.Duration) error

	// NewStreamedEncryptedSender returns an io.WriteCloser that encrypts data as it is written to the stream.
	//
	// Possible errors:
	//   - ErrAesKey: the aesKey is nil, meaning authentication is not enabled
	//   - ErrWrite: failed to send the nonce through the stream
	//   - ErrShortWrite: the nonce sent was shorter than expected
	//   - ErrStreamCipher: failed to create an AES-CTR stream
	//   - ErrCipher: invalid key size
	NewStreamedEncryptedSender(timeout time.Duration) (io.WriteCloser, error)

	// Receive reads data into buffer with timeout.
	//
	// Possible errors:
	//   - ErrRead: failed to receive data from the stream
	//   - ErrTimeout: timeout occurred when receiving data from the stream
	Receive(b []byte, timeout time.Duration) error

	// NewStreamedReceiver returns an io.ReadCloser that allows
	// the caller to directly read data from the stream and set a timeout.
	NewStreamedReceiver(timeout time.Duration) io.ReadCloser

	// ReceiveSerialized reads serialized data with length header.
	//
	// The buffer provided should be at least 8 bytes bigger than the data expected to receive,
	// so for example if you want to receive a 1024 bytes message you should provide a buffer with at least 1032 bytes.
	//
	// Possible errors:
	//   - ErrSmallBuffer: the buffer provided was too small to receive the incoming data
	//   - ErrRead: failed to receive data from the stream
	//   - ErrTimeout: timeout occurred when receiving data from the stream
	ReceiveSerialized(b []byte, timeout time.Duration) (uint64, error)

	// ReceiveEncrypted reads AES-GCM encrypted data.
	//
	// The buffer provided should be at least 28 bytes bigger than the data expected to receive.
	// so for example if you want to receive a 1024 bytes message you should provide a buffer with at least 1048 bytes.
	//
	// Possible errors:
	//   - ErrAesKey: the aesKey is nil, meaning authentication is not enabled
	//   - ErrSmallBuffer: the buffer provided was too small to receive the incoming data
	//   - ErrRead: failed to receive data from the stream
	//   - ErrCipher: invalid key size
	//   - ErrGCM: failed to create GCM
	//   - ErrShort: ciphertext received is malformed because it is too short
	//   - ErrDecrypt: failed to decrypt the received data
	//   - ErrTimeout: timeout occurred when receiving data from the stream
	ReceiveEncrypted(b []byte, timeout time.Duration) (uint64, error)

	// NewStreamedEncryptedReceiver returns an io.ReadCloser that decrypts data as it comes from the stream.
	//
	// Possible errors:
	//   - ErrAesKey: the aesKey is nil, meaning authentication is not enabled
	//   - ErrRead: failed to receive the nonce from the stream
	//   - ErrStreamCipher: failed to create an AES-CTR stream
	NewStreamedEncryptedReceiver(timeout time.Duration) (io.ReadCloser, error)

	// GetDieCh returns a readonly chan which can be readable when the stream is to be closed.
	GetDieCh() <-chan struct{}
}

var headerPool = sync.Pool{
	New: func() any {
		buf := make([]byte, 2)
		return &buf
	},
}
