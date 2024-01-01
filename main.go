package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"
)

const (
	lockID = "my_lock"
	ttl    = 10 // seconds
)

func main() {
	serverID := os.Args[1] // Unique ID for each instance

	baseConnectionString := os.Getenv("DB_BASE_CONNECTION_STRING")
	if baseConnectionString == "" {
		log.Fatalf("Base DB connection string is not set")
	}

	connectionString := fmt.Sprintf("%s?application_name=%s", baseConnectionString, serverID)

	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()

	logfile, err := os.OpenFile("log_"+serverID+".txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Error opening log file: %v", err)
	}
	defer logfile.Close()

	log.SetOutput(logfile)

	for {
		if acquireLock(db, serverID) {
			fmt.Println(serverID, "acquired the lock")
			log.Println("Acquired the lock")
			time.Sleep(10 * time.Second) // Simulate work
			releaseLock(db, serverID)
			fmt.Println(serverID, "released the lock")
			log.Println("Released the lock")
		}
		time.Sleep(1 * time.Second) // Wait before trying again
	}
}

func acquireLock(db *sql.DB, serverID string) bool {
	tx, err := db.Begin()
	if err != nil {
		log.Println("Error starting transaction")
		return false
	}

	// Ensure lock is not expired
	var existingServerID string
	var timestamp time.Time
	err = tx.QueryRow("SELECT server_id, timestamp FROM distributed_locks WHERE lock_id = $1", lockID).Scan(&existingServerID, &timestamp)
	if err == nil && time.Since(timestamp) < ttl*time.Second {
		tx.Rollback()
		return false // Lock is already held and not expired
	}

	// Acquire or update the lock
	_, err = tx.Exec("UPSERT INTO distributed_locks (lock_id, server_id, timestamp) VALUES ($1, $2, now())", lockID, serverID)
	if err != nil {
		tx.Rollback()
		log.Println("Error acquiring lock")
		return false
	}

	return tx.Commit() == nil
}

func releaseLock(db *sql.DB, serverID string) {
	_, err := db.Exec("DELETE FROM distributed_locks WHERE lock_id = $1 AND server_id = $2", lockID, serverID)
	if err != nil {
		log.Println("Error releasing lock")
	}
}
