package gem

import "testing"

func TestToHTML_ConvertTextLine(t *testing.T) {
	src := "this is a text line"
	test := ToHTML(src)
	ref := "<p>this is a text line</p>"

	assertEqual(t, test, ref)
}

func TestToHTML_ConvertWhitespace(t *testing.T) {
	src := ""
	test := ToHTML(src)
	ref := "<br />"

	assertEqual(t, test, ref)
}

func TestToHTML_ConvertLink(t *testing.T) {
	src := "=> example.com my website"
	test := ToHTML(src)
	ref := `<p><a href="example.com">my website</a></p>`

	assertEqual(t, test, ref)
}

func TestToHTML_ConvertLinkVariableWhitespace(t *testing.T) {
	src := "=>      example.com      my website"
	test := ToHTML(src)
	ref := `<p><a href="example.com">my website</a></p>`

	assertEqual(t, test, ref)
}

func TestToHTML_ConvertLinkNoLinkName(t *testing.T) {
	src := "=> example.com"
	test := ToHTML(src)
	ref := `<p><a href="example.com">example.com</a></p>`

	assertEqual(t, test, ref)
}

func TestToHTML_ConvertConsecutiveLinks(t *testing.T) {
	src := `=> example.com my website
=> website.gov a governmental agency
=> search.com`
	test := ToHTML(src)
	ref := `<p><a href="example.com">my website</a></p><p><a href="website.gov">a governmental agency</a></p><p><a href="search.com">search.com</a></p>`

	assertEqual(t, test, ref)
}

func TestToHTML_ConvertPreformatted(t *testing.T) {
	src := "```\n" + `
=> example.com this will not be converted

Hello! Whitespace    is preserved here
<p>fake html tags</p>
` + "```"
	test := ToHTML(src)
	ref := `<figure><pre><code>
=> example.com this will not be converted

Hello! Whitespace    is preserved here
<p>fake html tags</p>
</code></pre></figure>`

	assertEqual(t, test, ref)
}

func TestToHTML_ConvertPreformattedAltText(t *testing.T) {
	src := "```this is alt text\n" + `
some arbitrary code content
` + "```"
	test := ToHTML(src)
	ref := `<figure><figcaption>this is alt text</figcaption><pre><code>
some arbitrary code content
</code></pre></figure>`

	assertEqual(t, test, ref)
}
