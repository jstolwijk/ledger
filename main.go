package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
)

/*
	@ validate idempotency key should be unique per user/tenant (??)
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

type Ledger struct {
	Id      string `db:"id" json:"id"`
	Name    string `db:"name" json:"name"`
	SumFrom int64  `db:"sum_from" json:"sumFrom"`
	SumTo   int64  `db:"sum_to" json:"sumTo"`
	Version int64  `db:"version" json:"version"`
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

	r.GET("/ledgers", func(c *gin.Context) {
		var ledgers []Ledger
		if err := pgxscan.Select(context.Background(), db, &ledgers, `select * from ledger`); err == nil {
			c.JSON(200, ledgers)
		} else {
			c.JSON(500, gin.H{
				"error": err.Error(),
			})
		}
	})

	// @ DESCRIPTION: Create a new journal entry
	r.POST("/journal", func(c *gin.Context) {
		var request JournalEntryRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		journalEntryId, err := addToJournal(request)

		if err == nil {
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

// todo: transactional shite
func addToJournal(request JournalEntryRequest) (*string, error) {
	fromLedger, err := upsertLedger(request.From, request.Amount.Value, 0)

	if err != nil {
		return nil, err
	}

	toLedger, err := upsertLedger(request.To, 0, request.Amount.Value)

	if err != nil {
		return nil, err
	}

	query := `
	insert into journal_entry
	(id, idempotency_key, from_account, to_account, amount_value, amount_currency, metadata) values 
	($1, $2, $3, $4, $5, $6, $7)
	`

	insertLedgerEntry := `
		insert into ledger_entry
		(id, ledger_id, ledger_name, ledger_version, journal_entry_id, currency, from_amount, to_amount, ledger_sum_to, ledger_sum_from) values
		($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`

	fmt.Println(request)
	journalEntryId := uuid.New().String()
	ledgerEntryFromId := uuid.New().String()
	ledgerEntryToId := uuid.New().String()

	fmt.Println(journalEntryId)

	// TODO: handle error
	rawMetadata, _ := json.Marshal(request.Metadata)

	if _, err := db.Exec(context.Background(), query, journalEntryId, request.IdempotencyKey, request.From, request.To, request.Amount.Value, request.Amount.Currency, rawMetadata); err == nil {
		if _, err := db.Exec(context.Background(), insertLedgerEntry, ledgerEntryFromId, fromLedger.Id, fromLedger.Name, fromLedger.Version, journalEntryId, request.Amount.Currency, request.Amount.Value, nil, fromLedger.SumTo, fromLedger.SumFrom); err == nil {
			if _, err := db.Exec(context.Background(), insertLedgerEntry, ledgerEntryToId, toLedger.Id, toLedger.Name, toLedger.Version, journalEntryId, request.Amount.Currency, nil, request.Amount.Value, toLedger.SumTo, toLedger.SumFrom); err == nil {
				return &journalEntryId, nil
			} else {
				return nil, err
			}
		} else {
			return nil, err
		}
	} else {
		return nil, err
	}
}

func upsertLedger(ledgerName string, fromAmount int64, toAmount int64) (*Ledger, error) {
	query := `
		insert into ledger
		(id, name, version, sum_from, sum_to) values
		($1, $2, 1, $3, $4)
		on conflict (name) do update set
		version = ledger.version + 1,
		sum_from = ledger.sum_from + $3,
		sum_to = ledger.sum_to + $4
		returning *
	`

	var ledgers []*Ledger
	err := pgxscan.Select(context.Background(), db, &ledgers, query, uuid.New().String(), ledgerName, fromAmount, toAmount)

	if err != nil {
		return nil, err
	}

	if len(ledgers) == 1 {
		return ledgers[0], nil
	} else if len(ledgers) > 1 {
		return nil, errors.New("multiple ledgers found")
	} else {
		return nil, errors.New("no ledger found")
	}
}
