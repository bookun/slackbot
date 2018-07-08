package slack

import (
	"bytes"
	"log"
	"net/http"
)

type Slack struct {
	Webhook string
}

func NewSlack(webhook string) *Slack {
	slack := &Slack{}
	//slack.Webhook := "https://hooks.slack.com/services/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"
	slack.Webhook = webhook
	return slack
}

func (s *Slack) Send(message *bytes.Buffer) {
	//bot := NewBot("bookun", "名前も変えられる", "general")
	req, err := http.NewRequest("POST", s.Webhook, message)
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		log.Fatal(err)
	}
	client := &http.Client{}
	_, err = client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
}
