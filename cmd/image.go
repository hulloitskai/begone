package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	str "strings"

	"gopkg.in/alecthomas/kingpin.v2"

	"github.com/stevenxie/begone/mbot"
	ess "github.com/unixpickle/essentials"
)

func registerImageCmd(app *kingpin.Application) {
	imageCmd = app.Command("image", "Spam a target with an image file.")

	// Args:
	imageOpts.Path = imageCmd.Arg("filepath", "Path to the image file to send.").
		Required().String()

	imageOpts.ConvoID = imageCmd.Arg(
		"conversation ID",
		"The target conversation ID (last portion of a www.messenger.com link).",
	).Default("").String()
}

var (
	imageCmd *kingpin.CmdClause

	imageOpts struct {
		Path, ConvoID *string
	}
)

func image() error {
	// Validate file.
	info, err := os.Stat(*imageOpts.Path)
	if err != nil {
		return err
	}

	if info.IsDir() {
		return fmt.Errorf("path '%s' refers to directory, not an image file",
			*imageOpts.Path)
	}

	switch str.ToLower(filepath.Ext(*imageOpts.Path)) {
	case ".jpeg", ".jpg", ".png", ".gif":
	default:
		return fmt.Errorf("file '%s' is not an image",
			filepath.Base(*imageOpts.Path))
	}

	// Derive convoURL.
	runner := deriveBotRunner(nil)
	convoURL, err := deriveConvoURL(*imageOpts.ConvoID, runner.Prompter)
	if err != nil {
		return err
	}

	// Configure and run BotRunner.
	bcfg := deriveBotConfig()
	if err = runner.Configure(bcfg); err != nil {
		return ess.AddCtx("configuring BotRunner", err)
	}
	return runner.Run(func(bot *mbot.Bot) error {
		return bot.CycleImage(convoURL, *imageOpts.ConvoID)
	})
}
