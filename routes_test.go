package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetBlocks(t *testing.T) {
	// Initialize blockchain with a genesis block
	bc = NewBlockChain()

	// Mine a new block with some transactions
	bc.AddTransactions([]Transaction{
		{Sender: "Alice", Receiver: "Bob", Amount: 10},
		{Sender: "Bob", Receiver: "Charlie", Amount: 5},
	})

	// Create a new Gin router
	router := NewRouter()

	// Create a new HTTP request to get blocks
	req, err := http.NewRequest("GET", "/blocks", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a new recorder to record the response
	w := httptest.NewRecorder()

	// Perform the request
	router.ServeHTTP(w, req)

	// Check the status code is 200 OK
	assert.Equal(t, http.StatusOK, w.Code)

	// Check the response body contains the correct number of blocks
	var blocks []*Block
	err = json.Unmarshal(w.Body.Bytes(), &blocks)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, len(bc.Blocks), len(blocks))

	for i, block := range blocks {
		assert.Equal(t, bc.Blocks[i], block)
	}
}

func TestMineBlockWithMultipleTransactions(t *testing.T) {

	bc = NewBlockChain()

	router := NewRouter()

	transactions := []Transaction{
		{Sender: "Alice", Receiver: "Bob", Amount: 10},
		{Sender: "Bob", Receiver: "Charlie", Amount: 5},
		{Sender: "Charlie", Receiver: "David", Amount: 2.5},
	}

	transactionJSON, err := json.Marshal(transactions)
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("POST", "/mine", bytes.NewBuffer(transactionJSON))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var newBlock *Block
	err = json.Unmarshal(w.Body.Bytes(), &newBlock)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, 2, len(bc.Blocks))
	assert.Equal(t, newBlock, bc.Blocks[1])

	assert.Equal(t, len(transactions), len(newBlock.Transactions))
	for i, tx := range transactions {
		assert.Equal(t, tx, newBlock.Transactions[i])
	}
}
