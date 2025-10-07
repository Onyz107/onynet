package transfer

import "sync"

var headerPool = sync.Pool{
	New: func() any {
		buf := make([]byte, 8)
		return &buf
	},
}
