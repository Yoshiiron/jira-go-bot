package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/nats-io/nats.go"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/thanhpk/randstr"
	"github.com/vitaliy-ukiru/fsm-telebot"
	"github.com/vitaliy-ukiru/fsm-telebot/storages/memory"
	tele "gopkg.in/telebot.v3"
	"gopkg.in/telebot.v3/middleware"
	jiraFuncs "jira-go-bot/JiraFuncs"
	"jira-go-bot/Sender"
	"log"
	"strconv"
	"time"
)

var (
	JiraAccountInput = fsm.NewStateGroup("accountinput")
	SendedCode       = JiraAccountInput.New("code")
	InputedCode      = JiraAccountInput.New("incode")
	JiraLogin        = JiraAccountInput.New("login")
)

var debug = flag.Bool("debug", false, "log debug info")

var (
	ActivationBtn = tele.Btn{Text: "Активация"}
	returnBtn     = tele.Btn{Text: "Вернуться назад"}
)

type RecievedMessageFromHooks struct {
	To      string `json:"to"`
	Message string `json:"message"`
}

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

	nc, err := nats.Connect("nats://nats:4222")
	if err != nil {
		log.Fatal(err)
	}

	defer nc.Close()

	nc.Subscribe("foo", func(m *nats.Msg) {
		var hooksmessage RecievedMessageFromHooks
		//hooksmessage = json.Unmarshal(m.Data)
		//fmt.Printf("Received a message: %s\n", string(m.Data))

		json.Unmarshal(m.Data, &hooksmessage)

		fmt.Printf("Получено следующее сообщение: \nКому: %v\nСообщение: %v", hooksmessage.To, hooksmessage.Message)

		db, _ := leveldb.OpenFile("/DB", nil)
		jirakey, _ := db.Get([]byte(hooksmessage.To), nil)
		key, _ := strconv.Atoi(string(jirakey))
		db.Close()
		bot.Send(tele.ChatID(key), hooksmessage.Message)
		//fmt.Println(int64(jirakey1))
		//bot.Send(tele.ChatID(int64(jirakey1)), "Успешно.")

		//bot.Send(tele.ChatID(resultChan), "Успешно.")
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

func dbconn() *leveldb.DB {
	db, err := leveldb.OpenFile("/DB", nil)
	if err != nil {
		log.Println(err)
	}
	db.Close()
	return db
}

func OnActivation(c tele.Context, state fsm.Context) error {
	menu := &tele.ReplyMarkup{}
	menu.Reply(menu.Row(returnBtn))
	menu.ResizeKeyboard = true

	state.Set(SendedCode)
	return c.Send("Напиши мне свой логин от Jira.", menu)
}

func OnInputJiraLogin(c tele.Context, state fsm.Context) error {
	menu := &tele.ReplyMarkup{}
	menu.Reply(menu.Row(ActivationBtn))
	menu.ResizeKeyboard = true

	Login := c.Message().Text
	jirauser := jiraFuncs.UserReturner(Login)

	//defer state.Finish(true)
	result := "Произошла ошибка при регистрации. Пользователь не найден. Проверьте правильность введённого вами логина."

	if jirauser.Name != "" {

		code := randstr.String(6)
		Sender.SendMessage(code, jirauser.EmailAddress)

		go state.Update("login", Login)
		go state.Update("code", code)
		go state.Set(InputedCode)

		return c.Send(fmt.Sprintf("На вашу почту %v был направлен код, введите его.", jirauser.EmailAddress))
	}
	return c.Send(fmt.Sprintf(result), menu)
}

func OnInputCode(c tele.Context, state fsm.Context) error {
	menu := &tele.ReplyMarkup{}
	menu.Reply(menu.Row(ActivationBtn))
	menu.ResizeKeyboard = true

	defer state.Finish(true)
	var (
		sendedCode  string
		inputedCode string
		login       string
	)
	state.MustGet("code", &sendedCode)
	state.MustGet("login", &login)
	inputedCode = c.Message().Text

	fmt.Println(sendedCode, inputedCode)
	jirauser := jiraFuncs.UserReturner(login)
	db, _ := leveldb.OpenFile("/DB", nil)
	db.Put([]byte(jirauser.Key), []byte(fmt.Sprintf("%v", c.Chat().ID)), nil)
	db.Close()
	fmt.Println(sendedCode, inputedCode)

	result := "Произошла ошибка при регистрации. Проверьте правильность введённого вами кода."

	if sendedCode == inputedCode {
		result = "Вы были успешно зарегистрированы."
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
