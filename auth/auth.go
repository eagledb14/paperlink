package auth

import (
	db "github.com/eagledb14/paperlink/db"
	"crypto/rand"
	"encoding/base64"
)

type Auth struct {
	db *db.DbWrapper
	Cookies map[string]string
}

func NewAuth() *Auth {
	db, err := db.Open("./auth.db?_busy_timeout=10000")
	if err != nil {
		panic("Unable to create autht able")
	}

	createUserTable(db)

	newAuth := Auth {
		db: db,
		Cookies: make(map[string]string),
	}


	return &newAuth
}

func (a *Auth) GenerateCookie() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	cookie := base64.RawStdEncoding.EncodeToString(bytes)

	return cookie, nil
}
