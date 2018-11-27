package cmd

import (
	"os"

	"gopkg.in/alecthomas/kingpin.v2"

	ess "github.com/unixpickle/essentials"
)

var (
	// Version is the program version. To be generated upon compilation by:
	//   -ldflags "-X github.com/stevenxie/begone/cmd.Version=$(VERSION)"
	//
	// It should match the output of the following command:
	//   git describe --tags | cut -c 2-
	Version string

	// app (program entrypoint).
	app = kingpin.New(
		"begone",
		"A fully automatic spamming tool for FB Messenger.",
	).Version(Version)
)

func registerAppFlags(app *kingpin.Application) {
	opts.Debug = app.Flag("debug", "Enable debug mode.").Default("false").Bool()

	opts.Delay = app.Flag("delay", "Delay (in ms) between messages.").
		Short('d').Default("225").Int()
	opts.Cycles = app.Flag("cycles", "Number of spam cycles (-1 for infinite).").
		Short('c').Default("-1").Int()

	opts.SendFailDelay = app.Flag("send-fail-delay",
		"Delay (in ms) after a send fail.").Short('D').Default("1000").Int()
	opts.MaxSendFails = app.Flag("max-send-fails",
		"Max consecutive send fails before aborting.").Short('f').Default("3").Int()
}

var opts struct {
	Debug                       *bool
	Delay, Cycles               *int
	SendFailDelay, MaxSendFails *int
}

// Exec runs the root command. It is the application entrypoint.
func Exec(args []string) {
	var err error

	switch kingpin.MustParse(app.Parse(args)) {
	// "Attack" subcommands:
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
