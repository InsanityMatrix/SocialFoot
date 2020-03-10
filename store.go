package main
//Help from https://www.sohamkamani.com/blog/2017/10/18/golang-adding-database-to-web-application/
import (
    "database/sql"
    "strconv"
    "errors"
    "time"
    "os"
    "github.com/bdwilliams/go-jsonify/jsonify"
)
//Users DB : username (string), gender (bool), age (int), password (string), email (string)

type Store interface {
    CreateUser(user *User) error
    PostUserImage(publicity bool, caption string, tags string, userid int, extension string) int
    GetPublicPosts() ([]*ImagePost, error)
    LoginUser(user *User) (*User, error)
    GetUserInfo(user *User) *User
    GetUserSettings(user *User) *UserSettings
    CheckUserCredentials(user *User) bool
    UpdateSetting(user *User,setting string, value string) bool
    GetUserInfoById(id int) *User
    SetUserPublicity(userID int, mode bool) bool
}


type dbStore struct {
    db *sql.DB
}

type ImagePost struct {
	postid int `json:"postid"`
	userid int `json:"userid"`
	public bool `json:"publicity"`
	tags string `json:"tags"`
	caption string `json:"caption"`
	extension string `json:"extension"`
}

func (store *dbStore) CreateUser(user *User) error {
  var exists bool
  row := store.db.QueryRow("SELECT exists(SELECT * FROM users WHERE LOWER(username)=LOWER($1))", user.username)
  err := row.Scan(&exists)
  if err != nil {
    return err
  }
  if exists {
    return errors.New("User Exists")
  }
	rows, err := store.db.Query("INSERT INTO users(username,gender,age,password,email) VALUES ($1,$2,$3,$4,$5) RETURNING id;",user.username,user.gender,user.age,user.password,user.email)
  if err != nil {
    return err
  }
  defer rows.Close()
  account := User{}
  if rows.Next() {
    err = rows.Scan(&account.id)
  }
  if err != nil {
    return err
  }
  _, err = store.db.Query("INSERT INTO user_settings(userid,publicity) VALUES ($1,$2)",account.id,true)
  if err != nil {
    return err
  }
  _, err = store.db.Query("CREATE TABLE user" + strconv.Itoa(account.id) + "_following(relID serial,userid int UNIQUE, followed date, PRIMARY KEY(relID) );")
  if err != nil {
    return err
  }
  _, err = store.db.Query("CREATE TABLE user" + strconv.Itoa(account.id) + "_followers(relID serial,userid int UNIQUE, followed date, PRIMARY KEY(relID) );")
  if err != nil {
    return err
  }
	return err
}
func (store *dbStore) followUser(follower int, followed int) error {
  var exists bool
  err := store.db.QueryRow("SELECT exists(SELECT * FROM user" + strconv.Itoa(follower) + "_following WHERE userid=$1)", follower).Scan(&exists)
  if err != nil {
    return err
  }
  if exists {
    return nil
  }
  dt := time.Now()
  _, err = store.db.Query("INSERT INTO user" + strconv.Itoa(followed) + "_followers(userid, followed) VALUES ($1,$2)",follower, dt)
  if err != nil {
    return err
  }
  _, err = store.db.Query("INSERT INTO user" + strconv.Itoa(follower) + "_following(userid, followed) VALUES ($1, $2)",followed, dt)
  if err != nil {
    return err
  }
  return nil
}
func (store *dbStore) isUserFollowing(follower int, followed int) bool {
  var exists bool
  err := store.db.QueryRow("SELECT exists (SELECT * FROM user" + strconv.Itoa(follower) + "_following WHERE userid=$1)", followed).Scan(&exists)
  if err != nil {
     return false
  }
  return exists
}
func (store *dbStore) PostUserImage(publicity bool, caption string, tags string, userid int, extension string, t string) int {
  dt := time.Now()
  rows, err := store.db.Query("INSERT INTO posts(userid,publicity,tags,caption,type,posted,extension) VALUES ($1,$2,$3,$4,$5,$6,$7) RETURNING postid",userid,publicity,tags,caption,t,dt,extension)
  if err != nil {
    return 0
  }
  defer rows.Close()
  var postid int
  if rows.Next() {
    err = rows.Scan(&postid)
    if err != nil {
      return 0
    }
    return postid
  }
  return 0
}
func (store *dbStore) DeleteUserPost(postid int) {
  store.db.Query("DELETE FROM posts WHERE postid=$1",postid)
}

func (store *dbStore) LoginUser(user *User) (*User, error) {
  row := store.db.QueryRow("SELECT id,username,gender,age,password,email from users where LOWER(username)=LOWER($1)", user.username)
  account := &User{}
  switch err := row.Scan(&account.id, &account.username, &account.gender, &account.age, &account.password, &account.email); err {
  case sql.ErrNoRows:
    return nil, err
  case nil:
    return account, nil
  default:
    panic(err)
    return nil, err
  }
}
func (store *dbStore) UpdateSetting(user *User,setting string, value string) bool {
  if(setting == "publicity") {
    val, _ := strconv.ParseBool("value")
    _, err := store.db.Query("UPDATE TABLE user_settings SET publicity=$1 WHERE userid=$2",val,user.id)
    if err != nil {
      return false
    }
    return true
  }

  return false
}
func (store *dbStore) CheckUserCredentials(user *User) bool {
  row := store.db.QueryRow("SELECT password FROM users WHERE username=$1",user.username)
  account := User{}
  err := row.Scan(&account.password)
  if err != nil {
    return false
  }
  ps := []byte(user.password)
  return comparePasswords(account.password, ps)
}
//ONLY EVER USE ONCE YOU HAVE VALIDATED THEIR INFO FIRST
func (store *dbStore) GetUserInfo(user *User) *User {
  row := store.db.QueryRow("SELECT * FROM users WHERE username=$1",user.username)
  account := &User{}
  err := row.Scan(&account.id,&account.username,&account.gender,&account.age,&account.password,&account.email)
  if err != nil {
    return nil
  }
  return account
}
func (store *dbStore) GetUserInfoById(id int) *User {
  row := store.db.QueryRow("SELECT * FROM users WHERE id=$1",id)
  account := &User{}
  err := row.Scan(&account.id,&account.username,&account.gender,&account.age,&account.password,&account.email)
  if err != nil {
    return nil
  }
  return account
}
func (store *dbStore) GetUserSettings(user *User) *UserSettings {
  row := store.db.QueryRow("SELECT bio,location,publicity FROM user_settings WHERE userid=$1",user.id)
  settings := &UserSettings{}
  err := row.Scan(&settings.bio,&settings.location,&settings.publicity)

  if err != nil {
    return nil
  }
  return settings
}
func (store *dbStore) SetUserPublicity(userID int, mode bool) bool {
  _, err := store.db.Query("UPDATE user_settings SET publicity=$1 WHERE userid=$2",mode,userID)
  if err != nil {
    return false
  }
  return true
}
func (store *dbStore) ChangeUserEmail(userID int, email string) bool {
  _, err := store.db.Query("UPDATE users SET email=$1 WHERE id=$2",email, userID)
  if err != nil {
    return false
  }
  return true
}
func (store *dbStore) ChangeUserLocation(userID int, location string) bool {
  _, err := store.db.Query("UPDATE user_settings SET location=$1 WHERE userid=$2",location, userID)
  if err != nil {
    return false
  }
  return true
}
func (store *dbStore) ChangeUserBio(userID int, bio string) bool {
  _, err := store.db.Query("UPDATE user_settings SET bio=$1 WHERE userid=$2",bio,userID)
  if err != nil {
    return false
  }
  return true
}
func (store *dbStore) DeleteAccount(userID int) bool {
  _, err := store.db.Query("DELETE FROM users WHERE id=$1",userID)
  if err != nil {
    return false
  }
  _, err = store.db.Query("DELETE FROM user_settings WHERE userid=$1",userID)
  if err != nil {
    return false
  }
  _, err = store.db.Query("DELETE FROM posts WHERE userid=$1",userID)
  if err != nil {
    return false
  }
  return true
}
func (store *dbStore) GetUsers() ([]*User, error) {
	rows, err := store.db.Query("SELECT username,gender,age,password,email from users")

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := []*User{}

	for rows.Next() {
		user := &User{}

		if err := rows.Scan(&user.username, &user.gender, &user.age, &user.password, &user.email); err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	return users, nil
}
func (store *dbStore) SubmitBugReport(username string, content string) {
  dt := time.Now()
  _, err := store.db.Query("INSERT INTO bug_reports(username, content, submitted) VALUES ($1, $2, $3)",username,content,dt)
  if err != nil {
    return
  }
}



//JSON FUNCTIONS
func (store *dbStore) GetUserFollowing(userid int) []string {
  rows, err := store.db.Query("SELECT * FROM user" + strconv.Itoa(userid) + "_following ORDER BY followed DESC")
  if err != nil {
    return []string{"Failed","Failed"}
  }
  return jsonify.Jsonify(rows)
}
func (store *dbStore) GetUserFollowers(userid int) []string {
  rows, err := store.db.Query("SELECT * FROM user" + strconv.Itoa(userid) + "_followers ORDER BY followed DESC")
  if err != nil {
    return []string{"Failed","Failed"}
  }
  return jsonify.Jsonify(rows)
}
func (store *dbStore) GetFollowersAmount(userid int) int {
  rows, err := store.db.Query("SELECT * FROM user" + strconv.Itoa(userid) + "_followers;")
  if err != nil {
    return 0
  }
  count := 0
  for rows.Next() {
    count++
  }
  return count
}
func (store *dbStore) GetFollowingAmount(userid int) int {
  rows, err := store.db.Query("SELECT * FROM user" + strconv.Itoa(userid) + "_following;")
  if err != nil {
    return 0
  }
  count := 0
  for rows.Next() {
    count++
  }
  return count
}
func (store *dbStore) GetPublicPosts() []string {
  rows, err := store.db.Query("SELECT * FROM posts WHERE publicity=$1 ORDER BY postid DESC",true)

	if err != nil {
    var error []string
    error[0] = "{\"status\":\"error\"}"
		return error
	}
	defer rows.Close()

  return jsonify.Jsonify(rows)
}
func (store *dbStore) GetUsersPosts(userid int) []string {
  rows, err := store.db.Query("SELECT * FROM posts WHERE userid=$1 ORDER BY postid DESC", userid)
  if err != nil {
    var error []string
    error[0] = "{\"status\":\"error\"}"
		return error
  }
  defer rows.Close()

  return jsonify.Jsonify(rows)
}
func (store *dbStore) GetPostById(postid int) *Post {
  row := store.db.QueryRow("SELECT * FROM posts WHERE postid=$1",postid)
  postData := &Post{}
  err := row.Scan(&postData.Postid,&postData.Userid,&postData.Tags,&postData.Caption,&postData.Type,&postData.Posted,&postData.Extension,&postData.Publicity,&postData.Likes)
  if err != nil {
    return nil
  }
  return postData
}
func (store *dbStore) GetJSONUserByID(uid int) []string {
  rows, err := store.db.Query("SELECT id,username FROM users WHERE id=$1", uid)
  if err != nil {
    var error []string
    error[0] = "{\"status\":\"error\"}"
		return error
  }
  defer rows.Close()
  return jsonify.Jsonify(rows)
}
func (store *dbStore) GetJSONUsersByUsernames(username string) []string {
  newName := "%" + username + "%"
  rows, err := store.db.Query("SELECT id,username,gender FROM users WHERE username LIKE $1", newName)
  if err != nil {
    var error []string
    error[0] = "{\"status\":\"error\"}"
		return error
  }
  defer rows.Close()
  return jsonify.Jsonify(rows)
}
//{MESSAGES}
func (store *dbStore) CreateTwoWayConversation(user1 int, user2 int) error {
  dt := time.Now()
  rows, err := store.db.Query("INSERT INTO private_conversations(userOne, userTwo, created) VALUES ($1, $2, $3) RETURNING convoID;",user1, user2, dt)
  if err != nil {
    return err
  }
  defer rows.Close()

	var convoID int

	for rows.Next() {

		if err := rows.Scan(&convoID); err != nil {
			return err
		}
	}

  os.Mkdir("/root/go/src/github.com/InsanityMatrix/SocialFoot/messages/" + strconv.Itoa(convoID), 0755)
	_, err = store.db.Query("CREATE TABLE c" + strconv.Itoa(convoID) + "_pconv (messageid SERIAL, mfrom int, read BOOLEAN, PRIMARY KEY(messageid));")

  return err
}
//returns 0 if error, convoID will never equal 0
func (store *dbStore) GetConversationID(user1 int, user2 int) int {
  row := store.db.QueryRow("SELECT convoID FROM private_conversations WHERE (userOne=$1 AND userTwo=$2) OR (userOne=$2 AND userTwo=$1)",user1, user2)

  var convoID int

  err := row.Scan(&convoID)
  if err != nil {
    return 0;
  }
  return convoID
}
func (store *dbStore) SendMessage(uidFrom int, uidTo int, message string) error {
  convoID := store.GetConversationID(uidFrom, uidTo)
  if convoID == 0 {
    return errors.New("Conversation doesn't exist, can't send message!")
  }
  row := store.db.QueryRow("INSERT INTO c" + strconv.Itoa(convoID) + "_pconv (mfrom,read) VALUES ($1,$2) RETURNING messageid",uidFrom, true)
  var messageID int
  row.Scan(&messageID)
  encryptMessageFile(strconv.Itoa(convoID) + "/" + strconv.Itoa(messageID) + ".txt", []byte(message))
  return nil
}
func (store *dbStore) GetConversations(uid int) []Conversation {
  rows, err := store.db.Query("SELECT convoid,usertwo,created FROM private_conversations WHERE userOne=$1", uid)
  if err != nil {
    return nil
  }
  defer rows.Close()

  conversations := []Conversation{}
  for rows.Next() {
    conversation := Conversation{}
    rows.Scan(&conversation.ConvoID, &conversation.ParticipantID,&conversation.Created)

    conversations = append(conversations, conversation)
  }
  rows, err = store.db.Query("SELECT convoid,userone,created FROM private_conversations WHERE usertwo=$1", uid)
  if err != nil {
    return nil
  }
  defer rows.Close()
  for rows.Next() {
    conversation := Conversation{}
    rows.Scan(&conversation.ConvoID, &conversation.ParticipantID,&conversation.Created)

    conversations = append(conversations, conversation)
  }
  return conversations
}
func (store *dbStore) GetConversation(convoid int) []*Message {
  rows, _ := store.db.Query("SELECT * FROM c + strconv.Itoa(convoid) + "_pconv")
  defer rows.Close()
  messages := []*Message{}
  for rows.Next() {
    message := &Message{}
    rows.Scan(&message.MessageID, &message.From, &message.Read)

    messages = append(messages, message)
  }
  return messages
}
func (store *dbStore) IsUserInConversation(convoid int, userid int) bool {
  var exists bool
  row := store.db.QueryRow("SELECT exists (SELECT * FROM private_conversations WHERE convoid=$1 AND (userone=$2 OR usertwo=$2))", convoid, userid)
  err := row.Scan(&exists)
  if err != nil {
    return false
  }
  return exists
}
//ESSENTIALS:

var store dbStore

func InitStore(s dbStore) {
	store = s
}
