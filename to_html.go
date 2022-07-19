package gem

import (
	"regexp"
	"strings"
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
				output.WriteString(line + "\n")
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
			output.WriteString(convertHeading(line))
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
			output.WriteString(line)
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

	} else if strings.HasPrefix(line, "*") {
		lineType = lineUnorderedListItem

	} else if strings.HasPrefix(line, ">") {
		lineType = lineBlockQuote

	} else {
		lineType = lineParagraph
	}

	return lineType
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

	// output <img> or <a> based on url suffix
	match, _ := regexp.MatchString(".*.(gif|jpg|jpeg|png|svg|webp)", url)
	if match {
		output = output + `<img src="` + url + `" alt="`
		if len(name) > 0 {
			output = output + strings.ReplaceAll(name, "\"", "\\\"")
		}
		output = output + `" />`

	} else {
		output = output + `<a href="` + url + `">`
		if len(name) > 0 {
			output = output + name
		} else {
			output = output + url
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
		output = output + "<figcaption>" + alt + "</figcaption>"
	}

	output = output + "<pre><code>"
	return output
}

func convertPreformattedClosing(line string) string {
	return "</code></pre></figure>"
}

// converts a given Gemtext heading to HTML
func convertHeading(line string) string {
	heading := ""
	lineType := getLineType(line)

	if lineType == lineHeadingThree {
		headingText := strings.TrimPrefix(line, "###")
		heading = heading + "<h3>" + strings.TrimSpace(headingText) + "</h3>"

	} else if lineType == lineHeadingTwo {
		headingText := strings.TrimPrefix(line, "##")
		heading = heading + "<h2>" + strings.TrimSpace(headingText) + "</h2>"

	} else {
		headingText := strings.TrimPrefix(line, "#")
		heading = heading + "<h1>" + strings.TrimSpace(headingText) + "</h1>"
	}

	return heading
}

// converts a given Gemtext unordered list to HTML
func convertUnorderedListItem(line string) string {
	listitem := strings.TrimSpace(strings.TrimPrefix(line, "*"))
	return "<li>" + listitem + "</li>"
}

// given an array of strings, the current index, and a line type, returns true if
// the next line is not that type
func nextIsNotType(lines []string, index int, lineType lineType) bool {
	return len(lines) == index+1 || getLineType(lines[index+1]) != lineType
}

// converts a given Gemtext quote line to HTML
func convertBlockquote(line string) string {
	return strings.TrimSpace(strings.TrimPrefix(line, ">"))
}

// converts a given Gemtext quote lookahead line to HTML
func convertBlockquoteFiller(line string) string {
	if len(convertBlockquote(line)) == 0 {
		return "</p><p>"
	} else {
		return "<br />"
	}
}
