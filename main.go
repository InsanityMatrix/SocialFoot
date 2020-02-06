﻿//Using https://www.sohamkamani.com/blog/2017/09/13/how-to-build-a-web-application-in-golang/
//main.go
package main

import (
    "fmt"
    "net/http"
    "github.com/gorilla/mux"
)

type User struct {
    username string `json:"username"`
    gender boolean `json:"gender"`
    age int `json:"age"`
    password string `json:"password"`
    email string `json:"email"`
}



//Global variables
var users []User
func newRouter() *mux.Router {
    r := mux.NewRouter()
    r.HandleFunc("/user", getUserHandler).Methods("GET")
    r.HandleFunc("/user", createUserHandler).Methods("POST")
    //ALL PAGE FUNCTIONS HERE
    r.HandleFunc("/", handler).Methods("GET")

    //Declare static file directory
    staticFileDirectory := http.Dir("./assets/")

    staticFileHandler := http.StripPrefix("/assets/", http.FileServer(staticFileDirectory))
    
    r.PathPrefix("/assets/").Handler(staticFileHandler).Methods("GET")
    return r
}
func main() {
    router := newRouter
    http.ListenAndServe(":8080", router)
}

func handler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Hello World!")
}

func getUserHandler(w http.ResponseWriter, r *http.Request) {
    userListBytes, err := json.Marshal(users)
    if err != nil {
        fmt.Println(fmt.Errorf("Error: %v", err))
        w.WriteHeader(http.StatusInternalServerError)
        return
    }
    //Write the json list of users to response
    w.Write(userListBytes)
}

func createUserHandler(w http.ResponseWriter, r *http.Request) {
    user := User{}

    //Send all data as HTML form Data so parse form
    err := r.ParseForm()
    if err != nil {
        fmt.Println(fmt.Errorf("Error: %v", err))
        w.WriteHeader(http.StatusInternalServerError)
        return
    }

    //Get the information about the user from user info
    user.username = r.Form.Get("username")
    user.gender = r.Form.Get("gender")
    user.age = r.Form.Get("age")
    user.password = r.Form.Get("age")
    user.email = r.Form.Get("email")

    //Append existing list of users with a new entry
    users = append(users, user)
    //TODO: Create a save to a database json file somewhere

    http.Redirect(w, r, "/assets/", http.StatusFound)    
}