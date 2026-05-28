package scanner

import (
	"fmt"
	"strconv"

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
	s.tokens = append(s.tokens, token.New(token.EOF, "", nil, s.line))
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
	case '"':
		if err := s.string(); err != nil {
			return err
		}
	default:
		if isDigit(c) {
			if err := s.number(); err != nil {
				return err
			}
		} else if isAlpha(c) {
			if err := s.identifier(); err != nil {
				return err
			}
		} else {
			return fmt.Errorf("Unexpected character '%c' at line %d", c, s.line)
		}
	}
	return nil
}

func (s *Scanner) advance() byte {
	b := s.source[s.current]
	s.current++
	return b
}

func (s *Scanner) addToken(t token.Type) {
	s.addTokenLiteral(t, nil)
}

func (s *Scanner) addTokenLiteral(t token.Type, literal any) {
	text := s.source[s.start:s.current]
	s.tokens = append(s.tokens, token.New(t, text, literal, s.line))
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

func (s *Scanner) peekNext() byte {
	if s.current+1 >= len(s.source) {
		return 0
	}
	return s.source[s.current+1]
}

func (s *Scanner) string() error {
	for s.peek() != '"' && !s.isAtEnd() {
		if s.peek() == '\n' {
			s.line++
		}
		s.advance()
	}
	if s.isAtEnd() {
		return fmt.Errorf("Unterminated string at line %d", s.line)
	}
	s.advance() // Consume the closing "
	value := s.source[s.start+1 : s.current-1]
	s.addTokenLiteral(token.String, value)
	return nil
}

func isDigit(c byte) bool {
	return c >= '0' && c <= '9'
}

func (s *Scanner) number() error {
	for isDigit(s.peek()) {
		s.advance()
	}
	if s.peek() == '.' && isDigit(s.peekNext()) {
		s.advance() // Consume the '.'
		for isDigit(s.peek()) {
			s.advance()
		}
	}
	lexeme := s.source[s.start:s.current]
	value, err := strconv.ParseFloat(lexeme, 64)
	if err != nil {
		return fmt.Errorf("invalid number %q at line %d: %w", lexeme, s.line, err)
	}
	s.addTokenLiteral(token.Number, value)
	return nil
}

func isAlpha(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || c == '_'
}

func (s *Scanner) identifier() error {
	for isAlpha(s.peek()) || isDigit(s.peek()) {
		s.advance()
	}
	lexeme := s.source[s.start:s.current]
	t, ok := token.Keywords[lexeme]
	if !ok {
		t = token.Identifier
	}
	s.addToken(t)
	return nil
}
