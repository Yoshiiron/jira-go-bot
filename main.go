package main

import (
	"context"
	"fmt"
	"github.com/rusq/tbcomctl/v4"
	tele "gopkg.in/telebot.v3"
	"jira-go-bot/Sender"
	"log"
	"math/rand"
	"strings"
	"time"
)

func randomizer() string {
	rand.Seed(time.Now().UnixNano())
	chars := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ" +
		"abcdefghijklmnopqrstuvwxyz" +
		"0123456789")
	length := 8
	var build strings.Builder
	for i := 0; i < length; i++ {
		build.WriteRune(chars[rand.Intn(len(chars))])
	}
	str := build.String()

	return str
}

func UserRegistration(b *tele.Bot) {
	codeInput := tbcomctl.NewInputText("code", "Пожалуйста, введите код отправленный вам на mail.yandex.ru", processInput(b))
	form := tbcomctl.NewForm(codeInput)

	b.Handle("/jira_acc", form.Handler)
	str := randomizer()
	Sender.SendMessage(str)

	b.Handle(tele.OnText, func(c tele.Context) error {
		if c.Message().IsReply() {
			if c.Message().Text == str {
				c.Send("Прекрасно.")
			}
		}
		return nil
	})
}

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

	UserRegistration(b)
	//
	b.Handle("/start", func(m tele.Context) error {
		m.Send("Привет!\nЯ Джира-Бот, призван помочь твоей работе внутри Devim!\nЯ буду присылать тебе задачки из Jira прямо в телеграмм.")
		return nil
	})

	b.Start()
}

func processInput(b *tele.Bot) func(ctx context.Context, c tele.Context) error {
	return func(ctx context.Context, c tele.Context) error {
		val := c.Message().Text
		log.Println("msgCb function is called, input value:", val)
		switch val {
		case "error":
			return fmt.Errorf("error requested: %s", val)
		case "wrong":
			return tbcomctl.NewInputError("wrong input")
		}
		if ctrl, ok := tbcomctl.ControllerFromCtx(ctx); ok {
			log.Println("form values so far: ", ctrl.Form().Data(c.Sender()))
		}
		return nil
	}
}
