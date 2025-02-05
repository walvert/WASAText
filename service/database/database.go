/*
Package database is the middleware between the app database and the code. All data (de)serialization (save/load) from a
persistent database are handled here. Database specific logic should never escape this package.

To use this package you need to apply migrations to the database if needed/wanted, connect to it (using the database
data source name from config), and then initialize an instance of AppDatabase from the DB connection.

For example, this code adds a parameter in `webapi` executable for the database data source name (add it to the
main.WebAPIConfiguration structure):

	DB struct {
		Filename string `conf:""`
	}

This is an example on how to migrate the DB and connect to it:

	// Start Database
	logger.Println("initializing database support")
	db, err := sql.Open("sqlite3", "./foo.db")
	if err != nil {
		logger.WithError(err).Error("error opening SQLite DB")
		return fmt.Errorf("opening SQLite: %w", err)
	}
	defer func() {
		logger.Debug("database stopping")
		_ = db.Close()
	}()

Then you can initialize the AppDatabase and pass it to the api package.
*/
package database

import (
	"database/sql"
	"errors"
	"fmt"
	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/types"
)

// AppDatabase is the high level interface for the DB
type AppDatabase interface {
	Ping() error
	CreateUser(username string) (types.User, error)
	GetUser(id int) (types.User, error)
	EditUsername(id int, username string) (types.User, error)
	UsernameExists(username string) bool
	EditUserImage(id int, username string) error
	CreateChat(chat types.Chat) (types.Chat, error)
}

type appdbimpl struct {
	c *sql.DB
}

// New returns a new instance of AppDatabase based on the SQLite connection `db`.
// `db` is required - an error will be returned if `db` is `nil`.
func New(db *sql.DB) (AppDatabase, error) {
	if db == nil {
		return nil, errors.New("database is required when building a AppDatabase")
	}

	sqlStmt := `CREATE TABLE IF NOT EXISTS users (id INTEGER PRIMARY KEY AUTOINCREMENT, username TEXT NOT NULL UNIQUE);
				CREATE TABLE IF NOT EXISTS chats (id INTEGER PRIMARY KEY, chat_name TEXT, participants TEXT);
				CREATE TABLE IF NOT EXISTS image_paths (id INTEGER PRIMARY KEY, path TEXT NOT NULL,FOREIGN KEY (id) REFERENCES users (id) ON DELETE CASCADE);
				CREATE TABLE IF NOT EXISTS user_chats (user_id INTEGER NOT NULL, chat_id INTEGER NOT NULL, FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE, FOREIGN KEY (chat_id) REFERENCES chats(id) ON DELETE CASCADE, PRIMARY KEY (user_id, chat_id));
				CREATE TABLE IF NOT EXISTS messages (id INTEGER PRIMARY KEY, chat_id INTEGER,user_id INTEGER,content TEXT NOT NULL,created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, FOREIGN KEY (chat_id) REFERENCES chats(id) ON DELETE CASCADE, FOREIGN KEY (user_id) REFERENCES users(id));`

	_, err := db.Exec(sqlStmt)
	if err != nil {
		return nil, fmt.Errorf("error creating database structure: %w", err)
	}

	return &appdbimpl{
		c: db,
	}, nil
}

func (db *appdbimpl) Ping() error {
	return db.c.Ping()
}
