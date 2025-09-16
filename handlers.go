package main

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func GetPosts(c *gin.Context, db *pgxpool.Pool) {
	rows, err := db.Query(context.Background(), "SELECT id, title, content FROM posts")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var p Post
		if err := rows.Scan(&p.ID, &p.Title, &p.Content); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		posts = append(posts, p)
	}

	c.JSON(http.StatusOK, posts)
}

func CreatePost(c *gin.Context, db *pgxpool.Pool) {
	var p Post
	if err := c.BindJSON(&p); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := db.QueryRow(ctx,
		"INSERT INTO posts (title, content) VALUES ($1, $2) RETURNING id",
		p.Title, p.Content).Scan(&p.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, p)
}

func UpdatePost(c *gin.Context, db *pgxpool.Pool) {
	id, _ := strconv.Atoi(c.Param("id"))
	var p Post
	if err := c.BindJSON(&p); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := db.Exec(context.Background(),
		"UPDATE posts SET title=$1, content=$2 WHERE id=$3",
		p.Title, p.Content, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "updated"})
}

func DeletePost(c *gin.Context, db *pgxpool.Pool) {
	id, _ := strconv.Atoi(c.Param("id"))

	_, err := db.Exec(context.Background(), "DELETE FROM posts WHERE id=$1", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}
