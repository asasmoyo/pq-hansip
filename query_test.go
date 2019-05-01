package hansip

// func TestQuery_MasterOKSlaveOK(t *testing.T) {
// 	masterConn := &connection{
// 		connected: true,
// 		host:      "master",
// 		quitCh:    make(chan bool),
// 		s:         &dummySQL{},
// 	}
// 	slaveConns := []*connection{
// 		&connection{
// 			connected: true,
// 			host:      "slave1",
// 			quitCh:    make(chan bool),
// 			s:         &dummySQL{},
// 		},
// 	}
// 	cs = &connectionSet{
// 		masterConn: masterConn,
// 		slaveConns: slaveConns,
// 		asMutex:    new(sync.RWMutex),
// 		quitCh:     make(chan bool),
// 	}
// 	cs.updateActiveSlaves()

// 	Query(nil, "")
// 	assert.True(t, slaveConns[0].s.(*dummySQL).queryRun, "fail to run query on slave")
// }

// func TestQuery_MasterOKSlaveNotOK(t *testing.T) {
// 	masterConn := &connection{
// 		connected: true,
// 		host:      "master",
// 		quitCh:    make(chan bool),
// 		s:         &dummySQL{},
// 	}
// 	slaveConns := []*connection{
// 		&connection{
// 			connected: false,
// 			host:      "slave1",
// 			quitCh:    make(chan bool),
// 			s:         &dummySQL{},
// 		},
// 	}
// 	cs = &connectionSet{
// 		masterConn: masterConn,
// 		slaveConns: slaveConns,
// 		asMutex:    new(sync.RWMutex),
// 		quitCh:     make(chan bool),
// 	}
// 	cs.updateActiveSlaves()

// 	Query(nil, "")
// 	assert.True(t, masterConn.s.(*dummySQL).queryRun, "fail to run query on master")
// }

// func TestWriterExec(t *testing.T) {
// 	masterConn := &connection{
// 		connected: true,
// 		host:      "master",
// 		quitCh:    make(chan bool),
// 		s:         &dummySQL{},
// 	}
// 	slaveConns := []*connection{
// 		&connection{
// 			connected: true,
// 			host:      "slave1",
// 			quitCh:    make(chan bool),
// 			s:         &dummySQL{},
// 		},
// 	}
// 	cs = &connectionSet{
// 		masterConn: masterConn,
// 		slaveConns: slaveConns,
// 		asMutex:    new(sync.RWMutex),
// 		quitCh:     make(chan bool),
// 	}
// 	cs.updateActiveSlaves()

// 	WriterExec("")
// 	assert.True(t, masterConn.s.(*dummySQL).execRun, "fail to run exec on master")
// }

// func TestWriterQuery(t *testing.T) {
// 	masterConn := &connection{
// 		connected: true,
// 		host:      "master",
// 		quitCh:    make(chan bool),
// 		s:         &dummySQL{},
// 	}
// 	slaveConns := []*connection{
// 		&connection{
// 			connected: true,
// 			host:      "slave1",
// 			quitCh:    make(chan bool),
// 			s:         &dummySQL{},
// 		},
// 	}
// 	cs = &connectionSet{
// 		masterConn: masterConn,
// 		slaveConns: slaveConns,
// 		asMutex:    new(sync.RWMutex),
// 		quitCh:     make(chan bool),
// 	}
// 	cs.updateActiveSlaves()

// 	WriterQuery(nil, "")
// 	assert.True(t, masterConn.s.(*dummySQL).queryRun, "fail to run query on master")
// }

// func TestNewTransaction(t *testing.T) {
// 	masterConn := &connection{
// 		connected: true,
// 		host:      "master",
// 		quitCh:    make(chan bool),
// 		s:         &dummySQL{},
// 	}
// 	slaveConns := []*connection{
// 		&connection{
// 			connected: true,
// 			host:      "slave1",
// 			quitCh:    make(chan bool),
// 			s:         &dummySQL{},
// 		},
// 	}
// 	cs = &connectionSet{
// 		masterConn: masterConn,
// 		slaveConns: slaveConns,
// 		asMutex:    new(sync.RWMutex),
// 		quitCh:     make(chan bool),
// 	}
// 	cs.updateActiveSlaves()

// 	NewTransaction()
// 	assert.True(t, masterConn.s.(*dummySQL).newTransactionRun, "fail to create a new transaction on master")
// }
