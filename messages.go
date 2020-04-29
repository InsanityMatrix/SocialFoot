package main

import (
  "crypto/aes"
  "crypto/cipher"
  "crypto/md5"
  "crypto/rand"
  "encoding/hex"
  "os"
  "fmt"
  "io"
  "io/ioutil"
  "net/http"
  "strconv"
  "encoding/json"
)
//TODO: Better error handling

func getMessageHash() string {
  hasher := md5.New()
  key := os.Getenv("MESSAGEKEY")
  hasher.Write([]byte(key))
  return hex.EncodeToString(hasher.Sum(nil))
}

func encryptMessage(data []byte) []byte {
  block, _ := aes.NewCipher([]byte(getMessageHash()))
  gcm, err := cipher.NewGCM(block)
  if err != nil {
    fmt.Println(err.Error())
  }
  nonce := make([]byte, gcm.NonceSize())
  if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
    fmt.Println(err.Error())
  }
  ciphertext := gcm.Seal(nonce, nonce, data, nil)
  return ciphertext
}

func decryptMessage(data []byte) []byte {
  key := []byte(getMessageHash())
  block, err := aes.NewCipher(key)
  if err != nil {
    fmt.Println(err.Error())
  }
  gcm, err := cipher.NewGCM(block)
  if err != nil {
    fmt.Println(err.Error())
  }
  nonceSize := gcm.NonceSize()
  nonce, ciphertext := data[:nonceSize], data[nonceSize:]
  plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
  if err != nil {
    fmt.Println(err.Error())
  }
  return plaintext
}

func encryptMessageFile(filename string, data []byte) {
  f, _ := os.Create("/root/go/src/github.com/InsanityMatrix/SocialFoot/messages/" + filename)
  defer f.Close()
  f.Write(encryptMessage(data))
}
func decryptMessageFile(filename string) []byte {
  data, _ := ioutil.ReadFile("/root/go/src/github.com/InsanityMatrix/SocialFoot/messages/" + filename)
  return decryptMessage(data)
}

//HANDLERS
func getMessages(w http.ResponseWriter, r *http.Request) {
	//Output Messages with convoid?
	err := r.ParseForm()
	if err != nil {
		//TODO: Better Error Handling System
		fmt.Fprint(w, "[ {} ]")
		return
	}

	convoID, _ := strconv.Atoi(r.Form.Get("convo"))
	conversation := store.GetConversation(convoID)

	for index, message := range conversation {
		content := decryptMessageFile(strconv.Itoa(convoID) + "/" + strconv.Itoa(message.MessageID) + ".txt")
		conversation[index].Content = string(content);
	}
	data, err := json.Marshal(conversation)
	if err != nil {
		fmt.Fprint(w, "No Messages")
		return
	}
	w.Header().Set("Content-Type","application/json")
	fmt.Fprint(w, string(data))
}
func fromMsgTemplateHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	file, err := os.Open(TEMPLATES + "/messages/fromMsg.html")
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
func toMsgTemplateHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	file, err := os.Open(TEMPLATES + "/messages/toMsg.html")
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
func createPrivateMessageHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		fmt.Fprint(w, "Why are you here?")
		return
	}

	_, err = decryptCookie(r, "username")
	if err != nil {
		fmt.Fprint(w, NotLoggedIn())
		return
	}

	userid, _ := strconv.Atoi(r.Form.Get("userid"))
	profileid, _ := strconv.Atoi(r.Form.Get("profileid"))
  if userid == profileid {
    fmt.Fprint(w, MsgCreationErr())
    return
  }
	exists := store.GetConversationID(userid, profileid)
	if exists == 0 {
		err = store.CreateTwoWayConversation(userid, profileid)
		if err != nil {
			fmt.Fprint(w, MsgCreationErr())
			return
		}
		fmt.Fprint(w, store.GetConversationID(userid, profileid))
	}
	fmt.Fprint(w, exists)
}
func sendTextMessageHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		fmt.Fprint(w, err.Error())
		return
	}

	uidFrom, _ := strconv.Atoi(r.Form.Get("uidFrom"))
	uidTo, _ := strconv.Atoi(r.Form.Get("uidTo"))
	message := r.Form.Get("message")
	err = store.SendMessage(uidFrom, uidTo, message)
	if err != nil {
		fmt.Fprint(w, err.Error())
		return
	}
	fmt.Fprint(w, "Success")
}
func getConversationsHandler(w http.ResponseWriter, r *http.Request) {
  err := r.ParseForm()
  if err != nil {
    fmt.Fprint(w, err.Error())
    return
  }
  	accID, _ := strconv.Atoi(r.Form.Get("uid"))
  	conversations := store.GetConversations(accID)
    for index, convo := range conversations {
      account := store.GetUserInfoById(convo.ParticipantID)
      convo.ParticipantName = account.username
      conversations[index] = convo
    }
    jsonInfo, _ := json.Marshal(conversations)
    fmt.Fprint(w,string(jsonInfo))
}
