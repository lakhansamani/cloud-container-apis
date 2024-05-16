package utils

import (
	"bytes"
	"html/template"

	"github.com/lakhansamani/cloud-container/internal/global"
	"github.com/rs/zerolog/log"
	gomail "gopkg.in/mail.v2"
)

func getParsedTemplate(emailSubject, emailTemplate string, data map[string]interface{}) (string, string, error) {
	templ, err := template.New("template.tmpl").Parse(emailTemplate)
	if err != nil {
		return "", "", err
	}
	buf := &bytes.Buffer{}
	err = templ.Execute(buf, data)
	if err != nil {
		return "", "", err
	}
	templateString := buf.String()
	subject, err := template.New("subject.tmpl").Parse(emailSubject)
	if err != nil {
		return "", "", err
	}
	buf = &bytes.Buffer{}
	err = subject.Execute(buf, data)
	if err != nil {
		return "", "", err
	}
	subjectString := buf.String()
	return subjectString, templateString, nil
}

// SendEmail function to send mail
func SendMail(mailer *gomail.Dialer, emailSubject, emailTemplate string, to []string, data map[string]interface{}) error {
	subj, tmpl, err := getParsedTemplate(emailSubject, emailTemplate, data)
	if err != nil {
		log.Debug().Err(err).Msg("Failed to parse template")
		return err
	}
	m := gomail.NewMessage()
	m.SetAddressHeader("From", global.SMTPSenderEmail, global.SMTPSenderName)
	m.SetHeader("To", to...)
	m.SetHeader("Subject", subj)
	m.SetBody("text/html", tmpl)
	if err := mailer.DialAndSend(m); err != nil {
		return err
	}
	return nil
}
