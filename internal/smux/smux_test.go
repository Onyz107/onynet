package smux_test

import (
	"context"
	"crypto/rand"
	"net"
	"testing"

	"github.com/Onyz107/onynet/internal/kcp"
	intSmux "github.com/Onyz107/onynet/internal/smux"
	"github.com/xtaci/smux"
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

func establishSession(tb testing.TB) (serverManager *intSmux.Manager, clientManager *intSmux.Manager) {
	server := newServer(tb)
	client := newClient(tb)

	clientConn, err := server.Accept()
	if err != nil {
		tb.Fatal(err)
	}

	buf := make([]byte, 1)
	clientConn.Read(buf)

	serverSession, err := smux.Server(clientConn, nil)
	if err != nil {
		tb.Fatal(err)
	}

	clientSession, err := smux.Client(client, nil)
	if err != nil {
		tb.Fatal(err)
	}

	serverManager = intSmux.NewManager(serverSession, []byte("23456789abcdeffedcba987654321021"), tb.Context())
	clientManager = intSmux.NewManager(clientSession, []byte("23456789abcdeffedcba987654321021"), tb.Context())

	return serverManager, clientManager
}

func BenchmarkManager_Accept(b *testing.B) {
	serverManager, clientManager := establishSession(b)
	defer serverManager.Close()
	defer clientManager.Close()

	go func() {
		for {
			clientManager.Open("testStream", 0)
		}
	}()

	b.ResetTimer()
	for b.Loop() {
		serverManager.Accept("testStream", 0)
	}
}

func BenchmarkManager_Open(b *testing.B) {
	serverManager, clientManager := establishSession(b)
	defer serverManager.Close()
	defer clientManager.Close()

	go func() {
		for {
			serverManager.Accept("testStream", 0)
		}
	}()

	b.ResetTimer()
	for b.Loop() {
		clientManager.Open("testStream", 0)
	}
}

func BenchmarkStream_Send(b *testing.B) {
	serverManager, clientManager := establishSession(b)
	defer serverManager.Close()
	defer clientManager.Close()

	acceptDone := make(chan *intSmux.Stream, 1)
	go func() {
		stream, err := serverManager.Accept("benchmarkSend", 0)
		if err != nil {
			b.Error(err)
			return
		}
		acceptDone <- stream
	}()

	clientStream, err := clientManager.Open("benchmarkSend", 0)
	if err != nil {
		b.Fatal(err)
	}
	defer clientStream.Close()

	serverStream := <-acceptDone
	defer serverStream.Close()

	data := make([]byte, 1024)
	rand.Read(data)

	done := make(chan struct{})
	go func() {
		buf := make([]byte, 1052) // 1024 + 28
		for {
			select {
			case <-done:
				return
			default:
				if err := clientStream.Receive(buf, 0); err != nil {
					return
				}
			}
		}
	}()

	b.ResetTimer()
	for b.Loop() {
		if err := serverStream.Send(data, 0); err != nil {
			b.Fatal(err)
		}
	}
	close(done)
}

func BenchmarkStream_NewStreamedSender(b *testing.B) {
	serverManager, clientManager := establishSession(b)
	defer serverManager.Close()
	defer clientManager.Close()

	// Wait for Accept to complete before use
	acceptDone := make(chan *intSmux.Stream, 1)
	go func() {
		stream, err := serverManager.Accept("benchmarkSend", 0)
		if err != nil {
			b.Error(err)
			return
		}
		acceptDone <- stream
	}()

	clientStream, err := clientManager.Open("benchmarkSend", 0)
	if err != nil {
		b.Fatal(err)
	}
	defer clientStream.Close()

	serverStream := <-acceptDone
	defer serverStream.Close()

	s := serverStream.NewStreamedSender(0)
	defer s.Close()

	done := make(chan struct{})
	go func() {
		r := clientStream.NewStreamedReceiver(0)
		defer r.Close()

		buf := make([]byte, 1024)
		for {
			select {
			case <-done:
				return
			default:
				if _, err := r.Read(buf); err != nil {
					return
				}
			}
		}
	}()

	data := make([]byte, 1024)
	rand.Read(data)

	b.ResetTimer()
	for b.Loop() {
		if _, err := s.Write(data); err != nil {
			b.Fatal(err)
		}
	}
	close(done)
}

func BenchmarkStream_SendSerialized(b *testing.B) {
	serverManager, clientManager := establishSession(b)
	defer serverManager.Close()
	defer clientManager.Close()

	acceptDone := make(chan *intSmux.Stream, 1)
	go func() {
		stream, err := serverManager.Accept("benchmarkSend", 0)
		if err != nil {
			b.Error(err)
			return
		}
		acceptDone <- stream
	}()

	clientStream, err := clientManager.Open("benchmarkSend", 0)
	if err != nil {
		b.Fatal(err)
	}
	defer clientStream.Close()

	serverStream := <-acceptDone
	defer serverStream.Close()

	data := make([]byte, 1024)
	rand.Read(data)

	done := make(chan struct{})
	go func() {
		buf := make([]byte, 1052) // 1024 + 28
		for {
			select {
			case <-done:
				return
			default:
				if _, err := clientStream.ReceiveSerialized(buf, 0); err != nil {
					return
				}
			}
		}
	}()

	b.ResetTimer()
	for b.Loop() {
		if err := serverStream.SendSerialized(data, 0); err != nil {
			b.Fatal(err)
		}
	}
	close(done)
}

func BenchmarkStream_SendEncrypted(b *testing.B) {
	serverManager, clientManager := establishSession(b)
	defer serverManager.Close()
	defer clientManager.Close()

	acceptDone := make(chan *intSmux.Stream, 1)
	go func() {
		stream, err := serverManager.Accept("benchmarkSend", 0)
		if err != nil {
			b.Error(err)
			return
		}
		acceptDone <- stream
	}()

	clientStream, err := clientManager.Open("benchmarkSend", 0)
	if err != nil {
		b.Fatal(err)
	}
	defer clientStream.Close()

	serverStream := <-acceptDone
	defer serverStream.Close()

	data := make([]byte, 1024)
	rand.Read(data)

	done := make(chan struct{})
	go func() {
		buf := make([]byte, 1052) // 1024 + 28
		for {
			select {
			case <-done:
				return
			default:
				if _, err := clientStream.ReceiveEncrypted(buf, 0); err != nil {
					return
				}
			}
		}
	}()

	b.ResetTimer()
	for b.Loop() {
		if err := serverStream.SendEncrypted(data, 0); err != nil {
			b.Fatal(err)
		}
	}
	close(done)
}

func BenchmarkStream_NewStreamedEncryptedSender(b *testing.B) {
	serverManager, clientManager := establishSession(b)
	defer serverManager.Close()
	defer clientManager.Close()

	// Wait for Accept to complete before use
	acceptDone := make(chan *intSmux.Stream, 1)
	go func() {
		stream, err := serverManager.Accept("benchmarkSend", 0)
		if err != nil {
			b.Error(err)
			return
		}
		acceptDone <- stream
	}()

	clientStream, err := clientManager.Open("benchmarkSend", 0)
	if err != nil {
		b.Fatal(err)
	}
	defer clientStream.Close()

	serverStream := <-acceptDone
	defer serverStream.Close()

	s, err := serverStream.NewStreamedEncryptedSender(0)
	if err != nil {
		b.Fatal(err)
	}
	defer s.Close()

	done := make(chan struct{})
	go func() {
		r, err := clientStream.NewStreamedEncryptedReceiver(0)
		if err != nil {
			b.Error(err)
			return
		}
		defer r.Close()

		buf := make([]byte, 1024)
		for {
			select {
			case <-done:
				return
			default:
				if _, err := r.Read(buf); err != nil {
					return
				}
			}
		}
	}()

	data := make([]byte, 1024)
	rand.Read(data)

	b.ResetTimer()
	for b.Loop() {
		if _, err := s.Write(data); err != nil {
			b.Fatal(err)
		}
	}
	close(done)
}

func BenchmarkStream_Receive(b *testing.B) {
	serverManager, clientManager := establishSession(b)
	defer serverManager.Close()
	defer clientManager.Close()

	// synchronize accept/open
	acceptDone := make(chan *intSmux.Stream, 1)
	go func() {
		stream, err := serverManager.Accept("benchmarkSend", 0)
		if err != nil {
			b.Error(err)
			return
		}
		acceptDone <- stream
	}()

	clientStream, err := clientManager.Open("benchmarkSend", 0)
	if err != nil {
		b.Fatal(err)
	}
	defer clientStream.Close()

	serverStream := <-acceptDone
	defer serverStream.Close()

	done := make(chan struct{})
	go func() {
		data := make([]byte, 1024)
		rand.Read(data)
		for {
			select {
			case <-done:
				return
			default:
				_ = clientStream.Send(data, 0)
			}
		}
	}()

	buf := make([]byte, 1024)
	b.ResetTimer()
	for b.Loop() {
		if err := serverStream.Receive(buf, 0); err != nil {
			b.Fatal(err)
		}
	}
	close(done)
}

func BenchmarkStream_NewStreamedReceiver(b *testing.B) {
	serverManager, clientManager := establishSession(b)
	defer serverManager.Close()
	defer clientManager.Close()

	// Wait for the server side to accept before proceeding
	acceptDone := make(chan *intSmux.Stream, 1)
	go func() {
		stream, err := serverManager.Accept("benchmarkSend", 0)
		if err != nil {
			b.Fatal(err)
		}
		acceptDone <- stream
	}()

	clientStream, err := clientManager.Open("benchmarkSend", 0)
	if err != nil {
		b.Fatal(err)
	}
	defer clientStream.Close()

	serverStream := <-acceptDone
	defer serverStream.Close()

	s := serverStream.NewStreamedSender(0)
	defer s.Close()

	r := clientStream.NewStreamedReceiver(0)
	defer r.Close()

	done := make(chan struct{})
	go func() {
		data := make([]byte, 1024)
		rand.Read(data)
		for {
			select {
			case <-done:
				return
			default:
				if _, err := s.Write(data); err != nil {
					// stop cleanly if receiver is closed
					return
				}
			}
		}
	}()

	buf := make([]byte, 1024)
	b.ResetTimer()
	for b.Loop() {
		if _, err := r.Read(buf); err != nil {
			b.Fatal(err)
		}
	}
	close(done)
}

func BenchmarkStream_ReceiveSerialized(b *testing.B) {
	serverManager, clientManager := establishSession(b)
	defer serverManager.Close()
	defer clientManager.Close()

	// synchronize accept/open
	acceptDone := make(chan *intSmux.Stream, 1)
	go func() {
		stream, err := serverManager.Accept("benchmarkSend", 0)
		if err != nil {
			b.Error(err)
			return
		}
		acceptDone <- stream
	}()

	clientStream, err := clientManager.Open("benchmarkSend", 0)
	if err != nil {
		b.Fatal(err)
	}
	defer clientStream.Close()

	serverStream := <-acceptDone
	defer serverStream.Close()

	done := make(chan struct{})
	go func() {
		data := make([]byte, 1024)
		rand.Read(data)
		for {
			select {
			case <-done:
				return
			default:
				_ = clientStream.SendSerialized(data, 0)
			}
		}
	}()

	buf := make([]byte, 1032) // 1024 + 8
	b.ResetTimer()
	for b.Loop() {
		if _, err := serverStream.ReceiveSerialized(buf, 0); err != nil {
			b.Fatal(err)
		}
	}
	close(done)
}

func BenchmarkStream_ReceiveEncrypted(b *testing.B) {
	serverManager, clientManager := establishSession(b)
	defer serverManager.Close()
	defer clientManager.Close()

	acceptDone := make(chan *intSmux.Stream, 1)
	go func() {
		stream, err := serverManager.Accept("benchmarkSend", 0)
		if err != nil {
			b.Error(err)
			return
		}
		acceptDone <- stream
	}()

	clientStream, err := clientManager.Open("benchmarkSend", 0)
	if err != nil {
		b.Fatal(err)
	}
	defer clientStream.Close()

	serverStream := <-acceptDone
	defer serverStream.Close()

	done := make(chan struct{})
	go func() {
		data := make([]byte, 1024)
		rand.Read(data)
		for {
			select {
			case <-done:
				return
			default:
				_ = clientStream.SendEncrypted(data, 0)
			}
		}
	}()

	buf := make([]byte, 1052) // 1024 + 28
	b.ResetTimer()
	for b.Loop() {
		if _, err := serverStream.ReceiveEncrypted(buf, 0); err != nil {
			b.Fatal(err)
		}
	}
	close(done)
}

func BenchmarkStream_NewStreamedEncryptedReceiver(b *testing.B) {
	serverManager, clientManager := establishSession(b)
	defer serverManager.Close()
	defer clientManager.Close()

	// Wait for the server side to accept before proceeding
	acceptDone := make(chan *intSmux.Stream, 1)
	go func() {
		stream, err := serverManager.Accept("benchmarkSend", 0)
		if err != nil {
			b.Fatal(err)
		}
		acceptDone <- stream
	}()

	clientStream, err := clientManager.Open("benchmarkSend", 0)
	if err != nil {
		b.Fatal(err)
	}
	defer clientStream.Close()

	serverStream := <-acceptDone
	defer serverStream.Close()

	s, err := serverStream.NewStreamedEncryptedSender(0)
	if err != nil {
		b.Fatal(err)
	}
	defer s.Close()

	r, err := clientStream.NewStreamedEncryptedReceiver(0)
	if err != nil {
		b.Fatal(err)
	}
	defer r.Close()

	done := make(chan struct{})
	go func() {
		data := make([]byte, 1024)
		rand.Read(data)
		for {
			select {
			case <-done:
				return
			default:
				if _, err := s.Write(data); err != nil {
					// stop cleanly if receiver is closed
					return
				}
			}
		}
	}()

	buf := make([]byte, 1024)
	b.ResetTimer()
	for b.Loop() {
		if _, err := r.Read(buf); err != nil {
			b.Fatal(err)
		}
	}
	close(done)
}
