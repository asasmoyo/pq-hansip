package hansip

import (
	"errors"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

type ConnectionTestSuite struct {
	suite.Suite
}

func TestConnection(t *testing.T) {
	s := &ConnectionTestSuite{}
	suite.Run(t, s)
}

func (s *ConnectionTestSuite) TestPingOK() {
	c := &connection{
		host:           "dummy",
		quitChan:       make(chan struct{}),
		pingTimeout:    1 * time.Second,
		connCheckDelay: 100 * time.Millisecond,
		s:              &dummySQL{},
		pingFn: func() error {
			return nil
		},
	}
	s.Nil(c.ping())
}

func (s *ConnectionTestSuite) TestPingFail() {
	c := &connection{
		host:           "dummy",
		quitChan:       make(chan struct{}),
		pingTimeout:    1 * time.Second,
		connCheckDelay: 100 * time.Millisecond,
		s:              &dummySQL{},
		pingFn: func() error {
			return errors.New("something fails")
		},
	}
	s.EqualError(c.ping(), "something fails")
}

func (s *ConnectionTestSuite) TestPingTimeout() {
	c := &connection{
		host:        "dummy",
		quitChan:    make(chan struct{}),
		pingTimeout: 100 * time.Millisecond,
		pingFn: func() error {
			time.Sleep(200 * time.Millisecond)
			return errors.New("fail")
		},
	}
	s.Equal(errPingTimeout, c.ping())
}

func (s *ConnectionTestSuite) TestUpdateStatusWithPingOK() {
	c := &connection{
		host:           "dummy",
		quitChan:       make(chan struct{}),
		pingTimeout:    1 * time.Second,
		connCheckDelay: 100 * time.Millisecond,
		s:              &dummySQL{},
		pingFn: func() error {
			return nil
		},
	}
	c.updateStatus()
	s.True(c.getConnected())
}

func (s *ConnectionTestSuite) TestUpdateStatusWithPingError() {
	c := &connection{
		connected:      1,
		host:           "dummy",
		quitChan:       make(chan struct{}),
		pingTimeout:    1 * time.Second,
		connCheckDelay: 100 * time.Millisecond,
		s:              &dummySQL{},
		pingFn: func() error {
			return errors.New("something fails")
		},
	}
	c.updateStatus()
	s.False(c.getConnected())
}

func (s *ConnectionTestSuite) TestConnected() {
	c := &connection{
		host:           "dummy",
		quitChan:       make(chan struct{}),
		pingTimeout:    1 * time.Second,
		connCheckDelay: 100 * time.Millisecond,
		s:              &dummySQL{},
		pingFn: func() error {
			return errors.New("something fails")
		},
	}
	c.setConnected(true)
	s.True(c.getConnected())
	c.setConnected(false)
	s.False(c.getConnected())
}

func (s *ConnectionTestSuite) TestQuit() {
	var closed bool
	c := &connection{
		host:           "dummy",
		quitChan:       make(chan struct{}),
		pingTimeout:    1 * time.Second,
		connCheckDelay: 100 * time.Millisecond,
		s:              &dummySQL{},
		pingFn: func() error {
			return errors.New("fail")
		},
		closeFn: func() {
			closed = true
		},
	}
	go c.loop()

	c.quit()
	time.Sleep(300 * time.Millisecond) // wait quit to finish running

	s.True(closed)
	s.False(c.getConnected())
}

func (s *ConnectionTestSuite) TestLoopTriggersUpdateStatus() {
	var pingErrFlag int32

	c := &connection{
		host:           "dummy",
		quitChan:       make(chan struct{}),
		pingTimeout:    1 * time.Second,
		connCheckDelay: 100 * time.Millisecond,
		s:              &dummySQL{},
		pingFn: func() error {
			if atomic.LoadInt32(&pingErrFlag) == 0 {
				return nil
			}
			return errors.New("something fails")
		},
	}
	go c.loop()

	// wait loop to kick in
	time.Sleep(300 * time.Millisecond)
	s.True(c.getConnected())

	atomic.StoreInt32(&pingErrFlag, 1)
	time.Sleep(300 * time.Millisecond)
	s.False(c.getConnected())
}
