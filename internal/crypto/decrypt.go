package crypto

import (
	"errors"

	intErrors "github.com/Onyz107/onynet/errors"
)

// DecryptAESGCM decrypts AES-GCM encrypted data with the given key.
func DecryptAESGCM(data, key []byte) ([]byte, error) {
	gcm, err := getGCM(key)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return nil, intErrors.ErrShort
	}

	nonce := data[:nonceSize]
	ciphertext := data[nonceSize:]

	bufPtr := decryptionBufPool.Get().(*[]byte)
	defer decryptionBufPool.Put(bufPtr)
	buf := (*bufPtr)[:0]

	plaintext, err := gcm.Open(buf, nonce, ciphertext, nil)
	if err != nil {
		return nil, errors.Join(intErrors.ErrDecrypt, err)
	}

	return plaintext, nil
}
