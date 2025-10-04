package crypto

import "crypto/rand"

// GenerateAESKey creates a random AES key of specified bit length.
func GenerateAESKey(bits int) []byte {
	key := make([]byte, bits/8)
	rand.Read(key) // Docs say that this never returns an error
	return key
}
