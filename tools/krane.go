package tools

// Krane represents a ktools/krane tool.
type Krane struct {
	Ktools
}

// NewKrane creates a new Krane instance.
func NewKrane() (*Krane, error) {
	ktools, err := NewKtools("krane", "krane")
	if err != nil {
		return nil, err
	}
	return &Krane{
		Ktools: *ktools,
	}, nil
}
