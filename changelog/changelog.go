// Package changelog has been designed to parse CHANGELOG.md and give access to
// the information about existing releases. However, in the future, it may
// support manipulating and generating changelogs as well if that kind of
// behaviour will be needed in our CLI tools.
//
// It's based on the "Keep a Changelog" v1.0.0 specification:
// https://keepachangelog.com/en/1.0.0/
package changelog

import "github.com/yuin/goldmark/ast"

// Controller is the interface that wraps the Changelog methods.
type Controller interface {
	AddRelease(Release)
	HasReleases() bool
	FirstRelease() *Release
	LatestRelease() *Release
	FromGoldmarkNode([]byte, ast.Node) error
}

// Changelog represents the changelog itself.
type Changelog struct {
	// Releases holds a list of all releases.
	Releases []Release
}

// New creates a new Changelog instance.
func New() *Changelog {
	return &Changelog{}
}

// AddRelease adds a new release.
func (c *Changelog) AddRelease(r Release) {
	c.Releases = append(c.Releases, r)
}

// HasReleases checks if there are any releases.
func (c *Changelog) HasReleases() bool {
	return len(c.Releases) > 0
}

// FirstRelease returns the first Release which usually points to the "Initial
// release".
func (c *Changelog) FirstRelease() *Release {
	l := len(c.Releases)
	if l > 0 {
		return &c.Releases[l-1]
	}
	return nil
}

// LatestRelease returns the latest Release which usually points to the
// "Unreleased" one.
func (c *Changelog) LatestRelease() *Release {
	l := len(c.Releases)
	if l > 0 {
		return &c.Releases[0]
	}
	return nil
}

func (c *Changelog) FromGoldmarkNode(source []byte, node ast.Node) error { //nolint:funlen,gocyclo
	var release *Release
	var changes []ReleaseChange
	var changesType string

	err := ast.Walk(node, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		switch n.Kind() {
		case ast.KindHeading:
			block := n.(*ast.Heading)
			switch block.Level {
			case 2:
				if !entering {
					return ast.WalkContinue, nil
				}

				if release != nil {
					c.AddRelease(*release)
				}

				release = NewRelease()
				if err := release.fromGoldmarkHeadingNode(source, block); err != nil {
					return ast.WalkContinue, nil
				}
			case 3:
				if !entering {
					return ast.WalkContinue, nil
				}

				blockStr := block.Text(source)
				changesType = string(blockStr)
			}
		case ast.KindLink:
			if !entering || release == nil {
				return ast.WalkContinue, nil
			}

			block := n.(*ast.Link)
			if len(block.Destination) > 0 {
				release.Link = string(block.Destination)
			}
		case ast.KindList:
			if !entering {
				switch changesType {
				case "Added":
					release.Added = changes
				case "Changed":
					release.Changed = changes
				case "Deprecated":
					release.Deprecated = changes
				case "Removed":
					release.Removed = changes
				case "Fixed":
					release.Fixed = changes
				case "Security":
					release.Security = changes
				}

				changes = []ReleaseChange{}
				return ast.WalkContinue, nil
			}
			return ast.WalkContinue, nil
		case ast.KindListItem:
			if !entering {
				return ast.WalkContinue, nil
			}

			block := n.(*ast.ListItem)
			changes = append(changes, *NewReleaseChange(string(block.Text(source))))
		case ast.KindParagraph:
			if !entering {
				return ast.WalkContinue, nil
			}

			block := n.(*ast.Paragraph)
			ps := block.PreviousSibling()

			if release != nil && ps.Kind() == ast.KindHeading {
				psBlock := ps.(*ast.Heading)
				if psBlock.Level == 2 {
					release.Text = string(block.Text(source))
					c.AddRelease(*release)
				}
			}
		}

		return ast.WalkContinue, nil
	})

	return err
}
