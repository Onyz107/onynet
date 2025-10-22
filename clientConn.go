package onynet

import (
	"context"
	"errors"
	"net"
	"time"

	"github.com/Onyz107/onynet/internal/kcp"
	intSmux "github.com/Onyz107/onynet/internal/smux"
)

// ClientConn defines a client that is connected to a server.
// The difference between Client and ClientConn is that Client should
// only be used when performing operations on the Client, while ClientConn
// should only be used when performing operations on the Server.
type ClientConn struct {
	client    *kcp.ClientConn
	connected bool
	manager   *intSmux.Manager
	ctx       context.Context
}

// OpenStream opens a named stream to communicate with the client.
// The ctx argument defines the stream's deadline while timeout defines the handshake's deadline.
//
// Possible errors:
//   - ErrNameTooLong: name for stream is too long
//   - ErrCtxCancelled: context was cancelled while waiting for a stream to establish connection
//   - ErrTimeout: timeout occurred waiting for the stream to establish connection
//   - ErrOpenStream: failed to open a multiplexing stream
//   - ErrWrite: failed to send headers through the stream
//   - ErrShortWrite: headers sent were shorter than expected
//   - ErrRead: failed to receive headers from the stream
func (cn *ClientConn) OpenStream(name string, ctx context.Context, timeout time.Duration) (*intSmux.Stream, error) {
	return cn.manager.OpenStream(name, ctx, timeout)
}

// AcceptStream accepts an incoming named stream from the client.
// The ctx argument defines the stream's deadline while timeout defines the handshake's deadline.
//
// Possible errors:
//   - ErrNameTooLong: name for stream is too long
//   - ErrCtxCancelled: context was cancelled while waiting for a stream to establish connection
//   - ErrTimeout: timeout occurred waiting for the stream to establish connection
//   - ErrAcceptStream: failed to accept a multiplexing stream
//   - ErrRead: failed to receive headers from the stream
//   - ErrWrite: failed to send headers through the stream
func (cn *ClientConn) AcceptStream(name string, ctx context.Context, timeout time.Duration) (*intSmux.Stream, error) {
	return cn.manager.AcceptStream(name, ctx, timeout)
}

// IsConnected returns true if the client is currently connected to the server, and false otherwise.
// The connection status is tracked by a variable that is set to true when a connection is established,
// and set to false when the Close function is called (for example, after a heartbeat failure or a manual disconnect).
func (cn *ClientConn) IsConnected() bool {
	return cn.connected
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
	cn.connected = false
	var errs []error

	if err := cn.client.Close(); err != nil {
		errs = append(errs, err)
	}

	if err := cn.manager.Close(); err != nil {
		errs = append(errs, err)
	}

	return errors.Join(errs...)
}
