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

type Server struct {
	server     *kcp.Server
	privateKey *rsa.PrivateKey
	ctx        context.Context
}

// NewServer starts an OnyNet server listening on given address.
func NewServer(addr net.Addr, privateKey *rsa.PrivateKey, ctx context.Context) (*Server, error) {
	server, err := kcp.NewServer(addr, ctx)
	if err != nil {
		return nil, errors.Join(intErrors.ErrNewServer, err)
	}

	onynetServer := &Server{server: server, privateKey: privateKey, ctx: ctx}

	return onynetServer, nil
}

// Accept waits for a new client connection and performs authentication.
func (s *Server) Accept() (*ClientConn, error) {
	client, err := s.server.Accept()
	if err != nil {
		return nil, errors.Join(intErrors.ErrAcceptClient, err)
	}

	var aesKey []byte
	if s.privateKey != nil {
		aesKey, err = auth.AuthorizeClient(client, s.privateKey)
		if err != nil {
			client.Close()
			return nil, errors.Join(intErrors.ErrAuth, err)
		}

		if err := auth.AuthorizeSelfServer(client, s.privateKey); err != nil {
			client.Close()
			return nil, errors.Join(intErrors.ErrAuth, err)
		}
	}

	session, err := smux.Server(client, nil)
	if err != nil {
		client.Close()
		return nil, errors.Join(intErrors.ErrCreateSession, err)
	}
	manager := intSmux.NewManager(session, aesKey, s.ctx)

	onynetClientConn := &ClientConn{client: client, manager: manager, ctx: s.ctx}

	heartbeatStream, err := onynetClientConn.AcceptStream("heartbeatStream", 5*time.Second)
	if err != nil {
		onynetClientConn.Close()
		return nil, errors.Join(intErrors.ErrHeartbeatStream, err)
	}
	go func() {
		defer heartbeatStream.Close()
		if err := heartbeat.ReceiveHeartbeat(heartbeatStream, s.ctx); err != nil {
			logger.Log.Debugf("closing client because of heartbeat err: %v", err)
			onynetClientConn.Close()
		}
	}()

	return onynetClientConn, nil
}

// Close shuts down the server and all active connections.
func (s *Server) Close() error {
	return s.server.Close()
}
