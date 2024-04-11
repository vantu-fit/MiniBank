package mail

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/vantu-fit/master-go-be/utils"
)

func TestSendEmailWithGmail(t *testing.T) {
	if testing.Short() { 
		t.Skip()
	}
	config, err := utils.LoadConfig("..")
	require.NoError(t, err)

	sender := NewGmailSender(config.EmailSenderName, config.EmailSenderAddress, config.EmailSenderPassword)

	subject := "A test email"
	content := `
	<h1> hello world </h1>
	<p> this is a  test message from <a href="https://github.com/vantu-fit> Vantu-fit </a> </p>
	`
	to := []string{"dotu30257@gmail.com"}
	attachFile := []string{"../README.md"}

	err = sender.SendEmail(subject, content, to, nil, nil, attachFile)
	require.NoError(t, err)
}
