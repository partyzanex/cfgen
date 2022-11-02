package config

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/iancoleman/strcase"
)

type Flag struct {
	Name     string      `yaml:"flag"`
	Type     FlagType    `yaml:"type"`
	Enum     []string    `yaml:"enum"`
	Desc     string      `yaml:"desc"`
	Required bool        `yaml:"required"`
	Aliases  []string    `yaml:"aliases"`
	Env      interface{} `yaml:"env"`
	Value    interface{} `yaml:"value"`
}

func (flag *Flag) Args(env string) string {
	switch flag.Type {
	case FlagTypeString:
		return flag.stringArg(env)
	case FlagTypeTimestamp:
		return flag.timestampArg(env)
	case FlagTypeInt, FlagTypeInt64, FlagTypeUInt, FlagTypeUInt64:
		return flag.intArg(env)
	case FlagTypeEnum:
		return flag.enumArg(env)
	case FlagTypeFloat64:
		return flag.floatArg(env)
	case FlagTypeBool:
		return flag.boolArg(env)
	case FlagTypeDuration:
		return flag.durationArg(env)
	case FlagTypeFloat64Slice,
		FlagTypeStringSlice,
		FlagTypeIntSlice,
		FlagTypeInt64Slice,
		FlagTypeUIntSlice,
		FlagTypeUInt64Slice:
		return flag.sliceArg(env)
	default:
		panic(flag.errorf("Args: unsupported flag type"))
	}
}

const (
	enumMethodSet             = "Set"
	enumMethodSetDuration     = "SetDuration"
	enumMethodSetTimestamp    = "SetTimestamp"
	enumMethodSetStringSlice  = "SetStringSlice"
	enumMethodSetIntSlice     = "SetIntSlice"
	enumMethodSetInt64Slice   = "SetInt64Slice"
	enumMethodSetUIntSlice    = "SetUIntSlice"
	enumMethodSetUInt64Slice  = "SetUInt64Slice"
	enumMethodSetFloat64Slice = "SetFloat64Slice"
)

func (flag *Flag) ValueSetMethodName() string {
	switch flag.Type {
	case FlagTypeString,
		FlagTypeBool,
		FlagTypeEnum,
		FlagTypeInt,
		FlagTypeInt64,
		FlagTypeUInt,
		FlagTypeUInt64,
		FlagTypeFloat64:
		return enumMethodSet
	case FlagTypeDuration:
		return enumMethodSetDuration
	case FlagTypeTimestamp:
		return enumMethodSetTimestamp
	case FlagTypeStringSlice:
		return enumMethodSetStringSlice
	case FlagTypeIntSlice:
		return enumMethodSetIntSlice
	case FlagTypeInt64Slice:
		return enumMethodSetInt64Slice
	case FlagTypeUIntSlice:
		return enumMethodSetUIntSlice
	case FlagTypeUInt64Slice:
		return enumMethodSetUInt64Slice
	case FlagTypeFloat64Slice:
		return enumMethodSetFloat64Slice
	default:
		panic(flag.errorf("ValueSetMethodName: unknown flag type %q", flag.Type))
	}
}

func (flag *Flag) CLIFlagType() string {
	return flag.ValueType() + "Flag"
}

func (flag *Flag) ValueType() string {
	switch flag.Type {
	case FlagTypeString, FlagTypeEnum:
		return "String"
	case FlagTypeInt:
		return "Int"
	case FlagTypeInt64:
		return "Int64"
	case FlagTypeUInt:
		return "Uint"
	case FlagTypeUInt64:
		return "Uint64"
	case FlagTypeFloat64:
		return "Float64"
	case FlagTypeBool:
		return "Bool"
	case FlagTypeTimestamp:
		return "Timestamp"
	case FlagTypeDuration:
		return "Duration"
	case FlagTypeStringSlice:
		return "StringSlice"
	case FlagTypeIntSlice:
		return "IntSlice"
	case FlagTypeInt64Slice:
		return "Int64Slice"
	case FlagTypeUIntSlice:
		return "UIntSlice"
	case FlagTypeUInt64Slice:
		return "UInt64Slice"
	case FlagTypeFloat64Slice:
		return "Float64Slice"
	default:
		return ""
	}
}

func (flag *Flag) IsSlice() bool {
	switch flag.Type {
	case FlagTypeStringSlice,
		FlagTypeIntSlice,
		FlagTypeInt64Slice,
		FlagTypeUIntSlice,
		FlagTypeUInt64Slice,
		FlagTypeFloat64Slice:
		return true
	default:
		return false
	}
}

func (flag *Flag) RequiredField() string {
	return strconv.FormatBool(flag.Required)
}

func (flag *Flag) EnvVarsField(appName string) string {
	prefix := strcase.ToScreamingSnake(appName) + "_"

	names := make([]string, 1)
	names[0] = strconv.Quote(prefix + strcase.ToScreamingSnake(flag.Name))

	switch env := flag.Env.(type) {
	case []interface{}:
		for _, e := range env {
			s, ok := e.(string)
			if !ok {
				panic(flag.errorf("EnvVarsField: unsupported env value type %T", e))
			}

			names = append(names, strconv.Quote(strcase.ToScreamingSnake(s)))
		}
	case string:
		names = append(names, strconv.Quote(strcase.ToScreamingSnake(env)))
	case bool:
		return "nil"
	case nil:
		// nothing
	default:
		panic(flag.errorf("EnvVarsField: unknown env type %T", flag.Env))
	}

	return fmt.Sprintf("[]string{%s}", strings.Join(names, ", "))
}

func (flag *Flag) GoType() string {
	switch flag.Type {
	case FlagTypeString, FlagTypeEnum:
		return "string"
	case FlagTypeInt:
		return "int"
	case FlagTypeInt64:
		return "int64"
	case FlagTypeUInt:
		return "uint"
	case FlagTypeUInt64:
		return "uint64"
	case FlagTypeFloat64:
		return "float64"
	case FlagTypeBool:
		return "bool"
	case FlagTypeDuration:
		return "time.Duration"
	case FlagTypeStringSlice:
		return "[]string"
	case FlagTypeIntSlice:
		return "[]int"
	case FlagTypeInt64Slice:
		return "[]int64"
	case FlagTypeUIntSlice:
		return "[]uint"
	case FlagTypeUInt64Slice:
		return "[]uint64"
	case FlagTypeFloat64Slice:
		return "[]float64"
	default:
		panic(any("unsupported flag type"))
	}
}

func (flag *Flag) sliceArg(env string) string {
	var ss []interface{}

	switch v := flag.Value.(type) {
	case map[string]interface{}:
		var ok bool

		ss, ok = v[env].([]interface{})
		if !ok {
			return ``
		}
	case []interface{}:
		ss = v
	default:
		panic(flag.errorf("sliceArg: unsupported type %T", flag.Value))
	}

	r := make([]string, len(ss))

	for i, s := range ss {
		var n string

		switch v := s.(type) {
		case int:
			n = strconv.Itoa(v)
		case int64:
			n = strconv.FormatInt(v, 10)
		case uint:
			n = strconv.FormatUint(uint64(v), 10)
		case uint64:
			n = strconv.FormatUint(v, 10)
		case float64:
			n = strconv.FormatFloat(v, 'f', -1, 64)
		case string:
			n = `"` + s.(string) + `"`
		default:
			panic(flag.errorf("sliceArg: unsupported value type %T", s))
		}

		r[i] = n
	}

	return strings.Join(r, ",")
}

func (flag *Flag) durationArg(env string) string {
	var d int64

	switch v := flag.Value.(type) {
	case string:
		dur, err := time.ParseDuration(v)
		if err != nil {
			panic(flag.errorf("durationArg: %s", err))
		}

		d = dur.Nanoseconds()
	case map[string]interface{}:
		str, ok := v[env].(string)
		if ok {
			dur, err := time.ParseDuration(str)
			if err != nil {
				panic(flag.errorf("durationArg: %s", err))
			}

			d = dur.Nanoseconds()
		}
	default:
		panic(flag.errorf("durationArg: unsupported type %T", flag.Value))
	}

	return fmt.Sprintf("time.Duration(%d)", d)
}

func (flag *Flag) boolArg(env string) string {
	var b bool

	switch v := flag.Value.(type) {
	case bool:
		b = v
	case map[string]interface{}:
		a, ok := v[env].(bool)
		if ok {
			b = a
		}
	default:
		panic(any(fmt.Sprintf("unsupported type %T", flag.Value)))
	}

	return strconv.FormatBool(b)
}

func (flag *Flag) floatArg(env string) string {
	var f float64

	switch v := flag.Value.(type) {
	case float64:
		f = v
	case map[string]interface{}:
		d, ok := v[env].(float64)
		if ok {
			f = d
		}
	default:
		panic(flag.errorf("floatArg: unsupported type %T"))
	}

	return strconv.FormatFloat(f, 'f', -1, 64)
}

func (flag *Flag) enumArg(env string) string {
	var enum string

	switch v := flag.Value.(type) {
	case string:
		enum = strcase.ToCamel(flag.Name) + strcase.ToCamel(v)
	case map[string]interface{}:
		e, ok := v[env].(string)
		if ok {
			enum = strcase.ToCamel(flag.Name) + strcase.ToCamel(e)
		} else {
			panic(flag.errorf("undefined value for env %q", env))
		}
	default:
		panic(flag.errorf("unsupported type %T", flag.Value))
	}

	return enum
}

func (flag *Flag) intArg(env string) string {
	var (
		i   int64
		err error
	)

	switch flag.Value.(type) {
	case int, int64:
		i, err = toInt64(flag.Value)
		if err != nil {
			panic(flag.errorf("toInt64: %s", err))
		}
	case map[string]interface{}:
		m := flag.Value.(map[string]interface{})

		v, ok := m[env]
		if ok {
			i, err = toInt64(v)
			if err != nil {
				panic(flag.errorf("toInt64: %s", err))
			}
		}
	default:
		panic(flag.errorf("intArg: unsupported type %T", flag.Value))
	}

	return strconv.FormatInt(i, 10)
}

func (flag *Flag) timestampArg(env string) string {
	var (
		dt  time.Time
		err error
	)

	switch v := flag.Value.(type) {
	case string:
		dt, err = time.Parse(time.RFC3339, v)
		if err != nil {
			panic(flag.errorf("timestampArg: %s, value: %v", err, v))
		}
	case map[string]interface{}:
		s, ok := v[env].(string)
		if ok {
			dt, err = time.Parse(time.RFC3339, s)
			if err != nil {
				panic(flag.errorf("timestampArg: %s", err))
			}
		}
	default:
		panic(flag.errorf("timestampArg: unsupported type %T", flag.Value))
	}

	return strconv.Quote(dt.Format(time.RFC3339))
}

func (flag *Flag) stringArg(env string) string {
	var arg string

	switch v := flag.Value.(type) {
	case string:
		arg = v
	case int, int64, uint, uint64:
		arg = fmt.Sprintf(`%d`, v)
	case float64:
		arg = fmt.Sprintf("%f", v)
	case bool:
		arg = strconv.FormatBool(v)
	case map[string]interface{}:
		s, ok := v[env].(string)
		if ok {
			arg = s
		}
	default:
		panic(flag.errorf("stringArg: unsupported value type %T", flag.Value))
	}

	return strconv.Quote(arg)
}

func (flag *Flag) errorf(format string, args ...interface{}) any {
	msg := fmt.Sprintf(format, args...)
	msg += fmt.Sprintf(" [flag=%s type=%s]", flag.Name, flag.Type)

	return msg
}

type FlagType string

const (
	FlagTypeString       FlagType = "string"
	FlagTypeStringSlice  FlagType = "stringSlice"
	FlagTypeEnum         FlagType = "enum"
	FlagTypeBool         FlagType = "bool"
	FlagTypeInt          FlagType = "int"
	FlagTypeUInt         FlagType = "uint"
	FlagTypeInt64        FlagType = "int64"
	FlagTypeUInt64       FlagType = "uint64"
	FlagTypeIntSlice     FlagType = "intSlice"
	FlagTypeUIntSlice    FlagType = "uintSlice"
	FlagTypeInt64Slice   FlagType = "int64Slice"
	FlagTypeUInt64Slice  FlagType = "uint64Slice"
	FlagTypeFloat64      FlagType = "float64"
	FlagTypeFloat64Slice FlagType = "float64Slice"
	FlagTypeDuration     FlagType = "duration"
	FlagTypeTimestamp    FlagType = "timestamp"
)

func (ft FlagType) String() string {
	return string(ft)
}
