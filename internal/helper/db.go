package helper

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/Elanoran/d2go/pkg/data"
	"github.com/Elanoran/koolo/internal/config"
	_ "github.com/go-sql-driver/mysql"
)

// Function to open a database connection
func openDBConnection() (*sql.DB, error) {
	// Check if db is enabled in config
	if !config.Config.Db.Enabled {
		return nil, nil
	}

	// Construct the database connection string
	connStr := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", config.Config.Db.User, config.Config.Db.Pass, config.Config.Db.Host, config.Config.Db.Port, config.Config.Db.DbName)

	// Open a connection to the database
	db, err := sql.Open("mysql", connStr)
	if err != nil {
		return nil, fmt.Errorf("error opening database connection: %v", err)
	}

	// Check if the connection to the database is successful
	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("error connecting to the database: %v", err)
	}

	return db, nil
}

// DbConnectAndSendText connects to the database and sends text information
func DbLogText(textData string) error {
	// Open a database connection
	db, err := openDBConnection()
	if err != nil {
		return err
	}
	defer db.Close()

	// Perform database operation: insert the provided text into a table
	_, err = db.Exec(fmt.Sprintf("INSERT INTO logs (timestamp, message, bot_name, bot_char_name) VALUES (CURRENT_TIMESTAMP, ?, ?, ?)"),
		textData, config.Config.Db.BotName, config.Config.Db.BotCharName)
	if err != nil {
		return fmt.Errorf("DbLogText:error inserting data into the database: %v", err)
	}

	return nil
}

func DbupsertBotInfoStart(gameNum int, gameRuns string) error {
	// Open a database connection
	db, err := openDBConnection()
	if err != nil {
		return err
	}
	defer db.Close()

	// Begin a transaction
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("DbupsertBotInfo: error beginning transaction: %v", err)
	}
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p) // re-throw panic after rollback
		} else if err != nil {
			tx.Rollback() // rollback if there was an error
		} else {
			err = tx.Commit() // commit if everything went well
		}
	}()

	// Try to update the last seen timestamp for the bot
	result, err := db.Exec("UPDATE bots SET last_seen = CURRENT_TIMESTAMP, game_number = ?, start_time = CURRENT_TIMESTAMP, running_char = ?, run_list = ?, in_area = ? WHERE bot_name = ?",
		gameNum,
		config.Config.Db.BotCharName,
		gameRuns,
		"start",
		config.Config.Db.BotName,
	)
	if err != nil {
		return err
	}

	// game_number, start_time, running_char, run_list, in_area

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	// If no rows were affected, the bot doesn't exist, so insert a new one
	if rowsAffected == 0 {
		_, err := tx.Exec("INSERT INTO bots (bot_name, last_seen, game_number, start_time, running_char, run_list, in_area) VALUES (?, CURRENT_TIMESTAMP, 0, CURRENT_TIMESTAMP, 'unknown', 'unknown', 'unknown')", config.Config.Db.BotName)
		if err != nil {
			return err
		}
	}

	return nil
}

func DbupsertBotInfoRunning(d data.Data) error {
	// Open a database connection
	db, err := openDBConnection()
	if err != nil {
		return err
	}
	defer db.Close()

	// Begin a transaction
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("DbupsertBotInfo: error beginning transaction: %v", err)
	}
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p) // re-throw panic after rollback
		} else if err != nil {
			tx.Rollback() // rollback if there was an error
		} else {
			err = tx.Commit() // commit if everything went well
		}
	}()

	// Try to update the last seen timestamp for the bot
	result, err := tx.Exec("UPDATE bots SET last_seen = CURRENT_TIMESTAMP, in_area = ? WHERE bot_name = ?", d.PlayerUnit.Area, config.Config.Db.BotName)
	if err != nil {
		return err
	}

	// game_number, start_time, running_char, run_list, in_area

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	// If no rows were affected, the bot doesn't exist, so insert a new one
	if rowsAffected == 0 {
		_, err := tx.Exec("INSERT INTO bots (bot_name, last_seen, game_number, start_time, running_char, run_list, in_area) VALUES (?, CURRENT_TIMESTAMP, 0, CURRENT_TIMESTAMP, 'unknown', 'unknown', 'unknown')", config.Config.Db.BotName)
		if err != nil {
			return err
		}
	}

	return nil
}

func DbupsertBotRuns(d data.Data, numgames int) error {
	// Open a database connection
	db, err := openDBConnection()
	if err != nil {
		return err
	}
	defer db.Close()

	// Get the bot ID
	var botID int
	err = db.QueryRow("SELECT id FROM bots WHERE bot_name = ?", config.Config.Db.BotName).Scan(&botID)
	if err != nil {
		return fmt.Errorf("upsertCurrentCharacter: error querying bot ID: %v", err)
	}

	// Perform database operation: update the current character information
	result, err := db.Exec("UPDATE current_characters SET bot_char_name = ?, class = ?, game_number = ?, game_run_names = ?, current_area_name = ?, game_start_time = ? WHERE bot_id = ?",
		config.Config.Db.BotCharName,
		config.Config.Character.Class,
		numgames,
		config.Config.Game.Runs,
		d.PlayerUnit.Area,
		time.Now().Format("2006-01-02 15:04:05"),
		botID,
	)
	if err != nil {
		return fmt.Errorf("upsertCurrentCharacter: error updating data in the database: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("upsertCurrentCharacter: error getting rows affected: %v", err)
	}

	// If no rows were affected, the current character doesn't exist, so insert a new one
	if rowsAffected == 0 {
		_, err := db.Exec("INSERT INTO current_characters (bot_id, bot_char_name, class, game_number, game_run_names, current_area_name, game_start_time) VALUES (?, ?, ?, ?, ?, ?, ?)",
			botID,
			config.Config.Db.BotCharName,
			config.Config.Character.Class,
			numgames,
			config.Config.Game.Runs,
			d.PlayerUnit.Area,
			time.Now().Format("2006-01-02 15:04:05"),
		)
		if err != nil {
			return fmt.Errorf("upsertCurrentCharacter: error inserting data into the database: %v", err)
		}
	}

	//fmt.Println("Current character info inserted or updated successfully")
	return nil
}

// DbConnectAndSendText connects to the database and sends text information
func DbGetStatNameFromID(id int) (string, error) {
	// Open a database connection
	db, err := openDBConnection()
	if err != nil {
		return "", err
	}
	defer db.Close()

	// Replace "yourDB" and "stat_mapping" with your actual database name and table name
	query := "SELECT real_stat_string FROM stat_mapping WHERE id = ?"

	// Assuming you have a database connection stored in a variable db
	row := db.QueryRow(query, id)

	var realStatString string
	err2 := row.Scan(&realStatString)
	if err2 != nil {
		return "", err2
	}

	return realStatString, nil
}

func DbLogItem(item data.Item) error {
	// Open a database connection
	db, err := openDBConnection()
	if err != nil {
		return err
	}
	defer db.Close()

	// Initialize an empty slice to store the stats information
	var statsInfo []string

	// Iterate through each stat in the item's stats
	for statID, statData := range item.Stats {
		// Construct a string for each stat
		statInfo := fmt.Sprintf("%s %d", statID, statData.Value)

		// Add the statInfo to the statsInfo slice
		statsInfo = append(statsInfo, statInfo)
	}

	// Join the statsInfo slice into a single string
	statsString := strings.Join(statsInfo, ", ")

	// Marshal the item's stats into JSON
	statsJSON, err := json.Marshal(statsString)
	if err != nil {
		return fmt.Errorf("DbLogItem: error serializing stats to JSON: %v", err)
	}

	// Perform database operation: insert the item into the items table
	_, err = db.Exec("INSERT INTO items (item_name, item_quality, bot_char_name, stats) VALUES (?, ?, ?, ?)",
		item.Name, item.Quality.ToString(), config.Config.Db.BotCharName, statsJSON)
	if err != nil {
		return fmt.Errorf("DbLogItem: error inserting data into the database: %v", err)
	}

	return nil
}

// DbupsertPlayerInfoRunning updates player information in the database.
func DbupsertPlayerInfoRunning(d data.Data, level int) error {
	// Open a database connection
	db, err := openDBConnection()
	if err != nil {
		return err
	}
	defer db.Close()

	// Begin a transaction
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("DbupsertPlayerInfoRunning: error beginning transaction: %v", err)
	}
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p) // re-throw panic after rollback
		} else if err != nil {
			tx.Rollback() // rollback if there was an error
		} else {
			err = tx.Commit() // commit if everything went well
		}
	}()

	// Check if the character exists
	var characterExists int
	err = tx.QueryRow("SELECT COUNT(*) FROM characters WHERE character_name = ?", d.PlayerUnit.Name).Scan(&characterExists)
	if err != nil {
		return err
	}

	if characterExists > 0 {
		// Character exists, update the player information
		_, err := tx.Exec("UPDATE characters SET level = ?, spec = ?, class = ? WHERE character_name = ?",
			level, config.Config.Character.Build, config.Config.Character.Class, d.PlayerUnit.Name)
		if err != nil {
			return err
		}
	} else {
		// Character doesn't exist, insert a new one
		_, err := tx.Exec("INSERT INTO characters (character_name, class, level, spec) VALUES (?, ?, ?, ?)",
			d.PlayerUnit.Name, config.Config.Character.Class, level, config.Config.Character.Build)
		if err != nil {
			return err
		}
	}

	return nil
}
