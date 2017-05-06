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
)

// FacebookMessenger ...
type FacebookMessenger struct {
	Token string
}

// CallbackMessage ...
type CallbackMessage struct {
	Object string   `json:"object"`
	Entry  []*Entry `json:"entry"`
}

// Entry ...
type Entry struct {
	ID        int          `json:"id"`
	Time      int          `json:"time"`
	Messaging []*Messaging `json:"messaging"`
}

// Messaging ...
type Messaging struct {
	Sender    *ID       `json:"sender"`
	Recipient *ID       `json:"recipient"`
	Timestamp int       `json:"timestamp"`
	Message   *Message  `json:"message,omitempty"`
	Delivery  *Delivery `json:"delivery,omitempty"`
	string    `json:""`
}

// Message ...
type Message struct {
	Mid  string `json:"mid"`
	Seq  int    `json:"seq"`
	Text string `json:"text"`
}

// Delivery ...
type Delivery struct {
	Mids      []string `json:"mids"`
	Watermark int      `json:"watermark"`
	Seq       int      `json:"seq"`
}

// ID ...
type ID struct {
	ID int `json:"id"`
}

// Text ...
type Text struct {
	Text string `json:"text"`
}

// SendMessage ...
type SendMessage struct {
	Recipient *ID   `json:"recipient"`
	Message   *Text `json:"message"`
}


var debug bool
var fb *FacebookMessenger

func callbackHandler(w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("Something wrong: %s\n", err.Error())
		return
	}
	if debug {
		log.Println("RecievedMessage Body:", string(b))
	}

	m, _ := url.ParseQuery(r.URL.RawQuery)
	fmt.Println(m["hub.verify_token"])
	if len(m["hub.verify_token"]) > 0 && m["hub.verify_token"][0] == os.Getenv("VERIFY_TOKEN") && len(m["hub.challenge"]) > 0 {
		fmt.Fprintf(w, m["hub.challenge"][0])
		return
	}

	var msg CallbackMessage
	err = json.Unmarshal(b, &msg)
	if err != nil {
		fmt.Printf("Something wrong: %s\n", err.Error())
		return
	}

	for _, event := range msg.Entry[0].Messaging {
		sender := event.Sender.ID
		if event.Message != nil {
			fmt.Printf("Recieved Text: %s\n", event.Message.Text)
			err := fb.SendTextMessage(sender, event.Message.Text)
			if err != nil {
				fmt.Printf("Something wrong: %s\n", err.Error())
			}
		}
	}

}

// SendTextMessage ...
func (fb *FacebookMessenger) SendTextMessage(recipient int, text string) error {

	m := &SendMessage{
		Recipient: &ID{ID: recipient},
		Message:   &Text{Text: text},
	}
	b, err := json.Marshal(m)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", "https://graph.facebook.com/v2.6/me/messages?access_token="+fb.Token, bytes.NewBuffer(b))
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/json; charset=UTF-8")

	client := &http.Client{
		Timeout: time.Duration(30 * time.Second),
	}
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return err
	}
	log.Print("Response: ", result)
	return nil

}


func main() {
	debug = true
	fb = &FacebookMessenger{
		Token: os.Getenv("ACESS_TOKEN"),
	}

	http.HandleFunc("/fbbot/callback", callbackHandler)

	port := os.Getenv("PORT")
	addr := fmt.Sprintf(":%s", port)
	http.ListenAndServe(addr, nil)
}
