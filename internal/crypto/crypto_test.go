package crypto_test

import (
	"crypto/rand"
	"fmt"
	"testing"

	"github.com/Onyz107/onynet/internal/crypto"
)

// generateTestData creates test data of specified size
func generateTestData(size int) []byte {
	data := make([]byte, size)
	rand.Read(data)
	return data
}

// BenchmarkEncryptAESGCM tests AES-GCM encryption performance
func BenchmarkEncryptAESGCM(b *testing.B) {
	sizes := []int{64, 512, 1024, 4096, 16384, 65536, 262144} // Various payload sizes

	for _, size := range sizes {
		b.Run(fmt.Sprintf("size_%d", size), func(b *testing.B) {
			plaintext := generateTestData(size)
			key := crypto.GenerateAESKey(256)

			b.ResetTimer()
			b.ReportAllocs()

			for i := 0; i < b.N; i++ {
				_, err := crypto.EncryptAESGCM(plaintext, key)
				if err != nil {
					b.Fatal(err)
				}
			}

			// Report throughput
			b.SetBytes(int64(size))
		})
	}
}

// BenchmarkDecryptAESGCM tests AES-GCM decryption performance
func BenchmarkDecryptAESGCM(b *testing.B) {
	sizes := []int{64, 512, 1024, 4096, 16384, 65536, 262144}

	for _, size := range sizes {
		b.Run(fmt.Sprintf("size_%d", size), func(b *testing.B) {
			plaintext := generateTestData(size)
			key := crypto.GenerateAESKey(256)

			// Pre-encrypt the data
			ciphertext, err := crypto.EncryptAESGCM(plaintext, key)
			if err != nil {
				b.Fatal(err)
			}

			b.ResetTimer()
			b.ReportAllocs()

			for i := 0; i < b.N; i++ {
				_, err := crypto.DecryptAESGCM(ciphertext, key)
				if err != nil {
					b.Fatal(err)
				}
			}

			b.SetBytes(int64(size))
		})
	}
}
