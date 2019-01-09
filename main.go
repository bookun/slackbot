package main

import (
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"

	"github.com/kutsuzawa/slackbot/events"
	"github.com/kutsuzawa/slackbot/lib"

	"github.com/kutsuzawa/slackbot/handler"
	"github.com/kutsuzawa/slackbot/models"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func healthHandler(c echo.Context) error {
	return c.String(http.StatusOK, "OK")
}

func checkEnv() error {
	envs := []string{"CHANNEL", "SLACKWEBHOOK", "TOKEN"}
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
	slack := models.NewSlack(os.Getenv("SLACKWEBHOOK"), os.Getenv("CHANNEL"), os.Getenv("TOKEN"))
	eventMap := make(map[string]handler.Event)
	eventMap["pull_request"] = &events.PR{}
	controller := handler.NewEventController(eventMap, slack, &lib.Util{})

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.POST("/event", controller.EventHandler)
	e.GET("/", healthHandler)
	debugPort := os.Getenv("DEBUG_PORT")
	if debugPort == "" {
		debugPort = "6060"
	}
	go func() {
		e.Logger.Print(http.ListenAndServe(fmt.Sprintf(":%s", debugPort), nil))
	}()

	port := os.Getenv("PORT")
	if port == "" {
		port = "3030"
	}
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", port)))
}
