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
		"github.com/koyachi/go-nude"
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

type ViewPostData struct {
	Username string
	ProfileName string
	ProfileID int
	Caption string
	Image bool
	Video bool
	Posted string
	Extension string
	Postid int
	Tags template.HTML
}
type ConversationPage struct {
	Username string
	Conversations []Conversation
	ConvoID int
	Userid int
	Participant int
}

var HOME string
var TEMPLATES string
var GAMES string
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
		r.HandleFunc("/live/feed", customFeedHandler)
		r.HandleFunc("/live", liveIndexHandler)

		//Games
		r.HandleFunc("/games/snake", snakeGameHandler)
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
		r.HandleFunc("/search", searchHandler).Methods("POST")

		//MESSAGES FUNCTIONS
		r.HandleFunc("/messages/send/text", sendTextMessageHandler)
		r.HandleFunc("/messages/create/private", createPrivateMessageHandler)

		//JSON stuff
		r.HandleFunc("/json/user/id", HandleJSONUserById)
		r.HandleFunc("/json/messages/convo", getMessages)
		r.HandleFunc("/json/feed/custom", getCustomFeedPosts)
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
		GAMES = "/root/go/src/github.com/InsanityMatrix/SocialFoot/games"
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
		db.SetMaxOpenConns(25)
		db.SetMaxIdleConns(10)
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

func getPublicPostsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type","application/json")
	pubposts := store.GetPublicPosts()
  fmt.Fprint(w, pubposts)
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
	SetHeaders(w)
	w.Header().Set("Content-Type","text/html")
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
	caption := strings.Replace(strings.Replace(r.Form.Get("caption"),"<","&lt;",-1),">","&gt;",-1)

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
			SetHeaders(w)
			w.Header().Set("Content-Type","text/html")
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
			status := "This video is too large."
			tmpl.Execute(w, map[string]string{"username":username,"status":status})
			return
		}
	}
	//Display results
	isNude, err := nude.IsNude("/root/go/src/github.com/InsanityMatrix/SocialFoot/assets/uploads/imageposts/post" + idStr + extension)
	if err != nil {
		tmpl, err := template.ParseFiles(TEMPLATES + "/uploadSuccess.html")
		msg, err := decryptCookie(r, "username")
		var username string
		if err != nil {
			username = store.GetUserInfoById(id).username
		} else {
			username = msg
		}
		SetHeaders(w)
		w.Header().Set("Content-Type","text/html")
		status := "Your post has been created at http://www.socialfoot.me/live/view/post/" + userid + "." + idStr
		tmpl.Execute(w, map[string]string{"username":username,"status":status})
	}
	if isNude {
		store.DeleteUserPost(postid)
		os.Remove("/root/go/src/github.com/InsanityMatrix/SocialFoot/assets/uploads/imageposts/post" + idStr + extension)
		tmpl, err := template.ParseFiles(TEMPLATES + "/uploadSuccess.html")
		msg, err := decryptCookie(r, "username")
		var username string
		if err != nil {
			id, _ := strconv.Atoi(userid)
			username = store.GetUserInfoById(id).username
		} else {
			username = msg
		}
		SetHeaders(w)
		w.Header().Set("Content-Type","text/html")
		status := "Nudity was detected, if this was an error, please REPORT it."
		tmpl.Execute(w, map[string]string{"username":username,"status":status})
		return
	}
	tmpl, err := template.ParseFiles(TEMPLATES + "/uploadSuccess.html")
	msg, err := decryptCookie(r, "username")
	var username string
	if err != nil {
		username = store.GetUserInfoById(id).username
	} else {
		username = msg
	}
	SetHeaders(w)
	w.Header().Set("Content-Type","text/html")
	status := "Your post has been created at http://www.socialfoot.me/live/view/post/" + userid + "." + idStr
	tmpl.Execute(w, map[string]string{"username":username,"status":status})
}

func HandleJSONUserById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type","application/json")
	_ = r.ParseForm()
	id, _ := strconv.Atoi(r.Form.Get("userid"))
	userjson := store.GetJSONUserByID(id)
	fmt.Fprint(w, userjson)
}
func searchHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		w.WriteHeader(404)
		return
	}
	w.Header().Set("Content-Type","application/json")
	term := r.Form.Get("term")
	response := store.GetJSONUsersByUsernames(term)
	// TODO: Add posts to the search feature.
	// posts := store.GetPostsByTag(term)

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
	SetHeaders(w)
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
