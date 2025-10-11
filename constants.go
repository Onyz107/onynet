package onynet

import "github.com/Onyz107/onynet/internal/smux"

type Communicator interface {
	smux.Communicator
}

type Stream struct {
	*smux.Stream
}
