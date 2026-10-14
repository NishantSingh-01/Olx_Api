package database


import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"
	
	_ "github.com/jackc/pgx/v5/stdlib"
)

func Connect(db_Url string) (*sql.DB, error) {  // todo to look o it and other video of sql connection mention in  this there 
	db, err := sql.Open("pgx", db_Url)
	if err != nil {
		return nil, fmt.Errorf("SQL.Open : Database Connection failed : %w", err)
	}
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(10 *time.Minute)

	// fail test here:
	ctx,cancel := context.WithTimeout(context.Background(),5*time.Second)
	defer cancel() // release the resources 
	err = db.PingContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("db.Ping: Database Connection failed : %w", err)
	}

	log.Println("Database connection successful")
	return db, nil
}

