package cmd

import (
	"math/rand"
	"time"

	spin "github.com/briandowns/spinner"
)

var attackTexts = []string{
	"Launching an attack...",
	"Slaying thots...",
	"Justice rains from above...",
	"REEEEEEEEEEEEEEEEEE...",
}

// attackSpinner makes and starts a spinner with a random attack-related suffix
// text (chosen from attackTexts).
func attackSpinner(rng *rand.Rand) *spin.Spinner {
	if rng == nil {
		src := rand.NewSource(time.Now().UnixNano())
		rng = rand.New(src)
	}
	var (
		text    = attackTexts[rng.Intn(len(attackTexts))]
		spinner = cyanSpinner(text, "Finished attack.")
	)
	return spinner
}

// cyanSpinner makes and starts a cyan-colored spinner, with text
// as the suffix of the spinner.
func cyanSpinner(text, doneText string) *spin.Spinner {
	spinner := spin.New(spin.CharSets[14], 100*time.Millisecond)
	spinner.Color("cyan")
	spinner.Suffix = " " + text
	spinner.FinalMSG = doneText
	spinner.Start()
	return spinner
}
