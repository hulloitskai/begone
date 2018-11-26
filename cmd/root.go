package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func init() {
	// Run command-specific initializations.
	initLoginCmd()
	initEmojifyCmd()
	initRootCmd()
}

func initRootCmd() {
	cobra.EnableCommandSorting = false

	// Configure persistent flags.
	set := rootCmd.PersistentFlags()
	set.IntP("delay", "d", 225, "delay (in ms) between each message")
	set.IntP("send-fail-delay", "D", 1000, "delay (in ms) after a send fail")
	set.IntP("max-send-fails", "f", 3,
		"maximum consecutive send fails before aborting")
	set.IntP("cycles", "c", -1, "number of spam cycles (-1 for unlimited)")

	// Configure local flags.
	rootCmd.Flags().BoolP("version", "v", false, "show the version string")

	// Add subcommands.
	rootCmd.AddCommand(loginCmd)
	rootCmd.AddCommand(emojifyCmd)
	rootCmd.AddCommand(repeatCmd)
	rootCmd.AddCommand(imageCmd)
	rootCmd.AddCommand(fileCmd)
	rootCmd.AddCommand(completionCmd)
}

var (
	// Version is the program version. To be generated upon compilation by:
	//   -ldflags "-X github.com/stevenxie/begone/cmd.Version=$(VERSION)"
	//
	// It should match the output of the following command:
	//   git describe --tags | cut -c 2-
	Version string

	rootCmd = &cobra.Command{
		Use:          "begone",
		Long:         "Begone is a fully automatic spamming tool for FB Messenger.",
		Version:      Version,
		SilenceUsage: true,
		RunE:         begone,
	}
)

func begone(cmd *cobra.Command, _ []string) error {
	showVersion, err := cmd.Flags().GetBool("version")
	if err != nil {
		return err
	}

	if showVersion {
		fmt.Printf("begone version %s", cmd.Version)
		return nil
	}

	return cmd.Help()
}

// Execute runs the root program command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
