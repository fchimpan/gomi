package parser_test

import (
	"testing"

	"github.com/fchimpan/gomi/internal/ast"
	"github.com/fchimpan/gomi/internal/parser"
	"github.com/fchimpan/gomi/internal/scanner"
)

func TestParse_MultiplicationBindsTighterThanAddition(t *testing.T) {
	tokens, err := scanner.New("1 + 2 * 3").ScanTokens()
	if err != nil {
		t.Fatalf("scan: %v", err)
	}

	expr, err := parser.New(tokens).Parse()
	if err != nil {
		t.Fatalf("parse: %v", err)
	}

	got := ast.Print(expr)
	want := "(+ 1 (* 2 3))"
	if got != want {
		t.Errorf("Parse = %q, want %q", got, want)
	}
}
