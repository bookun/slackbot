package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/kutsuzawa/slackbot/pr"
	"github.com/kutsuzawa/slackbot/slack"
	"github.com/kutsuzawa/slackbot/user"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"

	"github.com/joho/godotenv"
)

// requestReviewHandle handle review requests from Github
func requestReviewHandler(c echo.Context) error {
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
			message, err := requestedPR.MakeJsonMessage("PullRequest", os.Getenv("CHANNEL"))
			if err != nil {
				log.Println(err)
				return err
			}
			slack := slack.NewSlack(os.Getenv("SLACKWEBHOOK"))
			if err = slack.Send(message); err != nil {
				log.Println(err)
				return err
			}
		}
	}
	return c.String(http.StatusOK, "succeed in sending PR to slack")
}

// userAddHandler handle user additional requests using a slash command from Slack
func userAddHandler(c echo.Context) error {
	commands := c.Request().FormValue("command")
	if commands == "/useradd" {
		user := user.NewUser(c.Request().FormValue("text"))
		user.Add()
	}
	return c.String(http.StatusOK, "succeed in adding user")
}

func welcomeHandler(c echo.Context) error {
	return c.String(http.StatusOK, "welcome to our slackbot")
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Not found .env file")
	}
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.POST("/commands", userAddHandler)
	e.POST("/pr", requestReviewHandler)
	e.GET("/", welcomeHandler)

	e.Logger.Fatal(e.Start(":"+os.Getenv("PORT")))
}
