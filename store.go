package main
//Help from https://www.sohamkamani.com/blog/2017/10/18/golang-adding-database-to-web-application/
import (
    "database/sql"
    "strconv"
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
	_, err := store.db.Query("INSERT INTO users(username,gender,age,password,email) VALUES ($1,$2,$3,$4,$5);",user.username,user.gender,user.age,user.password,user.email)
  if err != nil {
    return err
  }
  row := store.db.QueryRow("SELECT id FROM users WHERE username=",user.username)
  err := row.Scan(user.id)
  if err != nil {
    return err
  }
  _, err := store.db.Query("INSERT INTO user_settings(userid) VALUES ($1)",user.id)
	return err
}
func (store *dbStore) LoginUser(user *User) (*User, error) {
  row := store.db.QueryRow("SELECT username,gender,age,password,email from users where username=$1", user.username)
  account := &User{}
  switch err := row.Scan(&account.username, &account.gender, &account.age, &account.password, &account.email); err {
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
  row := store.db.QueryRow("SELECT username,gender,age,password,email FROM users WHERE username=$1",user.username)
  account := &User{}
  err := row.Scan(&account.username, &account.gender, &account.age, &account.password, &account.email)
  if err != nil {
    return false
  }
  ps := []byte(user.password)
  return comparePasswords(account.password, ps)
}
//ONLY EVER USE ONCE YOU HAVE VALIDATED THEIR INFO FIRST
func (store *dbStore) GetUserInfo(user *User) *User {
  row := store.db.QueryRow("SELECT id,username,gender,age,password,email FROM users WHERE username=$1",user.username)
  account := &User{}
  err := row.Scan(&account.id,&account.username,&account.gender,&account.age,&account.password,&account.email)
  if err != nil {
    return nil
  }
  return account
}
func (store *dbStore) GetUserSettings(user *User) *UserSettings {
  row := store.db.QueryRow("SELECT userid, bio, website, location, publicity FROM user_settings WHERE userid")
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


var store dbStore

func InitStore(s dbStore) {
	store = s
}
