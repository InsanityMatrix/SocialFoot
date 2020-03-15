package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"os"
	"io"
	"fmt"
	"net/http"
)

func getCookieHash() string {
	hasher := md5.New()
	key := os.Getenv("COOKIE_KEY")
	hasher.Write([]byte(key))
	return hex.EncodeToString(hasher.Sum(nil))
  }

func setEncryptedCookie(w http.ResponseWriter, name string, data []byte) {
	block, _ := aes.NewCipher([]byte(getCookieHash()))
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		fmt.Println(err.Error())
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		fmt.Println(err.Error())
	}
	ciphertext := gcm.Seal(nonce, nonce, data, nil)
	cookie := http.Cookie{
		Name: name,
		Value: hex.EncodeToString(ciphertext),
		Path: "/",
		MaxAge: 86400,
	}
	http.SetCookie(w, &cookie)
}

func decryptCookie(r *http.Request, name string) (string, error) {
	key := []byte(getCookieHash())
	block, err := aes.NewCipher(key)
	if err != nil {
		fmt.Println(err.Error())
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		fmt.Println(err.Error())
	}
	msg, err := r.Cookie(name)
	if err != nil {
		return  "", err
	}
	data, err := hex.DecodeString(msg.Value)
	if err != nil {
		return "", err
	}
	nonceSize := gcm.NonceSize()
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		fmt.Println(err.Error())
	}
	return string(plaintext), nil
}
