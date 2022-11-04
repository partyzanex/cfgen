package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValue_Bool(t *testing.T) {
	v := NewValue("test").
		Set("test", true).
		Set("test2", false)

	assert.Equal(t, true, v.Bool())
	assert.Equal(t, false, v.Env("test2").Bool())
	assert.Equal(t, false, v.Env("unknown").Bool())
}

func TestValue_String(t *testing.T) {
	v := NewValue("test").
		Set("test", "123").
		Set("docker", "test").
		Set("local", "321")

	assert.Equal(t, "123", v.String())
	assert.Equal(t, "321", v.Env("local").String())
	assert.Equal(t, "test", v.Env("docker").String())
	assert.Equal(t, "321", v.Env("unknown").String())
}
