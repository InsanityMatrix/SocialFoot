package main

import (
  "net/http"
  "html/template"
  "strconv"
  "encoding/json"
  "fmt"
)
type GameScore struct {
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
  user := store.GetUserInfo(&User{username: name})
  tmpl.Execute(w, map[string]string{"username": name, "userid":strconv.Itoa(user.id)})
}
type ScoreList struct {
  Scoreid int
  Username string
  Score int
}
func snakeScoresHandler (w http.ResponseWriter, r *http.Request) {
  SetHeaders(w)
  w.Header().Set("Content-Type","application/json")
  scores := store.GetTopSnakeScores()
  scoreData := []ScoreList{}
  for _, score := range scores {
    user := store.GetUserInfoById(score.Userid)
    listData := ScoreList{Scoreid: score.Scoreid, Username: user.username, Score: score.Score}

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

//2048
func Handler2048(w http.ResponseWriter, r *http.Request) {
  SetHeaders(w)
  w.Header().Set("Content-Type","text/html")
  name, err := decryptCookie(r, "username")
  if err != nil {
    http.Redirect(w, r, "/assets/login.html", http.StatusSeeOther)
    return
  }
  tmpl, _ := template.ParseFiles(GAMES + "/2048.html")
  user := store.GetUserInfo(&User{username: name})
  tmpl.Execute(w, map[string]string{"username": name, "userid":strconv.Itoa(user.id)})
}
func Handler2048Scores(w http.ResponseWriter, r *http.Request) {
  SetHeaders(w)
  w.Header().Set("Content-Type","application/json")
  scores := store.GetTop2048Scores()
  scoreData := []ScoreList{}
  for _, score := range scores {
    user := store.GetUserInfoById(score.Userid)
    listData := ScoreList{Scoreid: score.Scoreid, Username: user.username, Score: score.Score}

    scoreData = append(scoreData, listData)
  }

  data, _ := json.Marshal(scoreData)
  fmt.Fprint(w, string(data))
}
func update2048Score (w http.ResponseWriter, r *http.Request) {
  err := r.ParseForm()
  if err != nil {
    fmt.Fprint(w, "You shouldn't be here")
  }
  userid,_ := strconv.Atoi(r.Form.Get("userid"))
  score,_ := strconv.Atoi(r.Form.Get("score"))

  store.Update2048Score(userid, score)
}

//Galaga
func GalagaHandler(w http.ResponseWriter, r *http.Request) {
  SetHeaders(w)
  w.Header().Set("Content-Type","text/html")
  name, err := decryptCookie(r, "username")
  if err != nil {
    http.Redirect(w, r, "/assets/login.html", http.StatusSeeOther)
    return
  }
  tmpl, _ := template.ParseFiles(GAMES + "/galaga.html")
  user := store.GetUserInfo(&User{username: name})
  tmpl.Execute(w, map[string]string{"username": name, "userid":strconv.Itoa(user.id)})
}
