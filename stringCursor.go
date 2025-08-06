package main

import (
	"strings"
)

type StringCursor struct {
	pos  int
	ind  int
	text string
}

func simplifyReq(req string) string {
	sc := NewStringCursor(req)
	return sc.simlifyExpr()
}

func NewStringCursor(s string) *StringCursor {
	return &StringCursor{
		text: s,
	}
}

func (sc *StringCursor) Consume() {
	sc.pos = sc.ind
}

func (sc *StringCursor) isEOF() bool {
	sc.TrimWhites()
	return len(sc.text) == sc.ind
}

func (sc *StringCursor) TrimWhites() {
	for sc.pos < len(sc.text) && sc.text[sc.pos] == ' ' {
		sc.pos++
	}
}

func (sc *StringCursor) NextWord() string {
	sc.TrimWhites()

	if sc.isEOF() {
		return ""
	}
	sc.ind = sc.pos
	start := sc.ind
	for sc.ind < len(sc.text) && sc.text[sc.ind] != ' ' && sc.text[sc.ind] != '(' && sc.text[sc.ind] != ')' {
		sc.ind++
	}
	word := sc.text[start:sc.ind]
	return word
}

func (sc *StringCursor) HasPrefix(s string) bool {
	sc.TrimWhites()
	sc.ind = min(len(sc.text), sc.pos+len(s))
	return sc.text[sc.pos:sc.ind] == s
}

func (sc *StringCursor) simlifyExpr() string {
	sb := new(strings.Builder)
	sb.WriteString(sc.simlifyTerm())

	for !sc.isEOF() && (sc.HasPrefix("AND") || sc.HasPrefix("OR")) {
		sb.WriteByte(' ')
		if sc.HasPrefix("AND") {
			sb.WriteString("AND")
		} else {
			sb.WriteString("OR")
		}
		sc.Consume()

		sb.WriteByte(' ')
		sb.WriteString(sc.simlifyTerm())
	}

	return sb.String()
}

func (sc *StringCursor) simlifyTerm() string {
	sb := new(strings.Builder)
	if sc.isEOF() {
		return ""
	}

	if sc.HasPrefix("(") {
		sb.WriteRune('(')
		sc.Consume()

		sb.WriteString(sc.simlifyExpr())
		if !sc.HasPrefix(")") {
			panic("Expected right parenthesis")
		}
		sc.Consume()
		sb.WriteByte(')')
	} else if sc.HasPrefix(")") {
		return ""
	} else {
		sb.WriteString(sc.NextWord())
		sc.Consume()

		sb.WriteByte(' ')
		sb.WriteString(sc.NextWord())
		sc.Consume()

		bracketBalance := 0
		for !sc.isEOF() {
			if sc.HasPrefix("(") {
				// slog.Info("Found left parenthesis", "text", sc.NextWord())
				bracketBalance++
			} else if sc.HasPrefix(")") {
				// slog.Info("Found right parenthesis", "text", sc.NextWord())
				if bracketBalance > 0 {
					bracketBalance--
				} else {
					break
				}
			} else if sc.HasPrefix("AND") || sc.HasPrefix("OR") {
				break
			} else {
				sc.NextWord()
			}

			sc.Consume()
		}
	}

	return sb.String()
}
