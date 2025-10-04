package onynet

import (
	"context"
	"crypto/rsa"
	"errors"
	"net"
	"sync"
	"sync/atomic"
	"time"

	intErrors "github.com/Onyz107/onynet/errors"
	"github.com/Onyz107/onynet/internal/auth"
	"github.com/Onyz107/onynet/internal/heartbeat"
	"github.com/Onyz107/onynet/internal/kcp"
	"github.com/Onyz107/onynet/internal/logger"
	intSmux "github.com/Onyz107/onynet/internal/smux"
	"github.com/xtaci/smux"
)

// Server defines a server which will be listening for incoming connections.
type Server struct {
	server     *kcp.Server
	clients    map[int]*ClientConn
	mu         sync.RWMutex
	privateKey *rsa.PrivateKey
	ctx        context.Context
}

var clientCounter int64

// NewServer starts an OnyNet server listening on given address.
func NewServer(addr net.Addr, privateKey *rsa.PrivateKey, ctx context.Context) (*Server, error) {
	server, err := kcp.NewServer(addr, ctx)
	if err != nil {
		return nil, errors.Join(intErrors.ErrNewServer, err)
	}

	onynetServer := &Server{server: server, clients: make(map[int]*ClientConn), privateKey: privateKey, ctx: ctx}

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

	id := int(atomic.AddInt64(&clientCounter, 1))
	s.mu.Lock()
	s.clients[id] = onynetClientConn
	s.mu.Unlock()

	heartbeatStream, err := onynetClientConn.AcceptStream("heartbeatStream", 5*time.Second)
	if err != nil {
		delete(s.clients, id)
		onynetClientConn.Close()
		return nil, errors.Join(intErrors.ErrHeartbeatStream, err)
	}

	go func() {
		defer heartbeatStream.Close()
		if err := heartbeat.ReceiveHeartbeat(heartbeatStream, s.ctx); err != nil {
			logger.Log.Debugf("closing client because of heartbeat err: %v", err)
			delete(s.clients, id)
			onynetClientConn.Close()
		}
	}()

	return onynetClientConn, nil
}

// GetClients returns a map of all connected clients with id being the key and ClientConn being the value.
func (s *Server) GetClients() map[int]*ClientConn {
	s.mu.RLock()
	defer s.mu.RUnlock()

	copyMap := make(map[int]*ClientConn, len(s.clients))
	for k, v := range s.clients {
		copyMap[k] = v
	}
	return copyMap
}

// GetClient returns a ClientConn of a connected client with the id provided.
func (s *Server) GetClient(id int) *ClientConn {
	s.mu.RLock()
	defer s.mu.RUnlock()
	client, ok := s.clients[id]
	if ok {
		return client
	}
	return nil
}

// Close shuts down the server and all active connections.
func (s *Server) Close() error {
	return s.server.Close()
}
