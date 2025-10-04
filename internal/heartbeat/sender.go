package heartbeat

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"time"

	intErrors "github.com/Onyz107/onynet/errors"
	"github.com/Onyz107/onynet/internal/logger"
)

// SendHeartbeat sends periodic heartbeat messages and checks for responses.
func SendHeartbeat(conn net.Conn, ctx context.Context) error {
	defer conn.SetDeadline(time.Time{})

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	buf := bufPool.Get().([]byte)
	defer bufPool.Put(buf)

	for {
		select {

		case <-ctx.Done():
			return intErrors.ErrCtxCancelled

		case <-ticker.C:
			conn.SetDeadline(time.Now().Add(15 * time.Second))
			n, err := conn.Write([]byte(senderMsg))
			if err != nil {
				return errors.Join(intErrors.ErrWrite, err)
			}
			if n != len(senderMsg) {
				return intErrors.ErrShortWrite
			}
			logger.Log.Debugf("SendHeartbeat: sent heartbeat")

			if _, err := io.ReadFull(conn, buf); err != nil {
				return errors.Join(intErrors.ErrRead, err)
			}
			heartbeatMsg := string(buf)

			if heartbeatMsg != receiverMsg {
				return errors.Join(intErrors.ErrUnexpectedMsg, fmt.Errorf("expected: %s: got: %s", receiverMsg, heartbeatMsg))
			}
			logger.Log.Debugf("SendHeartbeat: heartbeat received")
			conn.SetDeadline(time.Time{})
		}
	}
}
