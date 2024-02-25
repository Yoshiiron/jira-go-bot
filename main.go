package main

import (
	"fmt"
	tele "gopkg.in/telebot.v3"
	jiraFuncs "jira-go-bot/JiraFuncs"
	"log"
	"strings"
	"time"
)

func main() {

	pref := tele.Settings{
		Token:  "6582917189:AAGBdLFPLNTFR09tjWGKaggk87RsEQTZ6fc",
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}

	b, err := tele.NewBot(pref)
	if err != nil {
		log.Fatal(err)
		return
	}

	b.Handle("/jira_acc", func(m tele.Context) error {
		message := m.Message().Text
		parts := strings.Fields(message)
		fmt.Println(parts[1])
		user := jiraFuncs.UserKeyReturner(jiraFuncs.JiraClient(), parts[1])
		if user == "" {
			m.Send("Такого пользователя не сущесвует, повторите попытку.")
		} else {
			m.Send(fmt.Sprintf("Key: %v, TelegramUID: %v", user, m.Sender().ID))
		}
		return nil
	})

	b.Handle("/start", func(m tele.Context) error {
		m.Send("Привет!\nЯ Джира-Бот, призван помочь твоей работе внутри Devim!\nЯ буду присылать тебе задачки из Jira прямо в телеграмм.")
		return nil
	})
	b.Start()
}
