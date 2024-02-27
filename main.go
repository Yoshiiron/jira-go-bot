package main

import (
	"context"
	"fmt"
	"github.com/rusq/tbcomctl/v4"
	"github.com/sethvargo/go-password/password"
	tele "gopkg.in/telebot.v3"
	"jira-go-bot/Sender"
	"log"
	"math/rand"
	"time"
)

var (
	Menu       = &tele.ReplyMarkup{ResizeKeyboard: true}
	SecondMenu = &tele.ReplyMarkup{ResizeKeyboard: true}

	btnHelp     = Menu.Text("ℹ Help")
	btnSettings = Menu.Text("⚙ Settings")

	SendCode  = Menu.Text("Отправить код на почту.")
	EnterCode = Menu.Text("Ввести код.")
)

func randomizer() string {
	var seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))
	length := 6
	charset := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	password := make([]byte, length)
	for i := range password {
		password[i] = charset[seededRand.Intn(len(charset))]
	}

	return string(password)
}

func UserRegistration(b *tele.Bot) {

	SecondMenu.Reply(
		SecondMenu.Row(SendCode),
		SecondMenu.Row(EnterCode),
	)

	codeInput := tbcomctl.NewInputText("code", "Пожалуйста, введите код отправленный вам на mail.yandex.ru: ", processInput(b))
	form := tbcomctl.NewForm(codeInput)

	str, err := password.Generate(6, 2, 0, false, false)
	if err != nil {
		return
	}

	b.Handle(&btnHelp, func(c tele.Context) error {
		c.Send("1", SecondMenu)
		return nil
	})

	b.Handle(&SendCode, func(c tele.Context) error {
		Sender.SendMessage(str)
		return nil
	})

	b.Handle(&EnterCode, form.Handler)

	b.Handle(tele.OnText, func(c tele.Context) error {
		if c.Message().IsReply() == true {
			if str == c.Message().Text {
				c.Send("Вы зарегистрированы.")
			}
		}
		return nil
	})
}

func main() {

	Menu.Reply(
		Menu.Row(btnHelp),
		Menu.Row(btnSettings),
	)

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
	//Sender.SendMessage("Hello")
	//
	b.Handle("/start", func(m tele.Context) error {
		m.Send("Привет!\nЯ Джира-Бот, призван помочь твоей работе внутри Devim!\nЯ буду присылать тебе задачки из Jira прямо в телеграмм.", Menu)
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
