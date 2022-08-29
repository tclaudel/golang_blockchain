package main

import (
	"fmt"
	"log"
)

type Block struct {
	nonce        int
	previousHash int64
	transaction  []string
}

func init() {
	log.SetPrefix("Blochcahin: ")
}

func main() {
	log.Println("test")
	fmt.Println("test")
}
