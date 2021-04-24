package auth

import (
	"context"
	"net/http"
	"testing"
	"time"

	"lutra/db"
	"lutra/model"
	"lutra/pages"
	"github.com/dgrijalva/jwt-go"
)

func createExpiredCookieWithAuthorization() (*http.Cookie, error) {
	creds := &Credentials{
		Username: "test@test.org",
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(12)).Unix(),
			Issuer:    "test@test.org",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, creds)
	signedToken, err := token.SignedString(jwtKeyAsByte)
	if err != nil {
		return nil, err
	}
	c := &http.Cookie{
		Name:    tokenCookieName,
		Expires: time.Now(),
		Value:   signedToken,
		Domain:  "*",
	}

	return c, nil
}

func createExpiredCookieWithInvalidAuthorization() (*http.Cookie, error) {
	creds := &Credentials{
		Username: "test@test.org",
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(12)).Unix(),
			Issuer:    "test@test.org",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, creds)
	signedToken, err := token.SignedString([]byte("not valid"))
	if err != nil {
		return nil, err
	}
	c := &http.Cookie{
		Name:    tokenCookieName,
		Expires: time.Now(),
		Value:   signedToken,
		Domain:  "*",
	}

	return c, nil
}

func setupDatabase() (*db.Db, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	db, err := db.NewConnection(ctx, "mongodb+srv://dev:VmV5lBWnSPgog5xn@devcluster.pqajf.mongodb.net/defaultDatabase?retryWrites=true&w=majority", "TestDatabase")
	if err != nil {
		return nil, err
	}
	return db, nil
}

func TestAuth_FindAuthorization(t *testing.T) {
	// Test 1
	r, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	cookie, err := createExpiredCookieWithAuthorization()
	if err != nil {
		t.Fatal(err)
	}
	r.AddCookie(cookie)

	_, res := FindAuthorization(context.Background(), r)
	if res != AuthorizationFound {
		t.Fatalf("FindAuthorization failed check Authorization Cookie %v", res)
	}

	// Test 2
	r2, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	cookie.Value = "ignore"
	r2.AddCookie(cookie)
	_, res = FindAuthorization(context.Background(), r2)
	if res == AuthorizationFound {
		t.Fatalf("FindAuthorization failed check Authorization Cookie %v", res)
	}

	// Test 3
	r3, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	_, res = FindAuthorization(context.Background(), r3)
	if res != NotAuthorized {
		t.Fatal("FindAuthorization found something where there is no Cookie in the Request")
	}

	// Test 4
	r4, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	cookie2, err := createExpiredCookieWithInvalidAuthorization()
	if err != nil {
		t.Fatal(err)
	}
	r4.AddCookie(cookie2)
	_, res = FindAuthorization(context.Background(), r4)
	if res != AuthorizationFailed {
		t.Fatal("FindAuthorization does not Fail Authorization with invalid Token")
	}
}

func TestCredentials_CheckUser(t *testing.T) {

	db, err := setupDatabase()
	if err != nil {
		t.Fatal(err)
	}
	user := &model.User{
		Name:     "Test User",
		Login:    "test@test.org",
		Disabled: false,
		Email:    "test@test.org",
		Terms:    make(map[string]bool),
	}

	err = pages.AddNewUser(context.Background(), db, user)
	if err != nil {
		t.Fatal(err)
	}
	r, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	cookie, err := createExpiredCookieWithAuthorization()
	if err != nil {
		t.Fatal(err)
	}
	r.AddCookie(cookie)

	f, res := FindAuthorization(context.Background(), r)
	if res != AuthorizationFound {
		t.Fatalf("Failed to find Authorization Cookie with %v", res)
	}

	if err = f.CheckUser(context.Background(), db); err != nil {
		t.Fatalf("Could not find User in Db %v", err)
	}
	p, err := db.Collection(context.Background(), "auth")
	if err != nil {
		t.Fatal(err)
	}
	p.DeleteDatabase(context.Background())
}
