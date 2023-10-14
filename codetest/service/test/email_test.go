package test

import (
	"crypto/tls"
	"github.com/jordan-wright/email"
	"net/smtp"
	"testing"
)

func TestSendEmial(t *testing.T) {
	e := email.NewEmail()
	e.From = "Jordan Wright <2577595258@qq.com>"
	e.To = []string{"2577595258@qq.com"}
	e.Subject = "验证码发送测试"
	e.HTML = []byte("你的验证码是：<b>123456</b>")
	err := e.SendWithTLS("smtp.qq.com:587",
		smtp.PlainAuth("", "2577595258@qq.com", "eyikrfleeqyudiag", "smtp.qq.com"),
		&tls.Config{InsecureSkipVerify: true, ServerName: "smtp.qq.com"})
	if err != nil {
		t.Fatal(err)
	}
}
