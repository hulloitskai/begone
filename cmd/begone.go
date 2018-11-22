package cmd

import (
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"runtime"
	"time"

	spin "github.com/briandowns/spinner"
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
	Version = "1.2.0"
	Author  = "Steven Xie <hello@stevenxie.me>"
	Usage   = "begone [global options] <command> [command options] [arguments...]" +
		"\n\n   For more information about a command, run:\n   \t\t" +
		"begone help <command>"
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
	var spinner *spin.Spinner
	switch runtime.GOOS {
	case "windows":
		fmt.Println("Logging into FB Messenger...")
	default:
		spinner = cyanSpinner("Logging in to FB Messenger...", "Login successful.")
		defer func() {
			spinner.Stop()
			fmt.Println()
		}()
	}

	if err = bot.Login(); err != nil {
		if runtime.GOOS != "windows" {
			spinner.FinalMSG = "Login failed."
		}
		return err
	}

	if runtime.GOOS != "windows" {
		spinner.Stop()
	}
	fmt.Println()

	// Begin attack.
	startText := randAtackText(nil)
	const endText = "Attack finished."
	switch runtime.GOOS {
	case "windows":
		fmt.Println(startText)

	default:
		spinner = cyanSpinner(startText, endText)

		// Subscribe to updates from bot.Counter.
		var killed bool
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

				// Wait for spinner to stop asynchronously.
				time.Sleep(time.Millisecond)
				os.Exit(0)
			}
		}()

		// Watch for interrupt signals.
		sigch := make(chan os.Signal, 1)
		signal.Notify(sigch, os.Interrupt)
		go func() {
			for range sigch {
				bot.Kill() // kill bot upon interrupt.
				killed = true
			}
		}()
	}

	if err = bot.Begone(convoID); err != nil {
		const msg = "Attack was interrupted."
		switch runtime.GOOS {
		case "windows":
			fmt.Println(msg)
		default:
			spinner.FinalMSG = msg
		}
		return err
	}

	if runtime.GOOS == "windows" {
		fmt.Println(endText)
	}
	return nil
}

var attackTexts = []string{
	"Launching an attack...",
	"No turning back now...",
	"Justice rains from above...",
	"Hitting that yeet...",
	"Demolishing message box...",
	"I guess they never miss, huh...",
	"Packing a punch...",
	"Makin' a mess...",
	"Letting it rip...",
}

func randAtackText(rng *rand.Rand) string {
	if rng == nil {
		src := rand.NewSource(time.Now().UnixNano())
		rng = rand.New(src)
	}
	return attackTexts[rng.Intn(len(attackTexts))]
}
