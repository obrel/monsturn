package task

import (
	"context"
	"errors"
	"regexp"
	"strings"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/obrel/go-lib/pkg/log"
	"github.com/obrel/go-lib/pkg/wrk"
	"github.com/obrel/monsturn/internal/pkg/pdb"
	"github.com/obrel/monsturn/internal/pkg/util"
)

const (
	statusRgx  = "^(turn\\/realm\\/)([a-z0-9\\-\\.]+)(\\/user\\/)([a-z0-9\\:]+)(\\/allocation\\/)([0-9]+)(\\/status)"
	trafficRgx = "^(turn\\/realm\\/)([a-z0-9\\-\\.]+)(\\/user\\/)([a-z0-9\\:]+)(\\/allocation\\/)([0-9]+)(\\/total_traffic)"
)

// Monitor :nodoc:
type Monitor struct {
	conn   redis.Conn
	topics []string
}

// NewMonitorTask
func NewMonitorTask(conn redis.Conn, topics []string) wrk.Task {
	return &Monitor{
		conn:   conn,
		topics: topics,
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
				if err := save(n.Channel, n.Data); err != nil {
					log.For("monitor", "do").Error(err)
					continue
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

func save(channel string, data []byte) error {
	// Check if status or traffic channel
	if strings.Contains(channel, "status") {
		rgx := regexp.MustCompile(statusRgx)
		match := rgx.Match([]byte(channel))

		if !match {
			return errors.New("Invalid status channel.")
		}

	} else if strings.Contains(channel, "total_traffic") {
		rgx := regexp.MustCompile(trafficRgx)
		match := rgx.Match([]byte(channel))

		if !match {
			return errors.New("Invalid traffic channel.")
		}

		ch, err := util.ExtractData(channel)
		if err != nil {
			return err
		}

		data, err := util.MessageParser(string(data[:]))
		if err != nil {
			return err
		}

		stat := pdb.NewStat(
			ch.Realm,
			ch.User,
			ch.Allocation,
			data.SentP,
			data.SentB,
			data.RecvP,
			data.RecvB,
		)

		err = stat.Insert()
		if err != nil {
			return err
		}
	} else {
		return errors.New("Invalid channel.")
	}

	return nil
}
