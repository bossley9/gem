package url

import (
	"testing"

	th "git.sr.ht/~bossley9/gem/pkg/testhelpers"
)

func TestGenerateID_Routine(t *testing.T) {
	test := `# New section, here! Now!
content line
`
	ref := "new-section-here-now"

	th.AssertEqual(t, GenerateID(test), ref)
}

func TestGenerateID_MaxLength(t *testing.T) {
	test := "# New section, here! Now! This is more of the title and will eventually be cut off at some point"
	ref := "new-section-here-now-this-is-mor"

	th.AssertEqual(t, GenerateID(test), ref)
}

func TestGenerateID_TrailingHyphen(t *testing.T) {
	test := "# New section, here! Now! This is it for you"
	ref := "new-section-here-now-this-is-it"

	th.AssertEqual(t, GenerateID(test), ref)
}
