// Copyright (c) 2015 Matt Borgerson
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package truncate

import (
	"fmt"
	"regexp"
	"strings"
	"unicode"
	"unicode/utf8"
)

var (
	tagExpr    = regexp.MustCompile("<(/?)([A-Za-z0-9]+).*?>")
	entityExpr = regexp.MustCompile("&#?[A-Za-z0-9]+;")

	// We will consider HTML or XHTML as valid input. The following elements,
	// called "Void Elements" need not conform to the XHTML <tag /> convention
	// of void elements and may appear simply as <tag>. Hence, if one of the
	// following is picked up by the tag expression as a start tag, do not add
	// it to the stack of tags that should be closed.
	voidElementTags = []string{
		"area",
		"base",
		"br",
		"col",
		"embed",
		"hr",
		"img",
		"input",
		"keygen",
		"link",
		"meta",
		"param",
		"source",
		"track",
		"wbr",
	}
)

// HTML will truncate a given byte slice to a maximum of maxlen visible
// characters and optionally append suffix. HTML tags are automatically closed
// generating valid truncated HTML.
func HTML(buf []byte, maxlen int, suffix string) ([]byte, error) {
	// Scan the input bytestream. While scanning, count the
	// number of visible characters--that is, characters which are not part of
	// markup tags. When a start tag is encountered, push the tag name onto a
	// stack. When visible character count >= maxlen, or the EOF is reached,
	// stop counting. Copy from the input stream the bytes from the start to the
	// current scanning pointer. Finally, pop each tag off the tag stack and
	// append it to the output stream in the form of a closing tag.

	// Check to see if no input was provided.
	if buf == nil || len(buf) == 0 || maxlen == 0 {
		return []byte{}, nil
	}

	tagStack := []string{}
	visible := 0
	bufPtr := 0

	for bufPtr < len(buf) && visible < maxlen {
		// Move to nearest tag and count visible characters along the way.
		offset := 0
		visibleCharacterMaxReached := false
		entityDetected := false

		for localOffset, runeValue := range string(buf[bufPtr:]) {
			offset = localOffset

			if runeValue == '<' {
				// Start of tag.
				break
			} else if runeValue == '&' {
				// Possible start of HTML Entity
				loc := entityExpr.FindIndex(buf[bufPtr+localOffset:])
				if loc != nil && loc[0] == 0 {
					// Entity found!
					entityDetected = true
					offset += loc[1] - 1 // Now pointing to ;
				}
				visible += 1
			} else if unicode.IsPrint(runeValue) && !unicode.IsSpace(runeValue) {
				// Printable, non-space character. Increment visible count.
				visible += 1
			}

			// Check if the limit of visible characters has been reached.
			if visible >= maxlen {
				visibleCharacterMaxReached = true
				break
			}

			if entityDetected {
				break
			}
		}

		// Increment bufPtr to end of scanned section
		bufPtr += offset

		// Stop scanning if the end of the buffer was reached or if the max
		// desired visible characters was reached
		if visibleCharacterMaxReached || bufPtr >= len(buf)-1 {
			break
		}

		// If an entity was detected, continue scanning for next tag
		if entityDetected {
			// Advance past the ;
			bufPtr += 1
			continue
		}

		// Now find the expression sub-matches
		matches := tagExpr.FindSubmatch(buf[bufPtr:])
		if matches == nil {
			bufPtr += 1
			continue
		}
		tagName := strings.ToLower(string(matches[2]))

		// Advance pointer to the end of the tag
		bufPtr += len(matches[0])

		// If this is a void element, do not count it as a start tag
		isVoidElement := false
		for _, voidElementTagName := range voidElementTags {
			if strings.EqualFold(tagName, voidElementTagName) {
				isVoidElement = true
				break
			}
		}
		if isVoidElement {
			continue
		}

		isStartTag := len(matches[1]) == 0

		if isStartTag {
			// This is a start tag. Push the tag to the stack.
			tagStack = append(tagStack, tagName)
		} else {
			// This is an end tag. First, check to make sure the end tag is
			// matches what's on top of the stack.
			if len(tagStack) == 0 || tagStack[len(tagStack)-1] != tagName {
				return nil, fmt.Errorf("unbalanced tag %q", tagName)
			}

			// Now, pop the tag stack.
			tagStack = tagStack[0 : len(tagStack)-1]
		}
	}

	// At this point, bufPtr points to the last rune that should be copied to
	// the output stream. Increment bufPtr past this rune, turning bufPtr into
	// the number of bytes that should be copied.
	_, size := utf8.DecodeRune(buf[bufPtr:])
	bufPtr += size

	// Copy the desired input to the output buffer.
	output := buf[0:bufPtr]

	output = append(output, []byte(suffix)...)

	// Finally, create a closing tag for each tag in the stack.
	for i := len(tagStack) - 1; i >= 0; i-- {
		output = append(output, []byte(fmt.Sprintf("</%s>", tagStack[i]))...)
	}

	return output, nil
}
