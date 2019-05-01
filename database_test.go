package hansip

// func TestShutdown(t *testing.T) {
// 	connTickDelay = 1 * time.Millisecond
// 	csTickDelay = 1 * time.Millisecond

// 	masterClosed := false
// 	masterConn := &connection{
// 		connected: true,
// 		host:      "master",
// 		quitCh:    make(chan bool),
// 		s:         &dummySQL{},
// 		pingFn: func() error {
// 			return nil
// 		},
// 		closeFn: func() {
// 			masterClosed = true
// 		},
// 	}
// 	go masterConn.loop()
// 	slaveClosed := false
// 	slaveConns := []*connection{
// 		&connection{
// 			connected: true,
// 			host:      "slave1",
// 			quitCh:    make(chan bool),
// 			s:         &dummySQL{},
// 			pingFn: func() error {
// 				return nil
// 			},
// 			closeFn: func() {
// 				slaveClosed = true
// 			},
// 		},
// 	}
// 	go slaveConns[0].loop()
// 	cs = &connectionSet{
// 		masterConn: masterConn,
// 		slaveConns: slaveConns,
// 		asMutex:    new(sync.RWMutex),
// 		quitCh:     make(chan bool),
// 	}
// 	go cs.loop()

// 	time.Sleep(5 * time.Millisecond)
// 	Shutdown()

// 	time.Sleep(5 * time.Millisecond)
// 	assert.True(t, masterClosed, "fail to run closeFn on master")
// 	assert.True(t, slaveClosed, "fail to run closeFn connection on slave")
// 	assert.False(t, masterConn.connected, "fail to close master")
// 	assert.False(t, slaveConns[0].connected, "fail to close slave connections")
// }
