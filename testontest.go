package main

import (
	"flag"
	"fmt"
	"github.com/nats-io/nats.go"
	"github.com/thanhpk/randstr"
	"github.com/vitaliy-ukiru/fsm-telebot"
	"github.com/vitaliy-ukiru/fsm-telebot/storages/memory"
	tele "gopkg.in/telebot.v3"
	"gopkg.in/telebot.v3/middleware"
	jiraFuncs "jira-go-bot/JiraFuncs"
	"jira-go-bot/Sender"
	"log"
	"time"
)

var (
	JiraAccountInput = fsm.NewStateGroup("accountinput")
	SendedCode       = JiraAccountInput.New("code")
	InputedCode      = JiraAccountInput.New("incode")
)

var debug = flag.Bool("debug", false, "log debug info")

var (
	ActivationBtn = tele.Btn{Text: "Активация"}
	returnBtn     = tele.Btn{Text: "Вернуться назад"}
)

func main() {
	flag.Parse()

	bot, err := tele.NewBot(tele.Settings{
		Token:     "6582917189:AAGBdLFPLNTFR09tjWGKaggk87RsEQTZ6fc",
		Poller:    &tele.LongPoller{Timeout: 10 * time.Second},
		ParseMode: tele.ModeHTML,
		Verbose:   *debug,
		OnError: func(err error, context tele.Context) {
			log.Printf("{ERR} %q chat=%s", err, context.Recipient())
		},
	})
	if err != nil {
		log.Fatalln(err)
	}
	storage := memory.NewStorage()
	defer storage.Close()

	manager := fsm.NewManager(bot, nil, storage, nil)

	bot.Use(middleware.AutoRespond())

	// commands
	bot.Handle("/start", OnStart(ActivationBtn))
	manager.Bind("/jiraAcc", fsm.DefaultState, OnActivation)
	manager.Bind("/return", fsm.AnyState, OnReturn)

	manager.Bind("/state", fsm.AnyState, func(c tele.Context, state fsm.Context) error {
		s, err := state.State()
		if err != nil {
			return c.Send(fmt.Sprintf("can't get state: %s", err))
		}
		return c.Send(s.GoString())
	})

	nc, err := nats.Connect("nats://192.168.1.121:4222")
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()
	nc.Subscribe("foo", func(m *nats.Msg) {
		fmt.Printf("Received a message: %s\n", string(m.Data))
		bot.Send(tele.ChatID(435902334), jiraFuncs.UserReturner(jiraFuncs.JiraClient(), fmt.Sprintf("%v", string(m.Data))).Key)
	})

	// buttons
	manager.Bind(&ActivationBtn, fsm.DefaultState, OnActivation)
	manager.Bind(&returnBtn, fsm.AnyState, OnReturn)

	// form
	manager.Bind(tele.OnText, SendedCode, OnInputJiraLogin)
	manager.Bind(tele.OnText, InputedCode, OnInputCode)

	bot.Start()
}

func OnStart(start tele.Btn) tele.HandlerFunc {
	menu := &tele.ReplyMarkup{}
	menu.Reply(menu.Row(start))
	menu.ResizeKeyboard = true

	return func(c tele.Context) error {
		log.Println("new user", c.Sender().ID)
		return c.Send(
			"Привет!"+
				"\nЯ Джира-Бот, призван помочь твоей работе внутри Devim!"+
				"\nЯ буду присылать тебе задачки из Jira прямо в телеграмм.", menu)
	}
}

func OnActivation(c tele.Context, state fsm.Context) error {
	menu := &tele.ReplyMarkup{}
	menu.Reply(menu.Row(returnBtn))
	menu.ResizeKeyboard = true

	state.Set(SendedCode)
	return c.Send("Напиши мне свой логин от Jira.", menu)
}

func OnInputJiraLogin(c tele.Context, state fsm.Context) error {
	Login := c.Message().Text
	jirauser := jiraFuncs.UserReturner(jiraFuncs.JiraClient(), Login)
	code := randstr.String(6)
	Sender.SendMessage(code, jirauser.EmailAddress)

	go state.Update("code", code)
	go state.Set(InputedCode)

	return c.Send(fmt.Sprintf("На вашу почту %v был направлен код, введите его.", jirauser.EmailAddress))
}

func OnInputCode(c tele.Context, state fsm.Context) error {
	menu := &tele.ReplyMarkup{}
	menu.Reply(menu.Row(ActivationBtn))
	menu.ResizeKeyboard = true

	defer state.Finish(true)
	var (
		sendedCode  string
		inputedCode string
	)
	state.MustGet("code", &sendedCode)
	inputedCode = c.Message().Text

	fmt.Println(sendedCode, inputedCode)

	result := "Произошла ошибка при регистрации. Проверьте правильность введённого вами кода."

	if sendedCode == inputedCode {
		result = "Вы были успешно зарегестрированы."
	}
	return c.Send(result, menu)
}

func OnReturn(c tele.Context, state fsm.Context) error {
	menu := &tele.ReplyMarkup{}
	menu.Reply(menu.Row(ActivationBtn))
	menu.ResizeKeyboard = true

	go state.Finish(true)
	return c.Send("Активация была отменена.", menu)
}
