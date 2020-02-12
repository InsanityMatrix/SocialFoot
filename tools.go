package main

import (
	"log"
	"net/smtp"
	"golang.org/x/crypto/bcrypt"
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
			from, []string{recipient}, []byte(content))
	if err != nil {
		log.Printf("smtp error: %s", err)
		return
	}
	log.Print("Sent")
}
