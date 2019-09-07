package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"consul_fabio_demo/utility"

	"github.com/gorilla/mux"
)

type userDetail struct {
	Email   string `json:"email,omitempty"`
	Fname   string `json:"firstname,omitempty"`
	Lname   string `json:"lastname,omitempty"`
	Org     string `json:"organization,omitempty"`
	Title   string `json:"title,omitempty"`
	Address string `json:"address,omitempty"`
}

var emailMap map[string]userDetail

func init() {
	emailMap = make(map[string]userDetail)
}

func main() {
	log.Println("Running Storage Service...")
	c := utility.GetConsulClient()
	tags := []string{"urlprefix-/storage"}
	(*c).Register("StorageService", 8891, &tags)

	router := mux.NewRouter()
	router.HandleFunc("/health", health).Methods(http.MethodGet)
	router.HandleFunc("/storage/user/{email}", getUser).Methods(http.MethodGet)
	router.HandleFunc("/storage/user", createUser).Methods(http.MethodPost)
	router.HandleFunc("/storage/user/{email}", deleteUser).Methods(http.MethodDelete)
	router.HandleFunc("/storage/user/org/{org}", getOrgUsers).Methods(http.MethodGet)
	router.HandleFunc("/storage/users", getUsers).Methods(http.MethodGet)

	log.Fatal(http.ListenAndServe(":8891", router))
}

func health(w http.ResponseWriter, r *http.Request) {
	return
}

func getUser(w http.ResponseWriter, r *http.Request) {
	log.Println("Request for user info")
	params := mux.Vars(r)
	email := params["email"]

	user, ok := emailMap[email]
	w.WriteHeader(http.StatusOK)
	if !ok {
		fmt.Fprintf(w, "user with email %s does not exist", email)
		return
	}
	res, _ := json.Marshal(user)
	fmt.Fprint(w, string(res))
}

func createUser(w http.ResponseWriter, r *http.Request) {
	log.Println("Request to create user")
	email := r.FormValue("email")
	fname := r.FormValue("fname")
	lname := r.FormValue("lname")
	org := r.FormValue("organization")
	title := r.FormValue("title")
	address := r.FormValue("address")

	_, ok := emailMap[email]
	w.WriteHeader(http.StatusOK)
	if ok {
		fmt.Fprintf(w, "User with email %s already exists", email)
		return
	}

	userDetail := userDetail{
		Email:   email,
		Fname:   fname,
		Lname:   lname,
		Org:     org,
		Title:   title,
		Address: address,
	}

	user, _ := json.Marshal(userDetail)
	emailMap[email] = userDetail
	fmt.Fprintf(w, "User created successfully %s", user)
}

func deleteUser(w http.ResponseWriter, r *http.Request) {
	log.Println("Request to delete user")
	params := mux.Vars(r)
	email := params["email"]

	_, ok := emailMap[email]
	w.WriteHeader(http.StatusOK)
	if !ok {
		fmt.Fprintf(w, "user with email %s does not exist", email)
		return
	}
	delete(emailMap, email)
	fmt.Fprintf(w, "user with email %s deleted", email)
}

func getOrgUsers(w http.ResponseWriter, r *http.Request) {
	log.Println("Request for organization users")
	params := mux.Vars(r)
	org := params["org"]

	users := make([]string, 0)

	for _, user := range emailMap {
		if org == user.Org {
			userStr, _ := json.Marshal(user)
			users = append(users, string(userStr))
		}
	}

	w.WriteHeader(http.StatusOK)
	res, _ := json.Marshal(users)
	fmt.Fprint(w, string(res))
}

func getUsers(w http.ResponseWriter, r *http.Request) {
	log.Println("Request of all users")
	users := make([]string, 0)

	for _, user := range emailMap {
		userStr, _ := json.Marshal(user)
		users = append(users, string(userStr))
	}

	w.WriteHeader(http.StatusOK)
	res, _ := json.Marshal(users)
	fmt.Fprint(w, string(res))
}
