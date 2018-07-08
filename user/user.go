package user

import (
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
	u.GithubName = strings.Replace(u.GithubName, "_", "-", -1)
	os.Setenv(u.GithubName, u.SlackName)
}
