package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/bookun/slackbot/pr"
	"github.com/bookun/slackbot/slack"
	"github.com/bookun/slackbot/user"

	"github.com/joho/godotenv"
)

// Adapter is function
type Adapter func(http.Handler) http.Handler

// Adapt function takes the handler you want to adapt
func Adapt(h http.Handler, adapters ...Adapter) http.Handler {
	for _, adapter := range adapters {
		h = adapter(h)
	}
	return h
}

// LoadEnv function load .env before handle function
func LoadEnv() Adapter {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			err := godotenv.Load()
			if err != nil {
				log.Println("Not found .env file")
			}
			h.ServeHTTP(w, r)
		})
	}
}

// requestReviewHandle handle review requests from Github
func requestReviewHandle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	event := r.Header.Get("X-Github-Event")
	if event == "pull_request" {
		var requestedPR pr.PR
		dec := json.NewDecoder(r.Body)
		if err := dec.Decode(&requestedPR); err != nil {
			log.Fatal(err)
		}
		if requestedPR.Action == "review_requested" {
			message, err := requestedPR.MakeJsonMessage("PullRequest", os.Getenv("CHANNEL"))
			if err != nil {
				log.Println(err)
			}
			slack := slack.NewSlack(os.Getenv("SLACKWEBHOOK"))
			if err = slack.Send(message); err != nil {
				log.Println(err)
			}
		}
	}
}

// userAddHandler handle user additional requests using a slash command from Slack
func userAddHandler(w http.ResponseWriter, r *http.Request) {
	commands := r.FormValue("command")
	if commands == "/useradd" {
		user := user.NewUser(r.FormValue("text"))
		user.Add()
	}
}

func main() {
	http.HandleFunc("/commands", userAddHandler)
	handler := http.HandlerFunc(requestReviewHandle)
	http.Handle("/", Adapt(handler, LoadEnv()))
	http.ListenAndServe(":"+os.Getenv("PORT"), nil)
}
