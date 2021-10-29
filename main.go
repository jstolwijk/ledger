package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4"
)

/*

 */

type JournalEntryRequest struct {
	IdempotencyKey string            `json:"idempotencyKey" binding:"required"`
	From           string            `json:"from" binding:"required"`
	To             string            `json:"to" binding:"required"`
	Amount         Amount            `json:"amount" binding:"required"`
	Metadata       map[string]string `json:"metadata" binding:"required"`
}

type Amount struct {
	Currency string `json:"currency" binding:"required"`
	Value    int64  `json:"value" binding:"required"`
}

func main() {
	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	var name string
	var weight int64
	err = conn.QueryRow(context.Background(), "select name, weight from widgets where id=$1", 42).Scan(&name, &weight)
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(name, weight)

	webserver()
}

func webserver() {
	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "up",
		})
	})

	// @ DESCRIPTION: Create a new journal entry
	r.POST("/journal", func(c *gin.Context) {
		var json JournalEntryRequest
		if err := c.ShouldBindJSON(&json); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		fmt.Println(json)

		c.JSON(200, gin.H{
			"status":         "added",
			"journalEntryId": "123",
		})
	})

	// @ DESCRIPTION: Post a ledger up to another ledger
	r.POST("/post", func(c *gin.Context) {
		// TODO: Implement

		c.JSON(200, gin.H{
			"status":         "added",
			"journalEntryId": "123",
		})
	})

	r.Run()
}
