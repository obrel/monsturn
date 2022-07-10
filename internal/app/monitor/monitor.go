package monitor

import (
	"context"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/obrel/monsturn/internal/app/worker"
)

// Monitor :nodoc:
type Monitor struct {
	conn     redis.Conn
	topics   []string
	callback func(channel string, data []byte) error
}

// NewMonitorTask
func NewMonitorTask(conn redis.Conn, topics []string, callback func(channel string, data []byte) error) worker.Task {
	return &Monitor{
		conn:     conn,
		topics:   topics,
		callback: callback,
	}
}

// Do run monitor task
func (t *Monitor) Do(ctx context.Context) error {
	const healthCheckPeriod = time.Minute

	psc := redis.PubSubConn{Conn: t.conn}

	var topics []interface{} = make([]interface{}, len(t.topics))
	for i, tp := range t.topics {
		topics[i] = tp
	}

	if err := psc.PSubscribe(topics...); err != nil {
		return err
	}

	done := make(chan error, 1)

	// Start a goroutine to receive notifications from the server.
	go func() {
		for {
			switch n := psc.Receive().(type) {
			case error:
				done <- n
				return
			case redis.Message:
				if err := t.callback(n.Channel, n.Data); err != nil {
					done <- err
					return
				}
			case redis.Subscription:
				switch n.Count {
				case len(t.topics):
					// Notify application when all channels are subscribed.
				case 0:
					// Return from the goroutine when all channels are unsubscribed.
					done <- nil
					return
				}
			}
		}
	}()

	ticker := time.NewTicker(healthCheckPeriod)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			// Send ping to test health of connection and server. If
			// corresponding pong is not received, then receive on the
			// connection will timeout and the receive goroutine will exit.
			if err := psc.Ping(""); err != nil {
				break
			}
		case <-ctx.Done():
			break
		case err := <-done:
			// Return error from the receive goroutine.
			return err
		}
	}

	// Signal the receiving goroutine to exit by unsubscribing from all channels.
	if err := psc.Unsubscribe(); err != nil {
		return err
	}

	// Wait for goroutine to complete.
	return <-done
}
