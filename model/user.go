package model

import (
	"fmt"
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

var store *sessions.CookieStore

func InitUserSessions(key []byte) {
	store = sessions.NewCookieStore(key)
	store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   0,
		HttpOnly: true,
	}
}

type User struct {
	Id           uint64 `db:"id"`
	Email        string `db:"email"`
	Username     string `db:"username"`
	PasswordHash string `db:"password_hash"`
	IsAdmin      bool   `db:"is_admin"`
}

func IsUserLoggedIn(r *http.Request) (User, bool, error) {
	session, err := store.Get(r, "user_session")
	if err != nil {
		return User{}, false, err
	}

	userId, ok := session.Values["user_id"]
	if !ok {
		return User{}, false, nil
	}

	user, err := GetUserById(userId.(uint64))
	if err != nil {
		return User{}, false, err
	}

	return user, true, nil
}

func Login(r *http.Request, w http.ResponseWriter, userId uint64) error {
	session, err := store.Get(r, "user_session")
	if err != nil {
		return err
	}

	session.Values["user_id"] = userId
	if err := session.Save(r, w); err != nil {
		return err
	}

	return nil
}

func Logout(r *http.Request, w http.ResponseWriter) error {
	session, err := store.Get(r, "user_session")
	if err != nil {
		return err
	}

	session.Options.MaxAge = -1
	if err := session.Save(r, w); err != nil {
		return err
	}

	return nil
}

func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))
	return err == nil
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func GetUserByEmail(email string) (User, error) {
	user := User{}

	query := `SELECT * FROM site_user WHERE email = $1`

	if err := db.Get(&user, query, email); err != nil {
		return User{}, err
	}

	return user, nil
}

func GetUserById(id uint64) (User, error) {
	user := User{}

	query := `SELECT * FROM site_user WHERE id = $1`

	if err := db.Get(&user, query, id); err != nil {
		return User{}, err
	}

	return user, nil
}

func CheckUserEmailExists(email string) (bool, error) {
	query := `SELECT count(1) > 0 FROM site_user WHERE email = $1`
	exists := false
	if err := db.Get(&exists, query, email); err != nil {
		return false, err
	}
	fmt.Println(exists)

	return exists, nil
}

func RegisterUser(email, username, password string) (uint64, error) {
	passwordHash, err := hashPassword(password)
	if err != nil {
		return 0, err
	}

	query := `
	INSERT INTO site_user (email, username, password_hash)
		VALUES ($1, $2, $3) RETURNING id;`

	var userId uint64
	err = db.Get(&userId, query, email, username, passwordHash)
	if err != nil {
		pgErr, ok := err.(*pq.Error)
		if ok {
			if pgErr.Code == "23505" {
				return 0, fmt.Errorf("Email already exists")
			}
		}

		return 0, err
	}

	return userId, nil
}
