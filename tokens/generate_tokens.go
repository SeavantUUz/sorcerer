package main

import (
	"fmt"
	"github.com/pote/philote-go"
)

const Secret = "123456"

func generateToken(read []string, write []string) string {
	token, _ := philote.NewToken(Secret, read, write)
	return token
}

func main() {
	tokenTest := generateToken([]string{"test-channel"}, []string{"test-channel"})
	fmt.Println(tokenTest)
	tokenContractRW := generateToken([]string{"contract-event"}, []string{"contract-event"})
	fmt.Println(tokenContractRW)
	tokenContractR := generateToken([]string{"contract-event"}, []string{})
	fmt.Println(tokenContractR)
}
