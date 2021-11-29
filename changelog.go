package main

import (
	"errors"
	"fmt"
	"strings"

	"github.com/dstmodders/mod-cli/changelog"
)

type Changelog struct {
	Changelog    *changelog.Changelog
	Count        bool
	First        bool
	Latest       bool
	List         bool
	ListVersions bool
}

func NewChangelog() *Changelog {
	return &Changelog{}
}

func (c *Changelog) printTitle(release changelog.Release, brackets bool) {
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

func (c *Changelog) printType(str string) {
	fmt.Printf("\n%s\n\n", strings.ToUpper(str))
}

func (c *Changelog) printList(list []changelog.ReleaseChange) {
	for _, change := range list {
		fmt.Printf("- %s\n", change.Value)
	}
}

func (c *Changelog) printRelease(release changelog.Release) {
	c.printTitle(release, true)

	if !release.HasChanges() && release.HasText() {
		fmt.Printf("\n%s\n", release.Text)
		return
	}

	if len(release.Added) > 0 {
		c.printType("Added")
		c.printList(release.Added)
	}

	if len(release.Changed) > 0 {
		c.printType("Changed")
		c.printList(release.Changed)
	}

	if len(release.Deprecated) > 0 {
		c.printType("Deprecated")
		c.printList(release.Deprecated)
	}

	if len(release.Removed) > 0 {
		c.printType("Removed")
		c.printList(release.Removed)
	}

	if len(release.Fixed) > 0 {
		c.printType("Fixed")
		c.printList(release.Fixed)
	}

	if len(release.Security) > 0 {
		c.printType("Security")
		c.printList(release.Security)
	}
}

func (c *Changelog) print() error {
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
			c.printTitle(release, false)
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
		c.printRelease(*c.Changelog.LatestRelease())
		fmt.Println()
		c.printRelease(*c.Changelog.FirstRelease())
		return nil
	}

	if c.Latest {
		c.printRelease(*c.Changelog.LatestRelease())
		return nil
	}

	if c.First {
		c.printRelease(*c.Changelog.FirstRelease())
		return nil
	}

	for i, release := range c.Changelog.Releases {
		c.printRelease(release)
		if i != l-1 {
			fmt.Println()
		}
	}

	return nil
}

func (c *Changelog) run(path string) error {
	c.Changelog = changelog.New()

	if err := c.Changelog.Load(path); err != nil {
		return err
	}

	if err := c.print(); err != nil {
		return err
	}

	return nil
}
