package logger

import (
	"os"
	"sync"

	"github.com/Onyz107/onylogger"
)

var (
	Log  *onylogger.OnyLogger
	once sync.Once
)

func init() {
	once.Do(func() {
		Log = onylogger.New(os.Stderr)
	})
}
