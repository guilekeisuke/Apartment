package mail

import (
	"api/config"
	"fmt"
	"os"

	_ "image/jpeg"
	"net/smtp"
	"time"

)

var confPath = "/go/api/config/mail.yml"

func SendMail(smtp string, , from string, to []string, msg string) {
	// config情報を取得
	confData, err := config.LoadConfigForYaml(confPath)
	if err != nil {
		fmt.Println(err)
	}

	// コンフィファイルから値を取得
	from := ConfData.Gmail.From
	to := ConfData.Gmail.To
	_smtp := ConfData.Gmail.Smtp
	pwd := ConfData.Gmail.Password
	auth := smtp.PlainAuth("", from, pwd, _smtp)
	fromName := ConfData.MovieUpcoming.FromName
	subject := ConfData.MovieUpcoming.Subject
	port := ConfData.Gmail.Port

	// 送信メール内容作成
	addr := smtp + ":" + port
	auth := smtp.PlainAuth("", from, pwd, smtp)
	msg_from := "From: <" + fromName + "> <" + from + ">\r\n"
	msg_to := "To: " + to + "\r\n" +
	msg_subject := "Subject: " + subject + "\r\n\r\n" +
	msgFull := []byte("" + 
		msg_from +
		msg_to +
		msg_subject +
		msg +
		"\r\n" +
		"")

	// メール送信
	sendMail := smtp.SendMail(addr, auth, from, []string{to}, msgFull)
	if sendMail != nil {
		fmt.Fprintf(os.Stderr, "エラー: %v\n", sendMail)
		return
	}
}


func (m mail) body() string {
    return "To: " + m.to + "\r\n" +
        "Subject: " + m.sub + "\r\n\r\n" +
        m.msg + "\r\n"
}