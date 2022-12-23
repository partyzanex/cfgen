package config

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestValue_Int(t *testing.T) {
	v := NewValue("test").
		Set("test", 1).
		Set("test2", 2)

	assert.Equal(t, 1, v.Int())
	assert.Equal(t, 2, v.Env("test2").Int())
	assert.Equal(t, 1, v.Env("unknown").Int())
}

func TestValue_Bool(t *testing.T) {
	v := NewValue("test").
		Set("test", true).
		Set("test2", false)

	assert.Equal(t, true, v.Bool())
	assert.Equal(t, false, v.Env("test2").Bool())
	assert.Equal(t, true, v.Env("unknown").Bool())
}

func TestValue_String(t *testing.T) {
	v := NewValue("test").
		Set("test", "123").
		Set("docker", "test").
		Set("local", "321")

	assert.Equal(t, "123", v.String())
	assert.Equal(t, "321", v.Env("local").String())
	assert.Equal(t, "test", v.Env("docker").String())
	assert.Equal(t, "123", v.Env("unknown").String())
}

func TestValue_Duration(t *testing.T) {
	v := NewValue("test").
		Set("test", time.Duration(1)).
		Set("local", time.Duration(1234)).
		Set("dev", time.Duration(1000000)).
		Set("prod", time.Duration(0))

	assert.Equal(t, time.Duration(1), v.Duration())
	assert.Equal(t, time.Duration(1234), v.Env("local").Duration())
	assert.Equal(t, time.Duration(1000000), v.Env("dev").Duration())
	assert.Equal(t, time.Duration(0), v.Env("prod").Duration())
}

func TestValue_Timestamp(t *testing.T) {
	v := NewValue("local").
		SetTimestamp("local", "2021-05-25T17:15:16Z").
		SetTimestamp("dev", "2022-05-25T17:15:16Z").
		SetTimestamp("prod", "2023-05-25T17:15:16Z")

	assert.Equal(t, "2021-05-25T17:15:16Z", v.Timestamp().Value().Format(time.RFC3339))
	assert.Equal(t, "2022-05-25T17:15:16Z", v.Env("dev").Timestamp().Value().Format(time.RFC3339))
	assert.Equal(t, "2023-05-25T17:15:16Z", v.Env("prod").Timestamp().Value().Format(time.RFC3339))
}
