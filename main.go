package main

import (
	"os"

	"github.com/stevenxie/begone/cmd"
)

func main() {
	// Execute root command, with args that exclude the caller name.
	cmd.Exec(os.Args[1:])
}
