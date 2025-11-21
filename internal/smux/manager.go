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
		case <-manager.ctx.Done():
			logger.Log.Debug("closing smux session because of context cancellation")
			manager.Close()
			return
		case <-manager.CloseChan():
			return
		}
	}()

	return manager
}

// AcceptStream waits for a stream with a given name.
func (m *Manager) AcceptStream(name string, ctx context.Context, timeout time.Duration) (*Stream, error) {
	if len(name) > 0xFFFF {
		return nil, intErrors.ErrNameTooLong
	}

	logger.Log.Debugf("smux/manager AcceptStream: timeout is: %f", timeout.Seconds())

	start := time.Now()
	for {
		time.Sleep(50 * time.Millisecond)
		select {

		case <-ctx.Done():
			return nil, intErrors.ErrCtxCancelled

		default:
			if timeout > 0 && time.Since(start) >= timeout {
				return nil, intErrors.ErrTimeout
			}
			stream, err := m.accept(name, ctx, time.Second)
			if err != nil {
				if errors.Is(err, intErrors.ErrTimeout) || errors.Is(err, intErrors.ErrNameMismatch) {
					logger.Log.Debugf("smux/manager AcceptStream: got non crticial error: %v continuing", err)
					continue
				}
				return nil, err
			}
			return stream, nil
		}
	}
}

func (m *Manager) accept(name string, ctx context.Context, timeout time.Duration) (*Stream, error) {
	if timeout > 0 {
		logger.Log.Debugf("smux/manager accept: session timeout is: %f", timeout.Seconds())
		m.session.SetDeadline(time.Now().Add(timeout))
		defer m.session.SetDeadline(time.Time{})
	}

	stream, err := m.session.AcceptStream()
	if err != nil {
		return nil, errors.Join(intErrors.ErrAcceptStream, err)
	}

	if timeout > 0 {
		logger.Log.Debugf("smux/manager accept: stream timeout is: %f", timeout.Seconds())
		stream.SetDeadline(time.Now().Add(timeout))
	}

	headerPtr := headerPool.Get().(*[]byte)
	defer headerPool.Put(headerPtr)
	header := *headerPtr

	logger.Log.Debugf("smux/manager accept: reading header")
	if _, err := io.ReadFull(stream, header); err != nil {
		stream.Close()
		return nil, errors.Join(intErrors.ErrRead, err)
	}
	length := binary.BigEndian.Uint16(header)
	logger.Log.Debugf("smux/manager accept: read header and found length: %d", length)

	logger.Log.Debugf("smux/manager accept: reading name")
	buf := make([]byte, length)
	if _, err := io.ReadFull(stream, buf); err != nil {
		stream.Close()
		return nil, errors.Join(intErrors.ErrRead, err)
	}
	logger.Log.Debugf("smux/manager accept: read name as: %s", string(buf))

	if string(buf) != name {
		logger.Log.Debugf("smux/manager accept: name mismatch")
		logger.Log.Debugf("smux/manager accept: writing no ok")
		stream.Write([]byte{0})
		stream.Close()
		return nil, intErrors.ErrNameMismatch
	}

	logger.Log.Debugf("smux/manager accept: writing ok")
	if _, err := stream.Write([]byte{1}); err != nil {
		stream.Close()
		return nil, errors.Join(intErrors.ErrWrite, err)
	}
	logger.Log.Debugf("smux/manager accept: wrote ok")

	stream.SetDeadline(time.Time{})

	wrapped := &Stream{stream: stream, aesKey: m.aesKey, ctx: ctx}
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

// OpenStream creates a new stream with a given name.
func (m *Manager) OpenStream(name string, ctx context.Context, timeout time.Duration) (*Stream, error) {
	if len(name) > 0xFFFF {
		return nil, intErrors.ErrNameTooLong
	}

	logger.Log.Debugf("smux/manager OpenStream: timeout is: %f", timeout.Seconds())
	start := time.Now()
	for {
		time.Sleep(50 * time.Millisecond)
		select {

		case <-ctx.Done():
			return nil, intErrors.ErrCtxCancelled

		default:
			if timeout > 0 && time.Since(start) >= timeout {
				return nil, intErrors.ErrTimeout
			}
			stream, err := m.open(name, ctx, time.Second)
			if err != nil {
				if errors.Is(err, intErrors.ErrTimeout) || errors.Is(err, intErrors.ErrNameMismatch) {
					logger.Log.Debugf("smux/manager OpenStream: got noncrticial error: %v continuing", err)
					continue
				}
				return nil, err
			}
			return stream, nil
		}
	}
}

func (m *Manager) open(name string, ctx context.Context, timeout time.Duration) (*Stream, error) {
	if timeout > 0 {
		logger.Log.Debugf("smux/manager open: session timeout is: %f", timeout.Seconds())
		m.session.SetDeadline(time.Now().Add(timeout))
		defer m.session.SetDeadline(time.Time{})
	}

	stream, err := m.session.OpenStream()
	if err != nil {
		return nil, errors.Join(intErrors.ErrOpenStream, err)
	}

	if timeout > 0 {
		logger.Log.Debugf("smux/manager open: stream timeout is: %f", timeout.Seconds())
		stream.SetDeadline(time.Now().Add(timeout))
	}

	headerPtr := headerPool.Get().(*[]byte)
	defer headerPool.Put(headerPtr)
	header := *headerPtr

	length := uint16(len(name))
	logger.Log.Debugf("smux/manager open: found length: %d", length)
	binary.BigEndian.PutUint16(header, length)

	logger.Log.Debugf("smux/manager open: writing header")
	n, err := stream.Write(header)
	if err != nil {
		stream.Close()
		return nil, errors.Join(intErrors.ErrWrite, err)
	}
	if n != len(header) {
		stream.Close()
		return nil, intErrors.ErrShortWrite
	}
	logger.Log.Debugf("smux/manager open: wrote header")

	logger.Log.Debugf("smux/manager open: writing: %s", name)
	n, err = stream.Write([]byte(name))
	if err != nil {
		stream.Close()
		return nil, errors.Join(intErrors.ErrWrite, err)
	}
	if n != len(name) {
		stream.Close()
		return nil, intErrors.ErrShortWrite
	}
	logger.Log.Debugf("smux/manager open: wrote: %s", name)

	logger.Log.Debug("smux/manager open: reading ok/no ok")
	buf := make([]byte, 1)
	if _, err := io.ReadFull(stream, buf); err != nil {
		stream.Close()
		return nil, errors.Join(intErrors.ErrRead, err)
	}

	if buf[0] != 1 {
		logger.Log.Debug("smux/manager open: received no ok")
		stream.Close()
		return nil, intErrors.ErrNameMismatch
	}

	stream.SetDeadline(time.Time{})

	wrapped := &Stream{stream: stream, aesKey: m.aesKey, ctx: ctx}
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

// CloseChan returns a channel when the session is terminated.
func (m *Manager) CloseChan() <-chan struct{} {
	return m.session.CloseChan()
}
