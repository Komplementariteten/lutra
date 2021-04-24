package pages

import "errors"

// PortalUrls Map to hold all Urls registered in the Portal
var PortalUrls Urls

var defaultUrl = "/auth/login"

type Urls struct {
	m map[string]string
}

// UrlNameAlreadyRegisteredError is the Error to mark that a URL was tried to be overwritten
var UrlNameAlreadyRegisteredError = errors.New("Url name already registered")

func init() {
	PortalUrls.m = make(map[string]string)
	PortalUrls.Add("LoginFailed", LoginFailedPage)
	PortalUrls.Add("LoginOk", LoginOkPage)
	PortalUrls.Add("RegisterOk", RegisterOkPage)
	PortalUrls.Add("Register", "/auth/register")
	PortalUrls.Add("RegisterConfirmPrefix", "http://localhost:61616/auth/register/confirm")
	PortalUrls.Add("Login", "/login")
	PortalUrls.Add("RegistrationComplete", "/welcome")
	PortalUrls.Add("GenericError", "/errors/generic")
	PortalUrls.Add("500", "/errors/500")
	PortalUrls.Add("400", "/errors/400")
	PortalUrls.Add("403", "/errors/403")
	PortalUrls.Add("404", "/errors/404")
	PortalUrls.Add("405", "/errors/405")
	PortalUrls.Add("410", "/errors/410")
}

// Add a Item to the Urls Map
func (u *Urls) Add(name string, url string) error {
	if _, ok := u.m[name]; ok {
		return UrlNameAlreadyRegisteredError
	}
	u.m[name] = url
	return nil
}

// Get returns the name of a given Url Name
func (u *Urls) Get(name string) string {
	if s, ok := u.m[name]; ok {
		return s
	}
	return defaultUrl
}
