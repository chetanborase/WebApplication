package util

import (
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}
const (
	alphabetPool = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	numberPool   = "1234567890"
)

func GetRandomString(size int) *string {

	bts := make([]byte, size)
	for i, _ := range bts {
		bts[i] = alphabetPool[rand.Intn(len(alphabetPool))]
	}
	s := string(bts)
	return &s
}
func GetNumber(size int) string {
	bts := make([]byte, size)
	for i, _ := range bts {
		bts[i] = numberPool[rand.Intn(len(numberPool))]
	}
	return string(bts)
}
func GetStringPointer(str string) *string {
	return &str
}
