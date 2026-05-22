package gomi

import (
	"bufio"
	"errors"
	"fmt"
	"io"
)

// ExitError carries a specific process exit code alongside an error.
// Codes follow <sysexits.h>, matching the jlox convention from
// Crafting Interpreters (64 = usage, 65 = compile, 70 = runtime).
type ExitError struct {
	Code int
	Err  error
}

func (e *ExitError) Error() string { return e.Err.Error() }
func (e *ExitError) Unwrap() error { return e.Err }

// Run is the testable entry point. main.go is a thin wrapper that wires
// os.Args / os.Stdin / os.Stdout / os.Stderr into here.
//
// args excludes the program name (pass os.Args[1:]).
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
	// TODO: REPL — implement after the Scanner is in place.

	sc := bufio.NewScanner(stdin)
	for {
		fmt.Fprint(stdout, "> ")
		if !sc.Scan() {
			break
		}
		line := sc.Text()
		if _, err := fmt.Fprintln(stdout, "TODO: evaluate", line); err != nil {
			fmt.Fprintln(stderr, "gomi: error writing to stdout:", err)
		}

	}
	return sc.Err()
}

func runFile(path string, stdout, stderr io.Writer) error {
	// TODO: read file & evaluate — implement after the Scanner is in place.
	return nil
}
