package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
)

const maxTransactionsPerBlock = 10

func ReadTransactions(filename string) ([]Transaction, error) {
	file, err := os.Open(filename)
	if err != nil {
		log.Printf("Error opening file: %s", err)
		return nil, fmt.Errorf("error opening file: %s", err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		log.Printf("Error reading file: %s", err)
		return nil, fmt.Errorf("error reading file: %s", err)
	}

	var transactions []Transaction
	err = json.Unmarshal(data, &transactions)
	if err != nil {
		log.Printf("Error unmarshalling JSON: %s", err)
		return nil, fmt.Errorf("error unmarshalling JSON: %s", err)
	}

	return transactions, nil
}

func SplitTransactionsIntoBlocks(transactions []Transaction) [][]Transaction {
	var blocks [][]Transaction
	for i := 0; i < len(transactions); i += maxTransactionsPerBlock {
		end := i + maxTransactionsPerBlock
		if end > len(transactions) {
			end = len(transactions)
		}
		blocks = append(blocks, transactions[i:end])
	}
	return blocks
}
