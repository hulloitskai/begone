package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	str "strings"

	"github.com/stevenxie/begone/pkg/mbot"
	ess "github.com/unixpickle/essentials"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

func registerImageCmd(app *kingpin.Application) {
	imageCmd = app.Command("image", "Spam a target with an image file.")

	// Args:
	imageCmd.Arg("filepath", "Path to the image file to send.").
		Required().StringVar(&imagePath)

	// Common options:
	registerCommonOpts(imageCmd)
}

var (
	imageCmd  *kingpin.CmdClause
	imagePath string
)

func image() error {
	// Validate file.
	info, err := os.Stat(imagePath)
	if err != nil {
		return err
	}

	if info.IsDir() {
		return fmt.Errorf("path '%s' refers to directory, not an image file",
			imagePath)
	}

	switch str.ToLower(filepath.Ext(imagePath)) {
	case ".jpeg", ".jpg", ".png", ".gif":
	default:
		return fmt.Errorf("file '%s' is not an image", filepath.Base(imagePath))
	}

	// Derive convoURL.
	runner := deriveBotRunner(nil)
	convoURL, err := deriveConvoURL(copts.ConvoID, runner.Prompter)
	if err != nil {
		return err
	}

	// Configure and run BotRunner.
	bcfg := deriveBotConfig()
	if err = runner.Configure(bcfg); err != nil {
		return ess.AddCtx("configuring BotRunner", err)
	}
	return runner.Run(func(bot *mbot.Bot) error {
		return bot.CycleImage(convoURL, imagePath)
	})
}
