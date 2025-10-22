package transfer

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"net"
	"time"

	intErrors "github.com/Onyz107/onynet/errors"
	"github.com/Onyz107/onynet/internal/crypto"
)

// Receive reads exactly len(buf) bytes from conn.
func Receive(conn net.Conn, buf []byte, timeout time.Duration) error {
	timedConn := getTimedReadWriteCloser(conn, timeout)

	if _, err := io.ReadFull(timedConn, buf); err != nil {
		var ne net.Error
		if errors.As(err, &ne) && ne.Timeout() {
			return errors.Join(intErrors.ErrTimeout, err)
		}
		return errors.Join(intErrors.ErrRead, err)
	}

	return nil
}

// NewStreamedReceiver returns an io.ReadCloser that allows
// the caller to directly read data from the stream and set a timeout.
func NewStreamedReceiver(conn net.Conn, timeout time.Duration) io.ReadCloser {
	return getTimedReadWriteCloser(conn, timeout)
}

// ReceiveSerialized reads length-prefixed serialized data.
func ReceiveSerialized(conn net.Conn, buf []byte, timeout time.Duration) (uint64, error) {
	headerPtr := headerPool.Get().(*[]byte)
	defer headerPool.Put(headerPtr)
	header := *headerPtr

	if err := Receive(conn, header, timeout); err != nil {
		return 0, err
	}
	length := binary.BigEndian.Uint64(header)

	maxDataSize := uint64(cap(buf))
	if length > maxDataSize {
		return 0, errors.Join(intErrors.ErrSmallBuffer, fmt.Errorf("buffer cap: %d: data length: %d", maxDataSize, length))
	}

	data := buf[:length]
	if err := Receive(conn, data, timeout); err != nil {
		return 0, err
	}

	return length, nil
}

// ReceiveEncrypted reads and decrypts AES-GCM data.
func ReceiveEncrypted(conn net.Conn, buf, aesKey []byte, timeout time.Duration) (uint64, error) {
	if aesKey == nil {
		return 0, intErrors.ErrAESKey
	}

	n, err := ReceiveSerialized(conn, buf, timeout)
	if err != nil {
		return 0, err
	}
	data := buf[:n]

	plaintext, err := crypto.DecryptAESGCM(data, aesKey)
	if err != nil {
		return 0, err
	}

	n = uint64(copy(buf[:len(plaintext)], plaintext))

	return n, nil
}

// NewStreamedEncryptedReceiver returns an io.ReadCloser that decrypts data as it comes from the stream.
func NewStreamedEncryptedReceiver(conn net.Conn, aesKey []byte, timeout time.Duration) (io.ReadCloser, error) {
	if aesKey == nil {
		return nil, intErrors.ErrAESKey
	}

	reader := NewStreamedReceiver(conn, timeout)

	nonce := make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(reader, nonce); err != nil {
		reader.Close()
		return nil, errors.Join(intErrors.ErrRead, err)
	}

	stream, err := crypto.NewStreamedCipher(aesKey, nonce)
	if err != nil {
		reader.Close()
		return nil, errors.Join(intErrors.ErrStreamCipher, err)
	}

	streamReader := &cipher.StreamReader{
		S: stream,
		R: reader,
	}

	return struct {
		io.Reader
		io.Closer
	}{streamReader, reader}, nil
}
