package user

import (
	"log"
	"os"
	"strings"
)

type User struct {
	GithubName string
	SlackName  string
}

func NewUser(text string) *User {
	items := strings.Split(text, " ")
	user := &User{}
	user.GithubName = items[0]
	user.SlackName = items[1]
	return user
}

func (u *User) Add() {
	u.GithubName = strings.Replace(u.GithubName, "-", "_", -1)
	log.Printf("gihubname: %s\n", u.GithubName)
	log.Printf("slackname: %s\n", u.SlackName)
	err := os.Setenv(u.GithubName, u.SlackName)
	if err != nil {
		log.Printf("set env error: %s\n", err)
	}
	log.Printf("export %s -> %s\n", u.GithubName, os.Getenv(u.GithubName))
}
