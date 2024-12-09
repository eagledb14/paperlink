package auth

import (
	"crypto/rand"
	"crypto/sha1"
	"encoding/base64"

	"golang.org/x/crypto/argon2"

	db "github.com/eagledb14/paperlink/db"
)

type User struct {
	Username string
	PassHash string
	Salt string
	Admin bool
}

func (a *Auth) NewUser(username string, password string, admin bool) (User, error) {
	user, err := a.GetUser(username)
	if user.Username != "" {
		return user, nil
	}

	hasher := sha1.New()
	salt, err := generateSalt()
	if err != nil {
		return User{}, err
	}
	
	hasher.Write([]byte(username))
	passHash := createHash(password, salt)

	newUser := User {
		Username: username,
		PassHash: passHash,
		Salt: salt,
		Admin: admin,
	}

	a.InsertUser(newUser)

	return newUser, nil
}

func (a *Auth) ValidateUser(username string, password string) (bool, error) {
	savedUser, err := a.GetUser(username)
	if err != nil {
		return false, err
	}

	validatePassHash := createHash(password, savedUser.Salt)

	if validatePassHash == savedUser.PassHash {
		return true, nil
	}

	return false, nil
}


func generateSalt() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	return base64.RawStdEncoding.EncodeToString(bytes), nil
}

func createHash(input string, salt string) string {
	time := uint32(10)
	memory := uint32(64 * 1024)
	threads := uint8(4)
	keyLength := uint32(32)

	hash := argon2.IDKey(
		[]byte(input), 
		[]byte(salt),
		time,
		memory,
		threads,
		keyLength,
		)
	return base64.RawStdEncoding.EncodeToString(hash)
}

func createUserTable(db *db.DbWrapper) error {
	return db.Exec(`CREATE TABLE IF NOT EXISTS users(
username TEXT PRIMARY KEY,
passHash TEXT,
salt TEXT,
admin BOOL
)`)
}

func (a *Auth) InsertUser(user User) error {
	return a.db.Exec(`INSERT INTO users(
username,
passHash,
salt,
admin
) VALUES (?, ?, ?, ?)`, user.Username, user.PassHash, user.Salt, user.Admin)
}

func (a *Auth) GetUser(username string) (User, error) {
	hasher := sha1.New()
	hasher.Write([]byte(username))
	userHash := base64.RawStdEncoding.EncodeToString(hasher.Sum(nil))

	row := a.db.QueryRow(`SELECT username, passHash, salt, admin FROM users WHERE username = ?`, userHash)
	newUser := User{}
	if err := row.Scan(&newUser.Username, &newUser.PassHash, &newUser.Salt, &newUser.Admin); err != nil {
		return newUser, err
	}

	return newUser, nil
}
