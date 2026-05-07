package mask_test

import (
	"strings"
	"testing"

	"github.com/yourorg/driftwatch/internal/mask"
)

func TestPipeline_Empty(t *testing.T) {
	p := mask.NewPipeline()
	input := map[string]string{"key": "value"}
	out := p.Run(input)
	if out["key"] != "value" {
		t.Errorf("expected value unchanged, got %q", out["key"])
	}
}

func TestPipeline_MaskStage(t *testing.T) {
	m := mask.New()
	p := mask.NewPipeline().Add(mask.MaskStage(m))
	input := map[string]string{"password": "secret", "host": "db"}
	out := p.Run(input)
	if out["password"] != mask.Placeholder {
		t.Errorf("expected password masked, got %q", out["password"])
	}
	if out["host"] != "db" {
		t.Errorf("expected host unchanged, got %q", out["host"])
	}
}

func TestPipeline_MultipleStages(t *testing.T) {
	upper := func(cfg map[string]string) map[string]string {
		out := make(map[string]string, len(cfg))
		for k, v := range cfg {
			out[k] = strings.ToUpper(v)
		}
		return out
	}
	m := mask.New()
	p := mask.NewPipeline().
		Add(mask.MaskStage(m)).
		Add(upper)
	input := map[string]string{"token": "abc", "env": "prod"}
	out := p.Run(input)
	if out["token"] != strings.ToUpper(mask.Placeholder) {
		t.Errorf("unexpected token value: %q", out["token"])
	}
	if out["env"] != "PROD" {
		t.Errorf("expected env=PROD, got %q", out["env"])
	}
}

func TestPipeline_NilInput(t *testing.T) {
	p := mask.NewPipeline()
	out := p.Run(nil)
	if out == nil {
		t.Error("expected non-nil map for nil input")
	}
}
