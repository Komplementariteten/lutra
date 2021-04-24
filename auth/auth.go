package auth

import (
	"context"
	"crypto/sha512"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/Komplementariteten/lutra/db"
	"github.com/Komplementariteten/lutra/model"
	"github.com/Komplementariteten/lutra/pages"
	"github.com/Komplementariteten/lutra/util"

	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson"
)

const tokenCookieName = "token"
const refreshCookieName = "r"

const defaultTokenLiveTime = 15 * time.Minute
const defaultRefreshLiveTime = 2 * time.Hour
const defaultCookieLiveHours = 24 * 30

const (
	// NotAuthorized means, Authorization Cookie not Found
	NotAuthorized = iota
	// CookieError reports trouble with the golang Request
	CookieError
	// AuthorizationExpired say Authorization Cookie is Ok, but Expired
	AuthorizationExpired
	// AuthorizationFound says that the Cookie is presend and needs to be validated
	AuthorizationFound
	// Authorized says that the User is Authorized via JWT Cookie
	Authorized
	// SignatureInvalid Error
	SignatureInvalid
	// AuthorizationFailed is a general Error in Authorizing a Request
	AuthorizationFailed
)

const (
	authDbDocUser = iota
)

var jwtKeyAsByte = []byte("adventureservice has a 548$%_L@ Key")

// ErrRefresh indicates that the Request Token handeling failed
var ErrRefresh = errors.New("refresh of Authorization failed")

// ErrPassword indicates that the given Password does not match its master in the DB
var ErrPassword = errors.New("Password do not match")

var ErrNoAuth = errors.New("Not Authorized")

const jwtAudiance = "123-456-789-abc"

// Credentials are decoded from the Token Cookie
type Credentials struct {
	Username string `json:"username"`
	Password string `json:"-", bson:"-"`
	pwHash   string `json:"pawhash"`
	salt     string `json:"salt"`
	jwt.StandardClaims
}

// RefreshToken is the dao Object for Sessions
type RefreshToken struct {
	SessionId string    `json:"session"`
	Token     string    `json:"token"`
	User      string    `json:"username"`
	Blocked   bool      `json:"blocked"`
	Started   time.Time `json:"started"`
}

func createJwtID() string {
	return util.GetRandomString(32)
}

func getCredentialsFromRefresh(token string) (c *Credentials, err error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return nil, fmt.Errorf("Refresh token has the wrong format")
	}
	hexBytes, err := hex.DecodeString(parts[1])
	if err != nil {
		return nil, err
	}
	c = &Credentials{
		Username: string(hexBytes),
	}
	return
}

// ChecksRefresh updates auth and refresh token
func CheckRefresh(w http.ResponseWriter, r *http.Request) (*Credentials, error) {
	refreshToken, err := r.Cookie(refreshCookieName)
	if err == http.ErrNoCookie {
		// http.Redirect(w, r, pages.PortalUrls.Get("Login"), http.StatusFound)
		return nil, ErrNoAuth
	}
	if err != nil {
		return nil, err
	}
	cred, status := FindAuthorization(r.Context(), r)
	if status == NotAuthorized {
		claims, err := getCredentialsFromRefresh(refreshToken.Value)
		if err != nil {
			return nil, err
		}
		claims.SetJWTCookie(w, r)
		refreshToken.Expires = time.Now().Local().Add(defaultRefreshLiveTime)
		refreshToken.Path = "/"
		refreshToken.HttpOnly = true
		refreshToken.SameSite = http.SameSiteLaxMode
		http.SetCookie(w, refreshToken)
		Environment.User = claims.Username
		return claims, nil
	}
	switch status {
	case AuthorizationFound:
		Environment.User = cred.Username
		return cred, nil
	case SignatureInvalid:
	case AuthorizationFailed:
		http.Redirect(w, r, pages.PortalUrls.Get("403"), http.StatusForbidden)
		return nil, ErrRefresh
	case CookieError:
		http.Redirect(w, r, pages.PortalUrls.Get("400"), http.StatusBadRequest)
		return nil, ErrRefresh
	default:
		http.Redirect(w, r, pages.PortalUrls.Get("500"), http.StatusInternalServerError)
		return nil, ErrRefresh
	}
	return nil, ErrRefresh
}

// CreateRefreshToken creates a new Refresh Token on given claims
func (c *Credentials) CreateRefreshToken(ctx context.Context, w http.ResponseWriter, r *http.Request, hndl *db.Db) error {
	unameBytes := []byte(c.Username)
	hexStr := hex.EncodeToString(unameBytes)
	token := &RefreshToken{
		SessionId: c.StandardClaims.Id,
		Token:     fmt.Sprintf("%s.%s.%s", util.GetRandomString(6), hexStr, c.StandardClaims.Id),
		User:      c.Username,
		Blocked:   false,
		Started:   time.Now().Local(),
	}
	dbPool, err := hndl.Collection(ctx, db.AuthDbPool)
	if err != nil {
		return err
	}

	tokenDoc, err := util.ToBsonD(token)
	if err != nil {
		return err
	}
	dbPool.Add(ctx, tokenDoc)

	cookie := &http.Cookie{
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
		Name:     refreshCookieName,
		Domain:   r.URL.Host,
		Path:     "/",
		Expires:  time.Now().Local().Add(defaultRefreshLiveTime),
		Value:    token.Token,
	}
	http.SetCookie(w, cookie)
	return nil
}

// SetJWTCookie sets a JWT Cookie for valid trusted Credentials
func (c *Credentials) SetJWTCookie(w http.ResponseWriter, r *http.Request) {
	expires := time.Now().Local().Add(defaultTokenLiveTime)
	c.StandardClaims.ExpiresAt = expires.Unix()
	c.StandardClaims.Audience = jwtAudiance
	c.StandardClaims.Issuer = r.URL.Hostname()
	c.StandardClaims.Subject = "Access"
	c.StandardClaims.NotBefore = time.Now().Local().Unix()
	c.StandardClaims.IssuedAt = time.Now().Local().Unix()
	c.StandardClaims.Id = createJwtID()

	k := GetEcdsaKey(Environment.Config)
	token := jwt.NewWithClaims(jwt.SigningMethodES512, c)
	signedToken, err := token.SignedString(k)
	if err != nil {
		panic(err)
	}
	cookie := &http.Cookie{
		Name:     tokenCookieName,
		HttpOnly: false,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
		Domain:   r.URL.Host,
		Path:     "/",
		Expires:  expires,
		Value:    signedToken}

	http.SetCookie(w, cookie)
}

// FinishRegistration finishes the Registration and activates the User
func FinishRegistration(ctx context.Context, hndl *db.Db, token string) error {
	if !hndl.IsConnected {
		return fmt.Errorf(db.DatabaseNotConnected)
	}
	dbPool, err := hndl.Collection(ctx, db.AuthDbPool)
	if err != nil {
		return err
	}
	var item bson.M
	item, err = dbPool.Get(ctx, &bson.M{model.RegisterTokenField: token, model.RegisterTimeField: bson.M{"$lt": time.Now().Local().Add(48 * time.Hour)}, model.AuthTypesField: model.AuthTypeRegister})
	if err != nil {
		return err
	}
	if item == nil {
		return fmt.Errorf("not found")
	}

	dbUser, err := dbPool.Get(ctx, &bson.M{model.AuthTypesField: model.AuthTypeUser, model.UserLoginDbField: item[model.RegisterIdField], model.UserConfirmedDbField: false})
	if err != nil {
		return err
	}

	updatemap := make(map[string]interface{})

	updatemap[model.UserConfirmedDbField] = true
	updatemap[model.UserDisableDbField] = false
	updatemap[model.UserUpdatedDateDbField] = time.Now().Local()
	updatemap[model.UserEnabledDateDbField] = time.Now().Local()

	filter := bson.M{"_id": dbUser["_id"]}
	updated, err := dbPool.Update(ctx, &filter, updatemap)
	if err != nil {
		return err
	}
	if updated != 1 {
		panic(fmt.Sprintf("Update hit the wrong number of results: %d", updated))
	}
	dbPool.Delete(ctx, dbPool.GetId(item))
	return nil
}

// FindAuthorization tries to Authorize a http Request via Cookie
func FindAuthorization(ctx context.Context, r *http.Request) (*Credentials, int) {
	c, err := r.Cookie(tokenCookieName)
	if err != nil {
		if err == http.ErrNoCookie {
			return nil, NotAuthorized
		}
		return nil, CookieError
	}
	tokenStr := c.Value
	claims := &Credentials{}
	k := GetEcdsaKey(Environment.Config)

	tkn, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		m := token.Method.Alg()
		if m != jwt.SigningMethodES512.Name {
			return nil, fmt.Errorf("%s is not a valid algorithms", m)
		}
		return &k.PublicKey, nil
	})

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return nil, SignatureInvalid
		}
		return nil, AuthorizationFailed
	}

	if !tkn.Valid {
		return nil, AuthorizationFailed
	}

	return claims, AuthorizationFound
}

func CleanExpiredRegistrations(ctx context.Context, hndl *db.Db) error {
	if !hndl.IsConnected {
		return fmt.Errorf(db.DatabaseNotConnected)
	}
	dbPool, err := hndl.Collection(ctx, db.AuthDbPool)
	if err != nil {
		return err
	}

	// Clean old unused Registrations
	found_registration, err := dbPool.Find(ctx, &bson.M{model.AuthTypesField: model.AuthTypeRegister, model.RegisterTimeField: bson.M{"$lt": time.Now().Local().Add(-12 * time.Hour)}})
	if err != nil {
		return err
	}
	for _, item := range found_registration {
		userMap, err := dbPool.Get(ctx, &bson.M{model.UserLoginDbField: item[model.RegisterIdField], model.UserConfirmedDbField: false})
		if err == nil {
			dbPool.Delete(ctx, dbPool.GetId(userMap))
		}
		dbPool.Delete(ctx, dbPool.GetId(item))
	}

	return nil
}

func CleanExpiredRefresh(ctx context.Context, hndl *db.Db) error {
	if !hndl.IsConnected {
		return fmt.Errorf(db.DatabaseNotConnected)
	}
	dbPool, err := hndl.Collection(ctx, db.AuthDbPool)
	if err != nil {
		return err
	}

	// Clean old Refresh Tokens
	found, err := dbPool.Find(ctx, &bson.M{"blocked": false, "started": bson.M{"$lt": time.Now().Local().Add(-72 * time.Hour)}})
	if err != nil {
		return err
	}
	for _, item := range found {
		dbPool.Delete(ctx, dbPool.GetId(item))
	}

	return nil
}

// CheckPassword checks the user Password entered from Form
func (c *Credentials) CheckPassword(pwText string) error {
	saltedPw := fmt.Sprintf("%s%s", c.salt, pwText)
	h := sha512.New384()
	pwhash := h.Sum([]byte(saltedPw))
	hashText := hex.EncodeToString(pwhash)
	if hashText == c.pwHash {
		return nil
	}
	return ErrPassword
}

// CheckUser Looks up if User is in DB
func (c *Credentials) CheckUser(ctx context.Context, hndl *db.Db) error {
	if !hndl.IsConnected {
		return fmt.Errorf(db.DatabaseNotConnected)
	}
	dbPool, err := hndl.Collection(ctx, db.AuthDbPool)
	if err != nil {
		return err
	}
	dbUser, err := dbPool.Get(ctx, &bson.M{model.UserLoginDbField: c.Username, model.UserDisableDbField: false})
	if err != nil {
		return err
	}
	c.salt = dbUser[model.UserSaltDbField].(string)
	c.pwHash = dbUser[model.UserHashDbField].(string)
	return nil
}
