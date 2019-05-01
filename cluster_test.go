package hansip

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type ClusterTestSuite struct {
	TestSuite
}

func TestCluster(t *testing.T) {
	s := &ClusterTestSuite{}
	suite.Run(t, s)
}

func (s *ClusterTestSuite) TestCreatesMasterAndSlaveConnection() {
	s.Nil(s.cluster.WriterExec("select 1;"))

	var temp int
	s.Nil(s.cluster.WriterQuery(&temp, "select 1;"))
	s.Equal(temp, 1)

	s.Nil(s.cluster.Query(&temp, "select 2;"))
	s.Equal(temp, 2)
}

func (s *ClusterTestSuite) TestKillConnectionsAfterShutdown() {
	s.cluster.Shutdown()
	s.Nil(s.cluster.getConnectionManager().writer())
	s.Nil(s.cluster.getConnectionManager().reader())
}

func (s *ClusterTestSuite) TestUseMasterWhenNoSlaveAvailable() {
	s.cluster.getConnectionManager().slaves = nil
	s.cluster.getConnectionManager().updateActiveSlaves()

	var temp int
	s.Nil(s.cluster.Query(&temp, "select 2;"))
	s.Equal(temp, 2)
}
