// Package ignore provides functionality to skip specific fields or services
// during drift detection based on user-defined rules.
package ignore

import (
	"os"

	"gopkg.in/yaml.v3"
)

// Rule defines a single ignore rule.
type Rule struct {
	Service string   `yaml:"service"`
	Fields  []string `yaml:"fields"`
}

// Rules is a collection of ignore rules.
type Rules struct {
	Ignore []Rule `yaml:"ignore"`
}

// LoadFromFile reads ignore rules from a YAML file at the given path.
func LoadFromFile(path string) (*Rules, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var rules Rules
	if err := yaml.Unmarshal(data, &rules); err != nil {
		return nil, err
	}
	return &rules, nil
}

// MatchesService returns true if the given service name matches a rule.
func (r *Rules) MatchesService(service string) bool {
	for _, rule := range r.Ignore {
		if rule.Service == "*" || rule.Service == service {
			return true
		}
	}
	return false
}

// IgnoredFields returns the set of fields to ignore for the given service.
func (r *Rules) IgnoredFields(service string) map[string]bool {
	fields := make(map[string]bool)
	for _, rule := range r.Ignore {
		if rule.Service == "*" || rule.Service == service {
			for _, f := range rule.Fields {
				fields[f] = true
			}
		}
	}
	return fields
}
