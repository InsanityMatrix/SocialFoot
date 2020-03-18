package main

import (
  "fmt"
  "net/http"
  "html/template"
)

func snakeGameHandler (w http.ResponseWriter, r *http.Request) {
  SetHeaders(w)
  w.Header().Set("Content-Type","text/html")
  name, err := decryptCookie(r)
  if err != nil {
    http.Redirect(w, r, "/assets/login.html", http.StatusSeeOther)
    return
  }
  tmpl, _ := template.ParseFiles(GAMES + "/snake.html")
  tmpl.Execute(w, map[string]string{"username": name})
}
