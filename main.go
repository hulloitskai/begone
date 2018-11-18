package main

import (
	"fmt"
	"log"

	"github.com/stevenxie/begone/loader"
	"github.com/stevenxie/begone/messenger"
	"github.com/urfave/cli"
)

func init() {
	log.SetFlags(0)
}

func main() {
	l, err := loader.NewLoader(Namespace)
	if err != nil {
		log.Fatalf("Failed to initialize driver loader: %v", err)
	}

	// Configure name, usage, metadata, etc.
	app := cli.NewApp()
	app.Name = Name
	app.Usage = "fully automatic spamming utility for FB Messenger"
	app.UsageText = `begone [--gen="<generator [args]>"] [--reps <repeats>] [convoID]
   begone [--help | -h]
   begone command [command options]

For command usage details, please see:
   begone help [command]

Available generators:
   emoji (default)
   repeat <string to repeat>`

	app.Version = Version
	app.HideVersion = true

	// Primary app action.
	app.Action = begoneAction(l)
	app.Flags = begoneFlags

	// App subcommands.
	app.Commands = []cli.Command{
		{
			Name:  "login",
			Usage: "saves FB Messenger login credentials",
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
			Action: func(ctx *cli.Context) error {
				var (
					cfg messenger.Config
					err error
				)
				cfg.Namespace = Namespace

				// If the remove flag is toggled, remove the config file (if it exists)
				// and exit.
				if ctx.Bool("clear") {
					fmt.Println("Removing config files containing user credentials...")

					removed, err := cfg.RemoveFile()
					if err != nil {
						log.Fatal(err)
					}

					if len(removed) == 0 {
						fmt.Println("No config files were found.")
					} else {
						fmt.Println("Done. The following files were removed:")
						for _, path := range removed {
							fmt.Printf("\t%s\n", path)
						}
					}

					return nil
				}

				// User wants to create a config file:
				if cfg.User, err = queryUser(); err != nil {
					log.Fatalf("Error while querying for username: %v", err)
				}
				if !ctx.Bool("user-only") {
					if cfg.Pass, err = queryPass(); err != nil {
						log.Fatalf("Error while querying for password: %v", err)
					}
				}

				path, err := cfg.Save()
				if err != nil {
					log.Fatalf("Error while saving config: %v", err)
				}

				fmt.Printf("Credentials saved to '%s'.\n", path)
				return nil
			},
		},
		{
			Name:  "clean",
			Usage: "removes temporary installation artifacts",
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "all, a",
					Usage: "remove all downloads, including the ChromeDriver executable",
				},
			},
			Action: func(ctx *cli.Context) error {
				all := ctx.Bool("all")
				if all {
					fmt.Println("Cleaning up all files, including the driver...")
				} else {
					fmt.Println("Cleaning up temporary files from failed downloads...")
				}

				removed, err := l.Clean(all)
				if err != nil {
					log.Fatalf("Error while cleaning up files: %v", err)
				}

				if len(removed) == 0 {
					fmt.Println("No matching files found; your system is clean.")
					return nil
				}

				fmt.Println("Done. The following files were removed:")
				for _, path := range removed {
					fmt.Printf("\t%s\n", path)
				}

				return nil
			},
		},
		{
			Name:  "where",
			Usage: "print out driver installation path",
			Action: func(ctx *cli.Context) error {
				fmt.Println(l.DriverPath)
				return nil
			},
		}, {
			Name:  "setup",
			Usage: "(for debug use) manually installs Selenium driver",
			Action: func(ctx *cli.Context) error {
				fmt.Println("Installing driver (if not yet installed)...")

				if err := l.Ensure(); err != nil {
					log.Fatalf("Driver installation error: %v", err)
				}

				fmt.Println("Selenium ChromeDriver is installed and ready-for-use.")
				return nil
			},
		},
	}

	app.RunAndExitOnError()
}
