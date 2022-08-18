package log

import "github.com/agflow/tools/notification/slack"

// NewSlackHook returns a hook for slack
func NewSlackHook(token, channel string) Hook {
	return func(info MetaInfo) error {
		var color string
		switch info.Lvl {
		case InfoLvl:
			color = slack.ColorGood
		case FatalLvl, PanicLvl, ErrorLvl:
			color = slack.ColorDanger
		default:
			color = slack.ColorWarning
		}
		msg := info.Msg + "\n"
		slackCli := slack.New(token, true)
		return slackCli.SendWithColor(channel, msg, color)
	}
}
