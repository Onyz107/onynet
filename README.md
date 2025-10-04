# OnyNet

A high-performance networking library for Go with support for encrypted streams, KCP transport, smux multiplexing, and heartbeat monitoring.

## Features
- **KCP transport**: Reliable and fast UDP-based protocol.
- **AES-GCM and/or AES-CTR encryption**: Secure data transfer.
- **Smux multiplexing**: Multiple streams over a single connection.
- **Authentication**: RSA-based client/server authentication.
- **Heartbeat**: Connection liveness monitoring.
- **Concurrent-safe**: Designed for multiple streams and clients.

OnyNet is perfect for Go projects that need low-latency, secure, reliable, multiplexed connections over unreliable networks.
Think game servers, real-time collaboration tools, or IoT device communication. It shines when you want multiple
independent streams over a single UDP connection with AES encryption and RSA authentication built-in. It’s less suited
for ultra-high-frequency microservices or web APIs where standard HTTP/gRPC is simpler, but for custom protocols
needing speed, reliability, and security, OnyNet is solid.

## Installation
```bash
go get github.com/Onyz107/onynet
````

## Usage

### Server
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

	clientConn, err := server.Accept()
	if err != nil {
		panic(err)
	}
	defer clientConn.Close()
	log.Printf("Accepted new client: %s\n", clientConn.RemoteAddr().String())

	stream, err := clientConn.AcceptStream("testStream", 0)
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

### Client
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

	stream, err := client.OpenStream("testStream", 0)
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

### Stream Operations

* `Send(data, timeout)` – Send raw bytes.
* `NewStreamedSender(timeout)` - Get an `io.WriterCloser` that supports timeouts.
* `SendSerialized(data, timeout)` – Send length-prefixed data.
* `SendEncrypted(data, timeout)` – Send AES-GCM encrypted data.
* `NewStreamedEncryptedSender` - Same as `NewStreamedSender` but with AES-CTR encryption.
* `Receive(buffer, timeout)` – Receive raw bytes.
* `NewStreamedReceiver(timeout)` - Get an `io.ReadCloser` that supports timeouts.
* `ReceiveSerialized(buffer, timeout)` – Receive length-prefixed data.
* `ReceiveEncrypted(buffer, timeout)` – Receive and decrypt AES-GCM data.
* `NewStreamedEncryptedReceiver(timeout)` - Same as `NewStreamedReceiver` but with AES-CTR encryption.
