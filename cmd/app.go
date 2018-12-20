package cmd

import (
	"os"

	"gopkg.in/alecthomas/kingpin.v2"

	ess "github.com/unixpickle/essentials"
)

var (
	// Version is the program version. To be generated upon compilation by:
	//     -ldflags "-X github.com/stevenxie/begone/cmd.Version=$(VERSION)"
	//
	// It should match the output of the following command:
	//     git describe --tags | cut -c 2-
	Version = "unknown"

	app = kingpin.New(
		"begone",
		"A fully automatic spamming tool for FB Messenger.",
	).Version(Version)
)

// Exec runs the root command. It is the application entrypoint.
func Exec(args []string) {
	var err error

	switch kingpin.MustParse(app.Parse(args)) {
	// Spamming subcommands:
	case emojifyCmd.FullCommand():
		err = emojify()
	case repeatCmd.FullCommand():
		err = repeat()
	case fileCmd.FullCommand():
		err = file()
	case imageCmd.FullCommand():
		err = image()

	// Other subcommands:
	case loginCmd.FullCommand():
		err = login()

	// No subcommand:
	default:
		app.Usage(args)
		os.Exit(0)
	}

	if err != nil {
		ess.Die("Error: " + err.Error())
	}
}
