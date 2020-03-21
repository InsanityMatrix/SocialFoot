package main

import (
  "net/http"
  "html/template"
  "fmt"
  "encoding/json"
  "strings"
  "regexp"
  "strconv"
)
//All Handlers that start with live
func liveIndexHandler(w http.ResponseWriter, r *http.Request) {
	//Handle Live page with html templates
	SetHeaders(w)
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
func customFeedHandler(w http.ResponseWriter, r *http.Request) {
  SetHeaders(w)
  w.Header().Set("Content-Type", "text/html")
  username, err := decryptCookie(r, "username")
  if err != nil {
    http.Redirect(w, r, "/assets/login.html", http.StatusSeeOther)
    return
  }
  tmpl, _ := template.ParseFiles(TEMPLATES + "/feed/custom.html")
  tmpl.Execute(w, map[string]string{"username": username})
}
func getCustomFeedPosts(w http.ResponseWriter, r *http.Request) {
  SetHeaders(w)
  w.Header().Set("Content-Type", "application/json")
  username, err := decryptCookie(r, "username")
  if err != nil {
    http.Redirect(w, r, "/assets/login.html", http.StatusSeeOther)
    return
  }
  user := store.GetUserInfo(&User{username: username})
  //Get all of the posts from users they follow
  following, err := store.GetFollowing(user.id)
  if err != nil {
    //They are not following anybody.
    fmt.Fprint(w, "[]")
  }
  //Get Feed from their following
  feed := []Post{}
  for _, usr := range following {
    userfeed := store.GetPosts(usr.id)

    feed = append(feed, userfeed...)
  }
  //SORT in Descending by PostID
  for i := 0; i < len(feed) - 1; i++ {
    //Really inefficient but we can come back to this later
    thisPost := feed[i]
    nextPost := feed[i + 1]
    if thisPost.Postid < nextPost.Postid {
      feed[i] = nextPost
      feed[i + 1] = thisPost
      i = 0
    }
  }
  //Now we have all posts we need, and it is sorted in a nice order we output the JSON
  postJSON, _ := json.Marshal(feed)
  fmt.Fprint(w, string(postJSON))
}
func profileHandler(w http.ResponseWriter, r *http.Request) {
	//Handle Live Profile settings
	SetHeaders(w)
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
  user := store.GetUserInfo(&User{username: name})

	tmpl.Execute(w, map[string]string{"username":name, "userid":strconv.Itoa(user.id)})
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
	SetHeaders(w)
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
	SetHeaders(w)
	w.Header().Set("Content-Type","text/html")
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
	if post == nil {
		http.Redirect(w, r, "/live", http.StatusSeeOther)
		return
	}
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
		TAGS += "<p class='postTag'>" + group[0] + "</p>"
	}

	tHTML := template.HTML(TAGS)
	Image := false
	Video := false
	if post.Type == "IMAGE" {
		Image = true
	}
	if post.Type == "VIDEO" {
		Video = true
	}
	data := ViewPostData{
		Username: thisUser.username,
		ProfileName: profileUser.username,
		ProfileID: post.Userid,
		Caption: post.Caption,
		Image: Image,
		Video: Video,
		Posted: post.Posted.Format("01/02/2006"),
		Extension: post.Extension,
		Postid: post.Postid,
		Tags: tHTML,
	}
	SetHeaders(w)
	w.Header().Set("Content-Type","text/html")
	tmpl.Execute(w, data)
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
func postHandler(w http.ResponseWriter, r *http.Request) {
	SetHeaders(w)
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
func loadMessages(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		name, err := decryptCookie(r, "username")
		if err != nil {
			http.Redirect(w,r,"/assets/login.html",http.StatusSeeOther)
			return
		}
		SetHeaders(w)
		w.Header().Set("Content-Type","text/html")
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
	SetHeaders(w)
	w.Header().Set("Content-Type","text/html")
	account := store.GetUserInfo(&User{username: name})
	conversations := store.GetConversations(account.id)

	pageData := MessagePage{Username: name, Conversations: conversations, Home: false}
	tmpl, _ := template.ParseFiles(TEMPLATES + "/messages/index.html", TEMPLATES + "/messages/sidebar.html")

	tmpl.Execute(w, pageData)
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

	tmpl, err := template.ParseFiles(TEMPLATES + "/messages/conversation.html")
	if err != nil {
		//TODO: Send with error code
		http.Redirect(w, r, "/live/messages", http.StatusSeeOther)
		return
	}

	conversations := store.GetConversations(user.id)
	for i, convo := range conversations {
		usr := store.GetUserInfoById(convo.ParticipantID)
		convo.ParticipantName = usr.username
		conversations[i] = convo
	}
	convParticipant := store.GetConvoParticipant(convoID, user.id)
	data := ConversationPage{Username: user.username, Conversations: conversations, ConvoID: convoID, Userid: user.id, Participant: convParticipant}
	SetHeaders(w)
	w.Header().Set("Content-Type","text/html")
	tmpl.Execute(w, data)
}
type UserProfileData struct {
  Userid int
  Profileid int
  SamePerson bool
  DiffPerson bool
  Location string
  Bio string
  Publicity string
  Gender string
  Followers int
  Following int
  Username string
  ViewingUsername string
  Age int
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
	pageData := UserProfileData{
    Userid: account.id,
    Profileid:userViewing.id,
    SamePerson: account.id == userViewing.id,
    DiffPerson: account.id != userViewing.id,
		Location: settings.location,
		Bio: settings.bio,
	  Publicity:publicity,
		Gender: gender,
		Followers:store.GetFollowersAmount(userViewing.id),
		Following:store.GetFollowingAmount(userViewing.id),
		Username: account.username,
		ViewingUsername: userViewing.username,
		Age: userViewing.age,
   }
	// TODO: MAKE & Parse Template
	tmpl, err := template.ParseFiles(TEMPLATES + "/user/profile.html")
	if err != nil {
		fmt.Fprint(w, "404: Page not found")
	}
	tmpl.Execute(w, pageData)
}
