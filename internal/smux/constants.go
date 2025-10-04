package smux

import "sync"

var headerPool = sync.Pool{
	New: func() interface{} {
		return make([]byte, 2)
	},
}
