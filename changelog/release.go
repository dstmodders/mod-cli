package changelog

import (
	"errors"
	"regexp"
	"time"

	"github.com/Masterminds/semver"
	"github.com/yuin/goldmark/ast"
)

const regexDate string = `\d{4}-(0[1-9]|1[0-2])-(0[1-9]|[12][0-9]|3[01])`

const regexSemVer string = `v?([0-9]+)(\.[0-9]+)?(\.[0-9]+)?` +
	`(-([0-9A-Za-z\-]+(\.[0-9A-Za-z\-]+)*))?` +
	`(\+([0-9A-Za-z\-]+(\.[0-9A-Za-z\-]+)*))?`

var dateRegex *regexp.Regexp
var versionRegex *regexp.Regexp

// ReleaseController is the interface that wraps the Release methods.
type ReleaseController interface {
	AddAdded(string)
	AddChanged(string)
	AddDeprecated(string)
	AddRemoved(string)
	AddFixed(string)
	AddSecurity(string)
	CountChanges() int
	HasChanges() bool
	HasText() bool
	DateString() string
}

// Release represents a single CHANGELOG.md release.
type Release struct {
	// Title is an original release title. Usually, it's either "Unreleased" or,
	// as an example, "[1.0.0] - 2017-06-20".
	Title string

	// Date is a release date that is parsed from Title. For example,
	// "2017-06-20".
	Date *time.Time

	// Link is a release link which usually can be found in the footnotes.
	Link string

	// Version is a release version itself in a semantic versioning format:
	// https://semver.org/
	Version *semver.Version

	// Text is a release text that doesn't fit into changes. For example, the text
	// like "Initial release" in the first release.
	Text string

	// Added holds a list of all "Added" changes.
	Added []ReleaseChange

	// Changed holds a list of all "Changed" changes.
	Changed []ReleaseChange

	// Deprecated holds a list of all "Deprecated" changes.
	Deprecated []ReleaseChange

	// Removed holds a list of all "Removed" changes.
	Removed []ReleaseChange

	// Fixed holds a list of all "Fixed" changes.
	Fixed []ReleaseChange

	// Security holds a list of all "Security" changes.
	Security []ReleaseChange
}

func init() {
	dateRegex = regexp.MustCompile(regexDate)
	versionRegex = regexp.MustCompile(regexSemVer)
}

// NewRelease creates a new Release instance.
func NewRelease() *Release {
	return &Release{}
}

func (r *Release) versionFromString(str string) error {
	match := versionRegex.FindString(str)
	if len(match) > 0 {
		ver, err := semver.NewVersion(match)
		if err != nil {
			return err
		}
		r.Version = ver
		return nil
	}
	return errors.New("no semantic version found")
}

func (r *Release) dateFromString(str string) error {
	match := dateRegex.FindString(str)
	if len(match) > 0 {
		t, err := time.Parse("2006-01-02", match)
		if err != nil {
			return err
		}
		r.Date = &t
		return nil
	}
	return errors.New("no date found")
}

func (r *Release) fromGoldmarkHeadingNode(buf []byte, node ast.Node) error {
	if node.Kind() != ast.KindHeading {
		return errors.New("not a heading node")
	}

	headingNode := node.(*ast.Heading)
	if headingNode.Level != 2 {
		return errors.New("not a level 2 heading node")
	}

	r.Title = string(headingNode.Text(buf))

	if err := r.versionFromString(r.Title); err != nil {
		return err
	}

	if err := r.dateFromString(r.Title); err != nil {
		return err
	}

	return nil
}

// AddAdded adds a new "Added" change.
func (r *Release) AddAdded(desc string) {
	r.Added = append(r.Added, *NewReleaseChange(desc))
}

// AddChanged adds a new "Changed" change.
func (r *Release) AddChanged(desc string) {
	r.Changed = append(r.Changed, *NewReleaseChange(desc))
}

// AddDeprecated adds a new "Deprecated" change.
func (r *Release) AddDeprecated(desc string) {
	r.Deprecated = append(r.Deprecated, *NewReleaseChange(desc))
}

// AddRemoved adds a new "Removed" change.
func (r *Release) AddRemoved(desc string) {
	r.Removed = append(r.Removed, *NewReleaseChange(desc))
}

// AddFixed adds a new "Fixed" change.
func (r *Release) AddFixed(desc string) {
	r.Fixed = append(r.Fixed, *NewReleaseChange(desc))
}

// AddSecurity adds a new "Fixed" change.
func (r *Release) AddSecurity(desc string) {
	r.Security = append(r.Security, *NewReleaseChange(desc))
}

// CountChanges counts the total number of changes.
func (r *Release) CountChanges() int {
	return len(r.Added) +
		len(r.Changed) +
		len(r.Deprecated) +
		len(r.Removed) +
		len(r.Fixed) +
		len(r.Security)
}

// HasChanges checks if a release has any changes.
func (r *Release) HasChanges() bool {
	return r.CountChanges() > 0
}

// HasText checks if a release has any text.
func (r *Release) HasText() bool {
	return len(r.Text) > 0
}

// DateString returns a string representation of a Date.
func (r *Release) DateString() string {
	return r.Date.Format("2006-01-02")
}

// ReleaseChange represents a single release change. Created to be extended in
// the future to hold values in different formats like Plain Text, Markdown and
// Steam Workshop.
type ReleaseChange struct {
	Value string
}

// NewReleaseChange creates a new ReleaseChange instance with the provided
// value.
func NewReleaseChange(value string) *ReleaseChange {
	return &ReleaseChange{Value: value}
}
