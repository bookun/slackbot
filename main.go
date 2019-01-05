package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/kutsuzawa/slackbot/pr"
	"github.com/kutsuzawa/slackbot/slack"
	"github.com/kutsuzawa/slackbot/user"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

type handler struct {
	slackWebhook string
	slackChannel string
}

// requestReviewHandle handle review requests from Github
func (h *handler) requestReviewHandler(c echo.Context) error {
	c.Response().Header().Set("Content-Type", "application/json; charset=utf-8")
	event := c.Request().Header.Get("X-Github-Event")
	if event == "pull_request" {
		var requestedPR pr.PR
		dec := json.NewDecoder(c.Request().Body)
		if err := dec.Decode(&requestedPR); err != nil {
			log.Println(err)
			return err
		}
		if requestedPR.Action == "review_requested" {
			message, err := requestedPR.MakeJsonMessage("PullRequest", h.slackChannel)
			if err != nil {
				log.Println(err)
				return err
			}
			slack := slack.NewSlack(h.slackWebhook)
			if err = slack.Send(message); err != nil {
				log.Println(err)
				return err
			}
		}
	}
	return c.String(http.StatusOK, "succeed in sending PR to slack")
}

// userAddHandler handle user additional requests using a slash command from Slack
func (h *handler) userAddHandler(c echo.Context) error {
	commands := c.Request().FormValue("command")
	if commands == "/useradd" {
		user := user.NewUser(c.Request().FormValue("text"))
		user.Add()
	}
	return c.String(http.StatusOK, "succeed in adding user")
}

func (h *handler) welcomeHandler(c echo.Context) error {
	return c.String(http.StatusOK, "welcome to our slackbot")
}

func checkEnv() error {
	envs := []string{"CHANNEL", "SLACKWEBHOOK", "PORT"}
	for _, v := range envs {
		if os.Getenv(v) == "" {
			err := fmt.Errorf("env variable %s is not defined", v)
			return err
		}
	}
	return nil
}

func main() {
	if err := checkEnv(); err != nil {
		log.Fatal(err)
	}
	handler := &handler{
		slackWebhook: os.Getenv("SLACKWEBHOOK"),
		slackChannel: os.Getenv("CHANNEL"),
	}
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.POST("/commands", handler.userAddHandler)
	e.POST("/pr", handler.requestReviewHandler)
	e.GET("/", handler.welcomeHandler)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", port)))
}
