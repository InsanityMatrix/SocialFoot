//Deploy App: https://devcenter.heroku.com/articles/getting-started-with-go#deploy-the-app
//Using https://www.sohamkamani.com/blog/2017/09/13/how-to-build-a-web-application-in-golang/
//main.go
package main

import (
	"database/sql"
	"encoding/json"
	"strconv"
	"os"
    "fmt"
    "net/http"
    "github.com/gorilla/mux"
	"time"
	_ "github.com/lib/pq"
)

type User struct {
    username string `json:"username"`
    gender bool `json:"gender"`
    age int `json:"age"`
    password string `json:"password"`
    email string `json:"email"`
}



//Global variables
func newRouter() *mux.Router {
    r := mux.NewRouter()
    r.HandleFunc("/user", getUserHandler).Methods("GET")
    r.HandleFunc("/user", createUserHandler).Methods("POST")
	r.HandleFunc("/forms/login", loginUserHandler).Methods("POST")
    //ALL PAGE FUNCTIONS HERE
    r.HandleFunc("/", handler).Methods("GET")

    //Declare static file directory
    staticFileDirectory := http.Dir("./assets/")

    staticFileHandler := http.StripPrefix("/assets/", http.FileServer(staticFileDirectory))

    r.PathPrefix("/assets/").Handler(staticFileHandler).Methods("GET")
    return r
}
func main() {
    router := newRouter()
    portEnv := os.Getenv("PORT")
    port := ":" + portEnv
    http.ListenAndServe(port, router)

	connString := "dbname=postgresql-amorphous-36239 sslmode=disable"
	db, err := sql.Open("postgres", connString)

	if err != nil {
		panic(err)
	}
	err = db.Ping()

	if err != nil {
		panic(err)
	}

	InitStore(&dbStore{db: db})
}

func handler(w http.ResponseWriter, r *http.Request) {
    http.Redirect(w, r, "/assets/", http.StatusSeeOther)
}

func getUserHandler(w http.ResponseWriter, r *http.Request) {
    users, err := store.GetUsers()

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
    user.gender, _ = strconv.ParseBool(r.Form.Get("gender"))
    user.age, _ = strconv.Atoi(r.Form.Get("age"))
    user.password = r.Form.Get("password")
    user.email = r.Form.Get("email")

    //Append existing list of users with a new entry
    err = store.CreateUser(&user)
	if err != nil {
		fmt.Println(err)
	}
    //TODO: Create a save to a database json file somewhere

    http.Redirect(w, r, "/assets/", http.StatusFound)
}
func loginUserHandler(w http.ResponseWriter, r *http.Request) {
	user := User{}

	err := r.ParseForm()
	if err != nil {
		fmt.Println(fmt.Errorf("Error: %v", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	user.username = r.Form.Get("username")
	user.password = r.Form.Get("password")
	bytePass := []byte(user.password)
	newPass := hashAndSalt(bytePass)
	user.password = newPass
}
func addCookie(w http.ResponseWriter, name string, value string) {
    expire := time.Now().AddDate(0, 0, 1)
    cookie := http.Cookie{
        Name:    name,
        Value:   value,
        Expires: expire,
    }
    http.SetCookie(w, &cookie)
}
