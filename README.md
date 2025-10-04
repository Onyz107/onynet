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

TODO

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
