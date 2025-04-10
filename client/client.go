package client

import (
	"errors"
	"net"
	"slices"

	"github.com/Anurag-Raut/smtp/client/parser"
	"github.com/Anurag-Raut/smtp/client/session"
)

type Client struct {
}

func GetClient() *Client {
	return &Client{}
}

func (c *Client) getMxRecords(from string) ([]*net.MX, error) {
	domain, err := parser.GetDomainFromEmail(from)
	if err != nil {
		return nil, err
	}

	mxRecords, err := net.LookupMX(domain)
	if err != nil {
		return nil, err
	}
	slices.SortFunc(mxRecords, func(a, b *net.MX) int {
		if a.Pref < b.Pref {
			return -1
		} else if a.Pref > b.Pref {
			return 1
		}
		return 0
	})

	return mxRecords, nil
}

func (c *Client) SendEmail(from string, to string, body *string) error {

	mxRecords, err := c.getMxRecords(from)
	if err != nil {
		return err
	}

	for _, mxRecord := range mxRecords {
		conn, err := net.Dial("tcp", mxRecord.Host)
		if err != nil {
			return err
		}

		session := session.NewSession(conn)
		err = session.Begin()
		if err != nil {
			return err
		}

	}

	return errors.New("Could resolve any MX records")

}
