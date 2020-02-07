package main
//Help from https://www.sohamkamani.com/blog/2017/10/18/golang-adding-database-to-web-application/
import (
    "database/sql"
)
//Users DB : username (string), gender (bool), age (int), password (string), email (string)

type Store interface {
    CreateUser(user *User) error
    GetUsers()([]*User, error)
    LoginUser(user *User) (*User, error)
}


type dbStore struct {
    db *sql.DB
}

func (store *dbStore) CreateUser(user *User) error {
	_, err := store.db.Query("INSERT INTO users(username,gender,age,password,email) VALUES ($1,$2,$3,$4,$5);",user.username,user.gender,user.age,user.password,user.email)
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
