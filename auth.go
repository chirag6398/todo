package main

import (
	"encoding/json"
	"net/http"
)

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// type AuthToken struct {
// 	AuthToken *jwtauth.JWTAuth
// }

// var tokenAuth AuthToken

var users []User

func register(w http.ResponseWriter, r *http.Request) {

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

func login(w http.ResponseWriter, r *http.Request) {
	var user User

	err := json.NewDecoder(r.Body).Decode(&user)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("invalid request please enter correct data"))
		return
	}

	for _, v := range users {
		if v.Password == user.Password && v.Username == user.Username {
			_, tokenString, _ := tokenAuth.Encode(map[string]interface{}{"username": user.Username})
			w.Write([]byte(tokenString))
			break
		} else {
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		}
	}

}
