package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/suite"
)

func TestUtilTestSuite(t *testing.T) {
	suite.Run(t, new(UtilTestSuite))
}

type UtilTestSuite struct {
	suite.Suite
}

func (s *UtilTestSuite) TestMustGetInt() {
	key := "FOO"
	_ = os.Setenv(key, "4")
	v := mustGetInt(key)
	s.Equal(4, v)
	_ = os.Unsetenv(key)
}

func (s *UtilTestSuite) TestMustGetString() {
	key := "FOO"
	_ = os.Setenv(key, "Lorem")
	v := mustGetString(key)
	s.Equal("Lorem", v)
	_ = os.Unsetenv(key)
}
