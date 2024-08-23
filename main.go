package main

import "fmt"

var bc *BlockChain

func main() {
	bc = NewBlockChain()
	router := NewRouter()
	filename := "block.json"
	transactions, err := ReadTransactions(filename)
	if err != nil {
		fmt.Println(err)
		return
	}

	go bc.AddTransactions(transactions)

	router.Run(":8080")
}
