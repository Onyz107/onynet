package auth

import (
	"sync"
)

const serverChallengeLength = 32

var serverChallengePool = sync.Pool{
	New: func() any {
		buf := make([]byte, serverChallengeLength)
		return &buf
	},
}

var headerPool = sync.Pool{
	New: func() any {
		buf := make([]byte, 8)
		return &buf
	},
}
