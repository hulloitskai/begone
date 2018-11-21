package cmd

import (
	"fmt"

	"github.com/stevenxie/begone/config"
	ess "github.com/unixpickle/essentials"
	"github.com/urfave/cli"
)

// LoginCmd saves FB Messenger login credentials.
var LoginCmd = cli.Command{
	Name:  "login",
	Usage: "save FB Messenger login credentials",
	Flags: []cli.Flag{
		cli.BoolFlag{
			Name:  "user-only, u",
			Usage: "don't save the password (which will be requested every time)",
		},
		cli.BoolFlag{
			Name:  "clear, c",
			Usage: "remove saved login credentials",
		},
	},
	Action: login,
}

func login(ctx *cli.Context) error {
	// Remove config file.
	if ctx.Bool("clear") {
		fmt.Println("Removing config file with saved credentials...")

		removed, err := config.Remove()
		if err != nil {
			ess.Die(err)
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

	// Create config file.
	var (
		cfg = new(config.Config)
		err error
	)
	if cfg.Username, err = queryUsername(); err != nil {
		ess.Die("Error while getting username:", err)
	}
	if !ctx.Bool("user-only") { // conditionally query for password
		if cfg.Password, err = queryPassword(); err != nil {
			ess.Die("Error while getting password:", err)
		}
	}

	path, err := config.Save(cfg)
	if err != nil {
		ess.Die("Error while saving file:", err)
	}

	fmt.Printf("Credentials saved to '%s'.\n", path)
	return nil
}
