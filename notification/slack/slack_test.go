package slack

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSend(t *testing.T) {
	slackCli := New("xoxb-2314993037-3480243399810-eKklHxCNGvH7E3dnGaknRpZL", true)
	require.Nil(t, slackCli.Send(DefaultChannel, "this is a test notification"))
	require.Nil(t, slackCli.SendWithColor(
		DefaultChannel,
		"this is a test notification with a color",
		ColorGood,
	),
	)
}
