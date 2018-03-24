package web

import (
	"errors"
	"github.com/9chain/nbcsrvd01/config"
	"net/smtp"
	"strings"
)

type Email struct {
	to      string "to"
	subject string "subject"
	msg     string "msg"
}

func NewEmail(to, subject, msg string) *Email {
	return &Email{to: to, subject: subject, msg: msg}
}

func SendEmail(email *Email) error {
	cfg := &config.Cfg.SMTP
	auth := smtp.PlainAuth("", cfg.User, cfg.Password, cfg.Host)
	sendTo := strings.Split(email.to, ";")

	if len(sendTo) == 0 {
		return errors.New("invalid email")
	}

	v := sendTo[0]

	str := strings.Replace("From: "+cfg.User+"~To: "+v+"~Subject: "+email.subject+"~~", "~", "\r\n", -1) + email.msg
	err := smtp.SendMail(
		cfg.ServerAddr,
		auth,
		cfg.User,
		[]string{v},
		[]byte(str),
	)
	if err != nil {
		return err
	}
	return nil
}
