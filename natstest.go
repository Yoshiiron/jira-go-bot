package main

import (
	"fmt"
	"github.com/nats-io/nats.go"
	jiraFuncs "jira-go-bot/JiraFuncs"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	nc, err := nats.Connect("nats://192.168.1.121:4222")
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()
	nc.Subscribe("foo", func(m *nats.Msg) {
		fmt.Printf("Received a message: %s\n", string(m.Data))
		fmt.Println(jiraFuncs.UserReturner(jiraFuncs.JiraClient(), fmt.Sprintf("%v", string(m.Data))).Key)
	})

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	<-ch
}
