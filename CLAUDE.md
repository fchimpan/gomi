# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project

`gomi` is the user's Go implementation of the interpreter from *Crafting Interpreters* by Robert Nystrom (craftinginterpreters.com). The target language is **Lox**, as defined in the book. The book has two halves: Part II (**jlox**, a tree-walking interpreter, originally Java) and Part III (**clox**, a bytecode VM, originally C). Which half the user is following at any point may need to be confirmed.

The repository is in a pre-implementation state: only `go.mod`, `LICENSE`, and `README.md` exist. When scaffolding new code, prefer Go-idiomatic translations of the book's Java/C patterns (e.g., type switches over visitor-pattern boilerplate, interfaces, generics where they help) rather than literal ports.

## Collaboration mode

The user has explicitly asked Claude to act as a 伴走者 (step-by-step companion) for this project. They want to understand and implement the interpreter themselves; Claude's role is to guide, not to ship finished files. In practice:

- Explain the concept (what / why / alternatives from the book) before showing code.
- Prefer minimal, focused snippets the user adapts themselves over large generated files.
- Suggest the next small step; don't jump multiple chapters ahead.
- When translating from the book's Java/C into Go, call out the Go-idiomatic difference explicitly.
- Japanese is fine — match the language the user is writing in.

## Toolchain

- Module path: `github.com/fchimpan/gomi`
- Go version: **1.26.2** (per `go.mod`). This is unusually new — confirm the local `go` toolchain supports it before running commands; otherwise `go` will auto-download the matching toolchain on first use.

## Common commands

Standard Go workflow (no Makefile or task runner is configured):

- Build everything: `go build ./...`
- Run all tests: `go test ./...`
- Run a single test: `go test -run '^TestName$' ./path/to/pkg`
- Run a single test with verbose output: `go test -v -run '^TestName$' ./path/to/pkg`
- Vet: `go vet ./...`
- Format: `gofmt -w .` (or `go fmt ./...`)
- Tidy modules after adding/removing deps: `go mod tidy`
