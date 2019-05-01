package hansip

import (
	"errors"
)

// error definitions
var (
	ErrTxFinished = errors.New("tx is already finished")
)

// NewCluster creates new cluster.
func NewCluster(conf *Config) *Cluster {
	if conf.MaxConnAttempt == 0 {
		conf.MaxConnAttempt = defaultMaxAttempt
	}
	if conf.ConnRetryDelay == 0 {
		conf.ConnRetryDelay = defaultConnRetryDelay
	}
	if conf.ConnCheckDelay == 0 {
		conf.ConnCheckDelay = defaultConnCheckDelay
	}
	if conf.ConnPingTimeout == 0 {
		conf.ConnPingTimeout = defaultConnPingTimeout
	}

	manager := newConnectionManager(conf.ConnCheckDelay)
	return &Cluster{
		manager: manager,
		conf:    conf,
	}
}
