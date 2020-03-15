package main

import (
	"log"
	"net/smtp"
	"golang.org/x/crypto/bcrypt"
	"strings"
	"net/http"
)

func hashAndSalt(pwd []byte) string {
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
    if err != nil {
        log.Println(err)
    }

	return string(hash)
}
func comparePasswords(hashedPwd string, plainPwd []byte) bool {
	byteHash := []byte(hashedPwd)

	err := bcrypt.CompareHashAndPassword(byteHash, plainPwd)
	if err != nil {
		return false
	}

	return true
}

func sendAuthMail(recipient string, content string) {
	from := "SocialFoot.noreply@gmail.com"
	password := "password"
	//Need to set OS Environment variable with mail password

	msg := "From: " + from + "\n" +
	       "To: " + recipient + "\n" +
	       "Subject: Email Authentication\n\n" +
	       content
	err := smtp.SendMail("smtp.gmail.com:587",
			smtp.PlainAuth("", from, password, "smtp.gmail.com"),
			from, []string{recipient}, []byte(msg))
	if err != nil {
		log.Printf("smtp error: %s", err)
		return
	}
	log.Print("Sent")
}
func isSupportedFile(extension string) (string, bool) {
	pictureExtensions := []string{".jpg", ".jpeg", ".jpe", ".jif",".jfif",".jfi",".png",".tiff",".tif",".raw",".arw",".cr2",".bmp",".webp", ".mp4", ".mov"}
	for _, ext := range pictureExtensions {
		if extension == ext {
			if extension == ".mp4" {
				return "VIDEO", true
			}
			if extension == ".mov" {
				return "VIDEO", true
			}
			return "IMAGE",true
		}
	}
	return "",false
}
func badReport(content string) bool {
	blackList := []string{
		"gay","gei","gae","gey",
	}
	spam := []string{"1", "gay", "bad", "Hi", "Hello", "Hey"}

	for _, blacklisted := range blackList {
	  blacklisted = " " + blacklisted + " "
		if strings.Contains(content, blacklisted) {
			return true
		}
	}
	for _, spamMessage := range spam {
		if content == spamMessage {
			return true
		}
	}
	return false
}
func SetHeaders(w http.ResponseWriter) {
	w.Header().Set("X-Frame-Options", "DENY")
}
