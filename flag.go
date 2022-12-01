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

const (
	enumValueTypeString       = "String"
	enumValueTypeInt          = "Int"
	enumValueTypeInt64        = "Int64"
	enumValueTypeUint         = "Uint"
	enumValueTypeUint64       = "Uint64"
	enumValueTypeFloat64      = "Float64"
	enumValueTypeBool         = "Bool"
	enumValueTypeTimestamp    = "Timestamp"
	enumValueTypeDuration     = "Duration"
	enumValueTypeStringSlice  = "StringSlice"
	enumValueTypeIntSlice     = "IntSlice"
	enumValueTypeInt64Slice   = "Int64Slice"
	enumValueTypeUintSlice    = "UintSlice"
	enumValueTypeUint64Slice  = "Uint64Slice"
	enumValueTypeFloat64Slice = "Float64Slice"
)

func (flag *Flag) ValueType() string {
	switch flag.Type {
	case FlagTypeString, FlagTypeEnum:
		return enumValueTypeString
	case FlagTypeInt:
		return enumValueTypeInt
	case FlagTypeInt64:
		return enumValueTypeInt64
	case FlagTypeUInt:
		return enumValueTypeUint
	case FlagTypeUInt64:
		return enumValueTypeUint64
	case FlagTypeFloat64:
		return enumValueTypeFloat64
	case FlagTypeBool:
		return enumValueTypeBool
	case FlagTypeTimestamp:
		return enumValueTypeTimestamp
	case FlagTypeDuration:
		return enumValueTypeDuration
	case FlagTypeStringSlice:
		return enumValueTypeStringSlice
	case FlagTypeIntSlice:
		return enumValueTypeIntSlice
	case FlagTypeInt64Slice:
		return enumValueTypeInt64Slice
	case FlagTypeUIntSlice:
		return enumValueTypeUintSlice
	case FlagTypeUInt64Slice:
		return enumValueTypeUint64Slice
	case FlagTypeFloat64Slice:
		return enumValueTypeFloat64Slice
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
		return nilStr
	case nil:
		// nothing
	default:
		panic(flag.errorf("EnvVarsField: unknown env type %T", flag.Env))
	}

	return fmt.Sprintf("[]string{%s}", strings.Join(names, ", "))
}

const (
	enumGoTypeString       = "string"
	enumGoTypeInt          = "int"
	enumGoTypeInt64        = "int64"
	enumGoTypeUint         = "uint"
	enumGoTypeUint64       = "uint64"
	enumGoTypeFloat64      = "float64"
	enumGoTypeBool         = "bool"
	enumGoTypeDuration     = "time.Duration"
	enumGoTypeStringSlice  = "[]string"
	enumGoTypeIntSlice     = "[]int"
	enumGoTypeInt64Slice   = "[]int64"
	enumGoTypeUintSlice    = "[]uint"
	enumGoTypeUint64Slice  = "[]uint64"
	enumGoTypeFloat64Slice = "[]float64"
)

func (flag *Flag) GoType() string {
	switch flag.Type {
	case FlagTypeString, FlagTypeEnum:
		return enumGoTypeString
	case FlagTypeInt:
		return enumGoTypeInt
	case FlagTypeInt64:
		return enumGoTypeInt64
	case FlagTypeUInt:
		return enumGoTypeUint
	case FlagTypeUInt64:
		return enumGoTypeUint64
	case FlagTypeFloat64:
		return enumGoTypeFloat64
	case FlagTypeBool:
		return enumGoTypeBool
	case FlagTypeDuration:
		return enumGoTypeDuration
	case FlagTypeStringSlice:
		return enumGoTypeStringSlice
	case FlagTypeIntSlice:
		return enumGoTypeIntSlice
	case FlagTypeInt64Slice:
		return enumGoTypeInt64Slice
	case FlagTypeUIntSlice:
		return enumGoTypeUintSlice
	case FlagTypeUInt64Slice:
		return enumGoTypeUint64Slice
	case FlagTypeFloat64Slice:
		return enumGoTypeFloat64Slice
	default:
		panic(flag.errorf("unsupported flag type %q", flag.Type))
	}
}

const nilStr = "nil"

func (flag *Flag) AliasesField() string {
	if len(flag.Aliases) == 0 {
		return nilStr
	}

	return fmt.Sprintf(`[]string{"%s"}`, strings.Join(flag.Aliases, `", "`))
}

func (flag *Flag) DescField() string {
	if flag.Type == FlagTypeEnum {

		if flag.Desc != "" {
			return fmt.Sprintf("%s, (variants: %s)",
				flag.Desc, strings.Join(flag.Enum, ", "),
			)
		}

		return fmt.Sprintf("variants: %s", strings.Join(flag.Enum, ", "))
	}

	return flag.Desc
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
	case nil:
		// nothing
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
			n = v
		default:
			panic(flag.errorf("sliceArg: unsupported value type %T", s))
		}

		if flag.Type == FlagTypeStringSlice {
			n = strconv.Quote(n)
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
	case nil:
		// nothing
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
	case nil:
		// nothing
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
	case nil:
		// nothing
	default:
		panic(flag.errorf("floatArg: unsupported type %T", flag.Value))
	}

	return fmt.Sprintf("float64(%s)", strconv.FormatFloat(f, 'f', -1, 64))
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
	case nil:
		// nothing
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
	case nil:
		// nothing
	default:
		panic(flag.errorf("intArg: unsupported type %T", flag.Value))
	}

	return fmt.Sprintf("%s(%s)", flag.GoType(), strconv.FormatInt(i, 10))
}

func (flag *Flag) timestampArg(env string) string {
	var (
		dt  time.Time
		err error
	)

	switch v := flag.Value.(type) {
	case time.Time:
		dt = v
	case string:
		dt, err = time.Parse(time.RFC3339, v)
		if err != nil {
			panic(flag.errorf("timestampArg: %s, value: %v", err, v))
		}
	case map[string]interface{}:
		s, ok := v[env]
		if !ok {
			break
		}

		switch d := s.(type) {
		case string:
			dt, err = time.Parse(time.RFC3339, d)
			if err != nil {
				panic(flag.errorf("timestampArg: %s", err))
			}
		case time.Time:
			dt = d
		case nil:
			// nothing
		}
	case nil:
		// nothing
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
	case nil:
		// nothing
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
