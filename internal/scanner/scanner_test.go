package scanner_test

import (
	"reflect"
	"testing"

	"github.com/fchimpan/gomi/internal/scanner"
	"github.com/fchimpan/gomi/internal/token"
)

func TestScanTokens(t *testing.T) {
	tests := []struct {
		name    string
		source  string
		want    []token.Token
		wantErr bool
	}{
		// --- Edge ---
		{
			name:   "empty source emits only EOF",
			source: "",
			want: []token.Token{
				{Type: token.EOF, Lexeme: "", Line: 1},
			},
		},

		// --- Single-character tokens (no match()) ---
		{
			name:   "single-character tokens",
			source: "(){},.;-+*/",
			want: []token.Token{
				{Type: token.LeftParen, Lexeme: "(", Line: 1},
				{Type: token.RightParen, Lexeme: ")", Line: 1},
				{Type: token.LeftBrace, Lexeme: "{", Line: 1},
				{Type: token.RightBrace, Lexeme: "}", Line: 1},
				{Type: token.Comma, Lexeme: ",", Line: 1},
				{Type: token.Dot, Lexeme: ".", Line: 1},
				{Type: token.Semicolon, Lexeme: ";", Line: 1},
				{Type: token.Minus, Lexeme: "-", Line: 1},
				{Type: token.Plus, Lexeme: "+", Line: 1},
				{Type: token.Star, Lexeme: "*", Line: 1},
				{Type: token.Slash, Lexeme: "/", Line: 1},
				{Type: token.EOF, Lexeme: "", Line: 1},
			},
		},

		// --- One-or-two character operators (match()) ---
		{
			name:   "bang and equal: with and without trailing =",
			source: "! != = ==",
			want: []token.Token{
				{Type: token.Bang, Lexeme: "!", Line: 1},
				{Type: token.BangEqual, Lexeme: "!=", Line: 1},
				{Type: token.Equal, Lexeme: "=", Line: 1},
				{Type: token.EqualEqual, Lexeme: "==", Line: 1},
				{Type: token.EOF, Lexeme: "", Line: 1},
			},
		},
		{
			name:   "less and greater: with and without trailing =",
			source: "< <= > >=",
			want: []token.Token{
				{Type: token.Less, Lexeme: "<", Line: 1},
				{Type: token.LessEqual, Lexeme: "<=", Line: 1},
				{Type: token.Greater, Lexeme: ">", Line: 1},
				{Type: token.GreaterEqual, Lexeme: ">=", Line: 1},
				{Type: token.EOF, Lexeme: "", Line: 1},
			},
		},

		// --- Literals ---
		{
			name:   "string literal: lexeme keeps quotes, literal drops them",
			source: `"hi"`,
			want: []token.Token{
				{Type: token.String, Lexeme: `"hi"`, Literal: "hi", Line: 1},
				{Type: token.EOF, Lexeme: "", Line: 1},
			},
		},
		{
			name:   "number literal becomes float64",
			source: "1.5",
			want: []token.Token{
				{Type: token.Number, Lexeme: "1.5", Literal: 1.5, Line: 1},
				{Type: token.EOF, Lexeme: "", Line: 1},
			},
		},
		{
			name:   "trailing dot is not part of number (peekNext justifies itself)",
			source: "1.",
			want: []token.Token{
				{Type: token.Number, Lexeme: "1", Literal: 1.0, Line: 1},
				{Type: token.Dot, Lexeme: ".", Line: 1},
				{Type: token.EOF, Lexeme: "", Line: 1},
			},
		},

		// --- Identifiers and keywords ---
		{
			name:   "keyword hits Keywords map, identifier falls back",
			source: "var foo",
			want: []token.Token{
				{Type: token.Var, Lexeme: "var", Line: 1},
				{Type: token.Identifier, Lexeme: "foo", Line: 1},
				{Type: token.EOF, Lexeme: "", Line: 1},
			},
		},
		{
			name:   "identifier allows underscore start and digits after first letter",
			source: "_foo bar123 _",
			want: []token.Token{
				{Type: token.Identifier, Lexeme: "_foo", Line: 1},
				{Type: token.Identifier, Lexeme: "bar123", Line: 1},
				{Type: token.Identifier, Lexeme: "_", Line: 1},
				{Type: token.EOF, Lexeme: "", Line: 1},
			},
		},
		{
			name:   "multiple keywords",
			source: "if while for return",
			want: []token.Token{
				{Type: token.If, Lexeme: "if", Line: 1},
				{Type: token.While, Lexeme: "while", Line: 1},
				{Type: token.For, Lexeme: "for", Line: 1},
				{Type: token.Return, Lexeme: "return", Line: 1},
				{Type: token.EOF, Lexeme: "", Line: 1},
			},
		},

		// --- Whitespace and comments ---
		{
			name:   "comment is skipped, newline advances line",
			source: "// hello\nfoo",
			want: []token.Token{
				{Type: token.Identifier, Lexeme: "foo", Line: 2},
				{Type: token.EOF, Lexeme: "", Line: 2},
			},
		},
		{
			name:   "whitespace is skipped, newlines advance line",
			source: " \t\n  \nfoo",
			want: []token.Token{
				{Type: token.Identifier, Lexeme: "foo", Line: 3},
				{Type: token.EOF, Lexeme: "", Line: 3},
			},
		},

		// --- Errors ---
		{
			name:    "unterminated string returns error",
			source:  `"abc`,
			wantErr: true,
		},
		{
			name:    "unexpected character returns error",
			source:  "@",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			s := scanner.New(tt.source)
			got, err := s.ScanTokens()

			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error, got nil (tokens=%v)", got)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("tokens mismatch\n got:  %v\n want: %v", got, tt.want)
			}
		})
	}
}
