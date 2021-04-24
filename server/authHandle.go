package server

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/Komplementariteten/lutra/auth"
	"github.com/Komplementariteten/lutra/db"
	"github.com/Komplementariteten/lutra/pages"
	"github.com/Komplementariteten/lutra/util"
)

// Testuser PW: abcABC123$

const sessionCookieName = "_SIC"

const ErrNotAllowedWhileLoggedIn = 460

type Session struct {
	Host          string
	Browser       string
	Created       time.Time
	Challenge     string
	LoginPage     *pages.LoginPage
	RegisterPager *pages.RegisterPage
}

type authHandle struct {
	registerPage        *pages.RegisterPage
	confirmRegisterPage *Page
	db                  *db.Db
	mail                *util.SmtpClient
	sessions            map[string]*Session
}

func (a *authHandle) ShouldBeExecuted(n int) bool {
	targetTime := 1
	return n%targetTime == 0
}
func (a *authHandle) Cleanup(ctx context.Context) {
	fmt.Println("authHandle Cleanup called")
	auth.CleanExpiredRefresh(ctx, a.db)
	auth.CleanExpiredRegistrations(ctx, a.db)
}

func (a *authHandle) registrationService(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodPost {
		if a.registerPage == nil {
			return
		}
		if auth.DoRegister(a.registerPage, a.mail, a.db, w, r) {
			return
		}
	}

	if a.registerPage == nil {
		a.registerPage = pages.GetRegisterPageValues()
	}
	var challenge string
	sessCookie, err := r.Cookie(sessionCookieName)
	if err == http.ErrNoCookie {
		sessid := util.GetRandomString(14)
		browser := r.Header.Get("user-agent")
		created := time.Now().Local()
		sess := &Session{
			Host:      r.RemoteAddr,
			Browser:   browser,
			Created:   created,
			Challenge: util.GetRandomString(36),
			LoginPage: pages.GetLoginPageValues(),
		}
		challenge = sess.Challenge
		a.sessions[sessid] = sess
		sessCookie := &http.Cookie{
			Name:     sessionCookieName,
			Path:     "/",
			SameSite: http.SameSiteLaxMode,
			HttpOnly: true,
			MaxAge:   43200,
			Expires:  time.Now().Local().Add(12 * time.Hour),
		}
		sessCookie.Value = sessid
		http.SetCookie(w, sessCookie)
	} else {
		browser := r.Header.Get("user-agent")
		sessid := sessCookie.Value
		sess := a.sessions[sessid]
		if sess == nil || browser != sess.Browser {
			sessCookie.Value = "delete"
			sessCookie.MaxAge = -1
			sessCookie.HttpOnly = true
			sessCookie.Path = "/"
			sessCookie.SameSite = http.SameSiteLaxMode
			sessCookie.Expires = time.Unix(0, 0)
			http.SetCookie(w, sessCookie)
			return
		}
		challenge = sess.Challenge
	}

	registerInfo := a.registerPage.RegisterForm
	registerInfo.Challenge = challenge
	err = AsJson(w, registerInfo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (a *authHandle) updateOrCreateSession(w http.ResponseWriter, r *http.Request) (*Session, error) {
	sessCookie, err := r.Cookie(sessionCookieName)
	if err == http.ErrNoCookie {
		sessid := util.GetRandomString(14)
		browser := r.Header.Get("user-agent")
		created := time.Now().Local()
		sess := &Session{
			Host:      r.RemoteAddr,
			Browser:   browser,
			Created:   created,
			Challenge: util.GetRandomString(36),
			LoginPage: pages.GetLoginPageValues(),
		}
		a.sessions[sessid] = sess
		sessCookie := &http.Cookie{
			Name:     sessionCookieName,
			Path:     "/",
			SameSite: http.SameSiteLaxMode,
			HttpOnly: true,
			MaxAge:   43200,
			Expires:  time.Now().Local().Add(12 * time.Hour),
		}
		sessCookie.Value = sessid
		http.SetCookie(w, sessCookie)
		return sess, nil
	} else {
		sessid := sessCookie.Value
		sess := a.sessions[sessid]
		if sess == nil {
			sessCookie.Value = "delete"
			sessCookie.MaxAge = -1
			sessCookie.HttpOnly = true
			sessCookie.Path = "/"
			sessCookie.SameSite = http.SameSiteLaxMode
			sessCookie.Expires = time.Unix(0, 0)
			http.SetCookie(w, sessCookie)
			return nil, fmt.Errorf("Session seems to be Hjacked")
		}
		sess.Browser = r.Header.Get("user-agent")
		return sess, nil
	}
	panic("Session was not handled")
}

func (a *authHandle) GetSession(w http.ResponseWriter, r *http.Request) (*Session, error) {
	sessCookie, err := r.Cookie(sessionCookieName)
	if err != nil {
		return nil, err
	}
	sessid := sessCookie.Value
	if sess, ok := a.sessions[sessid]; ok {
		return sess, nil
	}
	sessCookie.Value = "delete"
	sessCookie.MaxAge = -1
	sessCookie.HttpOnly = true
	sessCookie.Path = "/"
	sessCookie.SameSite = http.SameSiteLaxMode
	sessCookie.Expires = time.Unix(0, 0)
	http.SetCookie(w, sessCookie)
	return nil, fmt.Errorf("Session gone")
}

func (a *authHandle) sessionService(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		a.updateOrCreateSession(w, r)
		w.Write([]byte("OK"))
	}
}

func (a *authHandle) loginService(w http.ResponseWriter, r *http.Request, claims *auth.Credentials) {

	session, err := a.GetSession(w, r)
	if err == http.ErrNoCookie {
		http.Error(w, err.Error(), http.StatusOK)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if session == nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if session.LoginPage == nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	v := session.LoginPage.LoginForm

	if r.Method == http.MethodPost {
		if err := auth.DoLogin(&v, a.db, w, r); err != nil {
			http.Error(w, err.Error(), http.StatusForbidden)
			return
		}
		w.Write([]byte("OK"))
		return
	}
	v.Challenge = session.Challenge
	err = AsJson(w, v)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (a *authHandle) handlePostRegisterRequests(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		http.Redirect(w, r, "/auth/", http.StatusFound)
		return
	}
	v, ok := auth.DoRegisterDone(a.db, w, r)
	if !ok {
		return
	}
	p, err := OpenPage(v.TemplateName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = p.Render(w, v)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (a *authHandle) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	claims, err := auth.CheckRefresh(w, r)

	if strings.HasPrefix(r.URL.Path, "login") {
		a.loginService(w, r, claims)
		return
	}

	if r.URL.Path == "register" && err == auth.ErrNoAuth {
		a.registrationService(w, r)
		return
	}

	if r.URL.Path == "session" {
		a.sessionService(w, r)
		return
	}

	if strings.HasPrefix(r.URL.Path, "register/confirm") {
		a.handlePostRegisterRequests(w, r)
		return
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
	}

	if err == nil {
		http.Error(w, "Can't register while Logged in.", ErrNotAllowedWhileLoggedIn)
	}
}
