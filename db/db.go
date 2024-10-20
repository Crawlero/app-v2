package db

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"os"
	"sync"
)

var lock = &sync.Mutex{}
var dbpool *pgxpool.Pool

func GetDbPool() *pgxpool.Pool {
	if dbpool == nil {
		lock.Lock()
		defer lock.Unlock()

		if dbpool == nil {
			fmt.Println("Creating pool instance now.")
			newPool, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
			if err != nil {
				fmt.Fprintf(os.Stderr, "Unable to connection to database: %v\n", err)
				os.Exit(1)
			}

			dbpool = newPool
		} else {
			fmt.Println("Single instance already created.")
		}
	} else {
		fmt.Println("Single instance already created.")
	}

	return dbpool
}
