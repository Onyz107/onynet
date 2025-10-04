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
	intCrypto "github.com/Onyz107/onynet/internal/crypto"
	"github.com/Onyz107/onynet/internal/logger"
)

// AuthorizeSelfClient performs client-side authentication and returns an AES key.
func AuthorizeSelfClient(conn net.Conn, publicKey *rsa.PublicKey) (aesKey []byte, err error) {
	conn.SetDeadline(time.Now().Add(5 * time.Second))
	defer conn.SetDeadline(time.Time{})

	aesKey = intCrypto.GenerateAESKey(256)
	logger.Log.Debugf("AuthorizeSelfClient: aesKey: length: %d", len(aesKey))
	ciphertext, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, publicKey, aesKey, nil)
	if err != nil {
		return nil, errors.Join(intErrors.ErrPublickey, err)
	}

	header := headerPool.Get().([]byte)
	defer headerPool.Put(header)
	length := uint64(len(ciphertext))
	binary.BigEndian.PutUint64(header, length)

	logger.Log.Debugf("AuthorizeSelfClient: sending length: %d", length)
	n, err := conn.Write(header)
	if err != nil {
		return nil, errors.Join(intErrors.ErrWrite, err)
	}
	if n != len(header) {
		return nil, intErrors.ErrShortWrite
	}

	logger.Log.Debugf("AuthorizeSelfClient: sending ciphertext: length: %d", len(ciphertext))
	n, err = conn.Write(ciphertext)
	if err != nil {
		return nil, errors.Join(intErrors.ErrWrite, err)
	}
	if n != len(ciphertext) {
		return nil, intErrors.ErrShortWrite
	}

	return aesKey, nil
}

// AuthorizeServer performs server-side challenge verification for the client.
func AuthorizeServer(conn net.Conn, publicKey *rsa.PublicKey) error {
	conn.SetDeadline(time.Now().Add(5 * time.Second))
	defer conn.SetDeadline(time.Time{})

	challenge := make([]byte, serverChallengeLength)
	rand.Read(challenge) // never returns an error
	hash := sha256.Sum256(challenge)

	logger.Log.Debugf("AuthorizeServer: sending challenge: length: %d", len(challenge))
	n, err := conn.Write(challenge)
	if err != nil {
		return errors.Join(intErrors.ErrWrite, err)
	}
	if n != len(challenge) {
		return intErrors.ErrShortWrite
	}

	header := headerPool.Get().([]byte)
	defer headerPool.Put(header)
	if _, err := io.ReadFull(conn, header); err != nil {
		return errors.Join(intErrors.ErrRead, err)
	}
	length := binary.BigEndian.Uint64(header)
	logger.Log.Debugf("AuthorizeServer: received length: %d", length)

	signature := make([]byte, length)
	if _, err := io.ReadFull(conn, signature); err != nil {
		return errors.Join(intErrors.ErrRead, err)
	}
	logger.Log.Debugf("AuthorizeServer: received signature: length: %d", len(signature))

	if err := rsa.VerifyPKCS1v15(publicKey, crypto.SHA256, hash[:], signature); err != nil {
		return errors.Join(intErrors.ErrPublickey, err)
	}

	return nil
}
