package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var completionCmd = &cobra.Command{
	Use:   "completion",
	Short: "Generates a bash completion script",
	Long: "Completion generates a bash completion script to standard output." +
		"\n\nTo manually enable bash completion for this program, append the " +
		"following to your ~/.bashrc or ~/.bash_profile: \n  . <(begone " +
		"completion)",
	Args: cobra.NoArgs,
	RunE: completion,
}

func completion(*cobra.Command, []string) error {
	return rootCmd.GenBashCompletion(os.Stdout)
}
