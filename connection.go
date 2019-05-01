package hansip

import (
	"errors"
	"sync/atomic"
	"time"

	"github.com/go-pg/pg"
)

var errPingTimeout = errors.New("ping timeout")

// connection abstracts connection to a database server.
// it handles connection updates by pinging the server every connTickDelay.
type connection struct {
	host string
	s    sql

	pingTimeout    time.Duration
	connCheckDelay time.Duration

	pingFn      func() error
	pingRunning int32
	closeFn     func()

	// 1 for connected, 0 for not
	connected int32

	closed   bool
	quitChan chan struct{}
}

// create a new connection instance
// and start loop in background to update connection status
func newConnection(options *pg.Options, pingTimeout, connCheckDelay time.Duration) (*connection, error) {
	db := pg.Connect(options)
	conn := &connection{
		host: options.Addr,
		s: &gopgSQL{
			db: db,
		},
		pingTimeout:    pingTimeout,
		connCheckDelay: connCheckDelay,
		quitChan:       make(chan struct{}),
		pingFn: func() error {
			_, err := db.Exec("select 1;")
			return err
		},
		closeFn: func() {
			db.Close()
		},
	}

	// check if connection is working
	if err := conn.ping(); err != nil {
		return nil, err
	}
	conn.updateStatus()

	// start main loop
	go conn.loop()

	return conn, nil
}

func (c *connection) ping() error {
	if !atomic.CompareAndSwapInt32(&c.pingRunning, 0, 1) {
		return nil
	}
	defer atomic.StoreInt32(&c.pingRunning, 0)

	errChan := make(chan error)
	go func() {
		errChan <- c.pingFn()
	}()

	select {
	case <-time.After(c.pingTimeout):
		return errPingTimeout
	case err := <-errChan:
		return err
	}
}

func (c *connection) getConnected() bool {
	return atomic.LoadInt32(&c.connected) == 1
}

func (c *connection) setConnected(connected bool) {
	if connected {
		atomic.StoreInt32(&c.connected, 1)
	} else {
		atomic.StoreInt32(&c.connected, 0)
	}
}

func (c *connection) loop() {
	ticker := time.NewTicker(c.connCheckDelay)
	for {
		select {
		case <-ticker.C:
			c.updateStatus()
		case <-c.quitChan:
			return
		}
	}
}

func (c *connection) updateStatus() {
	connected := c.ping() == nil
	c.setConnected(connected)
}

func (c *connection) quit() {
	if c.closed {
		return
	}

	// stop loop
	c.quitChan <- struct{}{}

	if c.closeFn != nil {
		c.closeFn()
	}
	c.setConnected(false)
	c.closed = true
}
