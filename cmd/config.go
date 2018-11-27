package cmd

// Configure Kingpin.
func init() {
	// Customize help, version flag.
	app.HelpFlag.Short('h')
	app.VersionFlag.Short('v')

	// Register all the things.
	registerAppFlags(app)
	registerLoginCmd(app)

	registerEmojifyCmd(app)
	registerRepeatCmd(app)
	registerFileCmd(app)
	registerImageCmd(app)
}
