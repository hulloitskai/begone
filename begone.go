package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/stevenxie/begone/strgen"

	spin "github.com/briandowns/spinner"
	"github.com/stevenxie/begone/loader"
	"github.com/stevenxie/begone/messenger"
	"github.com/urfave/cli"
)

// begoneAction creates a cli.ActionFunc that spams people on FB messenger.
func begoneAction(l *loader.Loader) cli.ActionFunc {
	return func(ctx *cli.Context) error {
		cfg, err := loadConfig()
		if err != nil {
			log.Fatalf("Failed to load config from file: %v", err)
		}

		args := ctx.Args()
		if len(args) > 1 {
			log.Print("Invalid command invokation.\n\n")
			cli.ShowAppHelpAndExit(ctx, 1)
		}

		if len(args) == 1 {
			cfg.ConvoID = args[0]
		}
		if err = queryMissing(cfg); err != nil {
			log.Fatalf("Error while querying for missing configuration details: %v",
				err)
		}

		gen, err := buildGenerator(ctx.String("gen"))
		if err != nil {
			log.Fatalf("Couldn't build generator: %v", err)
		}

		// Show a CLI spinner (along with start-up text).
		fmt.Println("Starting the great Selenium engine...")

		e := messenger.NewEngine(cfg, l)

		spinner := makeAndStartSpinner("Logging into FB Messenger...")
		spinner.FinalMSG = "Logged into FB Messenger."
		if err = e.Login(); err != nil {
			log.Fatalf("Error while logging into FB Messenger: %v", err)
		}
		spinner.Stop()
		fmt.Println()

		spinner = makeAndStartSpinner("Performing attack routine...")
		spinner.FinalMSG = "Attack complete."
		if err = e.Spam(gen, ctx.Int("reps")); err != nil {
			log.Fatalf("Error during spamming routine: %v", err)
		}
		spinner.Stop()
		fmt.Println()

		return nil
	}
}

var begoneFlags = []cli.Flag{
	cli.StringFlag{
		Name:  "gen, g",
		Usage: "specify a generator and its arguments\n\t(default: emoji)",
		Value: "repeat begone",
	},
	cli.IntFlag{
		Name:  "reps, r",
		Usage: "limit the number of iterations to generate\n\t(-1 == infinite)",
		Value: -1,
	},
}

// makeAndStartSpinner makes and starts a cyan-colored spinner.
func makeAndStartSpinner(text string) *spin.Spinner {
	spinner := spin.New(spin.CharSets[14], 100*time.Millisecond)
	spinner.Color("cyan")
	spinner.Suffix = " " + text
	spinner.Start()
	return spinner
}

// buildGenerator builds a generator from a specifier string.
func buildGenerator(specifier string) (strgen.Generator, error) {
	if specifier == "" {
		return nil, fmt.Errorf("begone: generator name not given")
	}

	var (
		args  = strings.Split(specifier, " ")
		name  = args[0]
		rargs = args[1:]
	)
	switch name {
	case "repeat":
		if len(rargs) == 0 {
			return nil, fmt.Errorf("begone: repeat requires arguments (usage: " +
				"repeat <string to repeat>)")
		}
		return strgen.NewRepeater(strings.Join(rargs, " ")), nil

	case "emoji":
		var mode string
		if len(rargs) >= 1 {
			mode = rargs[0]
		}

		return strgen.NewEmojifier(mode)
	}

	return nil, fmt.Errorf("begone: no generator found with the name '%s'", name)
}
