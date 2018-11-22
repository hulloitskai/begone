package cmd

import (
	"time"

	spin "github.com/briandowns/spinner"
)

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
