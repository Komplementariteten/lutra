package util

import (
	"bytes"
	"fmt"
	"log"
	"net/smtp"
	"path/filepath"
	"text/template"
)

var templateDir = "/Users/me/Documents/Workspace/AdventureService/src/lutra/tmpl"

type SmtpClient struct {
	cram   smtp.Auth
	plain  smtp.Auth
	Server string
}

type MailTemplateDefaults struct {
	Subject       string
	Recipient     string
	From          string
	RecipientName string
	FromName      string
	Text          string
	LinkText      string
	ClickLink     string
}

func NewMailClient(smtpServer string, port int, user string, password string) (*SmtpClient, error) {
	a := smtp.CRAMMD5Auth(user, password)
	a_plain := smtp.PlainAuth("", user, password, smtpServer)
	c := &SmtpClient{
		cram:   a,
		plain:  a_plain,
		Server: fmt.Sprintf("%s:%d", smtpServer, port),
	}
	return c, nil
}

func (m *SmtpClient) send(from string, to []string, msg []byte) error {
	err := smtp.SendMail(m.Server, m.plain, from, to, msg)
	if err != nil {
		err = nil
		err = smtp.SendMail(m.Server, m.cram, from, to, msg)
		return err
	}
	return nil
}

func (m *SmtpClient) SendText(from string, to []string, t *template.Template, template_data interface{}) error {
	buffer := bytes.NewBufferString("")
	err := t.Execute(buffer, template_data)
	if err != nil {
		return err
	}
	log.Println(buffer)
	err = m.send(from, to, buffer.Bytes())
	return err
}

func (m *SmtpClient) SendWithTextTemplate(from string, to []string, name string, template_data interface{}) error {
	templates := template.Must(template.ParseGlob(filepath.Join(templateDir, "*.mail")))
	template := templates.Lookup(name)
	if template != nil {
		err := m.SendText(from, to, template, template_data)
		if err != nil {
			log.Print(err.Error())
			return err
		}
		return nil
	}
	return fmt.Errorf("Template %s not found", name)
}
