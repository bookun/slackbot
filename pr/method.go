package pr

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"

	"github.com/bookun/slackbot/util"
)

var utils = &util.Util{}

type Field struct {
	Title string `json:"title"`
	Value string `json:"value"`
	Short bool   `json:"short"`
}

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
	//ImageURL   string `json:"image_url"`
	ThumbURL   string `json:"thumb_url"`
	Footer     string `json:"footer"`
	FooterIcon string `json:"footer_icon"`
	Ts         int64  `json:"ts"`
}

type Text struct {
	Name        string       `json:"username"`
	Channel     string       `json:"channel"`
	LinkName    bool         `json:"link_names"`
	Attachments []Attachment `json:"attachments"`
}

func (p *PR) MakeJsonMessage(name, channel string) *bytes.Buffer {
	tp := &Text{}
	tp.Name = name
	tp.Channel = channel
	tp.LinkName = true
	//tp.Attachments.Fields = make([]struct, 1)
	attachment := &Attachment{}
	tp.Attachments = append(tp.Attachments, p.makeAttachment(attachment))
	return p.jsonEncode(*tp)
}

func (p *PR) makeAttachment(attachment *Attachment) Attachment {
	sender := p.PullRequest.User
	reviewer := p.PullRequest.RequestedReviewers[0]
	title := p.PullRequest.Title
	log.Printf("translate %s -> %s\n", reviewer.Login, utils.Translate(reviewer.Login))
	attachment.Pretext = fmt.Sprintf("%s -> %s\nPR: %s\n", sender.Login, utils.Translate(reviewer.Login), title)
	attachment.Fallback = attachment.Pretext
	attachment.Color = "good"
	attachment.AuthorName = sender.Login
	attachment.AuthorLink = sender.URL
	attachment.AuthorIcon = sender.AvatarURL
	attachment.Title = title
	attachment.TitleLink = p.PullRequest.URL
	attachment.Text = p.PullRequest.Body
	attachment.Markdown = true
	//attachment.Fields[0].Title = "assignee"
	//attachment.Fields[0].Value = sender.Login
	//attachment.Fields[0].Short = true
	attachment.Fields = append(attachment.Fields, Field{Title: "assignee", Value: sender.Login, Short: true})
	attachment.Fields = append(attachment.Fields, Field{Title: "reviewer", Value: utils.Translate(reviewer.Login), Short: true})
	attachment.ThumbURL = attachment.AuthorIcon
	attachment.Footer = "GitHub"
	attachment.FooterIcon = "https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcR7G9JTqB8z1AVU-Lq7xLy1fQ3RMO-Tt6PRplyhaw75XCAnYvAYxg"
	attachment.Ts = (p.PullRequest.UpdatedAt).Unix()
	return *attachment
}

func (p *PR) jsonEncode(text Text) *bytes.Buffer {
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	if err := enc.Encode(text); err != nil {
		log.Fatal(err)
	}
	return bytes.NewBuffer(buf.Bytes())
}
