package kcp

import (
	"context"
	"net"
	"sync"
	"time"

	intErrors "github.com/Onyz107/onynet/errors"
	"github.com/Onyz107/onynet/internal/logger"
	"github.com/xtaci/kcp-go/v5"
)

type ClientConn struct {
	conn *kcp.UDPSession
	ctx  context.Context
	done chan struct{}
	once sync.Once
}

func (c *ClientConn) Read(p []byte) (n int, err error) {
	select {
	case <-c.ctx.Done():
		return 0, intErrors.ErrCtxCancelled
	default:
		return c.conn.Read(p)
	}
}

func (c *ClientConn) Write(p []byte) (n int, err error) {
	select {
	case <-c.ctx.Done():
		return 0, intErrors.ErrCtxCancelled
	default:
		return c.conn.Write(p)
	}
}

func (c *ClientConn) Close() error {
	c.once.Do(func() { close(c.done) })
	logger.Log.Debugf("kcp.ClientConn Close: closing client connection")
	return c.conn.Close()
}

func (c *ClientConn) LocalAddr() net.Addr {
	return c.conn.LocalAddr()
}

func (c *ClientConn) RemoteAddr() net.Addr {
	return c.conn.RemoteAddr()
}

func (c *ClientConn) SetDeadline(t time.Time) error {
	return c.conn.SetDeadline(t)
}

func (c *ClientConn) SetReadDeadline(t time.Time) error {
	return c.conn.SetReadDeadline(t)
}

func (c *ClientConn) SetWriteDeadline(t time.Time) error {
	return c.conn.SetWriteDeadline(t)
}
