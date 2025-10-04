package kcp_test

import (
	"context"
	"net"
	"testing"

	"github.com/Onyz107/onylogger"
	"github.com/Onyz107/onynet/internal/kcp"
	"github.com/Onyz107/onynet/internal/logger"
)

func newServer(tb testing.TB) *kcp.Server {
	tb.Helper()

	addr, err := net.ResolveUDPAddr("udp", "127.0.0.1:9090")
	if err != nil {
		tb.Fatal(err)
	}

	server, err := kcp.NewServer(addr, context.Background())
	if err != nil {
		tb.Fatal(err)
	}

	return server
}

func newClient(tb testing.TB) *kcp.Client {
	tb.Helper()

	addr, err := net.ResolveUDPAddr("udp", "127.0.0.1:9090")
	if err != nil {
		tb.Fatal(err)
	}

	client, err := kcp.Dial(addr, context.Background())
	if err != nil {
		tb.Fatal(err)
	}
	client.Write([]byte("0")) // Need to write something for it to be accepted

	return client
}

func TestConnect(t *testing.T) {
	logger.Log.SetLevel(onylogger.LevelDebug)

	server := newServer(t)
	defer server.Close()
	t.Log("Server started")

	client := newClient(t)
	defer client.Close()
	t.Log("Client started")

	clientConn, err := server.Accept()
	if err != nil {
		t.Fatal(err)
	}
	defer clientConn.Close()

	t.Log("Client connected")
}
