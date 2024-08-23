package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func getBlocks(c *gin.Context) {
	c.JSON(http.StatusOK, bc.Blocks)
}

func mineBlock(c *gin.Context) {
	var transactions []Transaction

	if err := c.ShouldBindJSON(&transactions); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	lastBlock := bc.getLastBlock()
	newBlock := NewBlock(transactions, lastBlock.Hash, bc.calculateDifficulty(lastBlock))
	bc.AddBlock(newBlock)
	c.JSON(http.StatusOK, newBlock)
}
