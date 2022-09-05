package main

import (
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
}
