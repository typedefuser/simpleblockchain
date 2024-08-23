package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"
)

type Transaction struct {
	Sender   string
	Receiver string
	Amount   float64
}

type Block struct {
	Timestamp     int64
	Transactions  []Transaction
	PrevBlockHash string
	Hash          string
	Nonce         int
	Difficulty    int
	MerkleRoot    string
}

func (b *Block) CalculateHash() string {
	var buff bytes.Buffer

	err := binary.Write(&buff, binary.LittleEndian, b.Timestamp)
	if err != nil {
		fmt.Println("Error writing timestamp:", err)
		return "error"
	}

	for _, t := range b.Transactions {
		buff.WriteString(t.Sender)
		buff.WriteString(t.Receiver)
		err := binary.Write(&buff, binary.LittleEndian, t.Amount)
		if err != nil {
			fmt.Println("Error writing transaction:", err)
			return "error"
		}
	}

	buff.WriteString(b.PrevBlockHash)

	err = binary.Write(&buff, binary.LittleEndian, int64(b.Nonce))
	if err != nil {
		fmt.Println("Error writing nonce:", err)
		return "error"
	}

	hash := sha256.Sum256(buff.Bytes())
	return hex.EncodeToString(hash[:])
}

func (b *Block) MineBlock() {
	target := strings.Repeat("0", b.Difficulty)
	for {
		hash := b.CalculateHash()
		if hash == "" {
			fmt.Println("Error: Unable to calculate hash")
			return
		}
		if len(hash) < b.Difficulty {
			b.Nonce++
			continue
		}
		if hash[:b.Difficulty] == target {
			b.Hash = hash
			fmt.Println("Block mined:", hash)
			break
		} else {
			b.Nonce++
		}
	}
}
func NewBlock(transactions []Transaction, prevBlockHash string, difficulty int) *Block {
	block := &Block{
		Timestamp:     time.Now().Unix(),
		Transactions:  transactions,
		PrevBlockHash: prevBlockHash,
		Difficulty:    difficulty,
	}
	block.MerkleRoot = block.calculateMerkleRoot()
	block.MineBlock()
	return block
}
func (b *Block) calculateMerkleRoot() string {
	var transactionHashes []string
	for _, tx := range b.Transactions {
		hash := sha256.Sum256([]byte(tx.Sender + tx.Receiver + strconv.FormatFloat(tx.Amount, 'f', -1, 64)))
		transactionHashes = append(transactionHashes, hex.EncodeToString(hash[:]))
	}

	if len(transactionHashes) == 0 {
		emptyHash := sha256.Sum256([]byte(""))
		return hex.EncodeToString(emptyHash[:])
	}

	for len(transactionHashes) > 1 {
		var newLevel []string
		for i := 0; i < len(transactionHashes); i += 2 {
			if i+1 < len(transactionHashes) {
				combinedHash := transactionHashes[i] + transactionHashes[i+1]
				hash := sha256.Sum256([]byte(combinedHash))
				newLevel = append(newLevel, hex.EncodeToString(hash[:]))
			} else {
				newLevel = append(newLevel, transactionHashes[i])
			}
		}
		transactionHashes = newLevel
	}
	return transactionHashes[0]
}

type BlockChain struct {
	Blocks []*Block
}

func (b *BlockChain) getLastBlock() *Block {
	return b.Blocks[len(b.Blocks)-1]
}

func (b *BlockChain) getPreviousHash() string {
	return b.Blocks[len(b.Blocks)-1].Hash
}

func (bc *BlockChain) AddBlock(block *Block) {
	bc.Blocks = append(bc.Blocks, block)
}

func (bc *BlockChain) AddTransactions(transactions []Transaction) {
	blocks := SplitTransactionsIntoBlocks(transactions)
	prevBlockHash := bc.getPreviousHash()
	for _, block := range blocks {
		newBlock := NewBlock(block, prevBlockHash, bc.calculateDifficulty(bc.getLastBlock()))
		bc.AddBlock(newBlock)
		prevBlockHash = newBlock.Hash
	}
}

func (bc *BlockChain) calculateDifficulty(previousBlock *Block) int {
	timeTaken := time.Now().Unix() - previousBlock.Timestamp

	switch {
	case timeTaken < 10:
		return previousBlock.Difficulty + 1
	case timeTaken > 20:
		return int(math.Max(float64(previousBlock.Difficulty-1), 1))
	default:
		return previousBlock.Difficulty
	}
}

func NewBlockChain() *BlockChain {
	genesisBlock := NewBlock(
		[]Transaction{},
		"",
		1,
	)
	return &BlockChain{Blocks: []*Block{genesisBlock}}
}
