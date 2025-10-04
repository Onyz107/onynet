package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"errors"

	intErrors "github.com/Onyz107/onynet/errors"
)

// DecryptAESGCM decrypts AES-GCM encrypted data with the given key.
func DecryptAESGCM(data, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, errors.Join(intErrors.ErrCipher, err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, errors.Join(intErrors.ErrGCM, err)
	}

	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return nil, intErrors.ErrShort
	}

	nonce := data[:nonceSize]
	ciphertext := data[nonceSize:]

	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, errors.Join(intErrors.ErrDecrypt, err)
	}

	return plaintext, nil
}
