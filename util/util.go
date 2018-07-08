package util

import (
	"log"
	"os"
	"strings"
)

// Util is empty struct
type Util struct{}

// Translate function traslate githubName into slackName
func (u *Util) Translate(githubName string) string {
	replacedGithubName := strings.Replace(githubName, "-", "_", -1)
	log.Printf("in Translate method %s -> %s\n", replacedGithubName, os.Getenv(replacedGithubName))
	if slackName := os.Getenv(replacedGithubName); slackName != "" {
		return slackName
	}
	return githubName
}
