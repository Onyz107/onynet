package onynet

import (
	"context"
	"crypto/rsa"
	"errors"
	"net"
	"sync"
	"time"

	intErrors "github.com/Onyz107/onynet/errors"
	"github.com/Onyz107/onynet/internal/auth"
	"github.com/Onyz107/onynet/internal/heartbeat"
	"github.com/Onyz107/onynet/internal/kcp"
	"github.com/Onyz107/onynet/internal/logger"
	intSmux "github.com/Onyz107/onynet/internal/smux"
	"github.com/xtaci/smux"
)

type Client struct {
	client  *kcp.Client
	manager *intSmux.Manager
	ctx     context.Context
	done    chan struct{}
	once    sync.Once
}

// Dial connects to an OnyNet server, performs optional authentication, and returns a client.
// The publicKey argument can be nil but this will fail anything that calls *.Send/ReceiveEncrypted
func Dial(addr net.Addr, publicKey *rsa.PublicKey, ctx context.Context) (*Client, error) {
	client, err := kcp.Dial(addr, ctx)
	if err != nil {
		return nil, errors.Join(intErrors.ErrDial, err)
	}

	var aesKey []byte
	if publicKey != nil {
		aesKey, err = auth.AuthorizeSelfClient(client, publicKey)
		if err != nil {
			client.Close()
			return nil, errors.Join(intErrors.ErrAuth, err)
		}

		if err := auth.AuthorizeServer(client, publicKey); err != nil {
			client.Close()
			return nil, errors.Join(intErrors.ErrAuth, err)
		}
	}

	session, err := smux.Client(client, nil)
	if err != nil {
		client.Close()
		return nil, errors.Join(intErrors.ErrCreateSession, err)
	}

	manager := intSmux.NewManager(session, aesKey, ctx)

	onynetClient := &Client{client: client, manager: manager, ctx: ctx, done: make(chan struct{}, 1)}

	go func() {
		select {
		case <-onynetClient.ctx.Done():
			onynetClient.Close()
			return
		case <-onynetClient.done:
			return
		}
	}()

	heartbeatStream, err := onynetClient.OpenStream("heartbeatStream", 5*time.Second)
	if err != nil {
		onynetClient.Close()
		return nil, errors.Join(intErrors.ErrHeartbeatStream, err)
	}
	go func() {
		defer heartbeatStream.Close()
		if err := heartbeat.SendHeartbeat(heartbeatStream, ctx); err != nil {
			logger.Log.Debugf("closing client because of heartbeat err: %v", err)
			onynetClient.Close()
		}
	}()

	return onynetClient, nil
}

// OpenStream opens a named stream to communicate with the server.
func (c *Client) OpenStream(name string, timeout time.Duration) (*intSmux.Stream, error) {
	return c.manager.Open(name, timeout)
}

// AcceptStream accepts an incoming named stream.
func (c *Client) AcceptStream(name string, timeout time.Duration) (*intSmux.Stream, error) {
	return c.manager.Accept(name, timeout)
}

// Close gracefully closes client connections and streams.
func (c *Client) Close() error {
	c.once.Do(func() { close(c.done) })

	var errs []error

	if err := c.client.Close(); err != nil {
		errs = append(errs, err)
	}

	if err := c.manager.Close(); err != nil {
		errs = append(errs, err)
	}

	return errors.Join(errs...)
}
