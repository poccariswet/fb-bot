package types

// ReceivedMessage ...
type ReceivedMessage struct {
	Object string  `json:"object"`
	Entry  []Entry `json:"entry"`
}

// Entry ...
type Entry struct {
	ID        string      `json:"id"`
	Time      int         `json:"time"`
	Messaging []Messaging `json:"messaging"`
}

// Messaging ...
type Messaging struct {
	Sender    Sender    `json:"sender"`
	Recipient Recipient `json:"recipient"`
	Timestamp int       `json:"timestamp"`
	Message   Message   `json:"message"`
}

// Sender ...
type Sender struct {
	ID string `json:"id"`
}

// Recipient ...
type Recipient struct {
	ID string `json:"id"`
}

// Message ...
type Message struct {
	MID  string `json:"mid"`
	Seq  int    `json:"seq"`
	Text string `json:"text"`
}

// SendMessage ...
type SendMessage struct {
	Recipient Recipient `json:"recipient"`
	Message   struct {
		Text string `json:"text"`
	} `json:"message"`
}

//---talk type---
type TalkJson struct {
	Status  int             `json:"status"`
	Message string          `json:"message"`
	Results []talkApiResult `json:"results"`
}

type talkApiResult struct {
	Perplexity float64 `json:"perplexity"`
	Reply      string  `json:"reply"`
}
