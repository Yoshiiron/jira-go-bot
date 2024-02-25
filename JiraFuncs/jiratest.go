package jiraFuncs

import (
	"fmt"
	"github.com/andygrunwald/go-jira"
	"log"
)

func UserKeyReturner(client *jira.Client, username string) string {
	users, _, err := client.User.Find(username, jira.WithUsername(username))
	if err != nil {
		fmt.Println(err)
	}
	var key string
	for _, user := range users {
		key = user.Key
	}
	return key
}

func JiraClient() *jira.Client {
	tp := jira.PATAuthTransport{
		Token: "MjQ3NDcxMjQ2MDE2OrGi1v39Ob9i9T58zhBIFYO0Lx3E",
	}
	jiraClient, err := jira.NewClient(tp.Client(), "https://jira.yoshiiron.space")
	if err != nil {
		log.Fatalf("Failed to create JIRA client: %s", err)
	}

	return jiraClient
}
