package onynet

import (
	"github.com/Onyz107/onynet/internal/smux"
)

type Handler interface {
	smux.Handler
}

type Communicator interface {
	smux.Communicator
}
