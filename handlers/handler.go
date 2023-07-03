package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	mongodb "github.com/chirag6398/todoApp/database"
	"github.com/go-chi/chi"
	"github.com/go-chi/jwtauth"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Resource struct {
	TokenAuth      *jwtauth.JWTAuth
	Client         *mongodb.MongoClient
	Database       string
	UserCollection string
	ListCollection string
}
type List struct {
	Title       string `json:"title"`
	IsCompleted bool   `json:"isCompleted"`
	IsDeleted   bool   `json:"isDeleted"`
	UserName    string `json:"username"`
}

var lists map[string][]List = map[string][]List{}
var users []User

func (rs *Resource) Register(w http.ResponseWriter, r *http.Request) {
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("invalid request please enter correct data"))
		return
	}
	collection := rs.Client.Client.Database(rs.Database).Collection(rs.UserCollection)
	var result User
	err2 := collection.FindOne(context.Background(), bson.M{"username": user.Username}).Decode(&result)
	if err2 != nil {
		if err2 == mongo.ErrNoDocuments {
			fmt.Println("User does not exist")
		} else {
			log.Println(err2.Error())
		}
	} else {
		fmt.Println("User exists")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("user exists please enter unique username"))
		return

	}
	_, err1 := collection.InsertOne(context.Background(), user)
	if err1 != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusBadGateway)
		w.Write([]byte("internal error try later"))
		return

	}
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("registration successful"))
}

func (rs *Resource) Login(w http.ResponseWriter, r *http.Request) {
	var user User

	err := json.NewDecoder(r.Body).Decode(&user)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("invalid request please enter correct data"))
		return
	}
	var result User
	collection := rs.Client.Client.Database(rs.Database).Collection(rs.UserCollection)
	err1 := collection.FindOne(context.Background(), bson.M{"username": user.Username, "password": user.Password}).Decode(&result)

	if err1 != nil {
		if err1 == mongo.ErrNoDocuments {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("user does not exists please enter correct username & password"))
			return
		} else {
			w.WriteHeader(http.StatusBadGateway)
			w.Write([]byte("please try later"))
			return
		}
	}

	_, tokenString, _ := rs.TokenAuth.Encode(map[string]interface{}{"username": user.Username})
	w.Write([]byte(tokenString))

}
func (rs *Resource) GetAllList(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var lists []List
	lists = []List{}
	_, token, _ := jwtauth.FromContext(r.Context())
	isCompletedStr := r.URL.Query().Get("isCompleted")
	isDeletedStr := r.URL.Query().Get("isDeleted")
	if isCompletedStr == "" {
		isCompletedStr = "false"
	}
	if isDeletedStr == "" {
		isDeletedStr = "false"
	}
	// isCompleted, err := strconv.ParseBool(isCompletedStr)
	// if err != nil {
	// 	w.WriteHeader(http.StatusBadRequest)
	// 	w.Write([]byte(err.Error()))
	// 	return
	// }
	// isDeleted, err := strconv.ParseBool(isDeletedStr)
	// if err != nil {
	// 	w.WriteHeader(http.StatusBadRequest)
	// 	w.Write([]byte(err.Error()))
	// 	return
	// }

	collection := rs.Client.Client.Database(rs.Database).Collection(rs.ListCollection)
	cursor, err1 := collection.Find(context.Background(), bson.M{"username": token["username"].(string)})
	if err1 != nil {
		log.Println(err1.Error())
		w.WriteHeader(http.StatusBadGateway)
		w.Write([]byte("internal error try later"))
		return

	}
	log.Println(token["username"].(string))

	for cursor.Next(context.Background()) {
		var list List
		if err := cursor.Decode(&list); err != nil {
			panic(err)
		}
		log.Println(list)
		lists = append(lists, list)
	}

	json.NewEncoder(w).Encode(lists)
}

func (rs *Resource) AddList(w http.ResponseWriter, r *http.Request) {
	var list List
	err := json.NewDecoder(r.Body).Decode(&list)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	_, token, _ := jwtauth.FromContext(r.Context())

	list.IsDeleted = false
	list.UserName = token["username"].(string)
	collection := rs.Client.Client.Database(rs.Database).Collection(rs.ListCollection)

	_, err1 := collection.InsertOne(context.Background(), list)
	if err1 != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusBadGateway)
		w.Write([]byte("internal error try later"))
		return

	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(list)
}

func (rs *Resource) UpdateList(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	newId, _ := primitive.ObjectIDFromHex(id)

	var updatedTodo List
	err := json.NewDecoder(r.Body).Decode(&updatedTodo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var updatedList List
	collection := rs.Client.Client.Database(rs.Database).Collection(rs.ListCollection)
	err1 := collection.FindOneAndUpdate(context.Background(), bson.M{"_id": newId}, bson.M{"$set": bson.M{"title": updatedTodo.Title}}).Decode(&updatedList)
	if err1 != nil {
		if err1 == mongo.ErrNoDocuments {
			log.Println("document not found")
			http.NotFound(w, r)
			return
		}
		w.WriteHeader(http.StatusBadGateway)
		w.Write([]byte(err1.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (rs *Resource) DeleteList(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	newId, _ := primitive.ObjectIDFromHex(id)
	var deletedList List
	collection := rs.Client.Client.Database(rs.Database).Collection(rs.ListCollection)
	err1 := collection.FindOneAndUpdate(context.Background(), bson.M{"_id": newId}, bson.M{"$set": bson.M{"isdeleted": true}}).Decode(&deletedList)
	if err1 != nil {
		if err1 == mongo.ErrNoDocuments {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("no item exists"))
			return
		}
		w.WriteHeader(http.StatusNoContent)
		w.Write([]byte("item deleted successfully"))
		return
	}
}

// func filterList(isCompleted, isDeleted bool, username string) []List {
// 	filteredList := []List{}

// 	for _, v := range lists[username] {
// 		if v.IsDeleted == isDeleted && v.IsCompleted == isCompleted {
// 			filteredList = append(filteredList, v)
// 		}
// 	}
// 	return filteredList
// }
