package main

import (
	"log/slog"
	"strings"
)

type tokenType string

const (
	tokenSTR      tokenType = "STR"
	tokenAND      tokenType = "AND"
	tokenOR       tokenType = "OR"
	tokenLeftBra  tokenType = "LEFT_BRACKET"
	tokenRightBra tokenType = "RIGHT_BRACKET"
	tokenEOF      tokenType = "EOF"
)

type token struct {
	typ   tokenType
	value string
}

func Resolve(s string, mp map[string]struct{}) bool {
	if s == "" {
		return true
	}
	l := NewLexer(s)
	n := parseExpression(l)
	return n.Resolve(mp)
}

type node struct {
	token token
	left  *node
	right *node
}

func (n *node) Resolve(mp map[string]struct{}) bool {

	switch n.token.typ {
	case tokenSTR:
		_, exists := mp[n.token.value]
		return exists
	case tokenOR:
		return n.left.Resolve(mp) || n.right.Resolve(mp)
	case tokenAND:
		return n.left.Resolve(mp) && n.right.Resolve(mp)
	}
	panic("Unexpected token type: " + string(n.token.typ))
}

func parseExpression(l *lexer) *node {
	n := parseTerm(l)
	for t := l.NextToken(); t.typ == tokenAND || t.typ == tokenOR; t = l.NextToken() {
		l.Consume()
		tmp := &node{
			token: t,
		}
		tmp.left = n
		n = tmp
		n.right = parseTerm(l)
	}

	return n
}

func parseTerm(l *lexer) *node {
	tok := l.NextToken()
	switch tok.typ {
	case tokenLeftBra:
		l.Consume()
		node := parseExpression(l)
		if l.NextToken().typ != tokenRightBra {
			slog.Error("Expected right bracket", "text", l.text)
		} else {
			l.Consume()
		}
		return node
	case tokenSTR:
		l.Consume()
		return &node{token: tok}
	default:
		panic("Unexpected token: " + string(tok.typ))
	}
}

type lexer struct {
	pos  int
	ind  int
	text string
	last token
}

func NewLexer(s string) *lexer {
	return &lexer{
		ind:  0,
		pos:  0,
		text: s,
	}
}

func (l *lexer) Consume() {
	l.pos = l.ind
}

func (l *lexer) NextToken() token {
	l.ind = l.pos
	for l.ind < len(l.text) && l.text[l.ind] == ' ' {
		l.ind++
	}
	if l.ind == len(l.text) {
		return token{typ: tokenEOF}
	}
	if strings.HasPrefix(l.text[l.ind:], "AND") {
		l.ind += 3
		return token{typ: tokenAND, value: "AND"}
	} else if strings.HasPrefix(l.text[l.ind:], "OR") {
		l.ind += 2
		return token{typ: tokenOR, value: "OR"}
	} else if strings.HasPrefix(l.text[l.ind:], "(") {
		l.ind++
		return token{typ: tokenLeftBra, value: "("}
	} else if strings.HasPrefix(l.text[l.ind:], ")") {
		l.ind++
		return token{typ: tokenRightBra, value: ")"}
	} else {
		start := l.ind
		for l.ind < len(l.text) && l.text[l.ind] != ' ' {
			l.ind++
		}
		l.ind++
		for l.ind < len(l.text) && l.text[l.ind] != ' ' && l.text[l.ind] != '(' && l.text[l.ind] != ')' {
			l.ind++
		}
		val := strings.TrimSpace(l.text[start:l.ind])
		if val != "" {
			return token{typ: tokenSTR, value: val}
		}
		panic("Unexpected token")
	}
}
