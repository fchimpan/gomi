package gomi

import (
	"bufio"
	"errors"
	"fmt"
	"io"

	"github.com/fchimpan/gomi/internal/ast"
	"github.com/fchimpan/gomi/internal/parser"
	"github.com/fchimpan/gomi/internal/scanner"
)

// ExitError carries a specific process exit code alongside an error.
// Codes follow <sysexits.h>: 64 = usage, 65 = compile, 70 = runtime.
type ExitError struct {
	Code int
	Err  error
}

func (e *ExitError) Error() string { return e.Err.Error() }
func (e *ExitError) Unwrap() error { return e.Err }

func Run(args []string, stdin io.Reader, stdout, stderr io.Writer) error {
	switch len(args) {
	case 0:
		return runPrompt(stdin, stdout, stderr)
	case 1:
		return runFile(args[0], stdout, stderr)
	default:
		return &ExitError{Code: 64, Err: errors.New("usage: gomi [script]")}
	}
}

func runPrompt(stdin io.Reader, stdout, stderr io.Writer) error {
	sc := bufio.NewScanner(stdin)
	for {
		fmt.Fprint(stdout, "> ")
		if !sc.Scan() {
			break
		}
		line := sc.Text()
		if err := run(line, stdout); err != nil {
			fmt.Fprintln(stderr, err)
		}
	}
	return sc.Err()
}

func runFile(path string, stdout, stderr io.Writer) error {
	// TODO: read file & evaluate — implement after the Scanner is in place.
	return nil
}

func run(source string, out io.Writer) error {
	sc := scanner.New(source)
	tokens, err := sc.ScanTokens()
	if err != nil {
		return err
	}
	p := parser.New(tokens)
	expr, err := p.Parse()
	if err != nil {
		return err
	}
	fmt.Fprintln(out, ast.Print(expr))
	return nil
}
