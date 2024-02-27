package Sender

import (
	"log"
	"net/smtp"
)

func SendMessage(text string) {

	smtpHost := "smtp.yandex.ru"
	smtpPort := "587"
	smtpUsername := "yoshiiron@yandex.ru"
	smtpPassword := "udzkrlwvxgeyfptz"

	from := "yoshiiron@yandex.ru"

	to := []string{"denis.astafev@devim.team"}
	message := []byte("Subject: Тестовое сообщение\r\n" +
		"From: " + from + "\r\n" +
		"\r\n" +
		"Ваш код для активации Jira-Devim-Bot: " + text + "\r\n")

	auth := smtp.PlainAuth("", smtpUsername, smtpPassword, smtpHost)

	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, message)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Сообщение успешно отправлено")

}
