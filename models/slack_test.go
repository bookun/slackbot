package models

import (
	"fmt"
	"os"
	"testing"
)

var (
	flagGetIDs = true
	slack      = &Slack{}
)
var message = Message{
	Name:     "Pull Request",
	Channel:  "test",
	LinkName: true,
	Attachments: []Attachment{
		{
			Pretext:    fmt.Sprintf("%s -> %s\nPR: %s\n", "octocat", "octocat", "new-feature"),
			Fallback:   fmt.Sprintf("%s -> %s\nPR: %s\n", "octocat", "octocat", "new-feature"),
			Color:      "good",
			AuthorName: "octocat",
			AuthorIcon: "https://github.com/images/error/octocat_happy.gif",
			AuthorLink: "https://api.github.com/users/octocat",
			Title:      "Test of Slackbot",
			TitleLink:  "https://api.github.com/repos/octocat/Hello-World/pulls/1347",
			Text:       "Please pull these awesome changes",
			Markdown:   true,
			Fields: []Field{
				{Title: "assignee", Value: "<@sender1>", Short: true},
				{Title: "reviewer", Value: "<@reviewer1>\n<@reviewer2>\n", Short: true},
			},
			ThumbURL:   "https://github.com/images/error/octocat_happy.gif",
			Footer:     "GitHub",
			FooterIcon: "https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcR7G9JTqB8z1AVU-Lq7xLy1fQ3RMO-Tt6PRplyhaw75XCAnYvAYxg",
			Ts:         1296068472,
		},
	},
}

func TestMain(m *testing.M) {
	slack = NewSlack("https://hooks.slack.com/services/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",
		"test",
		os.Getenv("TOKEN"))
	m.Run()
}

func TestSlack_Send(t *testing.T) {
	if err := slack.Send(message); err != nil {
		t.Errorf("want %v, got %v\n", nil, err)
	}
}
