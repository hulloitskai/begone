package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/stevenxie/begone/internal/config"
	"github.com/stevenxie/begone/internal/interact"
	ess "github.com/unixpickle/essentials"
)

func initLoginCmd() {
	// Configure flags:
	loginCmd.LocalFlags().BoolP("user-only", "u", false,
		"only save username; password will be requested every time")

	// Add clearLoginCmd as a subcommand under loginCmd.
	loginCmd.AddCommand(clearLoginCmd)
}

var (
	loginCmd = &cobra.Command{
		Use:   "login",
		Short: "Save FB Messenger login credentials",
		Long: "Save FB Messenger login credentials into a JSON file. \nThe " +
			"password will be obfuscated to prevent accidental exposure.",
		Args: withUsage(cobra.NoArgs),
		RunE: login,
	}

	clearLoginCmd = &cobra.Command{
		Use:   "clear",
		Short: "remove saved login credentials",
		Args:  withUsage(cobra.NoArgs),
		RunE:  clearLogin,
	}
)

func login(cmd *cobra.Command, _ []string) error {
	userOnly, err := cmd.LocalFlags().GetBool("user-only")
	if err != nil {
		return err
	}

	var (
		p   = interact.NewPrompter()
		cfg = new(config.Config)
	)
	if err = p.QueryMissing(cfg, userOnly); err != nil {
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

func clearLogin(_ *cobra.Command, _ []string) error {
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
