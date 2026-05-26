package scanner

import (
	"fmt"

	"github.com/fchimpan/gomi/internal/token"
)

type Scanner struct {
	source  string
	tokens  []token.Token
	start   int
	current int
	line    int
}

func New(source string) *Scanner {
	return &Scanner{
		source: source,
		line:   1,
	}
}

func (s *Scanner) ScanTokens() ([]token.Token, error) {
	for !s.isAtEnd() {
		s.start = s.current
		if err := s.scanToken(); err != nil {
			return nil, err
		}
	}
	return s.tokens, nil
}

func (s *Scanner) isAtEnd() bool {
	return s.current >= len(s.source)
}

func (s *Scanner) scanToken() error {
	c := s.advance()
	switch c {
	case '(':
		s.addToken(token.LeftParen)
	case ')':
		s.addToken(token.RightParen)
	case '{':
		s.addToken(token.LeftBrace)
	case '}':
		s.addToken(token.RightBrace)
	case ',':
		s.addToken(token.Comma)
	case '.':
		s.addToken(token.Dot)
	case '-':
		s.addToken(token.Minus)
	case '+':
		s.addToken(token.Plus)
	case ';':
		s.addToken(token.Semicolon)
	case '*':
		s.addToken(token.Star)
	case '!':
		if s.match('=') {
			s.addToken(token.BangEqual)
		} else {
			s.addToken(token.Bang)
		}
	case '=':
		if s.match('=') {
			s.addToken(token.EqualEqual)
		} else {
			s.addToken(token.Equal)
		}
	case '<':
		if s.match('=') {
			s.addToken(token.LessEqual)
		} else {
			s.addToken(token.Less)
		}
	case '>':
		if s.match('=') {
			s.addToken(token.GreaterEqual)
		} else {
			s.addToken(token.Greater)
		}
	case '/':
		if s.match('/') {
			for s.peek() != '\n' && !s.isAtEnd() {
				s.advance()
			}
		} else {
			s.addToken(token.Slash)
		}
	case ' ', '\r', '\t':
	case '\n':
		s.line++
	default:
		fmt.Printf("Unexpected character: %c at line %d\n", c, s.line)
	}
	return nil
}

func (s *Scanner) advance() byte {
	b := s.source[s.current]
	s.current++
	return b
}

func (s *Scanner) addToken(t token.Type) {
	text := s.source[s.start:s.current]
	s.tokens = append(s.tokens, token.Token{Type: t, Lexeme: text, Line: s.line})
}

func (s *Scanner) match(expected byte) bool {
	if s.isAtEnd() {
		return false
	}
	if s.source[s.current] != expected {
		return false
	}
	s.current++
	return true
}

func (s *Scanner) peek() byte {
	if s.isAtEnd() {
		return 0
	}
	return s.source[s.current]
}
