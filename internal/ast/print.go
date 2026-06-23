package ast

import (
	"fmt"
	"strings"
)

// Print renders an expression tree as a parenthesized, Lisp-like string,
// as a debugging aid for inspecting parser output. The operator (or the
// word "group" for a parenthesised sub-expression) comes first, followed
// by its operands:
//
//	1 + 2 * 3     ->  (+ 1 (* 2 3))
//	-2 * 3        ->  (* (- 2) 3)
//	(1 + 2) * 3   ->  (* (group (+ 1 2)) 3)
//
// Dispatch is a type switch on the concrete node types, the same shape
// every later pass over the tree (the evaluator, the resolver) will use.
func Print(e Expr) string {
	switch e := e.(type) {
	case *Literal:
		if e.Value == nil {
			return "nil"
		}
		return fmt.Sprintf("%v", e.Value)
	case *Grouping:
		return parenthesize("group", e.Expression)
	case *Unary:
		return parenthesize(e.Operator.Lexeme, e.Right)
	case *Binary:
		return parenthesize(e.Operator.Lexeme, e.Left, e.Right)
	default:
		return fmt.Sprintf("<unknown %T>", e)
	}
}

// parenthesize writes "(name operand operand ...)", recursing into each
// operand via Print.
func parenthesize(name string, exprs ...Expr) string {
	var b strings.Builder
	b.WriteByte('(')
	b.WriteString(name)
	for _, e := range exprs {
		b.WriteByte(' ')
		b.WriteString(Print(e))
	}
	b.WriteByte(')')
	return b.String()
}
