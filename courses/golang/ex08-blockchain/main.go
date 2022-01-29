package main

import (
	"crypto/sha256"
	"fmt"
	"time"
)

type Block struct {
	timestamp    time.Time
	transactions []string
	prevhash     []byte
	Hash         []byte
}

func main() {
	abc := []string{" A sent 50 coins to BC"}
	xyz := Blocks(abc, []byte{})
	fmt.Println("This is First Block")
	Print(xyz)

	pqrs := []string{" PQ sent 230 coins to RS"}
	klmn := Blocks(pqrs, xyz.Hash)
	fmt.Println("This is Second Block")
	Print(klmn)
}

func Blocks(transactions []string, prevhash []byte) *Block {
	currentTime := time.Now()
	return &Block{
		timestamp:    currentTime,
		transactions: transactions,
		prevhash:     prevhash,
		Hash:         newHash(currentTime, transactions, prevhash),
	}
}

func newHash(time time.Time, transactions []string, prevhash []byte) []byte {
	input := append(prevhash, time.String()...)
	for transaction := range transactions {
		input = append(input, string(rune(transaction))...)
	}
	hash := sha256.Sum256(input)
	return hash[:]
}

func Print(block *Block) {
	fmt.Printf("\ttime: %s\n", block.timestamp.String())
	fmt.Printf("\tprevhash: %x\n", block.prevhash)
	fmt.Printf("\thash: %x\n", block.Hash)
}

func Transaction(block *Block) {
	fmt.Println("\tTransactions: ")
	for i, transaction := range block.transactions {
		fmt.Printf("\t\t%v: %q\n", i, transaction)
	}
}
