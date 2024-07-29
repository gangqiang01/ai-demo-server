package v1

import (
	"github.com/edgehook/ithings/common/dbm/model"
	"github.com/edgehook/ithings/common/global"
	"gopkg.in/gomail.v2"
	"k8s.io/klog/v2"
)

type Email struct {
	// From    string `json:"from"`
	To      string `json:"to"`
	Subject string `json:"subject"`
	Content string `json:"content"`
}

// Smtp 邮件服务

type Smtp struct {
	Address  string `form:"address" json:"address"`
	Port     *int   `form:"port" json:"port"`
	Username string `form:"username" json:"username"`
	Password string `form:"password" json:"password"`
}

// emailutil := v1.Email{
// 	To:      "gangqiang.sun@advantech.com.cn",
// 	Subject: "test",
// 	Content: "test  我是一个小小鸟",
// }

// err := emailutil.SendEmail()
// if err != nil {
// 	klog.Errorf("Error: %s", err.Error())
// }

func (email *Email) SendEmail() error {
	m := gomail.NewMessage()
	emailConfig, err := model.GetIthingsConfigByType(global.EmailType)
	if err != nil {
		return err
	}
	m.SetHeader("From", emailConfig.Username)
	m.SetHeader("To", email.To)
	m.SetHeader("Subject", email.Subject)
	m.SetBody("text/html", email.Content)
	gomail.SetCharset("UTF-8")

	d := gomail.NewDialer(emailConfig.Address, emailConfig.Port, emailConfig.Username, emailConfig.Password)
	// Send the email to Bob, Cora and Dan.
	if err := d.DialAndSend(m); err != nil {
		return err
	}
	return nil
}

func SendAlarmEmail(to string, subject string, content string) bool {
	html := "<div style=\"margin-top:30px;\">\n" +
		"            <i>Dear friend,</i>\n" +
		"            <p style=\"margin-top:30px; font-size:20px\">Welcome to ithings.</p>\n" +
		"            <p style=\"margin-top:20px;\">If it is not your request,please ignore this email.</p>\n" +
		"            <p style=\"margin-top:10px; margin-bottom:10px;\"> Details as follows:</p>\n" +
		"            <p style=\"margin-top:10px; margin-bottom:10px;\">" + content + "</p>\n" +
		// "            <p style=\"margin-top:10px; margin-bottom:50px;\">" + result + "</p>\n" +
		"            <p style=\"margin-top:30px; margin-bottom:50px;\">This is an automatically generated email,please do not reply.</p >\n" +
		"            <i>Best regards,</i><br>\n" +
		"            <i>ithings</i>\n" +
		"        </div>"
	emailutil := Email{
		To:      to,
		Subject: subject,
		Content: html,
	}

	err := emailutil.SendEmail()
	if err != nil {
		klog.Errorf("Error: %s", err.Error())
		return false
	}
	return true
}
