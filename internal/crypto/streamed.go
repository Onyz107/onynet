package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"errors"

	intErrors "github.com/Onyz107/onynet/errors"
)

// NewStreamedCipher creates a CTR stream cipher from AES key and nonce.
func NewStreamedCipher(aesKey, nonce []byte) (cipher.Stream, error) {
	block, err := aes.NewCipher(aesKey)
	if err != nil {
		return nil, errors.Join(intErrors.ErrCipher, err)
	}

	return cipher.NewCTR(block, nonce), nil
}
