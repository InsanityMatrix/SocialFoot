package main

import (
  "net/http"
  "html/template"
  "strconv"
  "encoding/json"
  "fmt"
)
type SnakeScore struct {
  Scoreid int
  Userid int
  Score int
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

//SNAKE
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
type SnakeScoreList struct {
  Scoreid int
  Username string
  Score int
}
func snakeScoresHandler (w http.ResponseWriter, r *http.Request) {
  SetHeaders(w)
  w.Header().Set("Content-Type","application/json")
  scores := store.GetTopSnakeScores()
  scoreData := []SnakeScoreList{}
  for _, score := range scores {
    user := store.GetUserInfoById(score.Userid)
    listData := SnakeScoreList{Scoreid: score.Scoreid, Username: user.username, Score: score.Score}

    scoreData = append(scoreData, listData)
  }

  data, _ := json.Marshal(scoreData)
  fmt.Fprint(w, string(data))
}
func updateSnakeScore (w http.ResponseWriter, r *http.Request) {
  err := r.ParseForm()
  if err != nil {
    fmt.Fprint(w, "You shouldn't be here")
  }
  userid,_ := strconv.Atoi(r.Form.Get("userid"))
  score,_ := strconv.Atoi(r.Form.Get("score"))

  store.UpdateSnakeScore(userid, score)
}
