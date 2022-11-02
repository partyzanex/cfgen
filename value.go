package config

import (
	"time"

	"github.com/urfave/cli/v2"
)

type Value struct {
	env  EnvName
	raw  map[EnvName]interface{}
	envs []EnvName
}

func (v *Value) setEnv(env EnvName) {
	if _, ok := v.raw[env]; !ok {
		v.envs = append(v.envs, env)
	}
}

func (v *Value) Set(env EnvName, value interface{}) *Value {
	v.setEnv(env)
	v.raw[env] = value

	return v
}

func (v *Value) SetStringSlice(env EnvName, values ...string) *Value {
	v.setEnv(env)
	v.raw[env] = cli.NewStringSlice(values...)

	return v
}

func (v *Value) SetIntSlice(env EnvName, values ...int) *Value {
	v.setEnv(env)
	v.raw[env] = cli.NewIntSlice(values...)

	return v
}

func (v *Value) SetInt64Slice(env EnvName, values ...int64) *Value {
	v.setEnv(env)
	v.raw[env] = cli.NewInt64Slice(values...)

	return v
}

func (v *Value) SetUIntSlice(env EnvName, values ...uint) *Value {
	v.setEnv(env)
	v.raw[env] = cli.NewUintSlice(values...)

	return v
}

func (v *Value) SetUInt64Slice(env EnvName, values ...uint64) *Value {
	v.setEnv(env)
	v.raw[env] = cli.NewUint64Slice(values...)

	return v
}

func (v *Value) SetFloat64Slice(env EnvName, values ...float64) *Value {
	v.setEnv(env)
	v.raw[env] = cli.NewFloat64Slice(values...)

	return v
}

func (v *Value) SetDuration(env EnvName, value time.Duration) *Value {
	v.setEnv(env)
	v.raw[env] = value

	return v
}

func (v *Value) SetTimestamp(env EnvName, value string) *Value {
	ts, err := time.Parse(time.RFC3339, value)
	if err != nil {
		panic(any(err))
	}

	v.setEnv(env)
	v.raw[env] = ts

	return v
}

func (v *Value) String() string {
	value := v.get()
	if value == nil {
		return ""
	}

	return value.(string)
}

func (v *Value) Bool() bool {
	value := v.get()
	if value == nil {
		return false
	}

	return value.(bool)
}

func (v *Value) Int() int {
	value := v.get()
	if value == nil {
		return 0
	}

	return value.(int)
}

func (v *Value) Int64() int64 {
	value := v.get()
	if value == nil {
		return 0
	}

	return value.(int64)
}

func (v *Value) Uint() uint {
	value := v.get()
	if value == nil {
		return 0
	}

	return value.(uint)
}

func (v *Value) Uint64() uint64 {
	value := v.get()
	if value == nil {
		return 0
	}

	return value.(uint64)
}

func (v *Value) Float64() float64 {
	value := v.get()
	if value == nil {
		return 0
	}

	return value.(float64)
}

func (v *Value) Duration() time.Duration {
	value := v.get()
	if value == nil {
		return 0
	}

	return value.(time.Duration)
}

func (v *Value) Timestamp() *cli.Timestamp {
	value := v.get()
	if value == nil {
		return nil
	}

	return cli.NewTimestamp(value.(time.Time))
}

func (v *Value) StringSlice() *cli.StringSlice {
	value := v.get()
	if value == nil {
		return nil
	}

	return value.(*cli.StringSlice)
}

func (v *Value) IntSlice() *cli.IntSlice {
	value := v.get()
	if value == nil {
		return nil
	}

	return value.(*cli.IntSlice)
}

func (v *Value) Int64Slice() *cli.Int64Slice {
	value := v.get()
	if value == nil {
		return nil
	}

	return value.(*cli.Int64Slice)
}

func (v *Value) UintSlice() *cli.UintSlice {
	value := v.get()
	if value == nil {
		return nil
	}

	return value.(*cli.UintSlice)
}

func (v *Value) Uint64Slice() *cli.Uint64Slice {
	value := v.get()
	if value == nil {
		return nil
	}

	return value.(*cli.Uint64Slice)
}

func (v *Value) Float64Slice() *cli.Float64Slice {
	value := v.get()
	if value == nil {
		return nil
	}

	return value.(*cli.Float64Slice)
}

func (v *Value) get() interface{} {
	value, ok := v.raw[v.env]
	if ok {
		return value
	}

	n := len(v.envs)
	if n == 0 {
		return nil
	}

	return v.raw[v.envs[n-1]]
}

func (v *Value) Env(env EnvName) *Value {
	newValue := *v
	newValue.env = env

	return &newValue
}

func NewValue(current EnvName) *Value {
	return &Value{env: current, raw: map[EnvName]interface{}{}}
}
