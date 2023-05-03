package relay

import (
	"net"
	"regexp"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sesv2"
	"github.com/aws/aws-sdk-go/service/sesv2/sesv2iface"
	"github.com/blueimp/aws-smtp-relay/internal/relay"
)

// Client implements the Relay interface.
type Client struct {
	sesAPI          sesv2iface.SESV2API
	setName         *string
	allowFromRegExp *regexp.Regexp
	denyToRegExp    *regexp.Regexp
}

// Send uses the client SESAPI to send email data
func (c Client) Send(
	origin net.Addr,
	from string,
	to []string,
	data []byte,
) error {
	allowedRecipients, deniedRecipients, err := relay.FilterAddresses(
		from,
		to,
		c.allowFromRegExp,
		c.denyToRegExp,
	)
	if err != nil {
		relay.Log(origin, &from, deniedRecipients, err)
	}
	if len(allowedRecipients) > 0 {
		_, err := c.sesAPI.SendEmail(&sesv2.SendEmailInput{
			ConfigurationSetName: c.setName,
			Content:              &sesv2.EmailContent{Raw: &sesv2.RawMessage{Data: data}},
			Destination:          &sesv2.Destination{ToAddresses: allowedRecipients},
			FromEmailAddress:     &from,
		})
		relay.Log(origin, &from, allowedRecipients, err)
		if err != nil {
			return err
		}
	}
	return err
}

// New creates a new client with a session.
func New(
	configurationSetName *string,
	allowFromRegExp *regexp.Regexp,
	denyToRegExp *regexp.Regexp,
) Client {
	return Client{
		sesAPI:          sesv2.New(session.Must(session.NewSession())),
		setName:         configurationSetName,
		allowFromRegExp: allowFromRegExp,
		denyToRegExp:    denyToRegExp,
	}
}
