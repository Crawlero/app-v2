package db

import (
	"context"
	"fmt"
	"os"
	"sync"

	"github.com/jackc/pgx/v5/pgxpool"
)

var lock = &sync.Mutex{}
var dbpool *pgxpool.Pool

func GetDbPool() *pgxpool.Pool {
	if dbpool == nil {
		lock.Lock()
		defer lock.Unlock()

		if dbpool == nil {
			fmt.Println("Creating pool instance now.")
			config, err := pgxpool.ParseConfig(os.Getenv("DATABASE_URL"))
			config.MaxConns = 10
			config.MinConns = 1

			newPool, err := pgxpool.NewWithConfig(context.Background(), config)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Unable to connection to database: %v\n", err)
				os.Exit(1)
			}

			dbpool = newPool
			fmt.Println("Pool instance created.")
		} else {
			fmt.Println("Single instance already created.")
		}
	} else {
		fmt.Println("Single instance already created.")
	}

	return dbpool
}
