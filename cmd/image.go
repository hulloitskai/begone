package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	str "strings"

	"github.com/spf13/cobra"
	"github.com/stevenxie/begone/internal/interact"
	"github.com/stevenxie/begone/mbot"
	ess "github.com/unixpickle/essentials"
)

var imageCmd = &cobra.Command{
	Use:   "image [flags]\n  begone image [flags] FILEPATH [CONVERSATION ID]",
	Short: "Spam a target with an image",
	Long:  "Image spams a target with an image file.",
	Args:  withUsage(cobra.RangeArgs(1, 2)),
	RunE:  image,
}

func image(cmd *cobra.Command, args []string) error {
	// Parse arguments.
	var (
		fpath, convoID string
		p              = interact.NewPrompter()
	)

	fpath = args[0]
	if len(args) > 1 {
		convoID = args[1]
	}

	// Validate file.
	info, err := os.Stat(fpath)
	if err != nil {
		return err
	}

	if info.IsDir() {
		return fmt.Errorf("path '%s' refers to directory, not an image file", fpath)
	}

	switch str.ToLower(filepath.Ext(fpath)) {
	case ".jpeg", ".jpg", ".png", ".gif":
	default:
		return fmt.Errorf("file '%s' is not an image", filepath.Base(fpath))
	}

	// Derive Bot config.
	bcfg, err := deriveBotConfig(cmd.Flags())
	if err != nil {
		return ess.AddCtx("deriving bot config", err)
	}

	// Derive convoURL.
	convoURL, err := deriveConvoURL(convoID, p)
	if err != nil {
		return err
	}

	// Configure and run BotRunner.
	runner := interact.NewBotRunnerWith(p)
	if err = runner.Configure(bcfg); err != nil {
		return ess.AddCtx("configuring BotRunner", err)
	}
	return runner.Run(func(bot *mbot.Bot) error {
		return bot.CycleImage(convoURL, fpath)
	})
}
