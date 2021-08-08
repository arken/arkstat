package mail

import (
	"context"
	"io/ioutil"
	"strings"
	"time"

	mailgun "github.com/mailgun/mailgun-go/v4"
	"gopkg.in/yaml.v3"
)

type Mailbox struct {
	mg     *mailgun.MailgunImpl
	sender string
}

type email struct {
	Subject string `yaml:"subject"`
	Body    string `yaml:"body"`
}

func Init(domain, privateKey, sender string) (mailbox *Mailbox, err error) {
	result := &Mailbox{}
	result.mg = mailgun.NewMailgun(domain, privateKey)
	result.sender = sender
	return result, nil
}

func (m *Mailbox) Send(templatePath, toAddress string) (err error) {
	// Initialize message structure to hold email message
	data := email{}

	// Read data from file
	content, err := ioutil.ReadFile(templatePath)
	if err != nil {
		return err
	}

	// Parse content into YAML data
	err = yaml.Unmarshal(content, &data)

	// Construct new email message.
	message := m.mg.NewMessage(
		m.sender,
		data.Subject,
		data.Body,
		cleanEmail(toAddress))

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	_, _, err = m.mg.Send(ctx, message)
	return err
}

// cleanEmail clean the email address input.
func cleanEmail(email string) string {
	return strings.TrimSuffix(strings.TrimPrefix(email, `"`), `"`)
}
