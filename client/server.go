package client

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/Anurag-Raut/smtp/logger"
)

func sendMail(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Only post is allowed on this req", http.StatusBadRequest)
		return
	}
	type body struct {
		From string   `json:"from"`
		To   []string `json:"to"`
		Body *string  `json:"body"`
	}
	var m body
	bodyBytes, err := io.ReadAll(r.Body)
	err = json.Unmarshal(bodyBytes, &m)
	if err != nil {
		http.Error(w, "Error whi;e parsing the body", http.StatusBadRequest)
		return
	}
	client := getClient(w)
	err = client.SendEmail(m.From, m.To, m.Body)
	if err != nil {
		http.Error(w, "Error while sending the mail: "+err.Error(), http.StatusBadRequest)
		return
	}
}

type clientServer struct {
	http.Server
}

func NewClientServer(address string, port string) *clientServer {
	return &clientServer{
		Server: http.Server{
			Addr: address + ":" + port,
		},
	}
}

func (c *clientServer) Listen() {
	mux := http.NewServeMux()

	mux.HandleFunc("/newRequest", sendMail)
	logger.ClientLogger.Println("Listenting on port ", c.Addr)
	c.Handler = mux
	go c.ListenAndServe()

}

func (c *clientServer) Close() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := c.Shutdown(ctx)
	if err != nil {
		logger.ClientLogger.Println("Error shutting down server", err)
	} else {
		logger.ClientLogger.Println("Server gracefully stopped")
	}
}
