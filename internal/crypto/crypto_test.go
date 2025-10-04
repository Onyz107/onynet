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

// generateNonce creates a random nonce for CTR mode
func generateNonce() []byte {
	nonce := make([]byte, 16) // AES block size
	rand.Read(nonce)
	return nonce
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

// BenchmarkEncryptDecryptCycle tests full encrypt-decrypt cycle
func BenchmarkEncryptDecryptCycle(b *testing.B) {
	sizes := []int{1024, 4096, 16384, 65536}

	for _, size := range sizes {
		b.Run(fmt.Sprintf("size_%d", size), func(b *testing.B) {
			plaintext := generateTestData(size)
			key := crypto.GenerateAESKey(256)

			b.ResetTimer()
			b.ReportAllocs()

			for i := 0; i < b.N; i++ {
				ciphertext, err := crypto.EncryptAESGCM(plaintext, key)
				if err != nil {
					b.Fatal(err)
				}

				decrypted, err := crypto.DecryptAESGCM(ciphertext, key)
				if err != nil {
					b.Fatal(err)
				}

				// Verify correctness in benchmark to ensure we're not optimizing incorrectly
				if len(decrypted) != len(plaintext) {
					b.Fatal("decrypted length mismatch")
				}
			}

			b.SetBytes(int64(size * 2)) // Account for both encrypt and decrypt
		})
	}
}

// BenchmarkNewStreamedCipher tests CTR stream cipher creation
func BenchmarkNewStreamedCipher(b *testing.B) {
	key := crypto.GenerateAESKey(256)
	nonce := generateNonce()

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		stream, err := crypto.NewStreamedCipher(key, nonce)
		if err != nil {
			b.Fatal(err)
		}
		_ = stream
	}
}

// BenchmarkStreamedCipherXOR tests CTR mode XOR operations
func BenchmarkStreamedCipherXOR(b *testing.B) {
	sizes := []int{64, 512, 1024, 4096, 16384, 65536, 262144}

	for _, size := range sizes {
		b.Run(fmt.Sprintf("size_%d", size), func(b *testing.B) {
			key := crypto.GenerateAESKey(256)
			nonce := generateNonce()
			data := generateTestData(size)

			stream, err := crypto.NewStreamedCipher(key, nonce)
			if err != nil {
				b.Fatal(err)
			}

			b.ResetTimer()
			b.ReportAllocs()

			for i := 0; i < b.N; i++ {
				// Create a copy for each iteration to avoid modifying original
				dataCopy := make([]byte, len(data))
				copy(dataCopy, data)

				stream.XORKeyStream(dataCopy, dataCopy)
			}

			b.SetBytes(int64(size))
		})
	}
}

// BenchmarkKeyGeneration tests key generation performance
func BenchmarkKeyGeneration(b *testing.B) {
	b.Run("AES256_key", func(b *testing.B) {
		b.ReportAllocs()

		for i := 0; i < b.N; i++ {
			key := make([]byte, 32)
			rand.Read(key)
		}
	})

	b.Run("nonce", func(b *testing.B) {
		b.ReportAllocs()

		for i := 0; i < b.N; i++ {
			nonce := make([]byte, 16)
			rand.Read(nonce)
		}
	})
}

// BenchmarkMemoryReuse tests performance with key/cipher reuse
func BenchmarkMemoryReuse(b *testing.B) {
	sizes := []int{1024, 4096, 16384}

	for _, size := range sizes {
		b.Run(fmt.Sprintf("new_key_each_time_%d", size), func(b *testing.B) {
			plaintext := generateTestData(size)

			b.ResetTimer()
			b.ReportAllocs()

			for i := 0; i < b.N; i++ {
				key := crypto.GenerateAESKey(256)
				_, err := crypto.EncryptAESGCM(plaintext, key)
				if err != nil {
					b.Fatal(err)
				}
			}

			b.SetBytes(int64(size))
		})

		b.Run(fmt.Sprintf("reuse_key_%d", size), func(b *testing.B) {
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

			b.SetBytes(int64(size))
		})
	}
}

// BenchmarkParallelEncryption tests concurrent encryption performance
func BenchmarkParallelEncryption(b *testing.B) {
	plaintext := generateTestData(4096)
	key := crypto.GenerateAESKey(256)

	b.ResetTimer()
	b.ReportAllocs()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, err := crypto.EncryptAESGCM(plaintext, key)
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	b.SetBytes(4096)
}

// BenchmarkParallelDecryption tests concurrent decryption performance
func BenchmarkParallelDecryption(b *testing.B) {
	plaintext := generateTestData(4096)
	key := crypto.GenerateAESKey(256)

	ciphertext, err := crypto.EncryptAESGCM(plaintext, key)
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	b.ReportAllocs()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, err := crypto.DecryptAESGCM(ciphertext, key)
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	b.SetBytes(4096)
}

// BenchmarkCompareStreamVsGCM compares CTR streaming vs GCM for large data
func BenchmarkCompareStreamVsGCM(b *testing.B) {
	size := 65536
	plaintext := generateTestData(size)
	key := crypto.GenerateAESKey(256)
	nonce := generateNonce()

	b.Run("GCM", func(b *testing.B) {
		b.ReportAllocs()

		for i := 0; i < b.N; i++ {
			_, err := crypto.EncryptAESGCM(plaintext, key)
			if err != nil {
				b.Fatal(err)
			}
		}

		b.SetBytes(int64(size))
	})

	b.Run("CTR_Stream", func(b *testing.B) {
		b.ReportAllocs()

		for i := 0; i < b.N; i++ {
			stream, err := crypto.NewStreamedCipher(key, nonce)
			if err != nil {
				b.Fatal(err)
			}

			data := make([]byte, len(plaintext))
			copy(data, plaintext)
			stream.XORKeyStream(data, data)
		}

		b.SetBytes(int64(size))
	})
}
