package util_test

import (
	"os"
	"path/filepath"
	"testing"
	"text/template"

	"lutra/model"
	"lutra/pages"
	"lutra/util"
)

func Test_NewClient(t *testing.T) {
	_, err := util.NewMailClient("smtp.mailbox.org", 587, "myadventure@mailbox.org", "D-#7BE9_3QLyYyM98UEF")

	if err != nil {
		t.Fatal(err)
	}

	_, err = util.NewMailClient("smtp.mailbox.org", 587, "myadventure@mailbox.org", "D-#s")

	if err != nil {
		t.Fatalf("Client does error with invalid Auth Info")
	}

	_, err = util.NewMailClient("smtp.s2e3.org", 587, "myadventure@mailbox.org", "D-#s")

	if err != nil {
		t.Fatalf("Client does error with invalid Hostname")
	}

}

func TestStmpClient_SendText(t *testing.T) {
	tmp, err := template.New("mailtest").Parse(`From: {{.FromName}} <{{.From}}>
Subject:{{.Subject}}
To: {{.RecipientName}} <{{.Recipient}}>
	
Hello {{.RecipientName}}
{{.Text}}

{{.LinkText}}
{{.ClickLink}}`)

	if err != nil {
		t.Fatal(err)
	}

	d := &util.MailTemplateDefaults{
		Subject:       "Test E-Mail",
		RecipientName: "Oliver",
		Recipient:     "osiegemund@gmail.com",
		FromName:      "Other Name",
		From:          "info@myadventure.space",
		Text:          "E-Mail Content",
		LinkText:      "Click here to go to google",
		ClickLink:     "https://www.google.de",
	}
	c, err := util.NewMailClient("smtp.mailbox.org", 587, "info@myadventure.space", "D-#7BE9_3QLyYyM98UEF")
	if err != nil {
		t.Fatal(err)
	}
	err = c.SendText(d.From, []string{d.Recipient}, tmp, d)
	if err != nil {
		t.Fatal(err)
	}
}

func TestStmpClient_SendWithTextTemplate(t *testing.T) {
	u := &model.User{
		Email: "osiegemund@gmail.com",
		Name:  "Oliver",
	}

	d := pages.GetRegisterFormMailData(u)

	c, err := util.NewMailClient("smtp.mailbox.org", 587, "info@myadventure.space", "D-#7BE9_3QLyYyM98UEF")
	if err != nil {
		t.Fatal(err)
	}
	err = c.SendWithTextTemplate(d.MailTemplateDefaults.From, []string{d.MailTemplateDefaults.Recipient}, "register", d)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_ParseGlobFindsTemplates(t *testing.T) {
	t.Log(os.Getwd())
	pattern := filepath.Join("/Users/me/Documents/Workspace/AdventureService/src/lutra/tmpl", "*.tmpl")
	_, err := template.ParseGlob(pattern)
	if err != nil {
		t.Fatal(err)
	}
}
