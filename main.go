package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/bookun/slackbot/pr"
	"github.com/bookun/slackbot/slack"
	"github.com/bookun/slackbot/user"
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

// SetHeader set "content-type" header
func SetHeader() Adapter {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			h.ServeHTTP(w, r)
		})
	}
}

func handle(w http.ResponseWriter, r *http.Request) {
	event := r.Header.Get("X-Github-Event")
	if event == "pull_request" {
		var requestedPR pr.PR
		dec := json.NewDecoder(r.Body)
		if err := dec.Decode(&requestedPR); err != nil {
			log.Fatal(err)
		}
		if requestedPR.Action == "review_requested" {
			//fmt.Fprintln(w, requestedPR.PullRequest.URL)
			//fmt.Fprintln(w, requestedPR.PullRequest.User.Login)
			//for _, v := range requestedPR.PullRequest.RequestedReviewers {
			//	fmt.Fprintln(w, v.Login)
			//}
			//fmt.Fprintln(w, requestedPR.PullRequest.Title)
			//fmt.Fprintln(w, requestedPR.PullRequest.Body)
			message := requestedPR.MakeJsonMessage("PullRequest", "general")
			slack := slack.NewSlack("https://hooks.slack.com/services/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")
			slack.Send(message)
		}
	} else {
		fmt.Fprintln(w, "not")
	}
}

func userAddHandler(w http.ResponseWriter, r *http.Request) {
	//token=gIkuvaNzQIHg97ATvDxqgjtO
	//&team_id=T0001
	//&team_domain=example
	//&enterprise_id=E0001
	//&enterprise_name=Globular%20Construct%20Inc
	//&channel_id=C2147483705
	//&channel_name=test
	//&user_id=U2147483697
	//&user_name=Steve
	//&command=/weather
	//&text=94070
	//&response_url=https://hooks.slack.com/commands/1234/5678
	//&trigger_id=13345224609.738474920.8088930838d88f008e0
	commands := r.FormValue("command")
	if commands == "/useradd" {
		//text := r.FormValue("text")
		//token := r.FormValue("token")
		user := user.NewUser(r.FormValue("text"))
		user.Add()
	}
}

func main() {
	handler := http.HandlerFunc(handle)
	http.HandleFunc("/commands", userAddHandler)
	http.Handle("/", Adapt(handler, SetHeader()))
	http.ListenAndServe(":"+os.Getenv("PORT"), nil)
}
