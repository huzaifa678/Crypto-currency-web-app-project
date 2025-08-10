package mail

import (
	"testing"

	"github.com/huzaifa678/Crypto-currency-web-app-project/utils"
	"github.com/stretchr/testify/require"
)



func TestSenderWithGmail(t *testing.T) {
	config, err := utils.LoadConfig("..")

	require.NoError(t, err)
	
	sender := NewGmailSender(config.SenderName, config.SenderEmail, config.SenderPassword)

	subject := "Test Email"
	content := "<h1>This is a test email</h1>"
	to := []string{"huzaifagill411@gmail.com"}
	attachFiles := []string{
		"../README.md",
	}

	err = sender.SendEmail(subject, content, to, nil, nil, attachFiles)
	require.NoError(t, err)
}