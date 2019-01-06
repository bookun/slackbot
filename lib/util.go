package lib

import (
	"os"
	"strings"
)

// Util is empty struct
type Util struct{}

// Translate function translate githubName into slackName
func (u *Util) Translate(githubNames ...string) []string {
	var results []string
	for _, name := range githubNames {
		replacedGithubName := strings.Replace(name, "-", "_", -1)
		if slackName := os.Getenv(replacedGithubName); slackName != "" {
			results = append(results, slackName)
		} else {
			results = append(results, name)
		}
	}
	return results
}
