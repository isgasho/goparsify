package goparsify

import (
	"bytes"
	"unicode/utf8"
)

func StringLit(allowedQuotes string) Parser {
	return NewParser("string", func(ps *State) Node {
		ps.AutoWS()

		for i := 0; i < len(allowedQuotes); i++ {
			if ps.Input[ps.Pos] == allowedQuotes[i] {

			}
		}
		if !stringContainsByte(allowedQuotes, ps.Input[ps.Pos]) {
			ps.ErrorHere(allowedQuotes)
			return Node{}
		}
		quote := ps.Input[ps.Pos]

		var end int = ps.Pos + 1

		inputLen := len(ps.Input)
		var buf *bytes.Buffer

		for end < inputLen {
			switch ps.Input[end] {
			case '\\':
				if end+1 >= inputLen {
					ps.ErrorHere(string(quote))
					return Node{}
				}

				if buf == nil {
					buf = bytes.NewBufferString(ps.Input[ps.Pos+1 : end])
				}

				c := ps.Input[end+1]
				if c == 'u' {
					if end+6 >= inputLen {
						ps.Error.Expected = "[a-f0-9]{4}"
						ps.Error.pos = end + 2
						return Node{}
					}

					r, ok := unhex(ps.Input[end+2 : end+6])
					if !ok {
						ps.Error.Expected = "[a-f0-9]"
						ps.Error.pos = end + 2
						return Node{}
					}
					buf.WriteRune(r)
					end += 6
				} else {
					buf.WriteByte(c)
					end += 2
				}
			case quote:
				if buf == nil {
					result := ps.Input[ps.Pos+1 : end]
					ps.Pos = end + 1
					return Node{Token: result}
				}
				ps.Pos = end + 1
				return Node{Token: buf.String()}
			default:
				r, w := utf8.DecodeRuneInString(ps.Input[end:])
				end += w
				if buf != nil {
					buf.WriteRune(r)
				}
			}
		}

		ps.ErrorHere(string(quote))
		return Node{}
	})
}

func stringContainsByte(s string, b byte) bool {
	for i := 0; i < len(s); i++ {
		if b == s[i] {
			return true
		}
	}
	return false
}

func unhex(b string) (v rune, ok bool) {
	for _, c := range b {
		v <<= 4
		switch {
		case '0' <= c && c <= '9':
			v |= c - '0'
		case 'a' <= c && c <= 'f':
			v |= c - 'a' + 10
		case 'A' <= c && c <= 'F':
			v |= c - 'A' + 10
		default:
			return 0, false
		}
	}

	return v, true
}