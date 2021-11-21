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

type ReleaseController interface {
	VersionFromString(str string) error
	DateFromString(str string) error
	FromGoldmarkHeadingNode(node ast.Node)
	CountChanges() int
	HasChanges() bool
	HasText() bool
	AddAdded(desc string)
	AddChanged(desc string)
	AddDeprecated(desc string)
	AddRemoved(desc string)
	AddFixed(desc string)
	AddSecurity(desc string)
	DateString() string
}

type Release struct {
	Title      string
	Date       *time.Time
	Link       string
	Version    *semver.Version
	Text       string
	Added      []ReleaseChange
	Changed    []ReleaseChange
	Deprecated []ReleaseChange
	Removed    []ReleaseChange
	Fixed      []ReleaseChange
	Security   []ReleaseChange
}

func init() {
	dateRegex = regexp.MustCompile(regexDate)
	versionRegex = regexp.MustCompile(regexSemVer)
}

func NewRelease() *Release {
	return &Release{}
}

func (r *Release) VersionFromString(str string) error {
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

func (r *Release) DateFromString(str string) error {
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

func (r *Release) FromGoldmarkHeadingNode(buf []byte, node ast.Node) error {
	if node.Kind() != ast.KindHeading {
		return errors.New("not a heading node")
	}

	headingNode := node.(*ast.Heading)
	if headingNode.Level != 2 {
		return errors.New("not a level 2 heading node")
	}

	r.Title = string(headingNode.Text(buf))

	if err := r.VersionFromString(r.Title); err != nil {
		return err
	}

	if err := r.DateFromString(r.Title); err != nil {
		return err
	}

	return nil
}

func (r *Release) CountChanges() int {
	return len(r.Added) +
		len(r.Changed) +
		len(r.Deprecated) +
		len(r.Removed) +
		len(r.Fixed) +
		len(r.Security)
}

func (r *Release) HasChanges() bool {
	return r.CountChanges() > 0
}

func (r *Release) HasText() bool {
	return len(r.Text) > 0
}

func (r *Release) AddAdded(desc string) {
	r.Added = append(r.Added, *NewReleaseChange(desc))
}

func (r *Release) AddChanged(desc string) {
	r.Changed = append(r.Changed, *NewReleaseChange(desc))
}

func (r *Release) AddDeprecated(desc string) {
	r.Deprecated = append(r.Deprecated, *NewReleaseChange(desc))
}

func (r *Release) AddRemoved(desc string) {
	r.Removed = append(r.Removed, *NewReleaseChange(desc))
}

func (r *Release) AddFixed(desc string) {
	r.Fixed = append(r.Fixed, *NewReleaseChange(desc))
}

func (r *Release) AddSecurity(desc string) {
	r.Security = append(r.Security, *NewReleaseChange(desc))
}

func (r *Release) DateString() string {
	return r.Date.Format("2006-01-02")
}

type ReleaseChange struct {
	Value string
}

func NewReleaseChange(value string) *ReleaseChange {
	return &ReleaseChange{Value: value}
}
