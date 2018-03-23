package web

import (
	"net/smtp"
	"strings"
	"fmt"
	"github.com/9chain/nbcsrvd01/config"
	"crypto/md5"
	"encoding/hex"
)

const (
	HOST        = "smtp.163.com"
	SERVER_ADDR = "smtp.163.com:25"
	USER        = "aquariusye@163.com" //发送邮件的邮箱
	PASSWORD    = "helloshiki"         //发送邮件邮箱的密码
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

	go func() {
		for _, v := range sendTo {
			str := strings.Replace("From: "+cfg.User+"~To: "+v+"~Subject: "+email.subject+"~~", "~", "\r\n", -1) + email.msg
			fmt.Println(str)
			err := smtp.SendMail(
				cfg.ServerAddr,
				auth,
				cfg.User,
				[]string{v},
				[]byte(str),
			)
			if err != nil {
				fmt.Printf("send mail fail %+v %+v %+v\n", auth, sendTo, err)
			}
		}
	}()

	return nil
}

func encryptWithSalt(plantext []byte,salt []byte) string {
	hash := md5.New()
	hash.Write(plantext)
	hash.Write(salt)
	sum := hash.Sum(nil)
	return hex.EncodeToString(sum)
}
