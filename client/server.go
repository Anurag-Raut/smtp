package client

import (
	"encoding/json"
	"io"
	"net/http"
)

func sendMail(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Only post is allowed on this req", http.StatusBadRequest)
		return
	}
	type body struct {
		from string
		to   []string
		body *string
	}
	var m body
	bodyBytes, err := io.ReadAll(r.Body)
	err = json.Unmarshal(bodyBytes, &m)
	if err != nil {
		http.Error(w, "Error whi;e parsing the body", http.StatusBadRequest)
		return
	}
	client := GetClient()

	err = client.SendEmail(m.from, m.to[0], m.body)
	if err != nil {
		http.Error(w, "Error while sending the mail", http.StatusBadRequest)
		return
	}
}

func main() {
	http.HandleFunc("/newRequest", sendMail)

	http.ListenAndServe(":8000", nil)
}
