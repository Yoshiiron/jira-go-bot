package main

import (
	"flag"
	"github.com/vitaliy-ukiru/fsm-telebot"
	"github.com/vitaliy-ukiru/fsm-telebot/storages/memory"
	tele "gopkg.in/telebot.v3"
	"gopkg.in/telebot.v3/middleware"
	"log"
	"time"
)

var (
	JiraAccountInput = fsm.NewStateGroup("accountinput")
	JiraAccountMail  = JiraAccountInput.New("mail")
	JiraCode         = JiraAccountInput.New("code")
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
	manager.Bind("/return", fsm.DefaultState, OnReturn)

	// buttons
	manager.Bind(&ActivationBtn, fsm.DefaultState, OnActivation)
	manager.Bind(&returnBtn, fsm.DefaultState, OnReturn)

	// form
	manager.Bind(tele.OnText, JiraAccountMail, )

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

	state.Set(JiraAccountMail)
	return c.Send("Напиши мне свой логин от Jira.", menu)
}

func


func OnReturn(c tele.Context, state fsm.Context) error {
	menu := &tele.ReplyMarkup{}
	menu.Reply(menu.Row(ActivationBtn))
	menu.ResizeKeyboard = true

	go state.Finish(true)
	return c.Send("Активация была отменена.", menu)
}

func DeleteAfterHandler(next tele.HandlerFunc) tele.HandlerFunc {
	return func(c tele.Context) error {
		defer func(c tele.Context) {
			if err := c.Delete(); err != nil {
				c.Bot().OnError(err, c)
			}
		}(c)
		return next(c)
	}
}
