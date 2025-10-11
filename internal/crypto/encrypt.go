package crypto

import (
	"crypto/rand"
)

// EncryptAESGCM encrypts plaintext using AES-GCM with the given key.
func EncryptAESGCM(plaintext, key []byte) ([]byte, error) {
	gcm, err := getGCM(key)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	rand.Read(nonce)

	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)
	return ciphertext, nil
}
