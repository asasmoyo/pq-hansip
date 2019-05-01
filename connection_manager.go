package hansip

import (
	"math/rand"
	"sync"
	"time"
)

// connectionManager abstracts all database connections we are currently possessing.
// there are one connection to master and n number connections to slaves.
type connectionManager struct {
	master       *connection
	slaves       []*connection
	activeSlaves []*connection
	mutex        sync.RWMutex

	connCheckDelay time.Duration

	closed   bool
	quitChan chan struct{}
}

func newConnectionManager(connCheckDelay time.Duration) *connectionManager {
	manager := &connectionManager{
		slaves:         []*connection{},
		connCheckDelay: connCheckDelay,
		quitChan:       make(chan struct{}),
	}
	go manager.loop()
	return manager
}

func (m *connectionManager) loop() {
	ticker := time.NewTicker(m.connCheckDelay)
	for {
		select {
		case <-ticker.C:
			m.updateActiveSlaves()
		case <-m.quitChan:
			return
		}
	}
}

func (m *connectionManager) getSlaves() []*connection {
	m.mutex.RLock()
	slaves := m.slaves
	m.mutex.RUnlock()
	return slaves
}

func (m *connectionManager) addSlave(conn *connection) {
	m.mutex.Lock()
	m.slaves = append(m.slaves, conn)
	m.mutex.Unlock()

	m.updateActiveSlaves()
}

func (m *connectionManager) getActiveSlaves() []*connection {
	m.mutex.RLock()
	slaves := m.activeSlaves
	m.mutex.RUnlock()
	return slaves
}

func (m *connectionManager) setActiveSlaves(slaves []*connection) {
	m.mutex.Lock()
	m.activeSlaves = slaves
	m.mutex.Unlock()
}

func (m *connectionManager) updateActiveSlaves() {
	current := m.getSlaves()
	if len(current) == 0 {
		return
	}

	slaves := make([]*connection, 0, len(current))
	for _, conn := range current {
		if conn.getConnected() {
			slaves = append(slaves, conn)
		}
	}
	m.setActiveSlaves(slaves)
}

func (m *connectionManager) reader() sql {
	current := m.getActiveSlaves()
	n := len(current)
	if n == 0 {
		return m.writer()
	}
	return current[rand.Intn(n)].s
}

func (m *connectionManager) writer() sql {
	if !m.master.getConnected() {
		return nil
	}
	return m.master.s
}

func (m *connectionManager) quit() {
	if m.closed {
		return
	}

	// stop loop
	m.quitChan <- struct{}{}

	if m.master != nil {
		m.master.quit()
	}
	for _, conn := range m.slaves {
		conn.quit()
	}
	m.updateActiveSlaves()
	m.closed = true
}
