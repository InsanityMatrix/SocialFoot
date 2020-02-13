//Deploy App: https://devcenter.heroku.com/articles/getting-started-with-go#deploy-the-app
//Using https://www.sohamkamani.com/blog/2017/09/13/how-to-build-a-web-application-in-golang/
//https://developers.google.com/web/fundamentals/native-hardware/fullscreen <- make web app
//main.go
package main

import (
	"database/sql"
	"os"
	//"io"
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

//Global variables
func newRouter() *mux.Router {
    r := mux.NewRouter()
    r.HandleFunc("/user", createUserHandler).Methods("POST")
		r.HandleFunc("/forms/login", loginUserHandler).Methods("POST")
		r.HandleFunc("/forms/signup", createUserHandler).Methods("POST")
		r.HandleFunc("/live/profile/settings", profileSettingsHandler).Methods("POST")
		r.HandleFunc("/live/profile", profileHandler)
		r.HandleFunc("/live/post", postHandler)
		r.HandleFunc("/live", liveIndexHandler)


		//Settings FUNCTIONS
		r.HandleFunc("/settings/user/publicity", changePublicityHandler)
		r.HandleFunc("/settings/user/email", changeEmailHandler)
		r.HandleFunc("/settings/user/location", changeLocationHandler)
		r.HandleFunc("/settings/user/bio", changeBioHandler)
		r.HandleFunc("/settings/user/delete", deleteUserHandler)
		r.HandleFunc("/settings/user/signout", signoutHandler)

		//r.HandleFunc("/user/post/imagepost", imagePostHandler)
		//report
		r.HandleFunc("/report", reportHandler)
		r.HandleFunc("/report/submit/bugreport", bugReportHandler)
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

    http.Redirect(w, r, "/live/", http.StatusFound)
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
	w.Header().Set("Content-Type", "text/html")
	msg, err := r.Cookie("username")
	if err != nil {
		http.Redirect(w,r,"/assets/",http.StatusSeeOther)
	}
	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, map[string]string{"username":msg.Value})

}
func profileHandler(w http.ResponseWriter, r *http.Request) {
	//Handle Live Profile settings
	w.Header().Set("Content-Type", "text/html")
	msg, err := r.Cookie("username")
	if err != nil {
		http.Redirect(w,r,"/assets/", http.StatusSeeOther)
	}
	tmpl, err := template.ParseFiles("templates/ProfileSettings.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, map[string]string{"username":msg.Value})
}
func postHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type","text/html")
	msg, err := r.Cookie("username")
	if err != nil {
		http.Redirect(w,r,"/assets/login.html", http.StatusSeeOther)
	}
	tmpl, err := template.ParseFiles("templates/post.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	user := User{}
	user.username = msg.Value
	account := store.GetUserInfo(&user)
	settings := store.GetUserSettings(account)

	publicity := "Private"
	xpublicity := "Public"
	if settings.publicity {
		publicity = "Public"
		xpublicity = "Private"
	}
	tmpl.Execute(w, map[string]string{"username":user.username, "publicity":publicity, "xpublicity":xpublicity})
}
func profileSettingsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
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
	//User is Verified
	account := store.GetUserInfo(&user)
	if account == nil {
		http.Redirect(w,r,"/live/profile", http.StatusSeeOther)
		return
	}
	verified := comparePasswords(account.password, []byte(user.password))
	if !verified {
		http.Redirect(w,r,"/live/profile", http.StatusSeeOther)
		return
	}
	settings := store.GetUserSettings(account)
	publicity := "Private"
	if settings.publicity {
		publicity = "Public"
	}
	tmpl, err := template.ParseFiles("templates/profile.html")
	if err != nil {
		http.Redirect(w, r, "/live", http.StatusInternalServerError)
	}
	idVal := strconv.Itoa(account.id)
	tmpl.Execute(w, map[string]string{"id": idVal, "username":account.username,"email":account.email, "publicity":publicity, "location":settings.location, "bio":settings.bio})
}

//SETTINGS FUNCTIONS
func changePublicityHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		fmt.Println(fmt.Errorf("Error: %v", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	userID, _ := strconv.Atoi(r.Form.Get("userID"))
	status := r.Form.Get("status")
	if status == "mPrivate" {
		//Make user private
		store.SetUserPublicity(userID, false)
		fmt.Fprint(w, "Private")
		return
	}
	store.SetUserPublicity(userID, true)
	fmt.Fprint(w, "Public")
}
func changeEmailHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		fmt.Println(fmt.Errorf("Error: %v", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	userID, _ := strconv.Atoi(r.Form.Get("userID"))
	email := r.Form.Get("email")
	status := store.ChangeUserEmail(userID, email)
	if !status {
		fmt.Fprint(w, "Failed")
		return
	}
	fmt.Fprint(w, "Success")
}
func changeLocationHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		fmt.Println(fmt.Errorf("Error: %v", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	userID, _ := strconv.Atoi(r.Form.Get("userID"))
	location := r.Form.Get("location")
	status := store.ChangeUserLocation(userID, location)
	if !status {
		fmt.Fprint(w, "Failed")
		return
	}
	fmt.Fprint(w, "Success")
}
func changeBioHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		fmt.Println(fmt.Errorf("Error: %v", err))
		fmt.Fprint(w, "Failed")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	userID, _ := strconv.Atoi(r.Form.Get("userID"))
	bio := r.Form.Get("bio")
	status := store.ChangeUserBio(userID, bio)
	if !status {
		fmt.Fprint(w, "Failed")
		return
	}
	fmt.Fprint(w, "Success")
}
func deleteUserHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		fmt.Println(fmt.Errorf("Error: %v", err))
		fmt.Fprint(w, "Failed")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	userID, _ := strconv.Atoi(r.Form.Get("userID"))
	status := store.DeleteAccount(userID)
	if !status {
		fmt.Fprint(w, "Failed")
		return
	}
	c, err := r.Cookie("username")
    if err != nil {
        panic(err.Error())
    }
  c.Value = "Anonymous"
  c.Expires = time.Unix(1414414788, 1414414788000)
	fmt.Fprint(w,"Success")
}
func signoutHandler(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("username")
    if err != nil {
        panic(err.Error())
    }
  c.Value = "Anonymous"
  c.Expires = time.Unix(1414414788, 1414414788000)
	fmt.Fprint(w,"Success")
}
func reportHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/bugReport.html")
	if err != nil {
		http.Redirect(w, r, "/live", http.StatusInternalServerError)
	}
	msg, err := r.Cookie("username")
	username := "Anonymous"
	if err == nil {
		username = msg.Value
	}
	tmpl.Execute(w, map[string]string{"username":username})
}
/*
func imagePostHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		panic(err.Error())
	}
	in, header, err := r.FormFile("upload")
	if err != nil {
		fmt.Fprint(w, "Could not parse upload")
		return
	}
	defer in.Close()


	caption := r.Form.Get("caption")
	tags := r.Form.Get("tags")
	username := r.Form.Get("username")
	userid := r.Form.Get("id")


	//Actually post image

	//postid := store.PostUserImage()
	out, err := os.OpenFile("./assets/uploads/imageposts/post", os.O_WRONLY, 0644)
	if err != nil {
		//handle error
		panic(err.Error())
	}
	defer out.Close()
	io.Copy(out, in)


}*/
func bugReportHandler(w http.ResponseWriter, r *http.Request) {
	msg, err := r.Cookie("username")
	username := "Anonymous"
	if err == nil {
		username = msg.Value
	}
	err = r.ParseForm()
	if err != nil {
		fmt.Println(fmt.Errorf("Error: %v", err))
		fmt.Fprint(w, "Failed")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	content := r.Form.Get("report")
	store.SubmitBugReport(username, content)
	http.Redirect(w, r, "/live", http.StatusSeeOther)
}
//Page functions to help with stuff
func addCookie(w http.ResponseWriter, name string, value string) {
    cookie := http.Cookie{
        Name:    name,
        Value:   value,
				Path: "/",
    }
    http.SetCookie(w, &cookie)
}
