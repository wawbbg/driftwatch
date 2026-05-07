package mask

// Stage is a function that transforms a config map.
type Stage func(map[string]string) map[string]string

// Pipeline applies a sequence of Stages to a config map in order.
// Each stage receives the output of the previous one.
type Pipeline struct {
	stages []Stage
}

// NewPipeline creates an empty Pipeline.
func NewPipeline() *Pipeline {
	return &Pipeline{}
}

// Add appends a Stage to the pipeline.
func (p *Pipeline) Add(s Stage) *Pipeline {
	p.stages = append(p.stages, s)
	return p
}

// Run executes all stages sequentially and returns the final map.
// If cfg is nil an empty map is returned.
func (p *Pipeline) Run(cfg map[string]string) map[string]string {
	if cfg == nil {
		cfg = map[string]string{}
	}
	out := cfg
	for _, s := range p.stages {
		out = s(out)
	}
	return out
}

// MaskStage returns a Stage that applies the given Masker.
func MaskStage(m *Masker) Stage {
	return m.Apply
}
