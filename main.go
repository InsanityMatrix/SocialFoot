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
	"encoding/json"
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
type Post struct {
	Postid int
	Userid int
	Tags string
	Caption string
	Type string
	Posted time.Time
	Extension string
	Publicity bool
	Likes int
}
type Follower struct {
	Relid int `json:"relid"`
	Userid int `json:"userid"`
	Followed string `json:"followed"`
}
type FollowerResult struct {
	FUsername string
	FFollowed string
}
type FollowersPageData struct {
	Activity string
	Userid int
	Username string
	Followers []FollowerResult
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
type Message struct {
	MessageID int
	Content string
	From int
	Read bool
}
type Conversation struct {
	ConvoID int
	ParticipantID int
	ParticipantName string
	Created time.Time
}
type MessagePage struct {
	Username string
	Conversations []Conversation
	Messages []Message
	Home bool
}
var HOME string
var TEMPLATES string
//Global variables
func newRouter() *mux.Router {
    r := mux.NewRouter()
		r.HandleFunc("/favicon.ico", faviconHandler)
    r.HandleFunc("/user", createUserHandler).Methods("POST")
		r.HandleFunc("/forms/login", loginUserHandler).Methods("POST")
		r.HandleFunc("/forms/signup", createUserHandler).Methods("POST")
		r.HandleFunc("/live/profile/settings", profileSettingsHandler).Methods("POST")
		r.HandleFunc("/live/profile", profileHandler)
		r.HandleFunc("/live/view/post/{postid}", viewPostHandler)
		r.HandleFunc("/live/post", postHandler)
		r.HandleFunc("/live/user/followers/{uid}", userFollowersHandler)
		r.HandleFunc("/live/user/following/{uid}", userFollowingHandler)
		r.HandleFunc("/live/user/posts", userPostHandler)
		r.HandleFunc("/live/messages", loadMessages)
		r.HandleFunc("/live/messages/{convoid}", conversationHandler)
		r.HandleFunc("/live/search",searchPageHandler)
		r.HandleFunc("/live/user/{uid}", userProfileHandler)
		r.HandleFunc("/live", liveIndexHandler)


		//Settings FUNCTIONS
		r.HandleFunc("/settings/user/publicity", changePublicityHandler)
		r.HandleFunc("/settings/user/email", changeEmailHandler)
		r.HandleFunc("/settings/user/location", changeLocationHandler)
		r.HandleFunc("/settings/user/bio", changeBioHandler)
		r.HandleFunc("/settings/user/delete", deleteUserHandler)
		r.HandleFunc("/settings/user/signout", signoutHandler)
		r.HandleFunc("/user/follow", followUserHandler)
		r.HandleFunc("/user/isfollowing", isFollowingUserHandler)
		r.HandleFunc("/user/post/imagepost", imagePostHandler).Methods("POST")
		r.HandleFunc("/posts/public", getPublicPostsHandler)
		r.HandleFunc("/search", searchUserHandler).Methods("POST")

		//MESSAGES FUNCTIONS
		r.HandleFunc("/messages/send/text", sendTextMessageHandler)
		r.HandleFunc("/messages/create/private", createPrivateMessageHandler)

		//JSON stuff
		r.HandleFunc("/json/user/id", HandleJSONUserById)
		r.HandleFunc("/json/messages/convo", getMessages)
		//report
		r.HandleFunc("/report", reportHandler)
		r.HandleFunc("/report/submit/bugreport", bugReportHandler)

		//TEMPLATES stuff
		r.HandleFunc("/templates/post", postTemplateHandler)
		r.HandleFunc("/templates/result", resultTemplateHandler)
		r.HandleFunc("/templates/tomsg", toMsgTemplateHandler)
		r.HandleFunc("/templates/frommsg", fromMsgTemplateHandler)
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
		return
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
			if err.Error() == "User Exists" {
				http.Redirect(w, r, "/assets/signup.html", http.StatusSeeOther)
				return
			}
			log.Println(err)
			fmt.Println(err)
			return
		}
  //Set Cookie with username
		setEncryptedCookie(w, "username", []byte(user.username))
		//Wait for like 1 second
    http.Redirect(w, r, "/live", http.StatusSeeOther)
}
func loginUserHandler(w http.ResponseWriter, r *http.Request) {
	user := User{}

	err := r.ParseForm()
	if err != nil {
		http.Redirect(w,r,"/assets/login.html", http.StatusSeeOther)
		return
	}

	user.username = r.Form.Get("username")
	user.password = r.Form.Get("password")
	bytePass := []byte(user.password)
	account, err := store.LoginUser(&user)
	if err != nil {
		//Username may not have been right
		http.Redirect(w,r,"/assets/login.html", http.StatusSeeOther)
		return
	}
	if comparePasswords(account.password, bytePass) {
		//Logged In
		setEncryptedCookie(w,"username", []byte(account.username))

		http.Redirect(w, r, "/live", http.StatusSeeOther)
		return
	} else {
		http.Redirect(w,r,"/assets/login.html", http.StatusSeeOther)
		return
	}

}
func liveIndexHandler(w http.ResponseWriter, r *http.Request) {
	//Handle Live page with html templates
	w.Header().Set("Content-Type", "text/html")
	name, err := decryptCookie(r, "username")
	if err != nil {
		http.Redirect(w,r,"/assets/",http.StatusSeeOther)
		return
	}
	tmpl, err := template.ParseFiles(TEMPLATES + "/index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, map[string]string{"username": name})
}
func getPublicPostsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type","application/json")
	pubposts := store.GetPublicPosts()
  fmt.Fprint(w, pubposts)
}
func profileHandler(w http.ResponseWriter, r *http.Request) {
	//Handle Live Profile settings
	w.Header().Set("Content-Type", "text/html")
	name, err := decryptCookie(r, "username")
	if err != nil {
		http.Redirect(w,r,"/assets/", http.StatusSeeOther)
		return
	}
	tmpl, err := template.ParseFiles(TEMPLATES + "/ProfileSettings.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, map[string]string{"username":name})
}
func postHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type","text/html")
	name, err := decryptCookie(r, "username")
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
	user.username = name
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
	name, err := decryptCookie(r, "username")
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
	user.username = name
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
  c.Value = ""
  c.MaxAge = -1
	fmt.Fprint(w,"Success")
}
func signoutHandler(w http.ResponseWriter, r *http.Request) {
	c := &http.Cookie{
		Name:    "username",
		Value: "",
		Path: "/",
		MaxAge: -1,
		HttpOnly:true,
	}
	http.SetCookie(w, c)
	fmt.Fprint(w,"Success")
}
func reportHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles(TEMPLATES + "/bugReport.html")
	if err != nil {
		http.Redirect(w, r, "/live", http.StatusInternalServerError)
	}
	msg, err := decryptCookie(r, "username")
	username := "Anonymous"
	if err == nil {
		username = msg
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
	ft, supported := isSupportedFile(extension)
	if !supported {
		tmpl, err := template.ParseFiles(TEMPLATES + "/uploadSuccess.html")
		msg, err := decryptCookie(r, "username")
		var username string
		if err != nil {
			id, _ := strconv.Atoi(userid)
			username = store.GetUserInfoById(id).username
		} else {
			username = msg
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
	postid := store.PostUserImage(publicity, caption, tags, id,extension, ft)
	if postid == 0 {
		//ERROR case
		fmt.Fprint(w, "Could not return post id or insert row")
	}
	idStr := strconv.Itoa(postid)
	var out *os.File
	if ft == "IMAGE" {
		out, _  = os.Create("/root/go/src/github.com/InsanityMatrix/SocialFoot/assets/uploads/imageposts/post" + idStr + extension)
	} else {
		out, _  = os.Create("/root/go/src/github.com/InsanityMatrix/SocialFoot/assets/uploads/videoposts/post" + idStr + extension)
	}

	defer out.Close()
	io.Copy(out, in)

	if ft == "VIDEO" {
		fInfo, _ := out.Stat()
		if fInfo.Size() > 21000000 {
			store.DeleteUserPost(postid)
			os.Remove("/root/go/src/github.com/InsanityMatrix/SocialFoot/assets/uploads/videoposts/post" + idStr + extension)
			tmpl, err := template.ParseFiles(TEMPLATES + "/uploadSuccess.html")
			msg, err := decryptCookie(r, "username")
			var username string
			if err != nil {
				id, _ := strconv.Atoi(userid)
				username = store.GetUserInfoById(id).username
			} else {
				username = msg
			}
			status := "This file type is too large."
			tmpl.Execute(w, map[string]string{"username":username,"status":status})
			return
		}
	}
	//Display results
	tmpl, err := template.ParseFiles(TEMPLATES + "/uploadSuccess.html")
	msg, err := decryptCookie(r, "username")
	var username string
	if err != nil {
		username = store.GetUserInfoById(id).username
	} else {
		username = msg
	}
	status := "Your post has been created at http://www.socialfoot.me/live/view/post/" + userid + "." + idStr
	tmpl.Execute(w, map[string]string{"username":username,"status":status})
}
func bugReportHandler(w http.ResponseWriter, r *http.Request) {
	msg, err := decryptCookie(r, "username")

	if err != nil {
		http.Redirect(w, r, "/assets/login.html", http.StatusSeeOther)
		return
	}

	username := msg
	err = r.ParseForm()
	if err != nil {
		fmt.Println(fmt.Errorf("Error: %v", err))
		fmt.Fprint(w, "Failed")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	content := r.Form.Get("report")
	if content == "" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	if badReport(content) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	store.SubmitBugReport(username, content)
	http.Redirect(w, r, "/live", http.StatusSeeOther)
}

func searchPageHandler(w http.ResponseWriter, r *http.Request) {
	msg, err := decryptCookie(r, "username")
	if err != nil {
		http.Redirect(w,r,"/assets/login.html",http.StatusSeeOther)
		return
	}
	tmpl, err := template.ParseFiles(TEMPLATES + "/search.html")
	if err != nil {
	  w.WriteHeader(http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, map[string]string{"username": msg})
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
func userFollowingHandler(w http.ResponseWriter, r *http.Request) {
	params := strings.Split(r.URL.Path, "/")
	userid, err := strconv.Atoi(params[len(params) - 1])
	if err != nil {
		http.Redirect(w, r, "/live/search", http.StatusSeeOther)
		return
	}
	w.Header().Set("Content-Type","text/html")
	userViewing := store.GetUserInfoById(userid)
	followersJSON := store.GetUserFollowing(userid)
	followersSTR := "["
	for _, data := range followersJSON {
		followersSTR += data
	}
	followersSTR += "]"

	var result []Follower
	json.Unmarshal([]byte(followersSTR),&result)
	var fResult []FollowerResult
	for _, fData := range result {
		fUser := store.GetUserInfoById(fData.Userid)
		followedDate, _ := time.Parse("2006-01-02T15:04:05Z", fData.Followed)
		newResult := FollowerResult{FUsername: fUser.username, FFollowed: followedDate.Format("01/02/2006")}
		fResult = append(fResult, newResult)
	}

	tmpl, _ := template.ParseFiles(TEMPLATES + "/user/followers.html")
	tmpl.Execute(w, FollowersPageData{Activity: "Follows", Userid: userid, Username: userViewing.username, Followers: fResult})
}
func userFollowersHandler(w http.ResponseWriter, r *http.Request) {
	params := strings.Split(r.URL.Path, "/")
	userid, err := strconv.Atoi(params[len(params) - 1])
	if err != nil {
		http.Redirect(w, r, "/live/search", http.StatusSeeOther)
		return
	}
	w.Header().Set("Content-Type","text/html")
	userViewing := store.GetUserInfoById(userid)
	followersJSON := store.GetUserFollowers(userid)
	followersSTR := "["
	for _, data := range followersJSON {
		followersSTR += data
	}
	followersSTR += "]"

	var result []Follower
	json.Unmarshal([]byte(followersSTR),&result)
	var fResult []FollowerResult
	for _, fData := range result {
		fUser := store.GetUserInfoById(fData.Userid)
		followedDate, _ := time.Parse("2006-01-02T15:04:05Z", fData.Followed)
		newResult := FollowerResult{FUsername: fUser.username, FFollowed: followedDate.Format("01/02/2006")}
		fResult = append(fResult, newResult)
	}

	tmpl, _ := template.ParseFiles(TEMPLATES + "/user/followers.html")
	tmpl.Execute(w, FollowersPageData{Activity: "Followers", Userid: userid, Username: userViewing.username, Followers: fResult})
}
func userProfileHandler(w http.ResponseWriter, r *http.Request) {
	params := strings.Split(r.URL.Path, "/")
	userid, err := strconv.Atoi(params[len(params) - 1])
	if err != nil {
		http.Redirect(w, r, "/live/search", http.StatusSeeOther)
		return
	}
	msg, err := decryptCookie(r, "username")
	if err != nil {
		http.Redirect(w, r, "/", http.StatusInternalServerError)
		return
	}
	account := store.GetUserInfo(&User{username: msg})
	//Use user id to get stuff.
	userViewing := store.GetUserInfoById(userid)
	//Now we have User so lets get user Settings
	settings := store.GetUserSettings(userViewing)
	var publicity string
	if settings.publicity {
		publicity = "Public"
	} else {
		publicity = "Private"
	}
	gender := "Female"
	if userViewing.gender {
		gender = "Male"
	}
	pageData := map[string]string{"userid": strconv.Itoa(account.id), "profileid": strconv.Itoa(userViewing.id),
		 "location": settings.location,
		 "bio": settings.bio,
		 "publicity":publicity,
		 "gender": gender,
		 "followers":strconv.Itoa(store.GetFollowersAmount(userViewing.id)),
		 "following":strconv.Itoa(store.GetFollowingAmount(userViewing.id)),
		 "username": account.username,
		 "viewingUsername": userViewing.username,
		 "age": strconv.Itoa(userViewing.age) }
	// TODO: MAKE & Parse Template
	tmpl, err := template.ParseFiles(TEMPLATES + "/user/profile.html")
	if err != nil {
		fmt.Fprint(w, "404: Page not found")
	}
	tmpl.Execute(w, pageData)
}
func faviconHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "/root/go/src/github.com/InsanityMatrix/SocialFoot/assets/images/favicon.ico")
}
func followUserHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		fmt.Fprint(w,err.Error())
		return
	}
	userid, _ := strconv.Atoi(r.Form.Get("userid"))
	profileid, _ := strconv.Atoi(r.Form.Get("profileid"))
	err = store.followUser(userid, profileid)
	if err != nil {
		fmt.Fprint(w,err.Error())
		return
	}
	fmt.Fprint(w, "Successfully followed this user!")
}
func isFollowingUserHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		fmt.Fprint(w,err.Error())
		return
	}
	userid, _ := strconv.Atoi(r.Form.Get("userid"))
	profileid, _ := strconv.Atoi(r.Form.Get("profileid"))
	isf := store.isUserFollowing(userid, profileid)
	if isf {
		fmt.Fprint(w, "1")
		return
	}
	fmt.Fprint(w, "0")
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
		return
	}
	defer file.Close()
	data, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Fprint(w, "Error")
		return
	}
	fmt.Fprint(w, string(data))
}
func userPostHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		fmt.Fprint(w, "Error")
		return
	}
	w.Header().Set("Content-Type","application/json")
	uid, err := strconv.Atoi(r.Form.Get("uid"))
	if err != nil {
		panic(err)
	}
	response := store.GetUsersPosts(uid)
	fmt.Fprint(w, response)
	return
}
type ViewPostData struct {
	Username string
	ProfileName string
	ProfileID int
	Caption string
	Type string
	Posted string
	Extension string
	Postid int
	Tags string
}
func viewPostHandler(w http.ResponseWriter, r *http.Request) {
	params := strings.Split(r.URL.Path, "/")
	identifier := params[len(params) - 1]
	params = strings.Split(identifier, ".")
	postid, _ := strconv.Atoi(params[len(params) - 1])
	profileid, _ := strconv.Atoi(params[len(params) - 2])
	//TODO: Finish View Post Handler
	profileUser := store.GetUserInfoById(profileid)

	username , err := decryptCookie(r, "username")
	if err != nil {
		http.Redirect(w, r, "/assets/login.html", http.StatusSeeOther)
		return
	}
	thisUser := store.GetUserInfo(&User{username: username})
	post := store.GetPostById(postid)
	tmpl, err := template.ParseFiles(TEMPLATES + "/feed/postPage.html")
	if err != nil {
		http.Redirect(w, r, "/live", http.StatusSeeOther)
		return
	}
	tags := post.Tags
	tagsRe := regexp.MustCompile("(#\\w+)")

	TAGS := ""
	matches := tagsRe.FindAllStringSubmatch(tags, -1)
	for _, group := range matches {
		TAGS += "<p class='postTag'>" + group[2] + "</p>"
	}
	data := ViewPostData{
		Username: thisUser.username,
		ProfileName: profileUser.username,
		ProfileID: post.Userid,
		Caption: post.Caption,
		Type: post.Type,
		Posted: post.Posted.Format("01/02/2006"),
		Extension: post.Extension,
		Postid: post.Postid,
		Tags: TAGS,
	}

	tmpl.Execute(w, data)
}

//{MESSAGES}
func sendTextMessageHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		fmt.Fprint(w, err.Error())
		return
	}

	uidFrom, _ := strconv.Atoi(r.Form.Get("uidFrom"))
	uidTo, _ := strconv.Atoi(r.Form.Get("uidTo"))
	message := r.Form.Get("message")
	err = store.SendMessage(uidFrom, uidTo, message)
	if err != nil {
		fmt.Fprint(w, err.Error())
		return
	}
	fmt.Fprint(w, "Success")
}
func loadMessages(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		name, err := decryptCookie(r, "username")
		if err != nil {
			http.Redirect(w,r,"/assets/login.html",http.StatusSeeOther)
			return
		}
		account := store.GetUserInfo(&User{username: name})
		conversations := store.GetConversations(account.id)
		for i, convo := range conversations {
			usr := store.GetUserInfoById(convo.ParticipantID)
			convo.ParticipantName = usr.username
			conversations[i] = convo
		}
		pageData := MessagePage{Username: name, Conversations: conversations}
		pageData.Home = true

		tmpl, _ := template.ParseFiles(TEMPLATES + "/messages/index.html", TEMPLATES + "/messages/sidebar.html")
		//TODO PARSE and combine sidebar
		tmpl.Execute(w, pageData)
		return
	}
	err := r.ParseForm()
	if err != nil {
		fmt.Fprint(w, err.Error())
	}
	name, err := decryptCookie(r, "username")
	if err != nil {
		http.Redirect(w,r,"/assets/login.html",http.StatusSeeOther)
		return
	}
	account := store.GetUserInfo(&User{username: name})
	conversations := store.GetConversations(account.id)

	pageData := MessagePage{Username: name, Conversations: conversations, Home: false}
	tmpl, _ := template.ParseFiles(TEMPLATES + "/messages/index.html", TEMPLATES + "/messages/sidebar.html")

	tmpl.Execute(w, pageData)
}
func createPrivateMessageHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		fmt.Fprint(w, "Why are you here?")
		return
	}

	_, err = decryptCookie(r, "username")
	if err != nil {
		fmt.Fprint(w, NotLoggedIn())
		return
	}

	userid, _ := strconv.Atoi(r.Form.Get("userid"))
	profileid, _ := strconv.Atoi(r.Form.Get("profileid"))
	exists := store.GetConversationID(userid, profileid)
	if exists == 0 {
		err = store.CreateTwoWayConversation(userid, profileid)
		if err != nil {
			fmt.Fprint(w, MsgCreationErr())
			return
		}
	}
	fmt.Fprint(w, "Conversation already exists")
}
func toMsgTemplateHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	file, err := os.Open(TEMPLATES + "/messages/toMsg.html")
	if err != nil {
		fmt.Fprint(w, "Error")
		return
	}
	defer file.Close()
	data, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Fprint(w, "Error")
		return
	}
	fmt.Fprint(w, string(data))
}
func fromMsgTemplateHandler(w http.ResponseWriter, r *http.Request) {

}
type ConversationPage struct {
	Username string
	Conversations []Conversation
	ConvoID int
	Userid int

}
func conversationHandler(w http.ResponseWriter, r *http.Request) {
	params := strings.Split(r.URL.Path, "/")
	convoID, _ := strconv.Atoi(params[len(params) - 1])

	//CHECK USER
	name, err := decryptCookie(r, "username")
	if err != nil {
		http.Redirect(w, r, "/assets/login.html", http.StatusSeeOther)
		return
	}
	user := store.GetUserInfo(&User{username: name})
	//CHECK IF USER IS IN THIS CONVERSATION
	if !store.IsUserInConversation(convoID, user.id) {
		http.Redirect(w, r, "/live/messages", http.StatusSeeOther)
		return
	}

	tmpl, err := template.ParseFiles(TEMPLATES + "/messages/conversation.html", TEMPLATES + "/messages/sidebar.html")
	if err != nil {
		//TODO: Send with error code
		http.Redirect(w, r, "/live/messages", http.StatusSeeOther)
		return
	}
	conversations := store.GetConversations(user.id)
	data := ConversationPage{Username: user.username, Conversations: conversations, ConvoID: convoID, Userid: user.id}

	tmpl.Execute(w, data)
}

func getMessages(w http.ResponseWriter, r *http.Request) {
	//Output Messages with convoid?
	err := r.ParseForm()
	if err != nil {
		//TODO: Better Error Handling System
		fmt.Fprint(w, "[ {} ]")
		return
	}

	convoID, _ := strconv.Atoi(r.Form.Get("convo"))
	conversation := store.GetConversation(convoID)

	data, err := json.Marshal(conversation)
	if err != nil {
		fmt.Fprint(w, "No Messages")
		return
	}
	fmt.Fprint(w, data)
}
