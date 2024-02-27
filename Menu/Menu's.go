package Menu

import (
	"context"
	"fmt"
	"github.com/rusq/tbcomctl/v4"
	tele "gopkg.in/telebot.v3"
	"log"
)

var (
	Menu       = &tele.ReplyMarkup{ResizeKeyboard: true}
	SecondMenu = &tele.ReplyMarkup{ResizeKeyboard: true}

	btnHelp     = Menu.Text("ℹ Help")
	btnSettings = Menu.Text("⚙ Settings")

	SendCode  = Menu.Text("Отправить код.")
	EnterCode = Menu.Text("Ввести код.")
)

func MenuReturner(b *tele.Bot) {

	Menu.Reply(
		Menu.Row(btnHelp),
		Menu.Row(btnSettings),
	)
	SecondMenu.Reply(
		SecondMenu.Row(SendCode),
		SecondMenu.Row(EnterCode),
	)

	code := tbcomctl.NewInputText("Code", "Пожалуйста, введите код: ", processInput(b))
	form := tbcomctl.NewForm(code)

	b.Handle(&btnHelp, func(c tele.Context) error {
		c.Send("1", SecondMenu)
		return nil
	})

	b.Handle(&EnterCode, form.Handler)

	iscode := "123"
	b.Handle(tele.OnText, func(c tele.Context) error {
		if c.Message().IsReply() == true {
			if iscode == c.Message().Text {
				c.Send("Вы зарегистрированы.")
			}
		}
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
