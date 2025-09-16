package main

import (
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	db := InitDB()
	defer db.Close()

	r := gin.Default()
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})
	r.GET("/posts", func(c *gin.Context) { GetPosts(c, db) })
	r.POST("/posts", func(c *gin.Context) { CreatePost(c, db) })
	r.PUT("/posts/:id", func(c *gin.Context) { UpdatePost(c, db) })
	r.DELETE("/posts/:id", func(c *gin.Context) { DeletePost(c, db) })

	log.Println("Server started on :8080")
	r.Run(":8080")
}
