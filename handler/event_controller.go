package handler

import (
	"github.com/kutsuzawa/slackbot/events"
	"github.com/kutsuzawa/slackbot/lib"
	"github.com/kutsuzawa/slackbot/models"
	"github.com/labstack/echo"
	"net/http"
)

type Event interface {
	GetSenderAndTargets() (string, []string)
	MakeMessage(channel string, senderID string, targetIDs []string) (models.Message, error)
}

type EventController struct {
	eventMap map[string]Event
	slack *models.Slack
	util *lib.Util
}

func NewEventController(slack *models.Slack, util *lib.Util) *EventController {
	return &EventController{
		eventMap: map[string]Event{
			"pull_request": &events.PR{},
		},
		slack: slack,
		util: util,
	}
}

func (e *EventController) EventHandler(c echo.Context) error {
	eventName := c.Request().Header.Get("X-Github-Event")
	event := e.eventMap[eventName]
	if err := c.Bind(event); err != nil {
		return err
	}
	senderID, targetsIDs, err := e.getSlackSenderIDAndTargetIDs(event.GetSenderAndTargets())
	if err != nil {
		return err
	}
	message, err := event.MakeMessage(e.slack.Channel, senderID, targetsIDs)
	if err := e.slack.Send(message); err != nil {
		return err
	}
	return c.String(http.StatusOK, "succeed at sending to slack")
}

func (e *EventController) getSlackSenderIDAndTargetIDs(sender string, targets []string) (string, []string, error) {
	slackSenderName := e.util.Translate(sender)
	slackSenderIDs, err := e.slack.GetIDs(slackSenderName)
	if err != nil {
		return "", nil, err
	}
	slackTargetNames := e.util.Translate(targets...)
	slackTargetIDs, err := e.slack.GetIDs(slackTargetNames)
	if err != nil {
		return "", nil, err
	}
	return slackSenderIDs[0], slackTargetIDs, nil
}
