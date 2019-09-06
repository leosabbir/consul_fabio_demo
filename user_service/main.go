package main

import (
	"consul_fabio_demo/utility"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/gorilla/mux"
)

const (
	fabiopath = "http://172.31.27.241:9999"
)

func main() {
	log.Println("Running User Service...")
	c := utility.GetConsulClient()
	tags := []string{"urlprefix-/service"}
	(*c).Register("UserService", 8890, &tags)

	router := mux.NewRouter()
	router.HandleFunc("/service/user/{email}", getUser).Methods(http.MethodGet)
	router.HandleFunc("/service/user", createUser).Methods(http.MethodPost)
	router.HandleFunc("/service/user/{email}", deleteUser).Methods(http.MethodDelete)
	router.HandleFunc("/service/user/org/{org}", getOrgUsers).Methods(http.MethodGet)
	router.HandleFunc("/service/users", getUsers).Methods(http.MethodGet)

	log.Fatal(http.ListenAndServe(":8890", router))
}

func getUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	email := params["email"]

	res, code, err := utility.SendRequest(fabiopath+"/user/"+email, http.MethodGet, nil, nil)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, err.Error())
		return
	}
	w.WriteHeader(code)
	fmt.Fprint(w, string(*res))
}

func createUser(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	fname := r.FormValue("fname")
	lname := r.FormValue("lname")
	org := r.FormValue("organization")
	title := r.FormValue("title")
	address := r.FormValue("address")

	data := url.Values{}
	data.Set("email", email)
	data.Set("fname", fname)
	data.Set("lname", lname)
	data.Set("org", org)
	data.Set("title", title)
	data.Set("address", address)

	res, code, err := utility.SendRequest(fabiopath+"/user/", http.MethodPost, strings.NewReader(data.Encode()), nil)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, err.Error())
		return
	}
	w.WriteHeader(code)
	fmt.Fprint(w, string(*res))
}

func deleteUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	email := params["email"]

	res, code, err := utility.SendRequest(fabiopath+"/user/"+email, http.MethodDelete, nil, nil)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, err.Error())
		return
	}
	w.WriteHeader(code)
	fmt.Fprint(w, string(*res))
}

func getOrgUsers(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	org := params["org"]

	res, code, err := utility.SendRequest(fabiopath+"/user/org/"+org, http.MethodGet, nil, nil)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, err.Error())
		return
	}
	w.WriteHeader(code)
	fmt.Fprint(w, string(*res))
}

func getUsers(w http.ResponseWriter, r *http.Request) {
	res, code, err := utility.SendRequest(fabiopath+"/users", http.MethodGet, nil, nil)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, err.Error())
		return
	}
	w.WriteHeader(code)
	fmt.Fprint(w, string(*res))
}
