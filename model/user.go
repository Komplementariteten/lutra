package model

import (
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/Komplementariteten/lutra/util"

	"go.mongodb.org/mongo-driver/bson"
)

const (
	UserLoginDbField        = "login"
	UserHashDbField         = "sec"
	UserDisableDbField      = "disabled"
	UserConfirmedDbField    = "confirmed"
	UserSaltDbField         = "salt"
	UserCreationDateDbField = "created"
	UserEnabledDateDbField  = "enabled"
	UserUpdatedDateDbField  = "updated"
	UserTermsDbField        = "terms"
	UserNameDbField         = "name"
	UserEmailDbField        = "email"
	RegisterIdField         = "id"
	RegisterTimeField       = "created"
	RegisterTokenField      = "token"
	AuthTypesField          = "type"
	AuthTypeUser            = "user"
	AuthTypeRegister        = "register"
)

// User represents a User DAO
type User struct {
	Name           string
	Login          string
	Email          string
	pwhash         string
	salt           string
	Disabled       bool
	EmailConfirmed bool
	Terms          map[string]bool
	Created        time.Time
	Enabled        time.Time
	Updated        time.Time
}

func (u *User) ToBson() *bson.D {
	return &bson.D{
		{AuthTypesField, AuthTypeUser},
		{UserNameDbField, u.Name},
		{UserLoginDbField, u.Login},
		{UserEmailDbField, u.Email},
		{UserHashDbField, u.pwhash},
		{UserDisableDbField, u.Disabled},
		{UserConfirmedDbField, u.EmailConfirmed},
		{UserSaltDbField, u.salt},
		{UserTermsDbField, u.Terms},
		{UserCreationDateDbField, u.Created},
		{UserEnabledDateDbField, u.Enabled},
		{UserUpdatedDateDbField, u.Updated},
	}
}

func (u *User) ToRegisterBson(token string) *bson.M {
	return &bson.M{
		AuthTypesField:     AuthTypeRegister,
		RegisterIdField:    u.Login,
		RegisterTimeField:  u.Created,
		RegisterTokenField: token,
	}
}

func (u *User) SetPassword(plainWord string) {
	if len(u.salt) == 0 {
		u.salt = util.GetRandomString(12)
	}
	saltedPw := fmt.Sprintf("%s%s", u.salt, plainWord)
	h := sha512.New384()
	pwhash := h.Sum([]byte(saltedPw))
	u.pwhash = hex.EncodeToString(pwhash)
}
