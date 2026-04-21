package policy

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// LoadFromFile reads a Policy from a JSON or YAML file at the given path.
// Only JSON is natively decoded; YAML support requires an external dependency,
// so this loader accepts .json files and returns an error for unsupported types.
func LoadFromFile(path string) (*Policy, error) {
	if path == "" {
		return &Policy{}, nil
	}

	ext := strings.ToLower(filepath.Ext(path))
	if ext != ".json" {
		return nil, fmt.Errorf("policy: unsupported file format %q (only .json supported)", ext)
	}

	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("policy: open %s: %w", path, err)
	}
	defer f.Close()

	var p Policy
	if err := json.NewDecoder(f).Decode(&p); err != nil {
		return nil, fmt.Errorf("policy: decode %s: %w", path, err)
	}

	if err := validate(&p); err != nil {
		return nil, err
	}

	return &p, nil
}

// validate checks that each rule has a non-empty field and a recognised level.
func validate(p *Policy) error {
	for i, r := range p.Rules {
		if strings.TrimSpace(r.Field) == "" {
			return fmt.Errorf("policy: rule[%d] has empty field", i)
		}
		switch r.Level {
		case LevelWarn, LevelError, LevelNone:
		default:
			return fmt.Errorf("policy: rule[%d] has unknown level %q", i, r.Level)
		}
	}
	return nil
}
