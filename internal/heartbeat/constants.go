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
		buf := make([]byte, 4)
		return &buf
	},
}
