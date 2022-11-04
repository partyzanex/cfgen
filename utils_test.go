package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetEnvName(t *testing.T) {
	env := GetEnvName("TEST_ENV", "one", "two", "three")
	assert.Equal(t, EnvName("one"), env)

	env = GetEnvName("TEST_ENV", "zero", "one", "two", "three")
	assert.Equal(t, EnvName("zero"), env)

	err := os.Setenv("TEST_ENV", "three")
	assert.NoError(t, err)

	env = GetEnvName("TEST_ENV", "one", "two", "three")
	assert.Equal(t, EnvName("three"), env)
}
