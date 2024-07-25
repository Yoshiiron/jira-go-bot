package jiraFuncs

import (
	"fmt"
	"github.com/andygrunwald/go-jira"
	"log"
)

func UserReturner(username string) jira.User {
	client := JiraClient()
	users, _, err := client.User.Find(username, jira.WithUsername(username))
	if err != nil {
		fmt.Println(err)
	}
	var user1 jira.User
	for _, user := range users {
		user1 = user
	}
	return user1
}

func JiraClient() *jira.Client {
	tp := jira.PATAuthTransport{
		Token: "---",
	}
	jiraClient, err := jira.NewClient(tp.Client(), "---")
	if err != nil {
		log.Fatalf("Failed to create JIRA client: %s", err)
	}

	return jiraClient
}
