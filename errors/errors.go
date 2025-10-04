package errors

import "errors"

// Miscellaneous error
var (
	ErrCtxCancelled    = errors.New("context cancelled")
	ErrWrite           = errors.New("write error")
	ErrRead            = errors.New("read error")
	ErrShortWrite      = errors.New("short write")
	ErrHeartbeatStream = errors.New("failed to open heartbeat stream")
)

// Crypto error
var (
	ErrCipher  = errors.New("invalid key size")
	ErrGCM     = errors.New("failed to create GCM")
	ErrShort   = errors.New("ciphertext too short")
	ErrDecrypt = errors.New("ciphertext corrupted")
)

// KCP error
var (
	ErrBadAddr = errors.New("invalid or unreachable address")
	ErrAccept  = errors.New("accept failed on KCP listener")
)

// Stream error
var (
	ErrAcceptStream = errors.New("failed to accept stream")
	ErrOpenStream   = errors.New("failed to open stream")
	ErrSetDeadline  = errors.New("failed to set deadline")
	ErrNameMismatch = errors.New("name mismatch")
	ErrTimeout      = errors.New("timeout")
	ErrNameTooLong  = errors.New("name too long")
)

// Transfer error
var (
	ErrSmallBuffer  = errors.New("buffer too small")
	ErrAESKey       = errors.New("invalid AES key")
	ErrStreamCipher = errors.New("failed to create cipher stream")
)

// Server error
var (
	ErrNewServer     = errors.New("failed to create new server")
	ErrAcceptClient  = errors.New("failed to accept client")
	ErrCreateSession = errors.New("failed to create session")
)

// Client error
var (
	ErrDial = errors.New("failed to dial")
	ErrAuth = errors.New("failed to authorize")
)

// Auth error
var (
	ErrPublickey  = errors.New("public key malformed or does not match private key")
	ErrPrivateKey = errors.New("private key malformed or does not match public key")
)

// Heartbeat error
var (
	ErrUnexpectedMsg = errors.New("unexpected message received")
)
