package mailer

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"time"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type SendGridMailer struct {
	fromEmail string
	apiKey    string
	client    *sendgrid.Client
}

func NewSendgridMailer(apiKey, fromEmail string) *SendGridMailer {
	if apiKey == "" {
		log.Fatal("SENDGRID_API_KEY is not set")
	}
	return &SendGridMailer{
		fromEmail: fromEmail,
		apiKey:    apiKey,
		client:    sendgrid.NewSendClient(apiKey),
	}
}

func (m *SendGridMailer) Send(templateFile, username, email string, data any, isSandbox bool) error {

	log.Println("sandbox", isSandbox)
	from := mail.NewEmail(FromName, m.fromEmail)
	to := mail.NewEmail(username, email)

	// template parsing and building
	tmpl, err := template.ParseFS(FS, "templates/"+templateFile)
	if err != nil {
		return err
	}
	subject := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(subject, "subject", data)
	if err != nil {
		return err
	}

	body := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(body, "body", data)
	if err != nil {
		return err
	}
	message := mail.NewSingleEmail(from, subject.String(), to, "", body.String())
	message.SetMailSettings(&mail.MailSettings{
		SandboxMode: &mail.Setting{
			Enable: &isSandbox,
		},
	})
	for i := 0; i < maxRetries; i++ {

		response, err := m.client.Send(message)
		if err != nil {
			log.Printf("Failed to send emai l to %s,attempt %d of %d", email, i+1, maxRetries)
			log.Printf("Error: %v", err)

			//Exponential backoff
			time.Sleep(time.Duration(i*1) * time.Second)
			continue
		}
		log.Printf("Email sent to %s, status code: %d", email, response.StatusCode)
		return nil
	}

	return fmt.Errorf("Failed to send email to %s after %d attempts", email, maxRetries)
}
