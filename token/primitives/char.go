package primitives

import (
	"fmt"
	"io"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/zimmski/tavor/log"
	"github.com/zimmski/tavor/token"
)

var simpleEscapes = map[rune]rune{
	'-':  '-',
	'\\': '\\',
	'f':  '\f',
	'n':  '\n',
	'r':  '\r',
	't':  '\t',
}

var characterClassEscapes = map[rune][]rune{
	'd': []rune{'0', '1', '2', '3', '4', '5', '6', '7', '8', '9'},
	's': []rune{' ', '\t', '\n', '\f', '\r'},
	'w': []rune{'0', '1', '2', '3', '4', '5', '6', '7', '8', '9', 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'N', 'M', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z', 'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'n', 'm', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z', '_'},
}

// CharacterClass implements a char token which holds a pattern of characters and character classes
// Each character class characters is added to the set of characters of the token. Every permutations chooses one character out of the available set of characters as the current value of the token.
type CharacterClass struct {
	chars       []rune
	charsLookup map[rune]struct{}
	charRanges  []characterRange

	permutations uint

	pattern string

	value rune
}

type characterRange struct {
	from, to rune
}

// NewCharacterClass returns a new instance of a CharacterClass token
func NewCharacterClass(pattern string) *CharacterClass {
	if pattern == "" {
		panic("pattern is empty")
	}

	/*

		TODO FIXME FIXME FIXME

		This part of the code is full of bad code, especially the error handling is not addressed at all, we simple panic!
		There is a time and place where this can be redone, but it is not now and here.

	*/

	var chars []rune
	var charRanges []characterRange
	var lastCharIsRangeChar = false
	var lastChar rune
	var isRange = false

	runes := strings.NewReader(pattern)

	c, _, err := runes.ReadRune()

	add := func(c rune) {
		if isRange {
			if lastChar > c {
				panic("Range to character is lower than range from character")
			}

			charRanges = append(charRanges, characterRange{
				from: lastChar,
				to:   c,
			})

			chars = chars[:len(chars)-1] // remove last character since it is now in a range

			isRange = false
			lastCharIsRangeChar = false
		} else {
			chars = append(chars, c)
		}
	}

	checkHex := func(c rune) bool {
		return (c >= '0' && c <= '9') || (c >= 'a' && c <= 'f') || (c >= 'A' && c <= 'F')
	}

PARSING:
	for err != io.EOF {
		if unicode.IsDigit(c) || unicode.IsLetter(c) || unicode.IsSpace(c) {
			add(c)
			lastChar = c
			lastCharIsRangeChar = true
		} else {
			switch c {
			case '-':
				if !lastCharIsRangeChar {
					panic("Range operator without range from character")
				}

				isRange = true
			case '\\':
				c, _, err = runes.ReadRune()
				if err == io.EOF {
					panic("early EOF for escaped character")
				} else if err != nil {
					break PARSING
				}

				switch c {
				case 'x':
					x, _, err := runes.ReadRune()
					if err == io.EOF {
						panic("early EOF for escaped character")
					} else if err != nil {
						break PARSING
					}

					var xses string

					if x == '{' {
						for {
							x, _, err = runes.ReadRune()
							if err == io.EOF {
								panic("early EOF for escaped character")
							} else if err != nil {
								break PARSING
							} else if x == '}' {
								break
							} else if !checkHex(x) {
								panic("x escaping needs HEX characters")
							}

							xses += string(x)
						}

						if len(xses) < 2 {
							panic("x escaping needs two HEX characters")
						}
					} else {
						if !checkHex(x) {
							panic("x escaping needs two HEX characters")
						}

						xses += string(x)

						x, _, err = runes.ReadRune()
						if err == io.EOF {
							panic("early EOF for escaped character")
						} else if err != nil {
							break PARSING
						} else if !checkHex(x) {
							panic("x escaping needs two HEX characters")
						}

						xses += string(x)
					}

					s, e := strconv.Unquote(`"\U` + strings.Repeat("0", 8-len(xses)) + xses + `"`)
					if e != nil {
						panic(e)
					}

					c, _ = utf8.DecodeRuneInString(s)

					add(c)
					lastChar = c
					lastCharIsRangeChar = true
				default:
					if simp, ok := simpleEscapes[c]; ok {
						add(simp)
						lastChar = simp
						lastCharIsRangeChar = true
					} else {
						if isRange {
							panic("Range operator without range to character")
						}

						esc, ok := characterClassEscapes[c]
						if !ok {
							panic(fmt.Sprintf("Unknown escape character %q", c))
						}

						for _, v := range esc {
							add(v)
						}

						lastCharIsRangeChar = false
					}
				}
			default:
				panic(fmt.Sprintf("Unknown character %q", c))
			}
		}

		c, _, err = runes.ReadRune()
	}

	if err != nil && err != io.EOF {
		panic(err)
	}

	if len(chars) == 0 && len(charRanges) == 0 {
		panic("empty character class is not allowed")
	}

	var first rune
	charsLookup := make(map[rune]struct{})

	if len(chars) != 0 {
		first = chars[0]

		for _, v := range chars {
			charsLookup[v] = struct{}{}
		}
	} else {
		first = charRanges[0].from
	}

	var permutations = uint(len(chars))

	for _, v := range charRanges {
		permutations += uint(v.to-v.from) + 1
	}

	return &CharacterClass{
		chars:       chars,
		charsLookup: charsLookup,
		charRanges:  charRanges,

		permutations: permutations,

		pattern: pattern,

		value: first,
	}
}

// Clone returns a copy of the token and all its children
func (c *CharacterClass) Clone() token.Token {
	chars := make([]rune, len(c.chars))

	copy(chars, c.chars)

	charsLookup := make(map[rune]struct{})

	for k := range c.charsLookup {
		charsLookup[k] = struct{}{}
	}

	charRanges := make([]characterRange, len(c.charRanges))
	for i := range c.charRanges {
		charRanges[i] = c.charRanges[i]
	}

	return &CharacterClass{
		chars:       chars,
		charsLookup: charsLookup,
		charRanges:  charRanges,

		permutations: c.permutations,

		pattern: c.pattern,

		value: c.value,
	}
}

// Parse tries to parse the token beginning from the current position in the parser data.
// If the parsing is successful the error argument is nil and the next current position after the token is returned.
func (c *CharacterClass) Parse(pars *token.InternalParser, cur int) (int, []error) {
	if cur+1 > pars.DataLen {
		return cur, []error{&token.ParserError{
			Message: fmt.Sprintf("expected %q but got early EOF", c.charsLookup),
			Type:    token.ParseErrorUnexpectedEOF,

			Position: pars.GetPosition(cur),
		}}
	}

	// TODO FIXME NOW we can see the need to put pars.Data into a reader... since we cannot do a readRune here!
	v := rune(pars.Data[cur])

	if _, ok := c.charsLookup[v]; !ok {
		found := false

		for _, r := range c.charRanges {
			if v >= r.from && v <= r.to {
				found = true

				break
			}
		}

		if !found {
			return cur, []error{&token.ParserError{
				Message: fmt.Sprintf("expected %q or %+v but got %q", c.charsLookup, c.charRanges, v),
				Type:    token.ParseErrorUnexpectedData,

				Position: pars.GetPosition(cur),
			}}
		}
	}

	c.value = v

	log.Debugf("Parsed %q", v)

	return cur + 1, nil
}

func (c *CharacterClass) permutation(i uint) {
	cl := uint(len(c.chars))

	if i < cl {
		c.value = c.chars[i]

		return
	}

	i -= cl

	for _, v := range c.charRanges {
		cl := uint(v.to-v.from) + 1

		if i < cl {
			c.value = rune(uint(v.from) + i)

			return
		}

		i -= cl
	}

	panic(fmt.Sprintf("TODO out of range with pattern %q and permutation %d", c.pattern, i))

}

// Permutation sets a specific permutation for this token
func (c *CharacterClass) Permutation(i uint) error {
	permutations := c.Permutations()

	if i < 1 || i > permutations {
		return &token.PermutationError{
			Type: token.PermutationErrorIndexOutOfBound,
		}
	}

	c.permutation(i - 1)

	return nil
}

// Permutations returns the number of permutations for this token
func (c *CharacterClass) Permutations() uint {
	return c.permutations
}

// PermutationsAll returns the number of all possible permutations for this token including its children
func (c *CharacterClass) PermutationsAll() uint {
	return c.Permutations()
}

func (c *CharacterClass) String() string {
	return string(c.value)
}
