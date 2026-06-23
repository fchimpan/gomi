package parser

import (
	"fmt"
	"slices"

	"github.com/fchimpan/gomi/internal/ast"
	"github.com/fchimpan/gomi/internal/token"
)

type Parser struct {
	tokens  []token.Token
	current int
}

func New(tokens []token.Token) *Parser {
	return &Parser{
		tokens:  tokens,
		current: 0,
	}
}

// Parse is the public entry point: it parses a single expression and
// converts any parser-internal panic back into an ordinary error.
//
// The results are named so the deferred recover can assign err.
func (p *Parser) Parse() (expr ast.Expr, err error) {
	defer func() {
		if r := recover(); r != nil {
			if parseErr, ok := r.(*parseError); ok {
				err = parseErr
			} else {
				panic(r) // re-panic if it's not a parse error
			}
		}
	}()
	return p.expression(), nil
}

// --- Grammar rules: one method per rule, lowest precedence first. ---

// expression → equality
func (p *Parser) expression() ast.Expr {
	return p.equality()
}

// equality → comparison ( ( "!=" | "==" ) comparison )*
func (p *Parser) equality() ast.Expr {
	expr := p.comparison()
	for p.match(token.BangEqual, token.EqualEqual) {
		op := p.previous()
		right := p.comparison()
		expr = &ast.Binary{Left: expr, Operator: op, Right: right}
	}
	return expr
}

// comparison → term ( ( ">" | ">=" | "<" | "<=" ) term )*
func (p *Parser) comparison() ast.Expr {
	expr := p.term()
	for p.match(token.Greater, token.GreaterEqual, token.Less, token.LessEqual) {
		op := p.previous()
		right := p.term()
		expr = &ast.Binary{Left: expr, Operator: op, Right: right}
	}
	return expr
}

// term → factor ( ( "-" | "+" ) factor )*
func (p *Parser) term() ast.Expr {
	expr := p.factor()
	for p.match(token.Minus, token.Plus) {
		op := p.previous()
		right := p.factor()
		expr = &ast.Binary{Left: expr, Operator: op, Right: right}
	}
	return expr
}

// factor → unary ( ( "/" | "*" ) unary )*
func (p *Parser) factor() ast.Expr {
	expr := p.unary()
	for p.match(token.Slash, token.Star) {
		op := p.previous()
		right := p.unary()
		expr = &ast.Binary{Left: expr, Operator: op, Right: right}
	}
	return expr
}

// unary → ( "!" | "-" ) unary | primary
func (p *Parser) unary() ast.Expr {
	if p.match(token.Bang, token.Minus) {
		op := p.previous()
		right := p.unary()
		return &ast.Unary{Operator: op, Right: right}
	}
	return p.primary()
}

// primary → NUMBER | STRING | "true" | "false" | "nil" | "(" expression ")"
func (p *Parser) primary() ast.Expr {
	switch p.peek().Type {
	case token.Number, token.String:
		return &ast.Literal{Value: p.advance().Literal}
	case token.True:
		p.advance()
		return &ast.Literal{Value: true}
	case token.False:
		p.advance()
		return &ast.Literal{Value: false}
	case token.Nil:
		p.advance()
		return &ast.Literal{Value: nil}
	case token.LeftParen:
		p.advance()
		expr := p.expression()
		p.consume(token.RightParen, "Expect ')' after expression.")
		return &ast.Grouping{Expression: expr}
	default:
		panic(&parseError{
			token: p.peek(),
			msg:   "Expect expression.",
		})
	}
}

// --- Token cursor helpers. ---

// consume checks whether the current token is of the expected type.
// If so, it consumes and returns that token. Otherwise it reports a
// parse error at the current token (panic with a parser-internal error
// value; recovered in Parse).
func (p *Parser) consume(t token.Type, msg string) token.Token {
	if p.check(t) {
		return p.advance()
	}
	panic(&parseError{
		token: p.peek(),
		msg:   msg,
	})
}

// match reports whether the current token has any of the given types.
// If so, it consumes the token and returns true; otherwise it leaves
// the current token alone and returns false.
func (p *Parser) match(types ...token.Type) bool {
	if slices.ContainsFunc(types, func(t token.Type) bool {
		return p.check(t)
	}) {
		p.advance()
		return true
	}
	return false
}

// check reports whether the current token is of type t. Unlike match,
// it never consumes the token, it only looks at it.
func (p *Parser) check(t token.Type) bool {
	if p.isAtEnd() {
		return false
	}
	return p.peek().Type == t
}

// advance consumes the current token and returns it. It never moves
// past the EOF sentinel.
func (p *Parser) advance() token.Token {
	if p.isAtEnd() {
		return p.peek()
	}
	p.current++
	return p.previous()
}

// isAtEnd reports whether the parser has run out of tokens to parse,
// i.e. the current token is the EOF sentinel.
func (p *Parser) isAtEnd() bool {
	return p.peek().Type == token.EOF
}

// peek returns the current token we have yet to consume.
func (p *Parser) peek() token.Token {
	return p.tokens[p.current]
}

// previous returns the most recently consumed token.
func (p *Parser) previous() token.Token {
	return p.tokens[p.current-1]
}

type parseError struct {
	token token.Token
	msg   string
}

func (e *parseError) Error() string {
	return fmt.Sprintf("[line %d] Error at '%s': %s", e.token.Line, e.token.Lexeme, e.msg)
}
