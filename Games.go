package main

import (
  "net/http"
  "html/template"
)
type SnakeScore struct {
  Username string
  Amount int
}
func GameHandler(w http.ResponseWriter, r *http.Request) {
  SetHeaders(w)
  w.Header().Set("Content-Type","text/html")
  name, err := decryptCookie(r, "username")
  if err != nil {
    http.Redirect(w, r, "/assets/login.html", http.StatusSeeOther)
    return
  }
  tmpl, _ := template.ParseFiles(GAMES + "/index.html")
  tmpl.Execute(w, map[string]string{"username": name})
}
func snakeGameHandler (w http.ResponseWriter, r *http.Request) {
  SetHeaders(w)
  w.Header().Set("Content-Type","text/html")
  name, err := decryptCookie(r, "username")
  if err != nil {
    http.Redirect(w, r, "/assets/login.html", http.StatusSeeOther)
    return
  }
  tmpl, _ := template.ParseFiles(GAMES + "/snake.html")
  tmpl.Execute(w, map[string]string{"username": name})
}
func snakeScoresHandler (w http.ResponseWriter, r *http.Request) {
  SetHeaders(w)
  w.Header().Set("Content-Type","application/json")
  //TODO: Set Up Score System

}
