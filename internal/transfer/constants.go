package transfer

import "sync"

var headerPool = sync.Pool{
	New: func() any {
		return make([]byte, 8)
	},
}
