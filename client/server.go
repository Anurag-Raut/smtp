package client

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"
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
	client := getClient()
	err = client.SendEmail(m.From, m.To[0], m.Body)
	if err != nil {
		http.Error(w, "Error while sending the mail", http.StatusBadRequest)
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
	log.Println("Listenting on port ", c.Addr)
	c.Handler = mux
	go c.ListenAndServe()

}

func (c *clientServer) Close() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := c.Shutdown(ctx)
	if err != nil {
		log.Println("Error shutting down server", err)
	} else {
		log.Println("Server gracefully stopped")
	}
}
