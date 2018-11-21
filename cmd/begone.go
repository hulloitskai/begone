package cmd

import (
	"fmt"
	"net/url"
	"os"
	"os/signal"
	"strings"

	"github.com/stevenxie/begone/config"
	"github.com/stevenxie/begone/mbot"
	"github.com/stevenxie/begone/strgen"
	ess "github.com/unixpickle/essentials"
	"github.com/urfave/cli"
)

// Constants for command begone.
const (
	Name    = "begone"
	Desc    = "a fully automatic spamming tool for FB Messenger"
	Version = "0.2.0"
	Author  = "Steven Xie <hello@stevenxie.me>"
	Usage   = `begone [global options] <command> [command options] [arguments...]


For more information about a command, run: begone help <command>`
)

// GlobalFlags are global cli flags for command begone.
var GlobalFlags = []cli.Flag{
	cli.IntFlag{
		Name:  "delay, d",
		Usage: "delay in milliseconds between each message",
		Value: 30,
	},
	cli.IntFlag{
		Name:  "cycles, c",
		Usage: "maximum number of spam cycles (-1 for unlimited)",
		Value: -1,
	},
}

// Begone creates a mbot.Bot, and launches an attack routine using 'gen'.
//
// Uses cli flags "cycles", "delay"
func Begone(ctx *cli.Context, gen strgen.Generator, convoID string) error {
	// Read config.
	cfg, err := config.Load()
	if err != nil {
		return err
	}
	if err = queryMissing(cfg); err != nil {
		return ess.AddCtx("cmd: requesting missing config values", err)
	}

	// Get conversation ID.
	if convoID == "" {
		if convoID, err = queryConvoID(); err != nil {
			return ess.AddCtx("cmd: requesting conversation ID", err)
		}
	}

	// Derive bot configuration, create bot.
	var (
		bcfg = botConfig(ctx, cfg)
		bot  = mbot.NewBot(bcfg, gen)
	)

	// Login to messenger.
	spinner := cyanSpinner("Logging in to FB Messenger...", "Login successful.")
	defer func() {
		spinner.Stop()
		fmt.Println()
	}()
	if err = bot.Login(); err != nil {
		spinner.FinalMSG = "Login failed."
		return err
	}
	spinner.Stop()
	fmt.Println()

	// Begin attack, watch out for kill signals.
	spinner = attackSpinner(nil)
	var killed bool

	// Subscribe to updates from bot.Counter.
	bot.Counter = make(chan int)
	go func() {
		sfx := spinner.Suffix
		for count := range bot.Counter {
			spinner.Suffix = fmt.Sprintf("%s (sent: %d)", sfx, count)
		}

		// In the event that an interrupt occurrred.
		if killed {
			spinner.FinalMSG = "Killed bot (caught interrupt signal)."
			spinner.Stop()
			os.Exit(0)
		}
	}()

	// Subscribe to interrupt signals.
	sigch := make(chan os.Signal, 1)
	signal.Notify(sigch, os.Interrupt)
	go func() {
		for range sigch {
			bot.Kill() // kill bot upon interrupt.
			killed = true
		}
	}()

	if err = bot.Begone(convoID); err != nil {
		spinner.FinalMSG = "Attack was interrupted."
		return err
	}
	return nil
}

// queryMissing fills missing fields of cfg with values read from stdin.
func queryMissing(cfg *config.Config) error {
	if cfg.Username == "" {
		uname, err := queryUsername()
		if err != nil {
			return err
		}
		cfg.Username = uname
	}

	if cfg.Password == "" {
		pw, err := queryPassword()
		if err != nil {
			return err
		}
		cfg.Password = pw
	}

	return nil
}

// queryConvoID requests and reads a conversation ID from stdin.
func queryConvoID() (string, error) {
	var convoID string
	for convoID == "" {
		fmt.Println("Enter the target conversation URL " +
			"(https://messenger.com/t/...):")

		var rawurl string
		if _, err := fmt.Scanf("%s", &rawurl); err != nil {
			return "", err
		}
		convoID = parseConvoURL(rawurl)
	}
	return convoID, nil
}

// parseConvoURL will parse the url of an FB Messenger conversation (rawurl)
// into an conversation ID.
//
// If parseConvoURL fails, it will print an error message to os.Stderr, and
// return an empty string.
func parseConvoURL(rawurl string) (convoID string) {
	u, err := url.Parse(rawurl)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Bad URL: %v\n", err)
		return ""
	}

	path := u.EscapedPath()
	slashIndex := strings.LastIndexByte(path, '/')
	if slashIndex == -1 {
		fmt.Fprintf(os.Stderr, "Bad URL: path does not contain any '/'\n")
		return ""
	}

	return path[slashIndex+1:]
}
