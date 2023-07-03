package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/go-chi/jwtauth"
)

type List struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	IsCompleted bool   `json:"is_completed"`
	IsDeleted   bool   `json:"is_deleted"`
}

var lists map[string][]List = map[string][]List{}

func getAllList(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	_, token, _ := jwtauth.FromContext(r.Context())

	isCompletedStr := r.URL.Query().Get("isCompleted")
	isDeletedStr := r.URL.Query().Get("isDeleted")
	if isCompletedStr == "" {
		isCompletedStr = "false"
	}
	if isDeletedStr == "" {
		isDeletedStr = "false"
	}
	isCompleted, err := strconv.ParseBool(isCompletedStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	isDeleted, err := strconv.ParseBool(isDeletedStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	newList := filterList(isCompleted, isDeleted, token["username"].(string))
	if len(newList) == 0 {
		newList = []List{}
	}
	json.NewEncoder(w).Encode(newList)
}

func addList(w http.ResponseWriter, r *http.Request) {
	var list List

	err := json.NewDecoder(r.Body).Decode(&list)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	_, token, _ := jwtauth.FromContext(r.Context())
	list.ID = len(lists[token["username"].(string)]) + 1
	list.IsDeleted = false
	lists[token["username"].(string)] = append(lists[token["username"].(string)], list)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(list)
}

func updateList(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	num, err := strconv.Atoi(id)
	if err != nil {
		w.Write([]byte("invalid id"))
		return
	}
	_, token, _ := jwtauth.FromContext(r.Context())
	for i, list := range lists[token["username"].(string)] {
		if list.ID == num {
			var updatedTodo List
			err := json.NewDecoder(r.Body).Decode(&updatedTodo)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			// Update the existing todo
			lists[token["username"].(string)][i].Title = updatedTodo.Title
			lists[token["username"].(string)][i].IsCompleted = updatedTodo.IsCompleted

			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(lists[token["username"].(string)])
			return
		}
	}

	http.NotFound(w, r)
}

func deleteList(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	num, err := strconv.Atoi(id)
	if err != nil {
		w.Write([]byte("invalid id"))
		return
	}
	_, token, _ := jwtauth.FromContext(r.Context())
	for i, list := range lists[token["username"].(string)] {
		if list.ID == num {
			// Delete the todo from the slice
			// lists = append(lists[:i], lists[i+1:]...)
			lists[token["username"].(string)][i].IsDeleted = true

			w.WriteHeader(http.StatusNoContent)
			w.Write([]byte("item deleted successfully"))
			return
		}
	}

	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("no item exists"))
}

func filterList(isCompleted, isDeleted bool, username string) []List {
	filteredList := []List{}

	for _, v := range lists[username] {
		if v.IsDeleted == isDeleted && v.IsCompleted == isCompleted {
			filteredList = append(filteredList, v)
		}
	}
	return filteredList
}
