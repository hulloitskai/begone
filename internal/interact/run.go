package interact

import (
	"fmt"
	"time"

	spin "github.com/briandowns/spinner"
	"github.com/stevenxie/begone/mbot"
	"github.com/unixpickle/essentials"
)

// BotFunc is a function that uses bot to perform some kind of action.
type BotFunc func(bot *mbot.Bot) error

// Run is a CLI-aware wrapper around bf.
func (br *BotRunner) Run(bf BotFunc) error {
	// Ensure br is ready.
	if br.Bot == nil {
		if err := br.Configure(nil); err != nil {
			return essentials.AddCtx("interact: configuring Bot", err)
		}
	}

	// Login to messenger.
	var (
		spinner *spin.Spinner

		msg      = "Logging into FB Messenger..."
		finalMsg = "Login successful."
		failMsg  = "Login failed."
	)

	switch br.Mode {
	case Standard:
		spinner = cyanSpinner(msg, finalMsg)
		defer func() {
			spinner.Stop()
			fmt.Println()
		}()
	case Reduced:
		br.Println(msg)
	}

	if err := br.Bot.Login(); err != nil {
		switch br.Mode {
		case Standard:
			spinner.FinalMSG = failMsg
		case Reduced:
			br.Println(failMsg)
		}
		return err
	}

	if br.Mode == Standard {
		spinner.Stop()
	}
	br.Println()

	// Begin attack.
	msg = br.runMessage()
	finalMsg = "Attack finished."
	failMsg = "Attack was interrupted."

	switch br.Mode {
	case Standard:
		spinner = cyanSpinner(msg, finalMsg)

		// Subscribe to updates from bot.Counter.
		br.Bot.Counter = make(chan int)
		go func() {
			sfx := spinner.Suffix
			for count := range br.Bot.Counter {
				spinner.Suffix = fmt.Sprintf("%s (sent: %d)", sfx, count)
			}
		}()

	case Reduced:
		br.Println(msg)
	}

	if err := bf(br.Bot); err != nil {
		switch br.Mode {
		case Standard:
			spinner.FinalMSG = failMsg
		case Reduced:
			br.Println(failMsg)
		}
		return err
	}

	if br.Mode == Reduced {
		br.Println(finalMsg)
	}
	return nil
}

var runMessages = []string{
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

// runMessage generates a random run message.
func (br *BotRunner) runMessage() string {
	return runMessages[br.rng.Intn(len(runMessages))]
}

// cyanSpinner makes and starts a cyan-colored spinner, with text
// as the suffix of the spinner.
func cyanSpinner(text, doneText string) *spin.Spinner {
	spinner := spin.New(spin.CharSets[14], 100*time.Millisecond)
	if err := spinner.Color("cyan"); err != nil {
		panic(err)
	}
	spinner.Suffix = " " + text
	spinner.FinalMSG = doneText
	spinner.Start()
	return spinner
}
