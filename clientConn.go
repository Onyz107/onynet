package onynet

import (
	"context"
	"errors"
	"net"
	"time"

	"github.com/Onyz107/onynet/internal/kcp"
	intSmux "github.com/Onyz107/onynet/internal/smux"
)

type ClientConn struct {
	client  *kcp.ClientConn
	manager *intSmux.Manager
	ctx     context.Context
}

// OpenStream opens a named stream to communicate with the server.
func (cn *ClientConn) OpenStream(name string, timeout time.Duration) (*intSmux.Stream, error) {
	return cn.manager.Open(name, timeout)
}

// AcceptStream accepts an incoming named stream from the client.
func (cn *ClientConn) AcceptStream(name string, timeout time.Duration) (*intSmux.Stream, error) {
	return cn.manager.Accept(name, timeout)
}

// LocalAddr returns the client's local address.
func (cn *ClientConn) LocalAddr() net.Addr {
	return cn.client.LocalAddr()
}

// RemoteAddr returns the client's remote address.
func (cn *ClientConn) RemoteAddr() net.Addr {
	return cn.client.RemoteAddr()
}

// Close closes the client connection and streams.
func (cn *ClientConn) Close() error {
	var errs []error

	if err := cn.client.Close(); err != nil {
		errs = append(errs, err)
	}

	if err := cn.manager.Close(); err != nil {
		errs = append(errs, err)
	}

	return errors.Join(errs...)
}
