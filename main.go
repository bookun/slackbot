package main

import (
	"fmt"
	"github.com/kutsuzawa/slackbot/events"
	"github.com/kutsuzawa/slackbot/lib"
	"log"
	"net/http"
	"os"

	"github.com/kutsuzawa/slackbot/handler"
	"github.com/kutsuzawa/slackbot/models"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func healthHandler(c echo.Context) error {
	return c.String(http.StatusOK, "OK")
}

func checkEnv() error {
	envs := []string{"CHANNEL", "SLACKWEBHOOK", "PORT", "TOKEN"}
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
	port := os.Getenv("PORT")
	slack := models.NewSlack(os.Getenv("SLACKWEBHOOK"), os.Getenv("CHANNEL"), os.Getenv("TOKEN"))
	eventMap := make(map[string]handler.Event)
	eventMap["pull_request"] = &events.PR{}
	controller := handler.NewEventController(eventMap, slack, &lib.Util{})

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.POST("/event", controller.EventHandler)
	e.GET("/", healthHandler)
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", port)))
}
