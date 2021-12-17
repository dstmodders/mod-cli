package tools

// Lint represents a linting result.
type Lint struct {
	Files []LintFile
}

// LintFile represents a single linting file.
type LintFile struct {
	// Path holds a file path.
	Path string

	// Issues holds the number of found issues.
	Issues int
}

// NewLint creates a new Lint instance.
func NewLint() *Lint {
	return &Lint{}
}

// NewLintFile creates a new LintFile instance.
func NewLintFile() *LintFile {
	return &LintFile{}
}
