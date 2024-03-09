package main

import (
	"fmt"
	jiraFuncs "jira-go-bot/JiraFuncs"
)

func main() {
	//tp := jira.PATAuthTransport{
	//	Token: "MjQ3NDcxMjQ2MDE2OrGi1v39Ob9i9T58zhBIFYO0Lx3E",
	//}
	//jiraClient, err := jira.NewClient(tp.Client(), "https://jira.yoshiiron.space")
	//if err != nil {
	//	log.Fatalf("Failed to create JIRA client: %s", err)
	//}

	//a, _, _ := jiraClient.User.Find("Yoshiiron", jira.WithUsername("Yoshiiron"))
	b := jiraFuncs.UserReturner("QQQ")
	if b.Name == "" {
		fmt.Println("Пользователя не существует.")
	}

}
