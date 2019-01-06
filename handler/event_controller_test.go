package handler

import (
	"bytes"
	"fmt"
	"github.com/kutsuzawa/slackbot/events"
	"github.com/kutsuzawa/slackbot/lib"
	"github.com/kutsuzawa/slackbot/models"
	"github.com/labstack/echo"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

type MockEvent struct {
	Name string `json:"name"`
}

func (e *MockEvent) GetSenderAndTargets() (string, []string) {
	return "sender", []string{"target1", "target2"}
}

func (e *MockEvent) MakeMessage(c string, sID string, tIDs []string) (models.Message, error) {
	return models.Message{}, nil
}

type MockSlack struct{}

func (s *MockSlack) GetIDs(names []string) ([]string, error) {
	return []string{"id1", "id2"}, nil
}

func (s *MockSlack) GetChannel() string {
	return "test_channel"
}

func (s *MockSlack) Send(message models.Message) error {
	return nil
}

type MockSlackError MockSlack

func (s *MockSlackError) GetIDs(names []string) ([]string, error) {
	return []string{"id1", "id2"}, nil
}

func (s *MockSlackError) GetChannel() string {
	return "test_channel"
}

func (s *MockSlackError) Send(message models.Message) error {
	return fmt.Errorf("send error")
}

func TestEventController_EventHandler(t *testing.T) {
	e := echo.New()
	eventMap := make(map[string]Event)
	eventMap["mock"] = &MockEvent{}
	eventMap["pull_request"] = &events.PR{}
	cases := []struct {
		name          string
		filepath string
		event string
		controller    *EventController
		expectedCode int
		expectedJSON  string
	}{
		{
			name:          "success",
			filepath: "../testdata/mock.json",
			event: "mock",
			controller:    NewEventController(eventMap, &MockSlack{}, &lib.Util{}),
			expectedCode: http.StatusOK,
			expectedJSON:  `{"event":"mock"}`,
		},
		{
			name:          "send error",
			event: "mock",
			filepath: "../testdata/mock.json",
			controller:    NewEventController(eventMap, &MockSlackError{}, &lib.Util{}),
			expectedCode: http.StatusInternalServerError,
			expectedJSON:  `{"error":"send error"}`,
		},
		{
			name:          "review_request",
			event: "pull_request",
			filepath: "../testdata/pull_request/review_requested.json",
			controller:    NewEventController(eventMap, &MockSlack{}, &lib.Util{}),
			expectedCode: http.StatusOK,
			expectedJSON:  `{"event":"pull_request"}`,
		},
	}
	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			data, err := ioutil.ReadFile(testCase.filepath)
			if err != nil {
				t.Errorf("can not open testdata")
			}
			req := httptest.NewRequest(http.MethodPost, "/event", bytes.NewBuffer(data))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			req.Header.Set("X-Github-Event", testCase.event)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			if assert.NoError(t, testCase.controller.EventHandler(c)) {
				assert.Equal(t, testCase.expectedCode, rec.Code)
				assert.Equal(t, testCase.expectedJSON, rec.Body.String())
			}
		})
	}
}
