package models

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// Slack have webhook and channel
type Slack struct {
	Channel string
	webhook string
	token   string
}

// NewSlack init slack
func NewSlack(webhook, channel, token string) *Slack {
	return &Slack{
		Channel: channel,
		webhook: webhook,
		token:   token,
	}
}

type Field struct {
	Title string `json:"title"`
	Value string `json:"value"`
	Short bool   `json:"short"`
}

// Attachment struct
type Attachment struct {
	Fallback   string  `json:"fallback"`
	Color      string  `json:"color"`
	Pretext    string  `json:"pretext"`
	AuthorName string  `json:"author_name"`
	AuthorLink string  `json:"author_link"`
	AuthorIcon string  `json:"author_icon"`
	Title      string  `json:"title"`
	TitleLink  string  `json:"title_link"`
	Text       string  `json:"text"`
	Fields     []Field `json:"fields"`
	Markdown   bool    `json:"markdown"`
	ThumbURL   string  `json:"thumb_url"`
	Footer     string  `json:"footer"`
	FooterIcon string  `json:"footer_icon"`
	Ts         int64   `json:"ts"`
}

type Message struct {
	Name        string       `json:"username"`
	Channel     string       `json:"channel"`
	LinkName    bool         `json:"link_names"`
	Attachments []Attachment `json:"attachments"`
}

type usersList struct {
	Ok      bool `json:"ok"`
	Members []struct {
		ID      string `json:"id"`
		TeamID  string `json:"team_id"`
		Name    string `json:"name"`
		Deleted bool   `json:"deleted"`
		Profile struct {
			Title                 string `json:"title"`
			Phone                 string `json:"phone"`
			Skype                 string `json:"skype"`
			RealName              string `json:"real_name"`
			RealNameNormalized    string `json:"real_name_normalized"`
			DisplayName           string `json:"display_name"`
			DisplayNameNormalized string `json:"display_name_normalized"`
			StatusText            string `json:"status_text"`
			StatusEmoji           string `json:"status_emoji"`
			StatusExpiration      int    `json:"status_expiration"`
			AvatarHash            string `json:"avatar_hash"`
			BotID                 string `json:"bot_id"`
			APIAppID              string `json:"api_app_id"`
			AlwaysActive          bool   `json:"always_active"`
			ImageOriginal         string `json:"image_original"`
			FirstName             string `json:"first_name"`
			LastName              string `json:"last_name"`
			Image24               string `json:"image_24"`
			Image32               string `json:"image_32"`
			Image48               string `json:"image_48"`
			Image72               string `json:"image_72"`
			Image192              string `json:"image_192"`
			Image512              string `json:"image_512"`
			Image1024             string `json:"image_1024"`
			StatusTextCanonical   string `json:"status_text_canonical"`
			Team                  string `json:"team"`
		} `json:"profile"`
		IsBot             bool   `json:"is_bot"`
		IsAppUser         bool   `json:"is_app_user"`
		Updated           int    `json:"updated"`
		Color             string `json:"color,omitempty"`
		RealName          string `json:"real_name,omitempty"`
		Tz                string `json:"tz,omitempty"`
		TzLabel           string `json:"tz_label,omitempty"`
		TzOffset          int    `json:"tz_offset,omitempty"`
		IsAdmin           bool   `json:"is_admin,omitempty"`
		IsOwner           bool   `json:"is_owner,omitempty"`
		IsPrimaryOwner    bool   `json:"is_primary_owner,omitempty"`
		IsRestricted      bool   `json:"is_restricted,omitempty"`
		IsUltraRestricted bool   `json:"is_ultra_restricted,omitempty"`
		Has2Fa            bool   `json:"has_2fa,omitempty"`
	} `json:"members"`
	CacheTs int `json:"cache_ts"`
}

// Send function send message to talk room in slack
func (s *Slack) Send(message Message) error {
	var requestBuffer bytes.Buffer
	if err := json.NewEncoder(&requestBuffer).Encode(message); err != nil {
		return err
	}
	req, err := http.NewRequest("POST", s.webhook, &requestBuffer)
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

func (s *Slack) GetIDs(names []string) ([]string, error) {
	var IDs []string
	usersList := usersList{}
	endpoint := fmt.Sprintf("https://slack.com/api/users.list?token=%s&pretty=1", s.token)
	resp, err := http.Get(endpoint)
	if err != nil {
		return nil, err
	}
	if err := json.NewDecoder(resp.Body).Decode(&usersList); err != nil {
		return nil, err
	}
	for _, name := range names {
		ID := ""
		for _, user := range usersList.Members {
			if user.Name == name {
				ID = user.ID
				IDs = append(IDs, ID)
			}
		}
		if ID == "" {
			IDs = append(IDs, name)
		}
	}
	return IDs, nil
}
