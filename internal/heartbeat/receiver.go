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

// ReceiveHeartbeat handles incoming heartbeat messages and responds.
func ReceiveHeartbeat(conn net.Conn, ctx context.Context) error {
	defer conn.SetDeadline(time.Time{})

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	bufPtr := bufPool.Get().(*[]byte)
	defer bufPool.Put(bufPtr)
	buf := *bufPtr

	for {
		select {

		case <-ctx.Done():
			return intErrors.ErrCtxCancelled

		case <-ticker.C:
			conn.SetDeadline(time.Now().Add(15 * time.Second))
			if _, err := io.ReadFull(conn, buf); err != nil {
				return errors.Join(intErrors.ErrRead, err)
			}
			heartbeatMsg := string(buf)

			if heartbeatMsg != senderMsg {
				return errors.Join(intErrors.ErrUnexpectedMsg, fmt.Errorf("expected: %s: got: %s", senderMsg, heartbeatMsg))
			}
			logger.Log.Debugf("ReceiveHeartbeat: heartbeat received")

			n, err := conn.Write([]byte(receiverMsg))
			if err != nil {
				return errors.Join(intErrors.ErrWrite, err)
			}
			if n != len(buf) {
				return intErrors.ErrShortWrite
			}
			logger.Log.Debugf("ReceiveHeartbeat: heartbeat sent")
			conn.SetDeadline(time.Time{})
		}
	}
}
