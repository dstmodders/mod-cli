package dir

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func assertIsPathIgnored(t *testing.T, d *Dir, testCases map[string]bool) {
	for testCase, result := range testCases {
		msg := `Path "%s" should be ignored`
		if result == false {
			msg = `Path "%s" should not be ignored`
		}
		assert.Equalf(t, d.IsPathIgnored(testCase), result, msg, testCase)
	}
}

func TestNew(t *testing.T) {
	d, err := New(".")
	assert.Nil(t, err)
	assert.IsType(t, Dir{}, *d)
	assert.Equal(t, ".", d.relPath)
	assert.Len(t, d.ignore, 0)
}

func TestWorkshop_IsPathIgnored(t *testing.T) {
	d, _ := New(".")

	// .git/
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

	d.ignore = []string{".git/"}
	assertIsPathIgnored(t, d, testCases)

	// .git
	d.ignore = []string{".git"}
	assertIsPathIgnored(t, d, testCases)

	// /.git/
	testCases = map[string]bool{
		".git":          true,
		".git/":         true,
		".git/objects/": true,

		"one/.git":                 false,
		"one/.git/":                false,
		"one/.git/test/":           false,
		"one/two/.git/test/":       false,
		"one/two/three/.git/test/": false,
	}

	d.ignore = []string{"/.git/"}
	assertIsPathIgnored(t, d, testCases)

	// /.git
	d.ignore = []string{"/.git"}
	assertIsPathIgnored(t, d, testCases)
}
