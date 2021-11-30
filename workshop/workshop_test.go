package workshop

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	// preparations
	src := "./"
	dest := "./test/"

	// test
	w, err := New(src, dest)
	assert.Nil(t, err)
	assert.IsType(t, Workshop{}, *w)
	assert.Equal(t, ".", w.relSrcPath)
	assert.Equal(t, "test", w.relDestPath)

	assert.ElementsMatch(t, []string{
		".*",
		"Makefile",
		"codecov.yml",
		"config.ld",
		"lcov.info",
		"luacov.*",
		"spec/",
	}, w.Ignore)
}

// IsPathIgnored checks if the provided path is ignored based on Ignore.
func TestWorkshop_IsPathIgnored(t *testing.T) {
	testCases := map[string]bool{
		".git":          true,
		".git/":         true,
		".git/objects/": true,

		"one/.git":                 true,
		"one/.git/":                true,
		"one/.git/test/":           true,
		"one/two/.git/test/":       true,
		"one/two/three/.git/test/": true,
	}

	w, _ := New(".", "test")

	for testCase, r := range testCases {
		assert.Equal(t, w.IsPathIgnored(testCase), r)
	}
}
