package main

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"math/rand"
	"strconv"
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

func hextoint(hexnumber string) int {
	d, _ := strconv.ParseInt("0x" + hexnumber , 0, 64)
	return int(d)
}

func upper_power_of_two(v int) int {
	v--;
	v |= v >> 1;
	v |= v >> 2;
	v |= v >> 4;
	v |= v >> 8;
	v |= v >> 16;
	v++;
	return v;

}


