package pages

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/Komplementariteten/lutra/model"

	"github.com/Komplementariteten/lutra/util"

	"github.com/Komplementariteten/lutra/db"
)

// RegisterPageTemplate is the Name of the Template file
const RegisterPageTemplate = "register"

// RegisterPage are the Values for the register Template
type RegisterPage struct {
	Defaults
	RegisterForm
	Headline    string
	Description string
}

type RegisterForm struct {
	NameField               string
	NameName                string
	PasswordFieldValidation string
	PasswordValidationName  string
	UsernameField           string
	PasswordField           string
	UsernameName            string
	PasswordName            string
	SubmitName              string
	PasswordNote            string
	PasswordValidationInfo  string
	RegisterMailSender      string
	PhoneField              string
	RegisterService         string
	Challenge               string
}

type RegisterMail struct {
	util.MailTemplateDefaults
	RegistrationID string
}

const RegisterOkPage = "/auth/user_registered"

func (rf *RegisterForm) GetUserFromJson(m map[string]interface{}) (*model.User, error) {

	if phone, ok := m[rf.PhoneField].(string); ok {
		if len(phone) > 0 {
			return nil, fmt.Errorf("Invalid registration information")
		}
	}

	u := &model.User{}
	allSet := 0
	if name, ok := m[rf.NameField]; ok {
		u.Name = name.(string)
		allSet++
	}

	if login, ok := m[rf.UsernameField]; ok {
		u.Login = login.(string)
		u.Email = login.(string)
		allSet++
	}

	if password, ok := m[rf.PasswordField]; ok {
		u.SetPassword(password.(string))
		allSet++
	}

	if allSet < 3 {
		return nil, fmt.Errorf("Not all need Userdata provided")
	}

	return u, nil
}

func (rf *RegisterForm) GetUserFromForm(v url.Values) *model.User {
	n := v.Get(rf.NameField)
	l := v.Get(rf.UsernameField)
	p := v.Get(rf.PasswordField)
	u := &model.User{
		Name:           n,
		Login:          l,
		Email:          l,
		Disabled:       true,
		EmailConfirmed: false,
	}
	u.SetPassword(p)
	return u
}

func GetRegisterPageValues() *RegisterPage {
	v := &RegisterPage{}
	v.Defaults.LangShort = "de"
	v.Defaults.Title = "Register a new User Account"
	v.TemplateName = RegisterPageTemplate
	v.RegisterForm.PasswordNote = "Wählen Sie ein Passwort das mindestens 10 Zeichen unterschiedlicher Art wie Groß- und Kleinbuchstaben sowie Zahlen enthält."
	v.RegisterForm.PasswordField = util.GetRandomString(12)
	v.RegisterForm.UsernameField = util.GetRandomString(12)
	v.RegisterForm.PasswordFieldValidation = util.GetRandomString(12)
	v.RegisterForm.NameField = util.GetRandomString(12)
	v.RegisterForm.PhoneField = util.GetRandomString(12)
	v.RegisterForm.NameName = "Anzeigename"
	v.RegisterForm.PasswordName = "Passwort"
	v.RegisterForm.PasswordValidationName = "Bestätigung"
	v.RegisterForm.PasswordValidationInfo = "Bestätigen sie ihr Passwort durch erneute Eingabe"
	v.RegisterForm.UsernameName = "Login/E-Mail"
	v.RegisterForm.SubmitName = "Anmelden"
	v.RegisterService = PortalUrls.Get("Register")
	return v
}

func GetRegisterFormMailData(u *model.User) *RegisterMail {
	rm := &RegisterMail{
		MailTemplateDefaults: util.MailTemplateDefaults{
			Subject:       "Account registration at my-adventurespace confirmation",
			Recipient:     u.Email,
			RecipientName: u.Name,
			From:          "no-reply@myadventure.space",
			FromName:      "no-reply@myadventure.space",
			ClickLink:     PortalUrls.Get("RegisterConfirmPrefix"),
		},
		RegistrationID: util.GetRandomString(48),
	}
	return rm
}

// AddNewUser adds a User in a given Format to the Database
func AddNewUser(ctx context.Context, hndl *db.Db, u *model.User) (string, error) {
	if !hndl.IsConnected {
		return "", fmt.Errorf(db.DatabaseNotConnected)
	}
	dbPool, err := hndl.Collection(ctx, db.AuthDbPool)
	if err != nil {
		return "", err
	}

	u.Created = time.Now().Local()
	u.Updated = time.Now().Local()
	u.Disabled = true
	u.EmailConfirmed = false
	trackedUser, err := dbPool.Add(ctx, u.ToBson())
	if err != nil {
		return "", err
	}

	registerToken := util.GetRandomString(32)
	registerBson := u.ToRegisterBson(registerToken)
	_, err = dbPool.AddM(ctx, registerBson)
	if err != nil {
		deleteErr := dbPool.Delete(ctx, trackedUser.InsertedID)
		if deleteErr != nil {
			return "", deleteErr
		}
		return "", err
	}

	return registerToken, nil
}
