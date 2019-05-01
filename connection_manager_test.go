package hansip

import (
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

type ConnectionManagerTestSuite struct {
	TestSuite
}

func TestConnectionManager(t *testing.T) {
	s := &ConnectionManagerTestSuite{}
	s.noCreateCluster = true
	suite.Run(t, s)
}

func (s *ConnectionManagerTestSuite) newIdleConnectionManager() *connectionManager {
	return &connectionManager{
		slaves:         []*connection{},
		connCheckDelay: 100 * time.Millisecond,
		quitChan:       make(chan struct{}),
	}
}

func (s *ConnectionManagerTestSuite) TestAddSlave() {
	manager := s.newIdleConnectionManager()
	manager.addSlave(&connection{
		connected: 1,
	})
	s.Len(manager.getActiveSlaves(), 1)
}

func (s *ConnectionManagerTestSuite) TestReader() {
	manager := s.newIdleConnectionManager()
	manager.addSlave(&connection{
		connected: 1,
		s:         &dummySQL{},
	})
	s.NotNil(manager.reader())

	// goes to master when no reader available
	manager.slaves[0].setConnected(false)
	manager.master = &connection{
		connected: 1,
		s:         &dummySQL{},
	}
	s.NotNil(manager.reader())
}

func (s *ConnectionManagerTestSuite) TestWriter() {
	manager := s.newIdleConnectionManager()
	manager.master = &connection{
		connected: 1,
		s:         &dummySQL{},
	}
	s.NotNil(manager.writer())
}

func (s *ConnectionManagerTestSuite) TestLoopAndQuit() {
	conn := &connection{
		connected:      1,
		pingTimeout:    1 * time.Second,
		connCheckDelay: 100 * time.Millisecond,
		pingFn: func() error {
			return nil
		},
		quitChan: make(chan struct{}),
	}
	go conn.loop()

	manager := s.newIdleConnectionManager()
	manager.addSlave(conn)
	go manager.loop()

	// wait loop to kick in
	time.Sleep(300 * time.Millisecond)

	s.Len(manager.getActiveSlaves(), 1)

	manager.quit()
	s.False(conn.getConnected())
	s.True(manager.closed)
}
