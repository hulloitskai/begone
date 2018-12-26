package main

import (
	"os"

	"github.com/stevenxie/begone/cmd"
)

func main() {
	// Execute app, with args that exclude the caller name.
	cmd.Exec(os.Args[1:])
}
