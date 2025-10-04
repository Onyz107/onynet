package transfer

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"net"
	"time"

	intErrors "github.com/Onyz107/onynet/errors"
	"github.com/Onyz107/onynet/internal/crypto"
)

// Send writes bytes to conn with optional timeout.
func Send(conn net.Conn, data []byte, timeout time.Duration) error {
	if timeout > 0 {
		conn.SetWriteDeadline(time.Now().Add(timeout))
	}
	defer conn.SetWriteDeadline(time.Time{})

	n, err := conn.Write(data)
	if err != nil {
		return errors.Join(intErrors.ErrWrite, err)
	}
	if n != len(data) {
		return errors.Join(intErrors.ErrShortWrite, fmt.Errorf("sent %d bytes instead of %d", n, len(data)))
	}

	return nil
}

// NewStreamedSender returns an io.WriteCloser that allows
// the caller to directly write data to the stream and set a timeout.
func NewStreamedSender(conn net.Conn, timeout time.Duration) io.WriteCloser {
	return &timedWriter{w: conn, timeout: timeout}
}

// SendSerialized sends length-prefixed data.
func SendSerialized(conn net.Conn, data []byte, timeout time.Duration) error {
	header := headerPool.Get().([]byte)
	defer headerPool.Put(header)
	length := uint64(len(data))
	binary.BigEndian.PutUint64(header, length)

	if err := Send(conn, header, timeout); err != nil {
		return err
	}

	if err := Send(conn, data, timeout); err != nil {
		return err
	}

	return nil
}

// SendEncrypted encrypts data and sends it.
func SendEncrypted(conn net.Conn, data, aesKey []byte, timeout time.Duration) error {
	if aesKey == nil {
		return intErrors.ErrAESKey
	}

	data, err := crypto.EncryptAESGCM(data, aesKey)
	if err != nil {
		return err
	}

	return SendSerialized(conn, data, timeout)
}

// NewStreamedEncryptedSender returns an io.WriteCloser that encrypts data as it is written to the stream.
func NewStreamedEncryptedSender(conn net.Conn, aesKey []byte, timeout time.Duration) (io.WriteCloser, error) {
	nonce := make([]byte, aes.BlockSize)
	rand.Read(nonce)

	streamWriter := NewStreamedSender(conn, timeout)

	n, err := streamWriter.Write(nonce)
	if err != nil {
		return nil, errors.Join(intErrors.ErrWrite, err)
	}
	if n != len(nonce) {
		return nil, errors.Join(intErrors.ErrShortWrite, fmt.Errorf("sent %d bytes instead of %d", n, len(nonce)))
	}

	stream, err := crypto.NewStreamedCipher(aesKey, nonce)
	if err != nil {
		return nil, errors.Join(intErrors.ErrStreamCipher, err)
	}

	encryptedStreamWriter := &cipher.StreamWriter{
		S: stream,
		W: streamWriter,
	}

	return encryptedStreamWriter, nil
}
