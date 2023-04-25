package emails

import (
	"bytes"
	"context"
	"crypto/tls"
	"html/template"
	"os"
	"strconv"

	gomail "gopkg.in/mail.v2"

	"github.com/sergio-abreu/note-taking-app-backend-golang/domain/notes"
)

func NewEMailer(host string, port int, username string, password string) EMailer {
	return EMailer{host: host, port: port, username: username, password: password}
}

func NewEMailerFromEnv() (EMailer, error) {
	host := os.Getenv("SMTP_HOST")
	portStr := os.Getenv("SMTP_PORT")
	user := os.Getenv("SMTP_EMAIL")
	password := os.Getenv("SMTP_PASSWORD")
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return EMailer{}, err
	}

	return NewEMailer(host, port, user, password), nil
}

type EMailer struct {
	host     string
	port     int
	username string
	password string
}

func (e EMailer) SendNoteReminder(ctx context.Context, to string, note notes.Note) error {
	m := gomail.NewMessage()
	m.SetHeader("From", e.username)
	m.SetHeader("To", to)
	m.SetHeader("Subject", "You've got a reminder")
	m.SetHeader("MIME-version", "1.0")
	m.SetHeader("Content-Type", "text/html")
	m.SetHeader("charset", "\"UTF-8\"")

	t, err := template.ParseFiles("./infrastructure/emails/templates/template.html")
	if err != nil {
		return err
	}

	var body bytes.Buffer
	err = t.Execute(&body, note)
	if err != nil {
		return err
	}

	m.SetBody("text/html", body.String())

	dialer := gomail.NewDialer(e.host, e.port, e.username, e.password)
	dialer.TLSConfig = &tls.Config{InsecureSkipVerify: false, ServerName: e.host}
	return dialer.DialAndSend(m)
}
