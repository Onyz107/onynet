package onynet

import (
	"context"
	"crypto/rsa"
	"errors"
	"net"
	"time"

	intErrors "github.com/Onyz107/onynet/errors"
	"github.com/Onyz107/onynet/internal/auth"
	"github.com/Onyz107/onynet/internal/heartbeat"
	"github.com/Onyz107/onynet/internal/kcp"
	"github.com/Onyz107/onynet/internal/logger"
	intSmux "github.com/Onyz107/onynet/internal/smux"
	"github.com/xtaci/smux"
)

// Client defines a client that is ready to connect to a server.
// The difference between Client and ClientConn is that Client should
// only be used when performing operations on the Client, while ClientConn
// should only be used when performing operations on the Server.
type Client struct {
	client  *kcp.Client
	manager *intSmux.Manager
	ctx     context.Context
}

// Dial connects to an OnyNet server, performs optional authentication (with publicKey being nil), and returns a client.
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

	onynetClient := &Client{client: client, manager: manager, ctx: ctx}

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
	var errs []error

	if err := c.client.Close(); err != nil {
		errs = append(errs, err)
	}

	if err := c.manager.Close(); err != nil {
		errs = append(errs, err)
	}

	return errors.Join(errs...)
}
