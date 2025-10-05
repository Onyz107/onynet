package smux

import (
	"context"
	"encoding/binary"
	"errors"
	"io"
	"time"

	intErrors "github.com/Onyz107/onynet/errors"
	"github.com/Onyz107/onynet/internal/logger"
	"github.com/xtaci/smux"
)

type Manager struct {
	session *smux.Session
	aesKey  []byte
	ctx     context.Context
}

// NewManager wraps a smux session with AES key and context.
func NewManager(session *smux.Session, aesKey []byte, ctx context.Context) *Manager {
	manager := &Manager{session: session, aesKey: aesKey, ctx: ctx}

	go func() {
		select {
		case <-ctx.Done():
			logger.Log.Debug("closing smux session because of context cancellation")
			manager.Close()
			return
		case <-manager.CloseChan():
			return
		}
	}()

	return manager
}

// Accept waits for a stream with a given name.
func (m *Manager) Accept(name string, timeout time.Duration) (*Stream, error) {
	if len(name) > 0xFFFF {
		return nil, intErrors.ErrNameTooLong
	}

	ctx := m.ctx
	if timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(m.ctx, timeout)
		defer cancel()
	}

	streamTimeout := timeout
	if timeout > 5*time.Second {
		streamTimeout = 5 * time.Second
	}

	for {
		time.Sleep(50 * time.Millisecond)
		select {

		case <-ctx.Done():
			if errors.Is(ctx.Err(), context.DeadlineExceeded) {
				return nil, intErrors.ErrTimeout
			}
			return nil, intErrors.ErrCtxCancelled

		default:
			stream, err := m.acceptStream(name, streamTimeout)
			if err != nil {
				if errors.Is(err, smux.ErrTimeout) || errors.Is(err, intErrors.ErrNameMismatch) {
					continue
				}
				return nil, err
			}
			return stream, nil
		}
	}
}

func (m *Manager) acceptStream(name string, timeout time.Duration) (*Stream, error) {
	stream, err := m.session.AcceptStream()
	if err != nil {
		return nil, errors.Join(intErrors.ErrAcceptStream, err)
	}

	stream.SetDeadline(time.Now().Add(timeout))

	header := headerPool.Get().([]byte)
	defer headerPool.Put(header)
	if _, err := io.ReadFull(stream, header); err != nil {
		stream.Close()
		return nil, errors.Join(intErrors.ErrRead, err)
	}
	length := binary.BigEndian.Uint16(header)

	buf := make([]byte, length)
	if _, err := io.ReadFull(stream, buf); err != nil {
		stream.Close()
		return nil, errors.Join(intErrors.ErrRead, err)
	}

	if string(buf) != name {
		stream.Write([]byte{0})
		stream.Close()
		return nil, intErrors.ErrNameMismatch
	}

	if _, err := stream.Write([]byte{1}); err != nil {
		stream.Close()
		return nil, errors.Join(intErrors.ErrWrite, err)
	}

	stream.SetDeadline(time.Time{})
	wrapped := &Stream{stream: stream, aesKey: m.aesKey, ctx: m.ctx}
	go func() {
		select {
		case <-wrapped.ctx.Done():
			logger.Log.Debug("closing smux stream because of context cancellation")
			wrapped.Close()
			return
		case <-wrapped.stream.GetDieCh():
			return
		}
	}()

	return wrapped, nil
}

// Open creates a new stream with a given name.
func (m *Manager) Open(name string, timeout time.Duration) (*Stream, error) {
	if len(name) > 0xFFFF {
		return nil, intErrors.ErrNameTooLong
	}

	ctx := m.ctx
	if timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(m.ctx, timeout)
		defer cancel()
	}

	streamTimeout := timeout
	if timeout > 5*time.Second {
		streamTimeout = 5 * time.Second
	}

	for {
		time.Sleep(50 * time.Millisecond)
		select {

		case <-ctx.Done():
			if errors.Is(ctx.Err(), context.DeadlineExceeded) {
				return nil, intErrors.ErrTimeout
			}
			return nil, intErrors.ErrCtxCancelled

		default:
			stream, err := m.openStream(name, streamTimeout)
			if err != nil {
				if errors.Is(err, smux.ErrTimeout) || errors.Is(err, intErrors.ErrNameMismatch) {
					continue
				}
				return nil, err
			}
			return stream, nil
		}
	}
}

func (m *Manager) openStream(name string, timeout time.Duration) (*Stream, error) {
	stream, err := m.session.OpenStream()
	if err != nil {
		return nil, errors.Join(intErrors.ErrOpenStream, err)
	}

	stream.SetDeadline(time.Now().Add(timeout))

	header := headerPool.Get().([]byte)
	defer headerPool.Put(header)
	length := uint16(len(name))
	binary.BigEndian.PutUint16(header, length)

	n, err := stream.Write(header)
	if err != nil {
		stream.Close()
		return nil, errors.Join(intErrors.ErrWrite, err)
	}
	if n != len(header) {
		stream.Close()
		return nil, intErrors.ErrShortWrite
	}

	n, err = stream.Write([]byte(name))
	if err != nil {
		stream.Close()
		return nil, errors.Join(intErrors.ErrWrite, err)
	}
	if n != len(name) {
		stream.Close()
		return nil, intErrors.ErrShortWrite
	}

	buf := make([]byte, 1)
	if _, err := io.ReadFull(stream, buf); err != nil {
		stream.Close()
		return nil, errors.Join(intErrors.ErrRead, err)
	}

	if buf[0] != 1 {
		stream.Close()
		return nil, intErrors.ErrNameMismatch
	}

	stream.SetDeadline(time.Time{})
	wrapped := &Stream{stream: stream, aesKey: m.aesKey, ctx: m.ctx}
	go func() {
		select {
		case <-wrapped.ctx.Done():
			logger.Log.Debug("closing smux stream because of context cancellation")
			wrapped.Close()
			return
		case <-wrapped.stream.GetDieCh():
			return
		}
	}()

	return wrapped, nil
}

// Close terminates the session.
func (m *Manager) Close() error {
	return m.session.Close()
}

// CloseChan returns channel closed when session is terminated.
func (m *Manager) CloseChan() <-chan struct{} {
	return m.session.CloseChan()
}
