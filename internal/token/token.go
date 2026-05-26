package token

import "fmt"

type Type int

const (
	// Single-character tokens.
	LeftParen Type = iota
	RightParen
	LeftBrace
	RightBrace
	Comma
	Dot
	Minus
	Plus
	Semicolon
	Slash
	Star

	// One- or two-character tokens.
	Bang
	BangEqual
	Equal
	EqualEqual
	Greater
	GreaterEqual
	Less
	LessEqual

	// Literals.
	Identifier
	String
	Number

	// Keywords.
	And
	Class
	Else
	False
	Fun
	For
	If
	Nil
	Or
	Print
	Return
	Super
	This
	True
	Var
	While

	EOF
)

var typeNames = map[Type]string{
	LeftParen: "LEFT_PAREN", RightParen: "RIGHT_PAREN",
	LeftBrace: "LEFT_BRACE", RightBrace: "RIGHT_BRACE",
	Comma: "COMMA", Dot: "DOT", Minus: "MINUS", Plus: "PLUS",
	Semicolon: "SEMICOLON", Slash: "SLASH", Star: "STAR",

	Bang: "BANG", BangEqual: "BANG_EQUAL",
	Equal: "EQUAL", EqualEqual: "EQUAL_EQUAL",
	Greater: "GREATER", GreaterEqual: "GREATER_EQUAL",
	Less: "LESS", LessEqual: "LESS_EQUAL",

	Identifier: "IDENTIFIER", String: "STRING", Number: "NUMBER",

	And: "AND", Class: "CLASS", Else: "ELSE", False: "FALSE",
	Fun: "FUN", For: "FOR", If: "IF", Nil: "NIL", Or: "OR",
	Print: "PRINT", Return: "RETURN", Super: "SUPER", This: "THIS",
	True: "TRUE", Var: "VAR", While: "WHILE",

	EOF: "EOF",
}

func (t Type) String() string {
	if name, ok := typeNames[t]; ok {
		return name
	}
	return fmt.Sprintf("Type(%d)", int(t))
}

type Token struct {
	Type    Type
	Lexeme  string
	Literal any
	Line    int
}

func New(typ Type, lexeme string, literal any, line int) Token {
	return Token{Type: typ, Lexeme: lexeme, Literal: literal, Line: line}
}

func (t Token) String() string {
	return fmt.Sprintf("%s %s %v", t.Type, t.Lexeme, t.Literal)
}
