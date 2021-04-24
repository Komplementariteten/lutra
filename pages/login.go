package pages

import (
	"github.com/Komplementariteten/lutra/util"
)

// LoginPageTemplate is the Name of the Template file
const LoginPageTemplate = "login"

// LoginPage are the Values for the Login Template
type LoginPage struct {
	Defaults
	LoginForm
	Headline    string
	Description string
}

type LoginForm struct {
	UsernameField string
	PasswordField string
	UsernameName  string
	PasswordName  string
	SubmitName    string
	Challenge     string
}

type PostRegisterPage struct {
	Defaults
	Message string
}

const LoginFailedPage = "/auth/login/failed"
const LoginOkPage = "/auth/login_ok"

// GetLoginPageValues provides Default Values for the Login Page
func GetLoginPageValues() *LoginPage {
	v := &LoginPage{}
	v.Defaults.LangShort = "de"
	v.Defaults.Title = "Login"
	v.TemplateName = "login"
	v.LoginForm.PasswordField = util.GetRandomString(12)
	v.LoginForm.UsernameField = util.GetRandomString(12)
	v.LoginForm.PasswordName = "Passwort"
	v.LoginForm.UsernameName = "Login"
	v.LoginForm.SubmitName = "Anmelden"
	return v
}

// GetPostRegisterValues provides default values for Register Page
func GetPostRegisterValues() *PostRegisterPage {
	v := &PostRegisterPage{}
	v.TemplateName = "registration_complete"
	v.LangShort = "de"
	v.Title = "Registration completed"
	v.Message = "Registration abgeschlossen, bitte melden Sie sich nun an um ad zu nutzen"
	return v
}
