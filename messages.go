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
