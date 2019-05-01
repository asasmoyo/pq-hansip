package hansip

import (
	"errors"
	"sync"
	"time"

	"github.com/go-pg/pg"
)

const (
	defaultMaxAttempt      = 3
	defaultConnRetryDelay  = 3 * time.Second
	defaultConnCheckDelay  = 3 * time.Second
	defaultConnPingTimeout = 1 * time.Second
)

// errors definition
var (
	ErrNoSlaveAvailable  = errors.New("no slave connection available")
	ErrNoMasterAvailable = errors.New("no master connection available")
)

// Config contains pg.Options for remote postgres
type Config struct {
	PrependQueryWithCaller bool
	MaxConnAttempt         int
	ConnRetryDelay         time.Duration
	ConnCheckDelay         time.Duration
	ConnPingTimeout        time.Duration
}

// Cluster abstracts database connections to remote postgres.
type Cluster struct {
	manager *connectionManager
	mutex   sync.Mutex

	conf *Config
}

func (c *Cluster) setConnectionManager(manager *connectionManager) {
	c.mutex.Lock()
	c.manager = manager
	c.mutex.Unlock()
}

func (c *Cluster) getConnectionManager() *connectionManager {
	c.mutex.Lock()
	manager := c.manager
	c.mutex.Unlock()

	return manager
}

// SetMaster creates a connection to given connection info and set it as master
func (c *Cluster) SetMaster(opts *pg.Options) error {
	conn, err := newConnection(opts, c.conf.ConnPingTimeout, c.conf.ConnCheckDelay)
	if err != nil {
		return err
	}
	c.getConnectionManager().master = conn
	return nil
}

// AddSlave creates a connection to given connection info and add it as slave
func (c *Cluster) AddSlave(opts *pg.Options) error {
	conn, err := newConnection(opts, c.conf.ConnPingTimeout, c.conf.ConnCheckDelay)
	if err != nil {
		return err
	}
	c.getConnectionManager().addSlave(conn)
	c.getConnectionManager().updateActiveSlaves()
	return nil
}

// Query runs query to one of randomly-picked slave connection.
// If there is no slave available, the query will be run on writer.
func (c *Cluster) Query(dest interface{}, query string, args ...interface{}) error {
	conn := c.getConnectionManager().reader()
	if conn == nil {
		return ErrNoSlaveAvailable
	}
	return conn.query(dest, query, args...)
}

// WriterExec runs a query to master connection.
func (c *Cluster) WriterExec(query string, args ...interface{}) error {
	conn := c.getConnectionManager().writer()
	if conn == nil {
		return ErrNoMasterAvailable
	}
	return conn.exec(query, args...)
}

// WriterQuery runs query to master connection.
func (c *Cluster) WriterQuery(dest interface{}, query string, args ...interface{}) error {
	conn := c.getConnectionManager().writer()
	if conn == nil {
		return ErrNoMasterAvailable
	}
	return conn.query(dest, query, args...)
}

// NewTransaction creates a new database transaction.
// This method guaratees that the transaction will be run on master connection.
func (c *Cluster) NewTransaction() (Transaction, error) {
	conn := c.getConnectionManager().writer()
	if conn == nil {
		return nil, ErrNoMasterAvailable
	}
	return c.getConnectionManager().writer().newTransaction()
}

// Shutdown kills all connections.
func (c *Cluster) Shutdown() {
	c.getConnectionManager().quit()
}
