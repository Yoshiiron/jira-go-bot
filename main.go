package main

import (
	"context"
	"fmt"
	"github.com/rusq/tbcomctl/v4"
	"github.com/thanhpk/randstr"
	tele "gopkg.in/telebot.v3"
	"jira-go-bot/Sender"
	"log"
	"time"
)

var (
	Menu       = &tele.ReplyMarkup{ResizeKeyboard: true}
	SecondMenu = &tele.ReplyMarkup{ResizeKeyboard: true}

	btnHelp = Menu.Text("Активация.")

	LoginJira = Menu.Text("Ввести логин Jira.")
	SendCode  = Menu.Text("Отправить код на почту.")
	EnterCode = Menu.Text("Ввести код.")
)

//func random() string {
//	rand.New(rand.NewSource())
//	chars := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ" +
//		"abcdefghijklmnopqrstuvwxyz" +
//		"0123456789")
//	length := 8
//	var b strings.Builder
//	for i := 0; i < length; i++ {
//		b.WriteRune(chars[rand.Intn(len(chars))])
//	}
//	str := b.String() // Например "ExcbsVQs"
//	return str
//}

func randomq() string {
	return randstr.String(6)
}

//
//func JiraLogining(b *tele.Bot) {
//	codeSend := tbcomctl.NewInputText("code", "Пожалуйста, введите ваш логин Jira:  ", processInput(b))
//	sendForm := tbcomctl.NewForm(codeSend)
//
//	b.Handle(&LoginJira, sendForm.Handler)
//	b.Handle(tele.OnText, sendForm.OnTextMiddleware(func(c tele.Context) error {
//		log.Printf("user data: %v", sendForm.Data(c.Sender())["code"])
//		log.Printf("%v", sendForm.Data(c.Sender()))
//		fmt.Println("Я тоже выполняюсь.")
//		return nil
//	}))
//
//}

func UserRegistration(b *tele.Bot) {

	SecondMenu.Reply(
		SecondMenu.Row(LoginJira),
		SecondMenu.Row(SendCode),
		SecondMenu.Row(EnterCode),
	)

	b.Handle(&btnHelp, func(c tele.Context) error {
		c.Send("1", SecondMenu)
		return nil
	})
	//
	var mail, code string

	b.Handle(&SendCode, func(c tele.Context) error {
		fmt.Println(c.Message().Text)
		code = randomq()
		Sender.SendMessage(code, mail)
		return nil
	})

}

func main() {

	Menu.Reply(
		Menu.Row(btnHelp),
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

	nameIp := tbcomctl.NewInputText("name", "Input your name:", processInput(b))
	ageIp := tbcomctl.NewInputText("age", "Input your age", processInput(b))

	form := tbcomctl.NewForm(nameIp, ageIp)
	b.Handle("/input", form.Handler)
	// b.Handle(tb.OnText, tbcomctl.NewMiddlewareChain(onText, nameIp.OnTextMw, ageIp.OnTextMw))
	b.Handle(tele.OnText, form.OnTextMiddleware(func(c tele.Context) error {
		log.Printf("onText is called: %q\nuser data: %v", c.Message().Text, form.Data(c.Sender()))
		return nil
	}))

	b.Start()

	////Sender.SendMessage("Hello")
	////
	//b.Handle("/start", func(m tele.Context) error {
	//	m.Send("Привет!\nЯ Джира-Бот, призван помочь твоей работе внутри Devim!\nЯ буду присылать тебе задачки из Jira прямо в телеграмм.", Menu)
	//	UserRegistration(b)
	//	return nil
	//})
	//
	//codeSend := tbcomctl.NewInputText("code", "Пожалуйста, введите ваш логин Jira:  ", processInput(b))
	//sendForm := tbcomctl.NewForm(codeSend)
	//b.Handle(&LoginJira, sendForm.Handler)
	//b.Handle(tele.OnText, func(c tele.Context) error {
	//	c.Send("Wha")
	//	fmt.Println(c.Message().IsReply())
	//	if c.Message().IsReply() == true {
	//		b.Handle(tele.OnText, sendForm.OnText(func(c tele.Context) error {
	//			log.Printf("user data: %v", sendForm.Data(c.Sender()))
	//			log.Printf("%v", sendForm.Data(c.Sender()))
	//			fmt.Println("Я тоже выполняюсь.")
	//			return nil
	//		}))
	//	}

	//return nil

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
