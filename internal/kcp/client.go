package kcp

import (
	"context"
	"errors"
	"net"
	"sync"
	"time"

	intErrors "github.com/Onyz107/onynet/errors"
	"github.com/Onyz107/onynet/internal/logger"
	"github.com/xtaci/kcp-go/v5"
)

type Client struct {
	conn *kcp.UDPSession
	ctx  context.Context
	done chan struct{}
	once sync.Once
}

func (c *Client) Read(b []byte) (n int, err error) {
	select {
	case <-c.ctx.Done():
		return 0, intErrors.ErrCtxCancelled
	default:
		return c.conn.Read(b)
	}
}

func (c *Client) Write(b []byte) (n int, err error) {
	select {
	case <-c.ctx.Done():
		return 0, intErrors.ErrCtxCancelled
	default:
		return c.conn.Write(b)
	}
}

func (c *Client) Close() error {
	c.once.Do(func() { close(c.done) })
	logger.Log.Debugf("kcp.Client Close: closing client connection")
	return c.conn.Close()
}

func (c *Client) LocalAddr() net.Addr {
	return c.conn.LocalAddr()
}

func (c *Client) RemoteAddr() net.Addr {
	return c.conn.RemoteAddr()
}

func (c *Client) SetDeadline(t time.Time) error {
	return c.conn.SetDeadline(t)
}

func (c *Client) SetReadDeadline(t time.Time) error {
	return c.conn.SetReadDeadline(t)
}

func (c *Client) SetWriteDeadline(t time.Time) error {
	return c.conn.SetWriteDeadline(t)
}

// Dial connects to a KCP server and returns a client wrapper.
func Dial(addr net.Addr, ctx context.Context) (*Client, error) {
	logger.Log.Debugf("kcp/client Dial: connecting to %s", addr.String())
	conn, err := kcp.DialWithOptions(addr.String(), nil, 0, 0)
	if err != nil {
		return nil, errors.Join(intErrors.ErrBadAddr, err)
	}
	logger.Log.Debugf("kcp/client Dial: connected to %s", addr.String())

	// Performance optimizations
	conn.SetWindowSize(512, 512)
	conn.SetNoDelay(1, 40, 2, 1)

	client := &Client{conn: conn, ctx: ctx, done: make(chan struct{}, 1)}

	go func() {
		select {
		case <-client.ctx.Done():
			logger.Log.Debug("kcp/client Dial: closing client because of context canceled")
			client.Close()
			return
		case <-client.done:
			return
		}
	}()

	return client, nil
}
