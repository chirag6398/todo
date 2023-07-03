package auth

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/jwtauth"
)

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type AuthToken struct {
	TokenAuth *jwtauth.JWTAuth
}

var users []User

func Register(w http.ResponseWriter, r *http.Request) {

	var user User

	err := json.NewDecoder(r.Body).Decode(&user)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("invalid request please enter correct data"))
		return
	}
	for _, v := range users {
		if user.Username == v.Username {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("username already exists"))
			return
		}
	}
	users = append(users, user)
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("registration successful"))

}

func (t *AuthToken) Login(w http.ResponseWriter, r *http.Request) {
	var user User

	err := json.NewDecoder(r.Body).Decode(&user)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("invalid request please enter correct data"))
		return
	}

	for _, v := range users {
		if v.Password == user.Password && v.Username == user.Username {
			_, tokenString, _ := t.TokenAuth.Encode(map[string]interface{}{"username": user.Username})
			w.Write([]byte(tokenString))
			break
		} else {
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		}
	}

}
