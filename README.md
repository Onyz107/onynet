# OnyNet

A high-performance, secure networking library for Go built on top of KCP and SMUX, featuring mutual TLS-style authentication, AES-GCM encryption, and multiplexed streams.

## Features

- **High Performance**: Built on KCP (ARQ protocol) for reliable UDP communication with optimized window sizes and no-delay settings
- **Secure by Default**: Optional RSA + AES-GCM encryption with mutual authentication
- **Stream Multiplexing**: Multiple logical streams over a single connection using SMUX
- **Named Streams**: Easy-to-use named stream API for organizing communication channels
- **Automatic Heartbeat**: Built-in connection health monitoring with automatic cleanup
- **Context-Aware**: Full context.Context support for graceful shutdown and cancellation
- **Flexible Data Transfer**: Multiple transfer modes including raw, serialized, encrypted, and streaming

## Installation

```bash
go get github.com/Onyz107/onynet
```

## Quick Start

### Server Example

```go
package main

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"log"
	"net"

	"github.com/Onyz107/onynet"
)

var privateKey = `-----BEGIN PRIVATE KEY-----
MIIEvAIBADANBgkqhkiG9w0BAQEFAASCBKYwggSiAgEAAoIBAQC+qrCgngPYwNri
m3CduCYFgTzgOsvBQf97IoEBtAHDKjTwl+fosYxmL7VwfFiMTnb+EGYOFl6wR8kv
ZGzdMhaqhDOiHEbmyPwQ72BMuhpXW5pYzFYUPtzxOrVcFGmnDS2FwFOr3eY4B2yn
i7+Its33pvHjKIu4tdfNRanEP29+/GoyLTfgXsazfk97ZzavYEGVZy96TDgLsHAy
w0b8qMRfBJ3FIuGsT9F5W6uz6/eTocnDnSIeQsbvHf+vJnZgIWUmaxcUawV4hOAy
XyPQa7TaE7j4l0FzbNVWuPgh6n3egOkeuwoJkl+qoxLoItAyQVz48zGt/a9CrrCv
lcfGrq7dAgMBAAECggEAFye1kZv/DZjGPFTyRUUy4OJEGVsqmYrBUxvqnOFWgXQj
v8BC+sBtM4/Bsip3unpg+xPUwQs4bqIZLbc9fVNy6zxo5NwYRDjOW+QmRssnHcrT
IMuX/Jdxz534nnkgJ5hzGdY0kx+8sLs+F87h8OF6fAP7RkILTeBPl+9I4btud/ZL
6Y/RBNulfksvviYeQvlNn4UXJUsrWV86hbTLx8fLMRpgnIKp7kFijbkNbd3i5g8k
DZoDKTsstTACHgBUJShiPFFEfi33C5z5j8h/7SX7UMGHYwcKS70qfL4/7BI499f4
a5PgSAi0G6pOcqy5KeOiBCUGFCe0SxZISJVQq8QCoQKBgQDebTdNgzhszd7nZh5X
dRaHJpk9QIgteWDzRvKr3FnZCaUmJLEgw6ro789GF9M76YlPlPNhFSti6tIx6XtK
xGcvsgCBSOr5wqJ39ah2cT+9U4+R8eHYRIJ8vItqBc0Shtu4oyOVxgQpLPf3seWq
XssTpiN7/u81CZLIwbG+hO2spQKBgQDbcj3JBPbtyvfkEV1xvE/y88BBNcykC3T2
opC1/1pe7/BY0QEhV3EhDCZ+aVSRNWZ5MYYPBwWjPcBseSJWHVWuVwBiBv2LtT2O
hWd27ZGHT40vqSCbHDZ5F3fuVjMxknoMP16XW87siDSCyIThbIJJdZSNJn3ZkP88
NJoc+KdL2QKBgDOQTCbLCdSncUphsR0DRuKz/whlImywW9pqEy3mWmnnQ4LxNKLs
2X1AwuNz+INGI4/wbu+Nsc4vs+TgRLXTjDxRXq6aEecuyO1YZOJ4ZJdmfL0PvxSc
5Uc3inZcu+rUmrFWGJTLIAHPq/ifJCf368o1VLqVIi1Ad+fUh3ksZdEFAoGAZA0a
DmStEI3Rp7IjII/zA5oOtayJuOFgWnKT9+aMlWxf8J6aHVF4ytB3XHs5i1sFdYwW
yxMwhtTIvqwb85c8UVhpXEhDoUbo4eoA2kBGcaLbhDdgHlgmnd8NVyUGAjv+WUcr
IWdCWKVhC5/QtdZ7MHLjX9eE2YU6WYDCIyNbY0ECgYApMJOoEWUHJ0wERC/2mcgr
7058V2oyqfrp23bwFH+lsgWPGpYAhbaGKHpbV762CMwG9fWTGjgqu1c9LrLuTuDJ
KIZSWO+X+K7YTZloU8aIxacSuUJa8YRv/cNkLvB6G1R5Q3qY3imK2yM7gRE4EpZa
O3l3AHVbdtj98riUs3Pmfg==
-----END PRIVATE KEY-----
`

func main() {
	addr, err := net.ResolveUDPAddr("udp", "127.0.0.1:8080")
	if err != nil {
		panic(err)
	}

	block, _ := pem.Decode([]byte(privateKey))
	if block == nil || block.Type != "PRIVATE KEY" {
		panic("failed to decode PEM block containing private key")
	}

	priv, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		panic(err)
	}

	rsaPriv, ok := priv.(*rsa.PrivateKey)
	if !ok {
		panic("not RSA private key")
	}

	log.Printf("Creating new server listening on: %s\n", addr.String())
	server, err := onynet.NewServer(addr, rsaPriv, context.Background())
	if err != nil {
		panic(err)
	}
	defer server.Close()

	clientConn, err := server.AcceptStream()
	if err != nil {
		panic(err)
	}
	defer clientConn.Close()
	log.Printf("Accepted new client: %s\n", clientConn.RemoteAddr().String())

	stream, err := clientConn.AcceptStream("testStream", context.Background(), 0)
	if err != nil {
		panic(err)
	}
	log.Printf("Accepted stream from client.\n")

	var msg string
	for {
		fmt.Print(">> ")
		fmt.Scanln(&msg)

		if err := stream.SendEncrypted([]byte(msg), 0); err != nil {
			panic(err)
		}
	}
}
```

### Client Example

```go
package main

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"log"
	"net"

	"github.com/Onyz107/onynet"
)

var publicKey = `-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAvqqwoJ4D2MDa4ptwnbgm
BYE84DrLwUH/eyKBAbQBwyo08Jfn6LGMZi+1cHxYjE52/hBmDhZesEfJL2Rs3TIW
qoQzohxG5sj8EO9gTLoaV1uaWMxWFD7c8Tq1XBRppw0thcBTq93mOAdsp4u/iLbN
96bx4yiLuLXXzUWpxD9vfvxqMi034F7Gs35Pe2c2r2BBlWcvekw4C7BwMsNG/KjE
XwSdxSLhrE/ReVurs+v3k6HJw50iHkLG7x3/ryZ2YCFlJmsXFGsFeITgMl8j0Gu0
2hO4+JdBc2zVVrj4Iep93oDpHrsKCZJfqqMS6CLQMkFc+PMxrf2vQq6wr5XHxq6u
3QIDAQAB
-----END PUBLIC KEY-----
`

func main() {
	addr, err := net.ResolveUDPAddr("udp", "127.0.0.1:8080")
	if err != nil {
		panic(err)
	}

	block, _ := pem.Decode([]byte(publicKey))
	if block == nil || block.Type != "PUBLIC KEY" {
		panic("failed to decode PEM block containing public key")
	}

	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		panic(err)
	}

	rsaPub, ok := pub.(*rsa.PublicKey)
	if !ok {
		panic("not RSA public key")
	}

	log.Printf("Connecting to server on: %s\n", addr.String())
	client, err := onynet.Dial(addr, rsaPub, context.Background())
	if err != nil {
		panic(err)
	}
	defer client.Close()

	stream, err := client.OpenStream("testStream", context.Background(), 0)
	if err != nil {
		panic(err)
	}
	defer stream.Close()
	log.Printf("Opened stream.")

	buf := make([]byte, 1024)
	for {
		n, err := stream.ReceiveEncrypted(buf, 0)
		if err != nil {
			panic(err)
		}

		msg := string(buf[:n])
		if msg == "exit" {
			return
		}

		fmt.Printf("Received message from server: %s\n", msg)
	}
}
```

## Authentication

OnyNet supports optional mutual authentication using RSA keys:

```go
// Generate RSA keys (do this once, store securely)
privateKey, _ := rsa.GenerateKey(rand.Reader, 2048)
publicKey := &privateKey.PublicKey

// Server side
server, _ := onynet.NewServer(addr, privateKey, ctx)

// Client side
client, _ := onynet.Dial(addr, publicKey, ctx)
```

When authentication is enabled:
1. Client encrypts an AES-256 key with the server's public key
2. Server decrypts the AES key using its private key
3. Server proves its identity by signing a challenge
4. All subsequent stream data can be encrypted with the shared AES key

## Stream Operations

### Named Streams

Streams are identified by names, making it easy to organize different types of communication:

```go
// Server accepts streams
dataStream, _ := client.AcceptStream("data", context.Background(), 5*time.Second)
controlStream, _ := client.AcceptStream("control", context.Background(), 5*time.Second)

// Client opens streams
dataStream, _ := client.OpenStream("data", context.Background(), 5*time.Second)
controlStream, _ := client.OpenStream("control", context.Background(), 5*time.Second)
```

### Transfer Methods

OnyNet provides multiple ways to transfer data:

#### 1. Raw Transfer
```go
// Send
data := []byte("raw data")
stream.Send(data, 10*time.Second)

// Receive
buf := make([]byte, 1024)
stream.Receive(buf, 10*time.Second)
```

#### 2. Serialized Transfer (with length prefix)
```go
// Send
stream.SendSerialized(data, 10*time.Second)

// Receive
n, _ := stream.ReceiveSerialized(buf, 10*time.Second)
actualData := buf[:n]
```

#### 3. Encrypted Transfer (AES-GCM)
```go
// Send (only works if authentication is enabled)
stream.SendEncrypted(data, 10*time.Second)

// Receive
n, _ := stream.ReceiveEncrypted(buf, 10*time.Second)
```

#### 4. Streaming Transfer
```go
// Send large data via streaming
writer := stream.NewStreamedSender(30*time.Second)
io.Copy(writer, largeFile)
writer.Close()

// Receive streamed data
reader := stream.NewStreamedReceiver(30*time.Second)
io.Copy(destination, reader)
reader.Close()
```

#### 5. Encrypted Streaming
```go
// Send encrypted stream
writer, _ := stream.NewStreamedEncryptedSender(30*time.Second)
io.Copy(writer, sensitiveFile)
writer.Close()

// Receive encrypted stream
reader, _ := stream.NewStreamedEncryptedReceiver(30*time.Second)
io.Copy(destination, reader)
reader.Close()
```

## Server Management

Get information about connected clients:

```go
// Get all clients
clients := server.GetClients()
for id, client := range clients {
    log.Printf("Client %d: %s", id, client.RemoteAddr())
}

// Get specific client
client := server.GetClient(5)
if client != nil {
    // Use client
}
```

## Architecture

OnyNet is built on three main layers:

1. **KCP Layer**: Provides reliable UDP transport with ARQ
2. **SMUX Layer**: Multiplexes multiple streams over a single KCP connection
3. **OnyNet Layer**: Adds named streams, authentication, encryption, and convenience methods

```
┌─────────────────────────────────────┐
│          Your Application           │
├─────────────────────────────────────┤
│  OnyNet (Named Streams, Auth, Enc)  │
├─────────────────────────────────────┤
│      SMUX (Stream Multiplexing)     │
├─────────────────────────────────────┤
│        KCP (Reliable UDP)           │
├─────────────────────────────────────┤
│              UDP                    │
└─────────────────────────────────────┘
```

## Error Handling

OnyNet provides structured errors in the `errors` package:

```go
import intErrors "github.com/Onyz107/onynet/errors"

_, err := client.OpenStream("test", context.Background(), 5*time.Second)
if errors.Is(err, intErrors.ErrTimeout) {
    // Handle timeout
} else if errors.Is(err, intErrors.ErrCtxCancelled) {
    // Handle cancellation
}
```

## Performance Tuning

KCP connections are pre-configured with optimized settings:
- Window size: 512/512
- No-delay mode: 1, 40ms interval, 2 resend, 1 no-congestion-control

These settings prioritize low latency and high throughput. Modify `internal/kcp/client.go` and `internal/kcp/server.go` if different settings are needed.

## Graceful Shutdown

Use context cancellation for graceful shutdown:

```go
ctx, cancel := context.WithCancel(context.Background())

server, _ := onynet.NewServer(addr, privateKey, ctx)

// Later, trigger shutdown
cancel()
// All connections and streams will be closed automatically
```

## Dependencies

- [kcp-go](https://github.com/xtaci/kcp-go) - KCP protocol implementation
- [smux](https://github.com/xtaci/smux) - Stream multiplexing
- [onylogger](https://github.com/Onyz107/onylogger) - Internal logging

## Thread Safety

- `Server.GetClients()` and `Server.GetClient()` are thread-safe
- Multiple goroutines can safely call `AcceptStream()` and `OpenStream()`
- Individual streams should not be used concurrently from multiple goroutines


### [Benchmarks](https://github.com/Onyz107/onynet/blob/main/benchmarks/benchmarks.md)