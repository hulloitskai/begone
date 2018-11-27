package cmd

import (
	"fmt"

	"gopkg.in/alecthomas/kingpin.v2"

	"github.com/stevenxie/begone/internal/config"
	"github.com/stevenxie/begone/internal/interact"
	ess "github.com/unixpickle/essentials"
)

func registerLoginCmd(app *kingpin.Application) {
	loginCmd = app.Command(
		"login",
		"Save FB Messenger login credentials (obfuscates password).",
	)

	// Register flags.
	loginOpts.UserOnly = loginCmd.Flag("user-only", "Only save the username.").
		Short('u').Default("false").Bool()

	loginOpts.Clear = loginCmd.Flag("clear", "Remove saved login credentials.").
		Short('U').Default("false").Bool()
}

var (
	loginCmd *kingpin.CmdClause

	loginOpts struct {
		UserOnly, Clear *bool
	}
)

func login() error {
	if *loginOpts.Clear {
		return clearLogin()
	}

	var (
		p   = interact.NewPrompter()
		cfg = new(config.Config)
	)
	if err := p.QueryMissing(cfg, *loginOpts.UserOnly); err != nil {
		return err
	}
	p.Println()

	path, err := config.Save(cfg)
	if err != nil {
		return ess.AddCtx("saving file", err)
	}

	p.Printf("Credentials saved to '%s'.\n", path)
	return nil
}

func clearLogin() error {
	fmt.Println("Removing config file with saved credentials...")

	removed, err := config.Remove()
	if err != nil {
		return err
	}

	if len(removed) == 0 {
		fmt.Println("Done; no config files were found.")
	} else {
		fmt.Println("Done; the following files were removed:")
		for _, path := range removed {
			fmt.Printf("\t%s\n", path)
		}
	}

	return nil
}
