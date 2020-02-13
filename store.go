package main
//Help from https://www.sohamkamani.com/blog/2017/10/18/golang-adding-database-to-web-application/
import (
    "database/sql"
    "strconv"
    "time"
)
//Users DB : username (string), gender (bool), age (int), password (string), email (string)

type Store interface {
    CreateUser(user *User) error
    GetUserInfo(user *User) *User
    GetUserSettings(user *User) *UserSettings
    CheckUserCredentials(user *User) bool
    LoginUser(user *User) (*User, error)
    UpdateSetting(user *User,setting string, value string) bool
}


type dbStore struct {
    db *sql.DB
}

func (store *dbStore) CreateUser(user *User) error {
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
    panic(err)
  }
	return err
}
func (store *dbStore) LoginUser(user *User) (*User, error) {
  row := store.db.QueryRow("SELECT id,username,gender,age,password,email from users where username=$1", user.username)
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

var store dbStore

func InitStore(s dbStore) {
	store = s
}
