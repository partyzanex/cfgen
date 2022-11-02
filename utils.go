package config

import (
	"fmt"
	"os"

	"github.com/pkg/errors"
)

func toInt64(v interface{}) (int64, error) {
	switch i := v.(type) {
	case int:
		return int64(i), nil
	case int64:
		return i, nil
	default:
		return 0, errors.Errorf("unsupported int type %T[%v]", v, v)
	}
}

func GetEnvName(envKey string, environments ...EnvName) EnvName {
	if len(environments) == 0 {
		panic(any("required environments"))
	}

	if env := os.Getenv(envKey); env != "" {
		for _, environment := range environments {
			if env == environment.String() {
				return environment
			}
		}

		panic(any(fmt.Sprintf("invalid environment %q", env)))
	}

	return environments[0]
}
