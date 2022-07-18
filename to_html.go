package gem

import (
	"regexp"
	"strings"
)

// converts a given string of Gemtext to basic HTML.
func ToHTML(gemtext string) string {
	const (
		stateDefault = iota
		statePreformatted
		stateUnorderedList
		stateBlockquote
	)

	var output strings.Builder
	var state = stateDefault

	lines := strings.Split(gemtext, "\n")

	for i, line := range lines {
		if state == statePreformatted {
			if strings.HasPrefix(line, "```") {
				// closing preformatted line
				output.WriteString(convertPreformattedClosing(line))
				state = stateDefault
			} else {
				// preformatted line
				output.WriteString(line + "\n")
			}
			continue
		}

		if len(line) == 0 {
			// whitespace line
			output.WriteString(convertWhitespace(line))
		} else if strings.HasPrefix(line, "=>") {
			// link line
			output.WriteString(convertLink(line))
		} else if strings.HasPrefix(line, "```") {
			// opening preformatted line
			state = statePreformatted
			output.WriteString(convertPreformattedOpening(line))
		} else if strings.HasPrefix(line, "#") {
			// heading line
			output.WriteString(convertHeading(line))
		} else if strings.HasPrefix(line, "*") {
			// unordered list line
			if state != stateUnorderedList {
				output.WriteString("<ul>")
				state = stateUnorderedList
			}
			output.WriteString(convertUnorderedListItem(line))
			if nextHasPrefix(lines, i, "*") {
				output.WriteString("</ul>")
				state = stateDefault
			}
		} else if strings.HasPrefix(line, ">") {
			// blockquote line
			if state != stateBlockquote {
				output.WriteString("<blockquote>")
				state = stateBlockquote
			}
			output.WriteString(convertBlockquote(line))
			if nextHasPrefix(lines, i, ">") {
				output.WriteString("</blockquote>")
				state = stateDefault
			}

		} else {
			// (default) text line
			output.WriteString(convertText(line))
		}
	}

	return output.String()
}

// converts a given Gemtext whitespace line to HTML
func convertWhitespace(line string) string {
	return "<br />"
}

// converts a given Gemtext text line to HTML
func convertText(line string) string {
	return "<p>" + line + "</p>"
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

	if strings.HasPrefix(line, "###") {
		// heading 3
		heading = heading + "<h3>"
		headingText := strings.TrimPrefix(line, "###")
		heading = heading + strings.TrimSpace(headingText)
		heading = heading + "</h3>"

	} else if strings.HasPrefix(line, "##") {
		// heading 2
		heading = heading + "<h2>"
		headingText := strings.TrimPrefix(line, "##")
		heading = heading + strings.TrimSpace(headingText)
		heading = heading + "</h2>"

	} else {
		// heading 1
		heading = heading + "<h1>"
		headingText := strings.TrimPrefix(line, "#")
		heading = heading + strings.TrimSpace(headingText)
		heading = heading + "</h1>"
	}

	return heading
}

// converts a given Gemtext unordered list to HTML
func convertUnorderedListItem(line string) string {
	listitem := strings.TrimSpace(strings.TrimPrefix(line, "*"))
	return "<li>" + listitem + "</li>"
}

// given an array of strings, the current index, and a prefix, returns true if
// the next line also has that prefix
func nextHasPrefix(lines []string, index int, prefix string) bool {
	return len(lines) == index+1 || !strings.HasPrefix(lines[index+1], prefix)
}

// converts a given Gemtext quote line to HTML
func convertBlockquote(line string) string {
	text := strings.TrimSpace(strings.TrimPrefix(line, ">"))
	quote := ""
	if len(text) == 0 {
		quote = convertWhitespace(text)
	} else {
		quote = convertText(text)
	}

	return quote
}
