package mail

import (
	"testing"

	"github.com/mishvets/WeatherForecast/util"
	"github.com/stretchr/testify/require"
)

func TestSendEmailWithGmail(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	config, err := util.LoadConfig("..")
	require.NoError(t, err)

	sender := newGmailSender(config.EmailSenderName, config.EmailSenderAdress, config.EmailSenderPassword)

	subject := "A test mail"
	content := `
	<h1>Hello, this is a test mail</h1>
	<p>This is a test msg with link to <a href="https://www.google.com.ua/?hl=uk">google</a>.</p>
	`
	to := []string{"weatherforecast099@gmail.com"}
	err = sender.SendEmail(subject, content, to)
	require.NoError(t, err)
}