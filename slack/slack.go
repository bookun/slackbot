package slack

import (
	"bytes"
	"net/http"
)

// Slack have webhook
type Slack struct {
	Webhook string
}

// NewSlack init Slack
func NewSlack(webhook string) *Slack {
	slack := &Slack{}
	slack.Webhook = webhook
	return slack
}

// Send function send message to tolk room in slack
func (s *Slack) Send(message *bytes.Buffer) error {
	//bot := NewBot("bookun", "名前も変えられる", "general")
	req, err := http.NewRequest("POST", s.Webhook, message)
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		return err
	}
	client := &http.Client{}
	_, err = client.Do(req)
	if err != nil {
		return err
	}
	return nil
}
