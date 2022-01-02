package tools

// Ktech represents a ktools/ktech tool.
type Ktech struct {
	Ktools
}

// NewKtech creates a new Ktech instance.
func NewKtech() (*Ktech, error) {
	ktools, err := NewKtools("ktech", "ktech")
	if err != nil {
		return nil, err
	}
	return &Ktech{
		Ktools: *ktools,
	}, nil
}
