package kcp

import (
	"context"
	"errors"
	"net"
	"sync"

	intErrors "github.com/Onyz107/onynet/errors"
	"github.com/Onyz107/onynet/internal/logger"
	"github.com/xtaci/kcp-go/v5"
)

type Server struct {
	listener *kcp.Listener
	ctx      context.Context
	done     chan struct{}
	once     sync.Once
}

// NewServer creates a KCP listener for accepting client connections.
func NewServer(addr net.Addr, ctx context.Context) (*Server, error) {
	logger.Log.Debugf("kcp/server NewServer: listening on %s", addr.String())
	listener, err := kcp.ListenWithOptions(addr.String(), nil, 0, 0)
	if err != nil {
		return nil, errors.Join(intErrors.ErrBadAddr, err)
	}

	server := &Server{listener: listener, ctx: ctx, done: make(chan struct{}, 1)}
	logger.Log.Debugf("kcp/server NewServer: created new server on %s", addr.String())

	go func() {
		select {
		case <-server.ctx.Done():
			logger.Log.Debugf("kcp/server NewServer: closing server because of context canceled")
			server.Close()
			return
		case <-server.done:
			return
		}
	}()

	return server, nil
}

// Accept waits for a new KCP client connection.
func (s *Server) Accept() (*ClientConn, error) {
	logger.Log.Debugf("kcp.Server Accept: accepting client connection")
	conn, err := s.listener.AcceptKCP()
	if err != nil {
		return nil, errors.Join(intErrors.ErrAccept, err)
	}
	logger.Log.Debug("kcp.Server Accept: accepted client connection")

	// Performance optimizations
	conn.SetWindowSize(512, 512)
	conn.SetNoDelay(1, 40, 2, 1)

	client := &ClientConn{conn: conn, ctx: s.ctx, done: make(chan struct{}, 1)}

	go func() {
		select {
		case <-client.ctx.Done():
			logger.Log.Debug("kcp.Server Accept: closing client connection because of context canceled")
			conn.Close()
			return
		case <-client.done:
			return
		}
	}()

	return client, nil
}

func (s *Server) Close() error {
	s.once.Do(func() { close(s.done) })
	logger.Log.Debugf("kcp.Server Close: closing server connection")
	return s.listener.Close()
}
