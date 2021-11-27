package main

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/dstmodders/mod-cli/changelog"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/text"
)

type Changelog struct {
	Changelog    *changelog.Changelog
	Count        bool
	First        bool
	Latest       bool
	List         bool
	ListVersions bool
	Node         ast.Node
	Source       []byte
}

func NewChangelog() *Changelog {
	return &Changelog{}
}

func (c *Changelog) NewChangelog() *Changelog {
	return &Changelog{}
}

func (c *Changelog) Load(path string) error {
	src, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	c.Source = src

	md := goldmark.New(goldmark.WithExtensions(extension.GFM))
	r := text.NewReader(src)
	node := md.Parser().Parse(r)

	c.Changelog = changelog.New()
	c.Node = node

	if err := c.Changelog.FromGoldmarkNode(src, node); err != nil {
		return err
	}

	return nil
}

func (c *Changelog) Print() error {
	if c.Changelog == nil {
		return errors.New("not loaded")
	}

	l := len(c.Changelog.Releases)
	if l == 0 {
		return errors.New("no releases")
	}

	if c.Count {
		fmt.Println(l)
		return nil
	}

	if c.List {
		for _, release := range c.Changelog.Releases {
			c.PrintTitle(release, false)
		}
		return nil
	}

	if c.ListVersions {
		for _, release := range c.Changelog.Releases {
			if release.Version != nil {
				fmt.Println(release.Version.String())
			} else {
				fmt.Println(strings.ToUpper(release.Title))
			}
		}
		return nil
	}

	if c.Latest && c.First {
		c.PrintRelease(*c.Changelog.LatestRelease())
		fmt.Println()
		c.PrintRelease(*c.Changelog.FirstRelease())
		return nil
	}

	if c.Latest {
		c.PrintRelease(*c.Changelog.LatestRelease())
		return nil
	}

	if c.First {
		c.PrintRelease(*c.Changelog.FirstRelease())
		return nil
	}

	for i, release := range c.Changelog.Releases {
		c.PrintRelease(release)
		if i != l-1 {
			fmt.Println()
		}
	}

	return nil
}

func (c *Changelog) PrintRelease(release changelog.Release) {
	c.PrintTitle(release, true)

	if !release.HasChanges() && release.HasText() {
		fmt.Printf("\n%s\n", release.Text)
		return
	}

	if len(release.Added) > 0 {
		c.PrintType("Added")
		c.PrintList(release.Added)
	}

	if len(release.Changed) > 0 {
		c.PrintType("Changed")
		c.PrintList(release.Changed)
	}

	if len(release.Deprecated) > 0 {
		c.PrintType("Deprecated")
		c.PrintList(release.Deprecated)
	}

	if len(release.Removed) > 0 {
		c.PrintType("Removed")
		c.PrintList(release.Removed)
	}

	if len(release.Fixed) > 0 {
		c.PrintType("Fixed")
		c.PrintList(release.Fixed)
	}

	if len(release.Security) > 0 {
		c.PrintType("Security")
		c.PrintList(release.Security)
	}
}

func (c *Changelog) PrintTitle(release changelog.Release, brackets bool) {
	title := release.Title
	if release.Version != nil {
		title = release.Version.String()
	}

	l, r := "", ""
	if brackets {
		l, r = "[", "]"
	}

	fmt.Printf("%s", l)

	if release.Date != nil && len(release.Link) > 0 {
		fmt.Printf("%s | %s | %s", title, release.DateString(), release.Link)
	} else if release.Date != nil {
		fmt.Printf("%s | %s", title, release.DateString())
	} else {
		fmt.Printf("%s", strings.ToUpper(title))
	}

	fmt.Printf("%s\n", r)
}

func (c *Changelog) PrintType(str string) {
	fmt.Printf("\n%s\n\n", strings.ToUpper(str))
}

func (c *Changelog) PrintList(list []changelog.ReleaseChange) {
	for _, change := range list {
		fmt.Printf("- %s\n", change.Value)
	}
}
