package client

type Client struct {
	port    string
	address string
}

type Config struct {
	port    *string
	address *string
}

func NewClient(cnf Config) *Client {
	client := Client{port: "8000", address: "127.0.0.1"}

	if cnf.port != nil {
		client.port = *cnf.port
	}

	if cnf.address != nil {
		client.address = *cnf.address
	}
	return &client
}
