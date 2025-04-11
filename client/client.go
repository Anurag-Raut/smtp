package client

import (
	"errors"
	"log"
	"net"
	"net/http"
	"slices"

	"github.com/Anurag-Raut/smtp/client/parser"
	"github.com/Anurag-Raut/smtp/client/session"
)

type Client struct {
	httpWriter http.ResponseWriter
}

func getClient(w http.ResponseWriter) *Client {
	return &Client{
		httpWriter: w,
	}
}

func (c *Client) getMxRecords(from string) ([]*net.MX, error) {
	domain, err := parser.GetDomainFromEmail(from)
	if err != nil {
		return nil, err
	}

	if domain == "localhost" {
		return []*net.MX{
			{Host: "127.0.0.1"},
		}, nil
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

func (c *Client) SendEmail(from string, to []string, body *string) error {
	log.Println("FROM", from)
	mxRecords, err := c.getMxRecords(from)
	if err != nil {
		log.Println(err)
		return err
	}

	for _, mxRecord := range mxRecords {
		conn, err := net.Dial("tcp", mxRecord.Host+":8000")
		if err != nil {
			return err
		}

		session := session.NewSession(conn, c.httpWriter)
		err = session.Begin()
		if err == nil {
			return nil
		} else {
			log.Println("err:", err.Error())
		}
		session.SendEmail(from, to, body)

	}

	return errors.New("Could resolve any MX records")

}
