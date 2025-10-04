package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"

	intErrors "github.com/Onyz107/onynet/errors"
)

// EncryptAESGCM encrypts plaintext using AES-GCM with the given key.
func EncryptAESGCM(plaintext, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, errors.Join(intErrors.ErrCipher, err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, errors.Join(intErrors.ErrGCM, err)
	}

	nonce := make([]byte, gcm.NonceSize())
	rand.Read(nonce)

	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)
	return ciphertext, nil
}
