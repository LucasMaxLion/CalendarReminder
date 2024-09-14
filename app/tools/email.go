package tools

import (
	"crypto/tls"
	"fmt"
	"github.com/spf13/viper"
	gomail "gopkg.in/mail.v2"
)

func SendEmail(to, subject, content string) error {

	m := gomail.NewMessage()
	m.SetHeader("From", m.FormatAddress(viper.GetString("email.address"), "萤火虫科技"))
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", fmt.Sprintf("<html><body>%s<br></body></html>", content))

	d := gomail.NewDialer("smtp.163.com", 25, viper.GetString("email.address"), viper.GetString("email.password"))
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	if err := d.DialAndSend(m); err != nil {
		return err
	}
	return nil
}
