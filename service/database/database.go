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
	GetUserById(id int) (types.User, error)
	GetUserByUsername(username string) (int, error)
	CreateUser(username string) (int, error)
	SetMyUsername(user types.User) error
	CreateChat(chatName string, isGroup bool) (int, error)
	SendMessage(chatID int, userID int, username string, msgType string, text string, mediaUrl string, isForward bool, replyTo int) (int, error)
	ValidateToken(token types.BearerToken) (bool, error)
	UpsertToken(token types.BearerToken) error
	GetPrivateChatID(user1ID int, user2ID int) (int, error)
	AddChatToUser(userID int, chatID int) error
	GetMessage(messageId int) (types.Message, error)
	AddPrivateChat(user1Id int, user2Id int, chatId int) error
	GetMyConversations(userID int) ([]types.Chat, error)
	GetConversation(userId int, chatID int) ([]types.Message, error)
	SetMyPhoto(userId int, path string) error
	GetMyPhoto(token string) (string, error)
	CommentMessage(messageID int, userID int) error
	DeleteComment(messageID int, userID int) error
	AddToGroup(chatID int, userID int) error
	SetGroupPhoto(chatId int, imagePath string) error
	LeaveGroup(chatId int, userId int) (bool, error)
	SetGroupName(chatId int, chatName string) error
	GetLastRead(chatID int) (int, error)
	SetLastRead(userId int, chatId int, lastRead int) error
	GetIdWithToken(token string) (int, error)
	GetChatInfo(chatId int) (types.Chat, error)
	GetUsernameByToken(token string) (string, error)
	GetUsernameById(id int) (string, error)
	GetUsers(userId int) ([]types.User, error)
	GetComments(messageId int) ([]string, error)
	DeleteMessage(messageId int) (bool, error)
	GetForwardInfo(messageId int) (types.ForwardInfo, error)
	GetGroupMembers(chatId int) ([]types.User, error)
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

	sqlStmt := `CREATE TABLE IF NOT EXISTS users (
					id INTEGER PRIMARY KEY AUTOINCREMENT,
					username TEXT NOT NULL UNIQUE,
					image_url TEXT DEFAULT ''
                );
				CREATE TABLE IF NOT EXISTS chats (
					id INTEGER PRIMARY KEY AUTOINCREMENT,
					chat_name TEXT DEFAULT '',
					chat_image TEXT DEFAULT '',
					is_group BOOLEAN DEFAULT FALSE,
					last_msg_id INTEGER DEFAULT 0,
					last_msg_username TEXT DEFAULT '',
					last_msg_text TEXT DEFAULT '',
					last_msg_time DATETIME DEFAULT NULL,
					last_msg_type TEXT DEFAULT ''
				);
				CREATE TABLE IF NOT EXISTS user_chats (
					user_id INTEGER NOT NULL,
					chat_id INTEGER NOT NULL,
					FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
					FOREIGN KEY (chat_id) REFERENCES chats(id) ON DELETE CASCADE,
					PRIMARY KEY (user_id, chat_id)
				);
				CREATE TABLE IF NOT EXISTS private_chats(
    				user1_id INTEGER,
    				user2_id INTEGER,
    				chat_id INTEGER NOT NULL UNIQUE,
    				PRIMARY KEY (user1_id, user2_id),
    				FOREIGN KEY (user1_id) REFERENCES users(id),
    				FOREIGN KEY (user2_id) REFERENCES users(id)
				);

				CREATE TABLE IF NOT EXISTS messages (
					id INTEGER PRIMARY KEY AUTOINCREMENT,
					chat_id INTEGER NOT NULL,
					sender_id INTEGER,
					username TEXT NOT NULL,
					type TEXT NOT NULL,
					text TEXT DEFAULT '',
					media_url TEXT DEFAULT '',
					is_forward BOOLEAN DEFAULT FALSE,
					reply_to INTEGER,
					created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
					FOREIGN KEY (chat_id) REFERENCES chats(id) ON DELETE CASCADE,
					FOREIGN KEY (sender_id) REFERENCES users(id)
				);
				CREATE TABLE IF NOT EXISTS tokens (
    				user_id  INTEGER PRIMARY KEY NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    				token TEXT NOT NULL
  				);
				CREATE TABLE IF NOT EXISTS comments (
    				message_id INTEGER NOT NULL REFERENCES messages(id) ON DELETE CASCADE,
    				user_id INTEGER NOT NULL,
				    PRIMARY KEY (message_id, user_id)
  				);
				CREATE TABLE IF NOT EXISTS last_read (
					user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
					chat_id INTEGER NOT NULL REFERENCES chats(id) ON DELETE CASCADE,
					message_id INTEGER NOT NULL REFERENCES messages(id) ON DELETE CASCADE,
					PRIMARY KEY (user_id, chat_id)
				);
				CREATE TRIGGER IF NOT EXISTS update_messages_username_on_user_delete
				BEFORE DELETE ON users
				FOR EACH ROW
				BEGIN
					UPDATE messages
					SET username = 'deleted account'
					WHERE sender_id = OLD.id;
				END`

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
