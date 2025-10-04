package heartbeat

import (
	"sync"
)

const (
	senderMsg   = "ping"
	receiverMsg = "pong"
)

var bufPool = sync.Pool{
	New: func() any {
		return make([]byte, 4)
	},
}
