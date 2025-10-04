package auth

import "sync"

const serverChallengeLength = 32

var serverChallengePool = sync.Pool{
	New: func() any {
		return make([]byte, serverChallengeLength)
	},
}

var headerPool = sync.Pool{
	New: func() any {
		return make([]byte, 8)
	},
}
