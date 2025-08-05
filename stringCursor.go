package main

type StringCursor struct {
	ind            int
	s              string
	t              string
	brackerBalance int
}

func NewStringCursor(s string) *StringCursor {
	return &StringCursor{
		ind: 0,
		s:   s,
		t:   "",
	}
}

func (sc *StringCursor) IsEOF() bool {
	return sc.ind >= len(sc.s)
}

func (sc *StringCursor) NextByte() byte {
	if sc.IsEOF() {
		return 0
	}
	return sc.s[sc.ind]
}

func (sc *StringCursor) isNextString(target string) bool {
	return sc.ind+len(target) <= len(sc.s) &&
		sc.s[sc.ind:sc.ind+len(target)] == target
}

func (sc *StringCursor) Consume() {
	if sc.IsEOF() {
		return
	}
	if sc.NextByte() == '(' {
		sc.brackerBalance++
	} else if sc.NextByte() == ')' {
		sc.brackerBalance--
	}
	sc.ind++
}

func (sc *StringCursor) ConsumeTill(targets []string) {
	for !sc.IsEOF() {
		if sc.brackerBalance > 0 && sc.NextByte() == ')' {
			sc.Consume()
			continue
		}
		for _, target := range targets {
			if sc.isNextString(target) {
				return
			}
		}
		sc.Consume()
	}
}

func (sc *StringCursor) Save() {
	if !sc.IsEOF() {
		sc.t += sc.s[sc.ind : sc.ind+1]
	}
}

func (sc *StringCursor) Add(c string) {
	sc.t += c
}

func (sc *StringCursor) SaveAndConsume() {
	sc.Save()
	sc.ind++
}

func (sc *StringCursor) NewString() string {
	return sc.t
}

func simplifyReq(s string) string {
	if len(s) == 0 {
		return s
	}

	sc := NewStringCursor(s)
	for !sc.IsEOF() {
		for sc.NextByte() == '(' {
			sc.SaveAndConsume()
		}

		for !sc.IsEOF() && sc.NextByte() != ' ' {
			sc.SaveAndConsume()
		}
		if !sc.IsEOF() && sc.NextByte() == ' ' {
			sc.SaveAndConsume()
		}

		for !sc.IsEOF() && sc.NextByte() != ' ' && sc.NextByte() != ')' {
			sc.SaveAndConsume()
		}
		if !sc.IsEOF() && sc.NextByte() == ' ' {
			sc.SaveAndConsume()
		}

		targets := []string{") OR ", ") AND ", "OR ", "AND "}
		sc.ConsumeTill(targets)
		for _, t := range targets {
			if sc.isNextString(t) {
				for range t {
					sc.SaveAndConsume()
				}
			}
		}
	}
	return sc.NewString()
}
