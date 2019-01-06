package events

import (
	"encoding/json"
	"fmt"
	"github.com/kutsuzawa/slackbot/models"
	"os"
	"reflect"
	"testing"
	"time"
)

var expectedMessage = models.Message{
	Name: "Pull Request",
	Channel: "test",
	LinkName: true,
	Attachments: []models.Attachment{
		{
			Pretext:    fmt.Sprintf("%s -> %s\nPR: %s\n", "octocat", "octocat", "new-feature"),
			Fallback:    fmt.Sprintf("%s -> %s\nPR: %s\n", "octocat", "octocat", "new-feature"),
			Color:      "good",
			AuthorName: "octocat",
			AuthorIcon: "https://github.com/images/error/octocat_happy.gif",
			AuthorLink:"https://api.github.com/users/octocat" ,
			Title:      "new-feature",
			TitleLink:  "https://api.github.com/repos/octocat/Hello-World/pulls/1347",
			Text:       "Please pull these awesome changes",
			Markdown:   true,
			Fields: []models.Field{
				{Title: "assignee", Value: "<@sender1>", Short: true},
				{Title: "reviewer", Value: "<@reviewer1>\n<@reviewer2>\n", Short: true},
			},
			ThumbURL:   "https://github.com/images/error/octocat_happy.gif",
			Footer:     "GitHub",
			FooterIcon: "https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcR7G9JTqB8z1AVU-Lq7xLy1fQ3RMO-Tt6PRplyhaw75XCAnYvAYxg",
			Ts:         timeParse("2011-01-26T19:01:12Z").Unix(),
		},
	},
}

func timeParse(timeStr string) time.Time{
	t, _ := time.Parse("2006-01-02T15:04:05Z", timeStr)
	return t
}

func TestPR_MakeMessage(t *testing.T) {
	cases := []struct {
		name string
		filepath string
		pr *PR
		expectedError error
		expectedMessage models.Message
	}{
		{
			name: "review_request",
			filepath: "../testdata/pull_request/review_requested.json",
			pr: &PR{},
			expectedError: nil,
			expectedMessage: expectedMessage,
		},
	}
	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			f, err := os.Open(testCase.filepath)
			if err != nil {
				t.Errorf("can not open testdata")
			}
			if err := json.NewDecoder(f).Decode(testCase.pr); err != nil {
				t.Errorf("can not decode")
			}
			message, err := testCase.pr.MakeMessage("test", "sender1", []string{"reviewer1", "reviewer2"})
			if err != testCase.expectedError {
				t.Errorf("want %v, got %v\n", testCase.expectedError, err)
			}
			if !reflect.DeepEqual(message, testCase.expectedMessage) {
				t.Errorf("want %v, got %v\n", testCase.expectedMessage, message)
			}
		})
	}

}

func TestPR_GetSenderAndTargets(t *testing.T) {
	cases := []struct{
		name string
		filepath string
		pr *PR
		expectedSender string
		expectedReviewers []string
	}{
		{
			name: "1 sender 3 reviewers",
			filepath: "../testdata/pull_request/review_requested.json",
			pr: &PR{},
			expectedSender: "octocat",
			expectedReviewers: []string{"octocat", "hubot", "other_user"},
		},
	}
	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			f, err := os.Open(testCase.filepath)
			if err != nil {
				t.Errorf("can not open testdata")
			}
			if err := json.NewDecoder(f).Decode(testCase.pr); err != nil {
				t.Errorf("can not decode")
			}
			sender, targets := testCase.pr.GetSenderAndTargets()
			if sender != testCase.expectedSender {
				t.Errorf("want %s, but got %s\n", testCase.expectedSender, sender)
			}
			if !reflect.DeepEqual(targets, testCase.expectedReviewers) {
				t.Errorf("want %v, but got %v\n", testCase.expectedReviewers, targets)
			}
		})
	}
}
