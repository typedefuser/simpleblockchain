package main

import (
	"github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine {
	router := gin.Default()
	router.GET("/blocks", getBlocks)
	router.POST("/mine", mineBlock)
	return router
}
