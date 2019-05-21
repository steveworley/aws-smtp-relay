package relay

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/service/pinpointemail"
	"github.com/aws/aws-sdk-go/service/pinpointemail/pinpointemailiface"
)

// Request object
type Request struct {
	Time  time.Time
	addr  net.Addr
	IP    string
	From  string
	To    []string
	Error string
	err   *error
}

// Log prints the Request data to Stdout/Stderr
func (req *Request) Log() {
	req.IP = req.addr.(*net.TCPAddr).IP.String()
	var out *os.File
	err := *req.err
	if err != nil {
		req.Error = err.Error()
		out = os.Stderr
	} else {
		out = os.Stdout
	}
	b, _ := json.Marshal(req)
	fmt.Fprintln(out, string(b))
}

// Send uses the given Pinpoint API to send email data
func Send(
	pinpointAPI pinpointemailiface.PinpointEmailAPI,
	origin net.Addr,
	from *string,
	to *[]string,
	data *[]byte,
	setName *string,
) {
	var err error
	req := &Request{
		Time: time.Now().UTC(),
		addr: origin,
		From: *from,
		To:   *to,
		err:  &err,
	}
	defer req.Log()
	destinations := []*string{}
	for k := range *to {
		destinations = append(destinations, &(*to)[k])
	}

	_, err = pinpointAPI.SendEmail(&pinpointemail.SendEmailInput{
		ConfigurationSetName: setName,
		FromEmailAddress: from,
		Destination: &pinpointemail.Destination{
			ToAddresses: destinations,
		},
		Content: &pinpointemail.EmailContent{
			Raw: &pinpointemail.RawMessage{
				Data: *data,
			},
		},
	})
}