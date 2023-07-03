package main

import (
	"log"
	"net/http"
	"os"

	"github.com/chirag6398/todoApp/auth"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/jwtauth"
	"github.com/joho/godotenv"
)

var tokenAuth *jwtauth.JWTAuth

func main() {
	router := chi.NewRouter()
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	port := os.Getenv("PORT")
	secretKey := os.Getenv("SECRET_KEY")
	tokenAuth = jwtauth.New("HS256", []byte(secretKey), nil)

	t := &auth.AuthToken{
		TokenAuth: tokenAuth,
	}

	router.Group(func(r chi.Router) {
		r.Use(middleware.Logger)
		r.Post("/register", auth.Register)
		r.Post("/login", t.Login)

		r.Group(func(r chi.Router) {
			r.Use(jwtauth.Verifier(tokenAuth))
			r.Use(jwtauth.Authenticator)

			r.Get("/todo", getAllList)
			r.Post("/todo", addList)
			r.Put("/todo/{id}", updateList)
			r.Delete("/todo/{id}", deleteList)
		})
	})

	log.Println("Server started on port ", port)
	log.Fatal(http.ListenAndServe(port, router))
}
