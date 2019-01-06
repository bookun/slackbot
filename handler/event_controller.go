package handler

import (
	"net/http"

	"github.com/kutsuzawa/slackbot/lib"
	"github.com/kutsuzawa/slackbot/models"
	"github.com/labstack/echo"
)


type Event interface {
	GetSenderAndTargets() (string, []string)
	MakeMessage(channel string, senderID string, targetIDs []string) (models.Message, error)
}

type Slack interface {
	GetIDs([]string) ([]string, error)
	GetChannel() string
	Send(message models.Message) error
}

type EventController struct {
	eventMap map[string]Event
	slack Slack
	util *lib.Util
}

func NewEventController(eventMap map[string]Event, slack Slack, util *lib.Util) *EventController {
	return &EventController{
		eventMap: eventMap,
		slack: slack,
		util: util,
	}
}

func (e *EventController) EventHandler(c echo.Context) error {
	eventName := c.Request().Header.Get("X-Github-Event")
	event := e.eventMap[eventName]
	if err := c.Bind(event); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	senderID, targetsIDs, err := e.getSlackSenderIDAndTargetIDs(event.GetSenderAndTargets())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	message, err := event.MakeMessage(e.slack.GetChannel(), senderID, targetsIDs)
	if err := e.slack.Send(message); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]string{"event": eventName})
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
