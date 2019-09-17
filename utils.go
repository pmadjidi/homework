package main

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"math/rand"
)

func RandomString(n int) string {
	var letter = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	b := make([]rune, n)
	for i := range b {
		b[i] = letter[rand.Intn(len(letter))]
	}
	return string(b)
}

func PrettyPrint(v interface{}) (err error) {
	b, err := json.MarshalIndent(v, "", "  ")
	if err == nil {
		fmt.Println(string(b))
	}
	return
}

func calcHash(name string) string {
	hasher := sha1.New()
	hasher.Write([]byte(name))
	return  fmt.Sprintf("%x", hasher.Sum(nil))
}
