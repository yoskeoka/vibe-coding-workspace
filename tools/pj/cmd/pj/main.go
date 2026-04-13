package main

import (
	"fmt"
	"os"

	"github.com/yoskeoka/vibe-coding-workspace/tools/pj/internal/pj"
)

func main() {
	if err := pj.Run(os.Args[1:], os.Stdout, os.Stderr); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}
