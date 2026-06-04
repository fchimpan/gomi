package ast

import "github.com/fchimpan/gomi/internal/token"

// Expr is the marker interface for every expression node in the AST.
//
// The unexported exprNode() method seals the interface: only types defined
// in this package can satisfy Expr. Operations over Expr are written as a
// type switch on the concrete pointer types (*Binary, *Grouping, *Literal,
// *Unary).
type Expr interface {
	exprNode()
}

// Binary holds a left/right operand pair joined by an operator token.
// Example source: `1 + 2`, `a == b`, `x < y`.
type Binary struct {
	Left     Expr
	Operator token.Token
	Right    Expr
}

func (*Binary) exprNode() {}

// Grouping wraps a parenthesised sub-expression. Example: `(1 + 2)`.
// The parentheses themselves are not stored; their structural meaning is
// preserved by this node existing in the tree.
type Grouping struct {
	Expression Expr
}

func (*Grouping) exprNode() {}

// Literal is a leaf node holding a literal value.
//
// Value is any because literals span multiple Go types:
//   - number  -> float64
//   - string  -> string
//   - true    -> bool
//   - false   -> bool
//   - nil     -> untyped nil
type Literal struct {
	Value any
}

func (*Literal) exprNode() {}

// Unary is a prefix-operator expression. Example: `-x`, `!ready`.
type Unary struct {
	Operator token.Token
	Right    Expr
}

func (*Unary) exprNode() {}
