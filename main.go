package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/soeyusuke/fb-bot/talk"
	"github.com/soeyusuke/fb-bot/types"
)

var accessToken = os.Getenv("ACCESS_TOKEN")
var verifyToken = os.Getenv("VERIFY_TOKEN")

const (
	EndPoint = "https://graph.facebook.com/v2.6/me/messages"
)

func main() {
	http.HandleFunc("/", TopPageHandler)
	http.HandleFunc("/webhook", webhookHandler)
	port := os.Getenv("PORT")
	address := fmt.Sprintf(":%s", port)
	http.ListenAndServe(address, nil)
}

func TopPageHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "This is go-bot application's top page.")
}

func webhookHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		verifyTokenAction(w, r)
	}
	if r.Method == "POST" {
		log.Print("ttttttttttttttttttttttttttttttttttttttttttttttttttttttt")
		webhookPostAction(w, r)
	}
}

func verifyTokenAction(w http.ResponseWriter, r *http.Request) {
	if r.URL.Query().Get("hub.verify_token") == verifyToken {
		log.Print("verify token success.")
		fmt.Fprintf(w, r.URL.Query().Get("hub.challenge"))
	} else {
		log.Print("Error: verify token failed.")
		fmt.Fprintf(w, "Error, wrong validation token")
	}
}

func webhookPostAction(w http.ResponseWriter, r *http.Request) {
	var receivedMessage types.ReceivedMessage
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Print(err)
	}
	log.Print("ffffffffffffffffffffffffffffffffffffffffffffffffffffffffff")
	if err = json.Unmarshal(body, &receivedMessage); err != nil {
		log.Print(err)
	}
	messagingEvents := receivedMessage.Entry[0].Messaging
	for _, event := range messagingEvents {
		senderID := event.Sender.ID
		if &event.Message != nil && event.Message.Text != "" {
			sendTextMessage(senderID, event.Message.Text)
		}
	}
	fmt.Fprintf(w, "Success")
}

func sendTextMessage(senderID string, text string) {
	recipient := new(types.Recipient)
	recipient.ID = senderID
	m := new(types.SendMessage)
	m.Recipient = *recipient
	m.Message.Text = talk.TalkApi(text)

	log.Print("-----------------------------------")
	log.Print(m.Message.Text)

	b, err := json.Marshal(m)
	if err != nil {
		log.Print(err)
	}

	req, err := http.NewRequest("POST", EndPoint, bytes.NewBuffer(b))
	if err != nil {
		log.Print(err)
	}

	values := url.Values{}
	values.Add("access_token", accessToken)
	req.URL.RawQuery = values.Encode()
	req.Header.Add("Content-Type", "application/json; charset=UTF-8")
	client := &http.Client{Timeout: time.Duration(30 * time.Second)}
	res, err := client.Do(req)
	if err != nil {
		log.Print(err)
	}

	defer res.Body.Close()
	var result map[string]interface{}
	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		log.Print(err)
	}

	if err := json.Unmarshal(body, &result); err != nil {
		log.Print(err)
	}
	log.Print(result)
}
