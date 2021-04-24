package auth

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/Komplementariteten/lutra/db"
	"github.com/Komplementariteten/lutra/model"
	"github.com/Komplementariteten/lutra/pages"
	"github.com/Komplementariteten/lutra/util"
)

// DoRegister handles the Register Requests
func DoRegister(rf *pages.RegisterPage, mail *util.SmtpClient, db *db.Db, w http.ResponseWriter, r *http.Request) bool {
	time.Sleep(500 * time.Millisecond)
	if err := r.ParseForm(); err != nil {
		return false
	}
	var u *model.User
	contentType := r.Header.Get("Content-Type")
	if contentType == "application/json" {
		j, err := util.GetJsonFromReader(r.Body)
		if err != nil {
			return false
		}
		u, err = rf.RegisterForm.GetUserFromJson(j)
		if err != nil {
			return false
		}
	} else {
		u = rf.RegisterForm.GetUserFromForm(r.PostForm)
	}
	if len(u.Login) == 0 {
		return false
	}
	token, err := pages.AddNewUser(r.Context(), db, u)
	if err == nil {
		time.Sleep(500 * time.Millisecond)
		data := pages.GetRegisterFormMailData(u)
		data.RegistrationID = token
		err = mail.SendWithTextTemplate(rf.RegisterForm.RegisterMailSender, []string{u.Email}, "register", data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return false
		}

		http.Redirect(w, r, pages.PortalUrls.Get("RegisterOk"), http.StatusFound)
		return true
	}
	// AddUser(r.Context(), db, )
	panic(err)
}

func DoLogin(lf *pages.LoginForm, db *db.Db, w http.ResponseWriter, r *http.Request) error {

	contentType := r.Header.Get("Content-Type")
	var m map[string]interface{}
	if contentType == "application/json" {
		j, err := util.GetJsonFromReader(r.Body)
		if err != nil {
			return err
		}
		m = j
	} else {
		return fmt.Errorf("Only json content is supported")
	}

	u := &Credentials{}
	if name, ok := m[lf.UsernameField]; ok {
		u.Username = name.(string)
	}
	err := u.CheckUser(r.Context(), db)
	if err != nil {
		time.Sleep(1 * time.Second)
		http.Error(w, "Forbidden", http.StatusForbidden)
		return fmt.Errorf("User not found")
	}
	pw := ""
	if password, ok := m[lf.PasswordField]; ok {
		pw = password.(string)
	}
	err = u.CheckPassword(pw)
	if err != nil {
		time.Sleep(1 * time.Second)
		http.Error(w, "Forbidden", http.StatusForbidden)
		return fmt.Errorf("Password do not match")
	}
	u.SetJWTCookie(w, r)
	u.CreateRefreshToken(r.Context(), w, r, db)
	return nil
}

// DoRegisterDone activates the User after registration
func DoRegisterDone(db *db.Db, w http.ResponseWriter, r *http.Request) (*pages.PostRegisterPage, bool) {
	time.Sleep(1 * time.Second)

	urlParts := strings.Split(r.URL.Path, "/")
	if urlParts[1] != "confirm" {
		return nil, false
	}

	err := FinishRegistration(r.Context(), db, urlParts[2])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	v := pages.GetPostRegisterValues()

	return v, true
}
