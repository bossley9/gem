package gem

import (
	"html"
	"regexp"
	"strconv"
	"strings"

	"git.sr.ht/~bossley9/gem/pkg/url"

	"github.com/gomarkdown/markdown"
)

type lineType int

const (
	lineParagraph lineType = iota
	lineWhitespace
	lineLink
	linePreformattedEdge
	lineHeadingOne
	lineHeadingTwo
	lineHeadingThree
	lineUnorderedListItem
	lineBlockQuote
)

// converts a given string of Gemtext to basic HTML.
func ToHTML(gemtext string) string {
	if len(gemtext) == 0 {
		return ""
	}
	idRefs := make(map[string]int, 0)

	const (
		stateDefault = iota
		stateParagraph
		statePreformatted
		stateUnorderedList
		stateBlockquote
	)

	var output strings.Builder
	var state = stateDefault

	lines := strings.Split(gemtext, "\n")

	for i, line := range lines {
		lineType := getLineType(line)

		if state == statePreformatted {
			if lineType == linePreformattedEdge {
				// closing preformatted line
				output.WriteString(convertPreformattedClosing(line))
				state = stateDefault
			} else {
				// preformatted line
				output.WriteString(html.EscapeString(line) + "\n")
			}
			continue
		}

		switch lineType {
		case lineLink:
			output.WriteString(convertLink(line))
		case linePreformattedEdge:
			// opening preformatted line
			state = statePreformatted
			output.WriteString(convertPreformattedOpening(line))
		case lineHeadingThree:
			fallthrough
		case lineHeadingTwo:
			fallthrough
		case lineHeadingOne:
			id := url.GenerateID(line)
			idCount, ok := idRefs[id]
			if ok {
				idRefs[id] = idCount + 1
				id = id + "-" + strconv.Itoa(idCount)
			} else {
				idRefs[id] = 1
			}
			output.WriteString(convertHeading(line, id))
		case lineUnorderedListItem:
			if state != stateUnorderedList {
				output.WriteString("<ul>")
				state = stateUnorderedList
			}
			output.WriteString(convertUnorderedListItem(line))
			if nextIsNotType(lines, i, lineUnorderedListItem) {
				output.WriteString("</ul>")
				state = stateDefault
			}
		case lineBlockQuote:
			// opening blockquote
			if state != stateBlockquote {
				output.WriteString("<blockquote><p>")
				state = stateBlockquote
			}
			quote := convertBlockquote(line)
			output.WriteString(quote)
			// closing blockquote
			if nextIsNotType(lines, i, lineBlockQuote) {
				output.WriteString("</p></blockquote>")
				state = stateDefault
			} else {
				// between blockquote lines
				if len(quote) > 0 {
					output.WriteString(convertBlockquoteFiller(lines[i+1]))
				}
			}
		case lineWhitespace:
			if state == stateParagraph {
				output.WriteString("<br />")
			}
		case lineParagraph:
			fallthrough
		default:
			// opening paragraph
			if state != stateParagraph {
				output.WriteString("<p>")
				state = stateParagraph
			}
			output.WriteString(convertText(line))
			// closing paragraph
			if nextIsNotType(lines, i, lineParagraph) {
				output.WriteString("</p>")
				state = stateDefault
			} else {
				output.WriteString("<br />")
			}
		}
	}

	return output.String()
}

// given a line of text, returns the Gemtext line type
func getLineType(line string) lineType {
	var lineType lineType

	if len(line) == 0 {
		lineType = lineWhitespace

	} else if strings.HasPrefix(line, "=>") {
		lineType = lineLink

	} else if strings.HasPrefix(line, "```") {
		lineType = linePreformattedEdge

	} else if strings.HasPrefix(line, "###") {
		lineType = lineHeadingThree

	} else if strings.HasPrefix(line, "##") {
		lineType = lineHeadingTwo

	} else if strings.HasPrefix(line, "#") {
		lineType = lineHeadingOne

	} else if strings.HasPrefix(line, "* ") {
		lineType = lineUnorderedListItem

	} else if strings.HasPrefix(line, ">") {
		lineType = lineBlockQuote

	} else {
		lineType = lineParagraph
	}

	return lineType
}

// converts a given text string to HTML
func convertText(text string) string {
	if len(text) == 0 {
		return text
	}

	md := string(markdown.ToHTML([]byte(text), nil, nil))
	// trim outer paragraph tag
	mdTrimmed := md[3 : len(md)-5]
	// prefer ambiguous quotations
	mdFormatted := strings.ReplaceAll(mdTrimmed, "&rdquo;", "&#34;")
	mdFormatted = strings.ReplaceAll(mdFormatted, "&ldquo;", "&#34;")

	return mdFormatted
}

// converts a given Gemtext link to HTML
func convertLink(line string) string {
	output := "<p>"
	link := strings.TrimSpace(strings.TrimPrefix(line, "=>"))
	split := strings.SplitN(link, " ", 2)

	url := ""
	name := ""

	if len(split) >= 1 {
		url = split[0]
	}
	if len(split) >= 2 {
		name = strings.TrimSpace(split[1])
	}

	regexImg := regexp.MustCompile(`.*\.(avif|bmp|gif|jpg|jpeg|png|svg|webp|xpm)`)
	regexAudio := regexp.MustCompile(`.*\.(m3u|m4a|mp3|ogg|wav)`)
	regexVideo := regexp.MustCompile(`.*\.(avi|divx|m4v|mkv|mov|mp4|mpeg|mpg|vob|webm|wmv)`)

	if regexImg.MatchString(url) {
		output = output + `<img src="` + html.EscapeString(url) + `" alt="`
		if len(name) > 0 {
			output = output + html.EscapeString(name)
		}
		output = output + `" />`

	} else if regexAudio.MatchString(url) {
		ext := url[strings.LastIndex(url, ".")+1:]

		output = output + "<audio controls>"
		output = output + `<source src="` + url + `" type="audio/` + ext + `" />`
		output = output + "Sorry, your browser doesn't support embedded audio."
		output = output + "</audio>"

	} else if regexVideo.MatchString(url) {
		ext := url[strings.LastIndex(url, ".")+1:]

		output = output + "<video controls>"
		output = output + `<source src="` + url + `" type="video/` + ext + `" />`
		output = output + "Sorry, your browser doesn't support embedded video."
		output = output + "</video>"

	} else {
		output = output + `<a href="` + html.EscapeString(url) + `">`
		if len(name) > 0 {
			output = output + convertText(name)
		} else {
			output = output + convertText(url)
		}
		output = output + `</a>`

	}

	output = output + "</p>"
	return output
}

// converts a given opening preformatted Gemtext line to HTML
func convertPreformattedOpening(line string) string {
	output := "<figure>"

	alt := strings.TrimSpace(strings.TrimPrefix(line, "```"))
	if len(alt) > 0 {
		output = output + "<figcaption>" + convertText(alt) + "</figcaption>"
	}

	output = output + "<pre><code>"
	return output
}

func convertPreformattedClosing(line string) string {
	return "</code></pre></figure>"
}

// converts a given Gemtext heading to HTML
func convertHeading(line string, id string) string {
	lineType := getLineType(line)
	headingNum := ""
	headingPrefix := ""

	if lineType == lineHeadingThree {
		headingNum = "3"
		headingPrefix = "###"

	} else if lineType == lineHeadingTwo {
		headingNum = "2"
		headingPrefix = "##"

	} else {
		headingNum = "1"
		headingPrefix = "#"
	}

	headingText := strings.TrimSpace(strings.TrimPrefix(line, headingPrefix))
	return "<h" + headingNum + ` id="` + id + `"><a href="#` + id + `">` + convertText(headingText) + "</a></h" + headingNum + ">"
}

// converts a given Gemtext unordered list to HTML
func convertUnorderedListItem(line string) string {
	listitem := strings.TrimSpace(strings.TrimPrefix(line, "*"))
	return "<li>" + convertText(listitem) + "</li>"
}

// given an array of strings, the current index, and a line type, returns true if
// the next line is not that type
func nextIsNotType(lines []string, index int, lineType lineType) bool {
	return len(lines) == index+1 || getLineType(lines[index+1]) != lineType
}

// converts a given Gemtext quote line to HTML
func convertBlockquote(line string) string {
	text := strings.TrimSpace(strings.TrimPrefix(line, ">"))
	return convertText(text)
}

// converts a given Gemtext quote lookahead line to HTML
func convertBlockquoteFiller(line string) string {
	if len(convertBlockquote(line)) == 0 {
		return "</p><p>"
	} else {
		return "<br />"
	}
}
