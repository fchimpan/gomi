package ast_test

import (
	"testing"

	"github.com/fchimpan/gomi/internal/ast"
	"github.com/fchimpan/gomi/internal/token"
)

// Compile-time assertions: every concrete expression type must satisfy
// the sealed Expr interface. If someone removes a type's exprNode()
// method, the package fails to build and the regression is caught here
// rather than at a distant call site.
var (
	_ ast.Expr = (*ast.Binary)(nil)
	_ ast.Expr = (*ast.Grouping)(nil)
	_ ast.Expr = (*ast.Literal)(nil)
	_ ast.Expr = (*ast.Unary)(nil)
)

// TestTypeSwitchDispatch verifies that operations over Expr can be
// written as a type switch on concrete pointer types. Every future
// operation (printer, evaluator, resolver) relies on this dispatch
// shape, so a silent break here would propagate widely. Building a
// small tree also exercises composition through the Expr interface.
func TestTypeSwitchDispatch(t *testing.T) {
	t.Parallel()

	plus := token.New(token.Plus, "+", nil, 1)
	star := token.New(token.Star, "*", nil, 1)
	minus := token.New(token.Minus, "-", nil, 1)

	// Build the tree for: 1 + (-2 * 3)
	tree := ast.Expr(&ast.Binary{
		Left:     &ast.Literal{Value: 1.0},
		Operator: plus,
		Right: &ast.Grouping{
			Expression: &ast.Binary{
				Left: &ast.Unary{
					Operator: minus,
					Right:    &ast.Literal{Value: 2.0},
				},
				Operator: star,
				Right:    &ast.Literal{Value: 3.0},
			},
		},
	})

	kind := func(e ast.Expr) string {
		switch e.(type) {
		case *ast.Binary:
			return "binary"
		case *ast.Grouping:
			return "grouping"
		case *ast.Literal:
			return "literal"
		case *ast.Unary:
			return "unary"
		default:
			return "unknown"
		}
	}

	root, ok := tree.(*ast.Binary)
	if !ok {
		t.Fatalf("root: want *ast.Binary, got %T", tree)
	}
	if got := kind(root.Left); got != "literal" {
		t.Errorf("root.Left: got %q, want %q", got, "literal")
	}
	if got := kind(root.Right); got != "grouping" {
		t.Errorf("root.Right: got %q, want %q", got, "grouping")
	}

	grp := root.Right.(*ast.Grouping).Expression.(*ast.Binary)
	if got := kind(grp.Left); got != "unary" {
		t.Errorf("grouped.Left: got %q, want %q", got, "unary")
	}
	if got := kind(grp.Right); got != "literal" {
		t.Errorf("grouped.Right: got %q, want %q", got, "literal")
	}
}
