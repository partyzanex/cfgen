package config

import (
	"sort"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

type EnvName string

func (en EnvName) String() string {
	return string(en)
}

type Source struct {
	App   App   `yaml:"app"`
	Flags Flags `yaml:"flags"`
}

type App struct {
	Name string    `yaml:"name"`
	Desc string    `yaml:"desc"`
	Env  []EnvName `yaml:"env"`
}

type Flags []*Flag

func (flags *Flags) UnmarshalYAML(node *yaml.Node) error {
	if node.Kind != yaml.MappingNode {
		return errors.Errorf("unsupported node kind: %d", node.Kind)
	}

	results := make([]*Flag, 0, len(node.Content))

	var key string

	for _, content := range node.Content {
		if content.Kind == yaml.ScalarNode {
			key = content.Value
		}

		if content.Kind == yaml.MappingNode {
			flag := new(Flag)

			err := content.Decode(flag)
			if err != nil {
				return errors.Wrap(err, "Decode")
			}

			if flag.Name == "" {
				flag.Name = key
			}

			results = append(results, flag)
		}
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].Name < results[j].Name
	})

	*flags = results

	return nil
}
