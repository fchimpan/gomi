package main

import (
	"fmt"
	"log"
	"os"

	"github.com/fchimpan/gomi/internal/gomi"
)

func main() {
	if err := gomi.Run(os.Args[1:], os.Stdin, os.Stdout, os.Stderr); err != nil {
		log.Fatal(fmt.Errorf("gomi: %w", err))
	}
}
