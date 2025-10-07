package kcp_test

import (
	"context"
	"net"
	"testing"

	"github.com/Onyz107/onynet/internal/kcp"
)

func BenchmarkDial(b *testing.B) {
	addr, err := net.ResolveUDPAddr("udp", "127.0.0.1:9090")
	if err != nil {
		b.Fatal(err)
	}

	ctx := context.Background()

	server, err := kcp.NewServer(addr, ctx)
	if err != nil {
		b.Fatal(err)
	}
	defer server.Close()

	b.ResetTimer()
	for b.Loop() {
		kcp.Dial(addr, ctx)
	}
}

func BenchmarkNewServer(b *testing.B) {
	ctx := context.Background()

	b.ResetTimer()
	for b.Loop() {
		addr, err := net.ResolveUDPAddr("udp", "127.0.0.1:0")
		if err != nil {
			b.Fatal(err)
		}
		if _, err := kcp.NewServer(addr, ctx); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkServer_Accept(b *testing.B) {
	addr, err := net.ResolveUDPAddr("udp", "127.0.0.1:5647")
	if err != nil {
		b.Fatal(err)
	}

	ctx := context.Background()

	server, err := kcp.NewServer(addr, ctx)
	if err != nil {
		b.Fatal(err)
	}
	defer server.Close()

	b.ResetTimer()
	for b.Loop() {
		go func() {
			client, err := kcp.Dial(addr, ctx)
			if err != nil {
				b.Fatal(err)
			}
			client.Write([]byte("1"))
		}()
		server.Accept()
	}
}
