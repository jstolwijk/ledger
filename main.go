package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
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

// pgxpool is a thread-safe connection pool for PostgreSQL.
var db *pgxpool.Pool

func main() {

	poolConfig, err := pgxpool.ParseConfig("postgres://postgres:postgres@localhost:5432/ledger")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to parse DATABASE_URL %v\n", err)
		os.Exit(1)
	}

	db, err = pgxpool.ConnectConfig(context.Background(), poolConfig)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer db.Close()

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
		var request JournalEntryRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		query := `
		insert into journal_entry
		(id, idempotency_key, from_account, to_account, amount_value, amount_currency, metadata) values 
		($1, $2, $3, $4, $5, $6, $7)
		`

		fmt.Println(request)
		journalEntryId := uuid.New().String()
		fmt.Println(journalEntryId)

		// TODO: handle error
		rawMetadata, _ := json.Marshal(request.Metadata)

		if _, err := db.Exec(context.Background(), query, journalEntryId, request.IdempotencyKey, request.From, request.To, request.Amount.Value, request.Amount.Currency, rawMetadata); err == nil {
			c.JSON(200, gin.H{
				"status":         "added",
				"journalEntryId": journalEntryId,
			})
		} else {
			c.JSON(500, gin.H{
				"status": "failed",
				"err":    err.Error(),
			})
		}

	})

	r.POST("/test", func(c *gin.Context) {
		if _, err := db.Exec(context.Background(), "select 1"); err == nil {
			c.JSON(200, gin.H{
				"status": "added",
			})
		} else {
			c.JSON(500, gin.H{
				"status": "failed",
			})
		}

	})

	// @ DESCRIPTION: Post a ledger up to another ledger
	// @ from ledger to ledger
	r.POST("/post", func(c *gin.Context) {
		// TODO: Implement

		c.JSON(200, gin.H{
			"status":         "added",
			"journalEntryId": "123",
		})
	})

	r.Run()
}
