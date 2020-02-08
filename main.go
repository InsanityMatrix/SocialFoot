//Deploy App: https://devcenter.heroku.com/articles/getting-started-with-go#deploy-the-app
//Using https://www.sohamkamani.com/blog/2017/09/13/how-to-build-a-web-application-in-golang/
//https://developers.google.com/web/fundamentals/native-hardware/fullscreen <- make web app
//main.go
package main

import (
	"database/sql"
	"os"
	"strconv"
	"html/template"
	"fmt"
	"time"
		"log"
    "net/http"
    "github.com/gorilla/mux"
	 _ "github.com/lib/pq"
)

type User struct {
		id int `json:"id"`
    username string `json:"username"`
    gender bool `json:"gender"`
    age int `json:"age"`
    password string `json:"password"`
    email string `json:"email"`
}

type UserSettings struct {
	id int
	bio string
	website string
	location string
	publicity bool
}

type UserInfo struct {
	 id int
	 username string
	 gender bool
	 age int
	 email string
	 bio string
	 website string
	 location string
	 publicity bool
}

var IndexHTML string
var ProfileHTML string
var ProfileSettingsHTML string
//Global variables
func newRouter() *mux.Router {
    r := mux.NewRouter()
    r.HandleFunc("/user", createUserHandler).Methods("POST")
		r.HandleFunc("/forms/login", loginUserHandler).Methods("POST")
		r.HandleFunc("/forms/signup", createUserHandler).Methods("POST")
		r.HandleFunc("/live/profile/settings", profileSettingsHandler).Methods("POST")
		r.HandleFunc("/live/profile", profileHandler)
		r.HandleFunc("/live", liveIndexHandler)
    //ALL PAGE FUNCTIONS HERE
    r.HandleFunc("/", handler)

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


		url := os.Getenv("DATABASE_URL")
		db, err := sql.Open("postgres", url)

		if err != nil {
			log.Fatalf("Connection error: %s", err.Error())
			panic(err)
		}
		defer db.Close()

		err = db.Ping()

		if err != nil {
			log.Fatalf("Ping error: %s", err.Error())
			panic(err)
		}
	//Set Connection Limit: https://www.alexedwards.net/blog/configuring-sqldb
		db.SetMaxOpenConns(15)
		db.SetMaxIdleConns(4)
		db.SetConnMaxLifetime(time.Hour)
		InitStore(dbStore{db: db})
		IndexHTML = initIndexHTML()
		ProfileHTML = initProfileHTML()
		ProfileSettingsHTML = initProfileSettingsHTML()
		http.ListenAndServe(port, router)
}

func handler(w http.ResponseWriter, r *http.Request) {
	_, err := r.Cookie("username")
	if err != nil {
		http.Redirect(w,r,"/assets/",http.StatusSeeOther)
	}
	http.Redirect(w,r,"/live",http.StatusSeeOther)
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
		cpassword := r.Form.Get("cpassword")
    user.email = r.Form.Get("email")

		if(user.password != cpassword) {
			http.Redirect(w, r, "/assets/signup.html", http.StatusSeeOther)
			return
		}
		user.password = hashAndSalt([]byte(user.password))
    //Append existing list of users with a new entry
    err = store.CreateUser(&user)
	if err != nil {
		log.Println(err)
		fmt.Println(err)
	}
  //Set Cookie with username
		addCookie(w, "username", user.username)

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
	account, err := store.LoginUser(&user)
	if err != nil {
		//Username may not have been right
		http.Redirect(w,r,"/assets/login.html", http.StatusSeeOther)
	}
	if comparePasswords(account.password, bytePass) {
		//Logged In
		addCookie(w,"username",account.username)

		http.Redirect(w, r, "/live", http.StatusFound)
	} else {
		http.Redirect(w,r,"/assets/login.html", http.StatusSeeOther)
	}

}
func liveIndexHandler(w http.ResponseWriter, r *http.Request) {
	//Handle Live page with html templates
	msg, err := r.Cookie("username")
	if err != nil {
		http.Redirect(w,r,"/assets/",http.StatusSeeOther)
	}
	tmpl, err := template.New("Index").Parse(IndexHTML)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, map[string]string{"username":msg.Value})

}
func profileHandler(w http.ResponseWriter, r *http.Request) {
	//Handle Live Profile settings
	msg, err := r.Cookie("username")
	if err != nil {
		http.Redirect(w,r,"/assets/", http.StatusSeeOther)
	}
	tmpl, err := template.New("Profile").Parse(ProfileHTML)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, map[string]string{"username":msg.Value})
}
func profileSettingsHandler(w http.ResponseWriter, r *http.Request) {
	msg, err := r.Cookie("username")
	if err != nil {
	 	http.Redirect(w,r,"/assets/", http.StatusSeeOther)
	}
	err = r.ParseForm()
	if err != nil {
		fmt.Println(fmt.Errorf("Error: %v", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	user := User{}
	user.username = msg.Value
	user.password = r.Form.Get("password")
	verified := store.CheckUserCredentials(&user)
	if !verified {
		http.Redirect(w, r, "/live/profile", http.StatusSeeOther)
	}
	//User is Verified
	account := store.GetUserInfo(&user)
	settings := store.GetUserSettings(&user)
	info := UserInfo{account.id, account.username, account.gender, account.age, account.email, settings.bio,settings.website, settings.location, settings.publicity}
	tmpl, err := template.New("ProfileSettings").Parse(ProfileHTML)
	if err != nil {
		http.Redirect(w, r, "/live", http.StatusInternalServerError)
	}
	tmpl.Execute(w, info)
}
func addCookie(w http.ResponseWriter, name string, value string) {
    cookie := http.Cookie{
        Name:    name,
        Value:   value,
				Path: "/",
    }
    http.SetCookie(w, &cookie)
}
