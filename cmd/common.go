package cmd

import "gopkg.in/alecthomas/kingpin.v2"

func registerCommonOpts(cmd *kingpin.CmdClause) {
	cmd.Flag("debug", "Enable debug mode.").BoolVar(&copts.Debug)

	cmd.Flag("delay", "Delay (in ms) between messages.").Short('d').
		Default("225").IntVar(&copts.Delay)
	cmd.Flag("cycles", "Number of spam cycles (-1 for infinite).").Short('c').
		Default("-1").IntVar(&copts.Cycles)

	cmd.Flag("send-fail-delay", "Delay (in ms) after a send fail.").Short('D').
		Default("1000").IntVar(&copts.SendFailDelay)
	cmd.Flag("max-send-fails", "Max consecutive send fails before aborting.").
		Short('f').Default("3").IntVar(&copts.MaxSendFails)

	cmd.Arg(
		"conversation ID",
		"The target conversation ID (last portion of a www.messenger.com link).",
	).StringVar(&copts.ConvoID)
}

var copts struct {
	Debug                       bool
	Delay, Cycles               int
	SendFailDelay, MaxSendFails int
	ConvoID                     string
}
