package auth

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/binary"
	"errors"
	"io"
	"net"
	"time"

	intErrors "github.com/Onyz107/onynet/errors"
	"github.com/Onyz107/onynet/internal/logger"
)

// AuthorizeSelfServer signs a challenge sent by a client to prove server identity.
func AuthorizeSelfServer(conn net.Conn, privateKey *rsa.PrivateKey) error {
	conn.SetDeadline(time.Now().Add(5 * time.Second))
	defer conn.SetDeadline(time.Time{})

	challenge := serverChallengePool.Get().([]byte)
	defer serverChallengePool.Put(challenge)

	if _, err := io.ReadFull(conn, challenge); err != nil {
		return errors.Join(intErrors.ErrRead, err)
	}

	hash := sha256.Sum256(challenge)
	signature, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, hash[:])
	if err != nil {
		return errors.Join(intErrors.ErrPrivateKey, err)
	}

	header := headerPool.Get().([]byte)
	defer headerPool.Put(header)
	length := uint64(len(signature))
	binary.BigEndian.PutUint64(header, length)

	logger.Log.Debugf("AuthorizeSelfServer: sending length: %d", length)
	n, err := conn.Write(header)
	if err != nil {
		return errors.Join(intErrors.ErrWrite, err)
	}
	if n != len(header) {
		return intErrors.ErrShortWrite
	}

	logger.Log.Debugf("AuthorizeSelfServer: sending signature: length: %d", len(signature))
	n, err = conn.Write(signature)
	if err != nil {
		return errors.Join(intErrors.ErrWrite, err)
	}
	if n != len(signature) {
		return intErrors.ErrShortWrite
	}

	return nil
}

// AuthorizeClient reads an AES key encrypted by a client and decrypts it.
func AuthorizeClient(conn net.Conn, privateKey *rsa.PrivateKey) (aesKey []byte, err error) {
	conn.SetDeadline(time.Now().Add(5 * time.Second))
	defer conn.SetDeadline(time.Time{})

	header := headerPool.Get().([]byte)
	defer headerPool.Put(header)

	logger.Log.Debug("AuthorizeClient: receiving header")
	if _, err := io.ReadFull(conn, header); err != nil {
		return nil, errors.Join(intErrors.ErrRead, err)
	}
	length := binary.BigEndian.Uint64(header)
	logger.Log.Debugf("AuthorizeClient: received header with length: %d", length)

	logger.Log.Debug("AuthorizeClient: receiving challenge")
	challenge := make([]byte, length)
	if _, err := io.ReadFull(conn, challenge); err != nil {
		return nil, errors.Join(intErrors.ErrRead, err)
	}
	logger.Log.Debugf("AuthorizeClient: received challenge: length: %d", len(challenge))

	logger.Log.Debug("AuthorizeClient: computing aesKey")
	aesKey, err = rsa.DecryptOAEP(sha256.New(), rand.Reader, privateKey, challenge, nil)
	if err != nil {
		return nil, errors.Join(intErrors.ErrPrivateKey, err)
	}
	logger.Log.Debugf("AuthorizeClient: computed aesKey: length: %d", len(aesKey))

	return aesKey, nil
}
