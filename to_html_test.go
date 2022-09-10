package gem

import (
	"testing"

	th "git.sr.ht/~bossley9/gem/pkg/testhelpers"
)

// text lines

func TestToHTML_ConvertParagraph(t *testing.T) {
	src := "this is a text line"
	test := ToHTML(src)
	ref := "<p>this is a text line</p>"

	th.AssertEqual(t, test, ref)
}

func TestToHTML_ConvertWhitespace(t *testing.T) {
	src := ""
	test := ToHTML(src)
	ref := ""

	th.AssertEqual(t, test, ref)
}

func TestToHTML_ConvertMultipleParagraphs(t *testing.T) {
	src := `this is the story
of a man who
did something very bad`
	test := ToHTML(src)
	ref := "<p>this is the story<br />of a man who<br />did something very bad</p>"

	th.AssertEqual(t, test, ref)
}

func TestToHTML_ConvertSpaceBetweenParagraphs(t *testing.T) {
	src := `this is a paragraph.

this is another paragraph
that is multiline`
	test := ToHTML(src)
	ref := "<p>this is a paragraph.</p>" + "<p>this is another paragraph<br />that is multiline</p>"

	th.AssertEqual(t, test, ref)
}

// links

func TestToHTML_ConvertLink(t *testing.T) {
	src := "=> example.com my website"
	test := ToHTML(src)
	ref := `<p><a href="example.com">my website</a></p>`

	th.AssertEqual(t, test, ref)
}

func TestToHTML_ConvertLinkVariableWhitespace(t *testing.T) {
	src := "=>      example.com      my website"
	test := ToHTML(src)
	ref := `<p><a href="example.com">my website</a></p>`

	th.AssertEqual(t, test, ref)
}

func TestToHTML_ConvertLinkNoLinkName(t *testing.T) {
	src := "=> example.com"
	test := ToHTML(src)
	ref := `<p><a href="example.com">example.com</a></p>`

	th.AssertEqual(t, test, ref)
}

func TestToHTML_ConvertConsecutiveLinks(t *testing.T) {
	src := `=> example.com my website
=> website.gov a governmental agency
=> search.com`
	test := ToHTML(src)
	ref := `<p><a href="example.com">my website</a></p><p><a href="website.gov">a governmental agency</a></p><p><a href="search.com">search.com</a></p>`

	th.AssertEqual(t, test, ref)
}

func TestToHTML_ConvertLinkImages(t *testing.T) {
	src := `=> car.png my car is what I would call "cool"`
	test := ToHTML(src)
	ref := `<p><img src="car.png" alt="my car is what I would call &#34;cool&#34;" /></p>`

	th.AssertEqual(t, test, ref)
}

func TestToHTML_ConvertLinkFakeImageExtension(t *testing.T) {
	src := `=> carpng my car is what I would call "cool"`
	test := ToHTML(src)
	ref := `<p><a href="carpng">my car is what I would call &#34;cool&#34;</a></p>`

	th.AssertEqual(t, test, ref)
}

func TestToHTML_ConvertLinkAudio(t *testing.T) {
	src := `=> sounds.mp3 this is what my car sounds like`
	test := ToHTML(src)
	ref := "<p><audio controls>" + `<source src="sounds.mp3" type="audio/mp3" />` + "Sorry, your browser doesn't support embedded audio." + "</audio></p>"

	th.AssertEqual(t, test, ref)
}

func TestToHTML_ConvertLinkVideo(t *testing.T) {
	src := `=> vid.mp4 a documentary on cron`
	test := ToHTML(src)
	ref := "<p><video controls>" + `<source src="vid.mp4" type="video/mp4" />` + "Sorry, your browser doesn't support embedded video." + "</video></p>"

	th.AssertEqual(t, test, ref)
}

// preformatted text

func TestToHTML_ConvertPreformatted(t *testing.T) {
	src := "```\n" + `
=> example.com this will not be converted

Hello! Whitespace    is preserved here
<p>fake html tags</p>
` + "```"
	test := ToHTML(src)
	ref := `<figure><pre><code>
=&gt; example.com this will not be converted

Hello! Whitespace    is preserved here
&lt;p&gt;fake html tags&lt;/p&gt;
</code></pre></figure>`

	th.AssertEqual(t, test, ref)
}

func TestToHTML_ConvertPreformattedAltText(t *testing.T) {
	src := "```this is alt text\n" + `
some arbitrary code content
` + "```"
	test := ToHTML(src)
	ref := `<figure><figcaption>this is alt text</figcaption><pre><code>
some arbitrary code content
</code></pre></figure>`

	th.AssertEqual(t, test, ref)
}

// headings

func TestToHTML_ConvertHeadingOne(t *testing.T) {
	src := "#   heading here   "
	test := ToHTML(src)
	ref := `<h1 id="heading-here">heading here</h1>`

	th.AssertEqual(t, test, ref)
}

func TestToHTML_ConvertHeadingTwo(t *testing.T) {
	src := "##heading here"
	test := ToHTML(src)
	ref := `<h2 id="heading-here">heading here</h2>`

	th.AssertEqual(t, test, ref)
}

func TestToHTML_ConvertHeadingThree(t *testing.T) {
	src := "### heading here"
	test := ToHTML(src)
	ref := `<h3 id="heading-here">heading here</h3>`

	th.AssertEqual(t, test, ref)
}

func TestToHTML_ConvertIdenticalHeadingIDs(t *testing.T) {
	src := `# reference
## reference
### reference
### reference
`
	test := ToHTML(src)
	ref := `<h1 id="reference">reference</h1>` +
		`<h2 id="reference-1">reference</h2>` +
		`<h3 id="reference-2">reference</h3>` +
		`<h3 id="reference-3">reference</h3>`

	th.AssertEqual(t, test, ref)
}

// unordered lists

func TestToHTML_ConvertList(t *testing.T) {
	src := "* eggs"
	test := ToHTML(src)
	ref := "<ul><li>eggs</li></ul>"

	th.AssertEqual(t, test, ref)
}

func TestToHTML_ConvertListMultiItem(t *testing.T) {
	src := `* eggs
* milk
* white bread
* greens`
	test := ToHTML(src)
	ref := "<ul>" + "<li>eggs</li>" + "<li>milk</li>" + "<li>white bread</li>" + "<li>greens</li>" + "</ul>"

	th.AssertEqual(t, test, ref)
}

// blockquotes

func TestToHTML_ConvertBlockquote(t *testing.T) {
	src := "> quote here"
	test := ToHTML(src)
	ref := "<blockquote><p>quote here</p></blockquote>"

	th.AssertEqual(t, test, ref)
}

func TestToHTML_ConvertBlockquoteMultiline(t *testing.T) {
	src := `> this is a
>   multiline spanning
>
>blockquote`
	test := ToHTML(src)
	ref := "<blockquote><p>this is a<br />multiline spanning</p><p>blockquote</p></blockquote>"

	th.AssertEqual(t, test, ref)
}
