//Deploy App: https://devcenter.heroku.com/articles/getting-started-with-go#deploy-the-app
//Using https://www.sohamkamani.com/blog/2017/09/13/how-to-build-a-web-application-in-golang/
//https://developers.google.com/web/fundamentals/native-hardware/fullscreen <- make web app
//main.go
package main

import (
	"database/sql"
	"os"
	"os/exec"
	"strings"
	"path/filepath"
	"io"
	"strconv"
	"html/template"
	"fmt"
	"io/ioutil"
	"regexp"
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

type LiveImagePost struct {
	User string
	imageLink string
}

type FeedData struct {
	username string
	Feed []LiveImagePost
}
var HOME string
var TEMPLATES string
//Global variables
func newRouter() *mux.Router {
    r := mux.NewRouter()
    r.HandleFunc("/user", createUserHandler).Methods("POST")
		r.HandleFunc("/forms/login", loginUserHandler).Methods("POST")
		r.HandleFunc("/forms/signup", createUserHandler).Methods("POST")
		r.HandleFunc("/live/profile/settings", profileSettingsHandler).Methods("POST")
		r.HandleFunc("/live/profile", profileHandler)
		r.HandleFunc("/live/post", postHandler)
		r.HandleFunc("/live/search",searchPageHandler)
		r.HandleFunc("/live", liveIndexHandler)


		//Settings FUNCTIONS
		r.HandleFunc("/settings/user/publicity", changePublicityHandler)
		r.HandleFunc("/settings/user/email", changeEmailHandler)
		r.HandleFunc("/settings/user/location", changeLocationHandler)
		r.HandleFunc("/settings/user/bio", changeBioHandler)
		r.HandleFunc("/settings/user/delete", deleteUserHandler)
		r.HandleFunc("/settings/user/signout", signoutHandler)

		r.HandleFunc("/user/post/imagepost", imagePostHandler).Methods("POST")
		r.HandleFunc("/posts/public", getPublicPostsHandler)
		r.HandleFunc("/search", searchUserHandler).Methods("POST")


		//JSON stuff
		r.HandleFunc("/json/user/id", HandleJSONUserById)
		//report
		r.HandleFunc("/report", reportHandler)
		r.HandleFunc("/report/submit/bugreport", bugReportHandler)

		//TEMPLATES stuff
		r.HandleFunc("/templates/post", postTemplateHandler)
		r.HandleFunc("/templates/result", resultTemplateHandler)
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
		HOME = filepath.Join(os.Getenv("HOME"), "/go/src/github.com/InsanityMatrix/SocialFoot")
		TEMPLATES = "/root/go/src/github.com/InsanityMatrix/SocialFoot/templates"
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
		cmd := exec.Command("(sleep 5; $HOME/go/src/github.com/InsanityMatrix/SocialFoot/SocialFoot &) &")
		_ = cmd.Run()
		return
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
		w.WriteHeader(200)
    //Get the information about the user from user info
    user.username = r.Form.Get("username")
    user.gender, _ = strconv.ParseBool(r.Form.Get("gender"))
    user.age, _ = strconv.Atoi(r.Form.Get("age"))
    user.password = r.Form.Get("password")
		cpassword := r.Form.Get("cpassword")
    user.email = r.Form.Get("email")

		//Email validation
		re := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

		if !re.MatchString(user.email) {
			fmt.Fprint(w, "This email is not valid.")
			return
		}

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

    http.Redirect(w, r, "/live", http.StatusFound)
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
		return
	}
	tmpl, err := template.ParseFiles(TEMPLATES + "/index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, map[string]string{"username":msg.Value})
}
func getPublicPostsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type","application/json")
	pubposts := store.GetPublicPosts()
  fmt.Fprint(w, pubposts)
}
func profileHandler(w http.ResponseWriter, r *http.Request) {
	//Handle Live Profile settings
	w.Header().Set("Content-Type", "text/html")
	msg, err := r.Cookie("username")
	if err != nil {
		http.Redirect(w,r,"/assets/", http.StatusSeeOther)
		return
	}
	tmpl, err := template.ParseFiles(TEMPLATES + "/ProfileSettings.html")
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
		return
	}
	tmpl, err := template.ParseFiles(TEMPLATES + "/post.html")
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
	tmpl.Execute(w, map[string]string{"username":user.username, "publicity":publicity, "xpublicity":xpublicity, "userid":strconv.Itoa(account.id)})
}
func profileSettingsHandler(w http.ResponseWriter, r *http.Request) {
	msg, err := r.Cookie("username")
	if err != nil {
	 	http.Redirect(w,r,"/assets/", http.StatusSeeOther)
		return
	}
	err = r.ParseForm()
	if err != nil {
		fmt.Println(fmt.Errorf("Error: %v", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html")
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
	tmpl, err := template.ParseFiles(TEMPLATES + "/profile.html")
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
	cookie := http.Cookie{
		Name : "username",
		Value: "",
		MaxAge: 0,
		Path: "/",
	}
	http.SetCookie(w, &cookie)
	fmt.Fprint(w,"Success")
}
func reportHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles(TEMPLATES + "/bugReport.html")
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

func imagePostHandler(w http.ResponseWriter, r *http.Request) {
	//Setup and get image
	in, header, err := r.FormFile("upload")

	if err != nil {
		fmt.Fprint(w, err.Error())
		return
	}


	err = r.ParseForm()
	if err != nil {
		panic(err.Error())
	}

	defer in.Close()
	userid := r.Form.Get("id")
	//VALIDATE FILE TYPE
	extension := strings.ToLower(filepath.Ext(header.Filename))
	if !isPictureFile(extension) {
		tmpl, err := template.ParseFiles(TEMPLATES + "/uploadSuccess.html")
		msg, err := r.Cookie("username")
		var username string
		if err != nil {
			id, _ := strconv.Atoi(userid)
			username = store.GetUserInfoById(id).username
		} else {
			username = msg.Value
		}
		status := "This file type is not supported. Try again with a different file type."
		tmpl.Execute(w, map[string]string{"username":username,"status":status})
		return
	}
	//Get Form Values
	caption := r.Form.Get("caption")
	tags := r.Form.Get("tags")

	tagsRe := regexp.MustCompile("((#\\w+),?\\s?)")

	TAGS := ""
	matches := tagsRe.FindAllStringSubmatch(tags, -1)
	for _, group := range matches {
		TAGS += group[2]
	}
	publicityText := r.Form.Get("type")


	//Set publicity
	publicity := true
	if publicityText == "Private" {
		publicity = false
	}


	//Actually post image
	id, _ := strconv.Atoi(userid)
	postid := store.PostUserImage(publicity, caption, tags, id,extension)
	if postid == 0 {
		//ERROR case
		fmt.Fprint(w, "Could not return post id or insert row")
	}
	idStr := strconv.Itoa(postid)
	out, err := os.Create("/root/go/src/github.com/InsanityMatrix/SocialFoot/assets/uploads/imageposts/post" + idStr + extension)
	if err != nil {
		//handle error
		panic(err.Error())
	}
	defer out.Close()
	io.Copy(out, in)
	//Display results
	tmpl, err := template.ParseFiles(TEMPLATES + "/uploadSuccess.html")
	msg, err := r.Cookie("username")
	var username string
	if err != nil {
		username = store.GetUserInfoById(id).username
	} else {
		username = msg.Value
	}
	status := "Your post has been created at http://www.socialfoot.me/live/view/post/" + userid + "." + idStr
	tmpl.Execute(w, map[string]string{"username":username,"status":status})


}
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

func searchPageHandler(w http.ResponseWriter, r *http.Request) {
	msg, err := r.Cookie("username")
	if err != nil {
		http.Redirect(w,r,"/assets/login.html",http.StatusSeeOther)
		return
	}
	tmpl, err := template.ParseFiles(TEMPLATES + "/search.html")
	if err != nil {
	  w.WriteHeader(http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, map[string]string{"username": msg.Value})
}

func HandleJSONUserById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type","application/json")
	_ = r.ParseForm()
	id, _ := strconv.Atoi(r.Form.Get("userid"))
	userjson := store.GetJSONUserByID(id)
	fmt.Fprint(w, userjson)
}
func searchUserHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		w.WriteHeader(404)
		return
	}
	w.Header().Set("Content-Type","application/json")
	term := r.Form.Get("term")
	response := store.GetJSONUsersByUsernames(term)
	fmt.Fprint(w, response)
}
func postTemplateHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	file, err := os.Open(TEMPLATES + "/feed/post.html")
	if err != nil {
		fmt.Fprint(w, "Error")
	}
	defer file.Close()
	data, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Fprint(w, "Error")
	}
	fmt.Fprint(w, string(data))
}
func resultTemplateHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	file, err := os.Open(TEMPLATES + "/search/result.html")
	if err != nil {
		fmt.Fprint(w, "Error")
	}
	defer file.Close()
	data, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Fprint(w, "Error")
	}
	fmt.Fprint(w, string(data))
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
