package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	mongodb "github.com/chirag6398/todoApp/database"
	handler "github.com/chirag6398/todoApp/handlers"
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
	var Database = os.Getenv("DATABASE")
	var UserCollection = os.Getenv("COLLECTION_USER")
	var ListCollection = os.Getenv("COLLECTION_LIST")

	tokenAuth = jwtauth.New("HS256", []byte(secretKey), nil)
	mongouri := os.Getenv("MONGO_URI")

	mongoClient, err := mongodb.ConnectMongoDb(mongouri)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	rs := &handler.Resource{
		TokenAuth:      tokenAuth,
		Client:         mongoClient,
		Database:       Database,
		UserCollection: UserCollection,
		ListCollection: ListCollection,
	}

	router.Group(func(r chi.Router) {
		r.Use(middleware.Logger)
		r.Post("/register", rs.Register)
		r.Post("/login", rs.Login)

		r.Group(func(r chi.Router) {
			r.Use(jwtauth.Verifier(tokenAuth))
			r.Use(jwtauth.Authenticator)

			r.Get("/todo", rs.GetAllList)
			r.Post("/todo", rs.AddList)
			r.Put("/todo/{id}", rs.UpdateList)
			r.Delete("/todo/{id}", rs.DeleteList)
		})
	})

	log.Println("Server started on port ", port)
	log.Fatal(http.ListenAndServe(port, router))
}
