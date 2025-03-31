package email

import (
	"testing"

	"github.com/phongnd2802/simple-bank/util"
	"github.com/stretchr/testify/require"
)

func TestSendEmailWithGmail(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	config, err := util.LoadConfig("local", "../config")
	require.NoError(t, err)

	sender := NewGmailSender(config.Email.EmailSenderName, config.Email.EmailSenderAddress, config.Email.EmailSenderPassword)

	subject := "A test email"
	content := `
		<h1> Hello World </h1>
		<p> This is a test messsage </p>
	`
	to := []string{"phongnguyen28022004@gmail.com"}
	attachFiles := []string{"../README.md"}

	err = sender.SendEmail(subject, content, to, nil, nil, attachFiles)
	require.NoError(t, err)
}