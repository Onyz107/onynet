package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"errors"
	"sync"

	intErrors "github.com/Onyz107/onynet/errors"
)

var decryptionBufPool = sync.Pool{
	New: func() any {
		buf := make([]byte, 0, 4096)
		return &buf
	},
}

var gcmCache sync.Map // key string -> cipher.AEAD

func getGCM(key []byte) (cipher.AEAD, error) {
	if val, ok := gcmCache.Load(string(key)); ok {
		return val.(cipher.AEAD), nil
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, errors.Join(intErrors.ErrCipher, err)
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, errors.Join(intErrors.ErrGCM, err)
	}
	gcmCache.Store(string(key), gcm)
	return gcm, nil
}
