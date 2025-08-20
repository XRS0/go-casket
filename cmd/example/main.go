package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/XRS0/go-casket/pkg/unpack"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: example <archive> <destDir> [password]")
		os.Exit(2)
	}
	archive := os.Args[1]
	dest := os.Args[2]

	opts := unpack.DefaultOptions()
	opts.Timeout = 0

	// Password via optional 3rd arg
	if len(os.Args) >= 4 {
		opts.Password = os.Args[3]
	}

	if !unpack.HasSupportedTool() {
		log.Fatal("No 7z/7zz/unrar found in PATH. Install p7zip or unrar.")
	}

	if err := unpack.Unpack(archive, dest, opts, nil, nil); err != nil {
		log.Fatalf("Unpack failed: %v", err)
	}
	fmt.Println("Done at", time.Now().Format(time.RFC3339))
}
