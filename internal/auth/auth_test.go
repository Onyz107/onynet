package auth_test

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"net"
	"sync"
	"testing"

	"github.com/Onyz107/onylogger"
	"github.com/Onyz107/onynet/internal/auth"
	"github.com/Onyz107/onynet/internal/kcp"
	"github.com/Onyz107/onynet/internal/logger"
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

func parsePrivateKey(tb testing.TB, privateKey string) *rsa.PrivateKey {
	tb.Helper()

	block, _ := pem.Decode([]byte(privateKey))
	if block == nil {
		tb.Fatal("failed to parse PEM block")
	}

	key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		tb.Fatal(err)
	}

	rsaPriv, ok := key.(*rsa.PrivateKey)
	if !ok {
		tb.Fatal("failed to assert type")
	}
	return rsaPriv
}

func parsePublicKey(tb testing.TB, publicKey string) *rsa.PublicKey {
	tb.Helper()

	block, _ := pem.Decode([]byte(publicKey))
	if block == nil {
		tb.Fatal("failed to parse PEM block")
	}

	key, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		tb.Fatal(err)
	}

	rsaPub, ok := key.(*rsa.PublicKey)
	if !ok {
		tb.Fatal("failed to assert type")
	}
	return rsaPub
}

func TestAuth(t *testing.T) {
	server := newServer(t)
	defer server.Close()

	client := newClient(t)
	defer client.Close()

	clientConn, err := server.Accept()
	if err != nil {
		t.Fatal(err)
	}
	defer clientConn.Close()

	buf := make([]byte, 1)
	clientConn.Read(buf)

	serverAuth := func(tb testing.TB, c *kcp.ClientConn) {
		priv := parsePrivateKey(t, privateKey)
		if _, err := auth.AuthorizeClient(c, priv); err != nil {
			tb.Fatal(err)
		}
		tb.Log("Authorized client")
		if err := auth.AuthorizeSelfServer(c, priv); err != nil {
			tb.Fatal(err)
		}
		tb.Log("Authorized self server")
	}

	clientAuth := func(tb testing.TB, c *kcp.Client) {
		pub := parsePublicKey(t, publicKey)
		if _, err := auth.AuthorizeSelfClient(c, pub); err != nil {
			tb.Fatal(err)
		}
		tb.Log("Authorized self client")
		if err := auth.AuthorizeServer(c, pub); err != nil {
			tb.Fatal(err)
		}
		tb.Log("Authorized server")
	}

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		serverAuth(t, clientConn)
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		clientAuth(t, client)
	}()
	wg.Wait()
}

func TestOnlyServerAuth(t *testing.T) {
	logger.Log.SetLevel(onylogger.LevelDebug)

	server := newServer(t)
	defer server.Close()

	client := newClient(t)
	defer client.Close()

	clientConn, err := server.Accept()
	if err != nil {
		t.Fatal(err)
	}
	defer clientConn.Close()

	buf := make([]byte, 1)
	clientConn.Read(buf)

	serverAuth := func(tb testing.TB, c *kcp.ClientConn) {
		priv := parsePrivateKey(t, privateKey)
		if _, err := auth.AuthorizeClient(c, priv); err == nil {
			tb.Fatal("no error returned")
		}
		if err := auth.AuthorizeSelfServer(c, priv); err == nil {
			tb.Fatal("no error returned")
		}
	}

	serverAuth(t, clientConn)
}
