package services

import (
	"log"
	"strconv"

	"gopkg.in/gomail.v2"
)

type MailJob struct {
	To      string
	Subject string
	Body    string
}

type Mailer struct {
	dialer *gomail.Dialer
}

func NewMailer(host string, port string, username, password string) (*Mailer, error) {
	p ,err := strconv.Atoi(port)
	if err != nil {
		return nil, err
	}

	d := gomail.NewDialer(host, p, username, password)

	return &Mailer{
		dialer: d,
	}, nil
}

func (m *Mailer) SendMail(job *MailJob) error {
	msg := gomail.NewMessage()
	msg.SetHeader("From", m.dialer.Username)
	msg.SetHeader("To", job.To)
	msg.SetHeader("Subject", job.Subject)
	msg.SetBody("text/html", job.Body)

	if err := m.dialer.DialAndSend(msg); err != nil {
		log.Printf("Failed to send email to %s: %v", job.To, err)
		return err
	}
	log.Printf("Email sent to %s successfully", job.To)
	return nil
}