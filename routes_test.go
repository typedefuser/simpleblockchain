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

	// Check the response body contains the genesis block
	var blocks []*Block
	err = json.Unmarshal(w.Body.Bytes(), &blocks)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, 1, len(blocks))
	assert.Equal(t, bc.Blocks[0], blocks[0])
}

func TestMineBlock(t *testing.T) {
	// Initialize blockchain with a genesis block
	bc = NewBlockChain()

	// Create a new Gin router
	router := NewRouter()

	// Create a sample transaction
	transaction := Transaction{
		Sender:   "Alice",
		Receiver: "Bob",
		Amount:   10,
	}

	// Marshal the transaction to JSON
	transactionJSON, err := json.Marshal([]Transaction{transaction})
	if err != nil {
		t.Fatal(err)
	}

	// Create a new HTTP request to mine a block
	req, err := http.NewRequest("POST", "/mine", bytes.NewBuffer(transactionJSON))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Create a new recorder to record the response
	w := httptest.NewRecorder()

	// Perform the request
	router.ServeHTTP(w, req)

	// Check the status code is 200 OK
	assert.Equal(t, http.StatusOK, w.Code)

	// Check the response body contains the new block
	var newBlock *Block
	err = json.Unmarshal(w.Body.Bytes(), &newBlock)
	if err != nil {
		t.Fatal(err)
	}

	// Check the new block is added to the blockchain
	assert.Equal(t, 2, len(bc.Blocks))
	assert.Equal(t, newBlock, bc.Blocks[1])
}
