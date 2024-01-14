package helper

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/hectorgimenez/koolo/internal/config"
)

// DbConnectAndSendText connects to the database and sends text information
func DbConnectAndSendText(tableName string, textData string) error {

	// Check if db is enabled in config
	if !config.Config.Db.Enabled {
		return nil
	}

	// Replace these values with your actual database connection details
	dbHost := config.Config.Db.Host     // "db-host"
	dbPort := config.Config.Db.Port     // "db-port"
	dbUser := config.Config.Db.User     // "db-username"
	dbPassword := config.Config.Db.Pass // "db-password"
	dbName := config.Config.Db.Db       // "db-name"

	// Construct the database connection string
	connStr := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPassword, dbHost, dbPort, dbName)

	// Open a connection to the database
	db, err := sql.Open("mysql", connStr)
	if err != nil {
		return fmt.Errorf("error opening database connection: %v", err)
	}
	defer db.Close()

	// Check if the connection to the database is successful
	err = db.Ping()
	if err != nil {
		return fmt.Errorf("error connecting to the database: %v", err)
	}

	//fmt.Println("Connected to the database")

	// Perform database operation: insert the provided text into a table
	_, err = db.Exec(fmt.Sprintf("INSERT INTO %s (data) VALUES ($1)", tableName), textData)
	if err != nil {
		return fmt.Errorf("error inserting data into the database: %v", err)
	}

	//fmt.Printf("Text sent to the table '%s' in the database\n", tableName)
	return nil
}
