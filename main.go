package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

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
	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "up",
		})
	})

	/*
		{
			"idempotencyKey": "183",
			"from": "jesse",
			"to": "jan",
			"amount": {
				"value": 100,
				"currency": "EUR"
			},
			"metadata": {
				"orderReference": "39123912"
			}
		}
	*/

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

	r.Run()

}
