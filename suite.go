package hansip

import (
	"os"

	"github.com/go-pg/pg"
	"github.com/stretchr/testify/suite"
)

// TestSuite custom test suite that provides working cluster instance
type TestSuite struct {
	suite.Suite
	cluster         *Cluster
	noCreateCluster bool
}

func (s *TestSuite) SetupTest() {
	if !s.noCreateCluster {
		s.setupCluster()
	}
}

func (s *TestSuite) TearDownTest() {
	if !s.noCreateCluster {
		s.cluster.Shutdown()
	}
}

func (s *TestSuite) setupCluster() {
	s.cluster = NewCluster(&Config{})
	s.Nil(s.cluster.SetMaster(s.getMasterConnectionInfo()))
	s.Nil(s.cluster.AddSlave(s.getSlave1ConnectionInfo()))
	s.Nil(s.cluster.AddSlave(s.getSlave2ConnectionInfo()))
}

func (s *TestSuite) getMasterConnectionInfo() *pg.Options {
	url := os.Getenv("DB_MASTER_URL")
	s.NotZero(url)

	opts, err := pg.ParseURL(url)
	s.Nil(err)

	return opts
}

func (s *TestSuite) getSlave1ConnectionInfo() *pg.Options {
	url := os.Getenv("DB_SLAVE1_URL")
	s.NotZero(url)

	opts, err := pg.ParseURL(url)
	s.Nil(err)

	return opts
}

func (s *TestSuite) getSlave2ConnectionInfo() *pg.Options {
	url := os.Getenv("DB_SLAVE2_URL")
	s.NotZero(url)

	opts, err := pg.ParseURL(url)
	s.Nil(err)

	return opts
}
