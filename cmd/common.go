package cmd

import kingpin "gopkg.in/alecthomas/kingpin.v2"

func registerCommonOpts(cmd *kingpin.CmdClause) {
	cmd.Flag("delay", "Delay (in ms) between messages.").Short('d').
		Default("250").IntVar(&copts.Delay)
	cmd.Flag("cycles", "Number of spam cycles (-1 for infinite).").Short('c').
		Default("-1").IntVar(&copts.Cycles)

	cmd.Flag("send-fail-delay", "Delay (in ms) after a send fail.").Short('D').
		Default("1250").IntVar(&copts.SendFailDelay)
	cmd.Flag("max-send-fails", "Max consecutive send fails before aborting.").
		Short('f').Default("3").IntVar(&copts.MaxSendFails)

	cmd.Flag("assume-user", "Treat numeric convo IDs as belonging to users, "+
		"not groups.").
		Short('u').BoolVar(&copts.AssumeUser)
	cmd.Flag("debug", "Enable debug mode.").BoolVar(&copts.Debug)
	cmd.Flag("no-fancy", "Disable fancy terminal graphics.").
		BoolVar(&copts.NoFancy)

	cmd.Arg(
		"conversation ID",
		"The target conversation ID (last portion of a www.messenger.com link).",
	).StringVar(&copts.ConvoID)
}

var copts struct {
	Delay, Cycles               int
	SendFailDelay, MaxSendFails int
	Debug, NoFancy, AssumeUser  bool
	ConvoID                     string
}
