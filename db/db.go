package db

import (
	_ "context"
	"database/sql"
	"log"

	_ "github.com/ncruces/go-sqlite3/driver"
	_ "github.com/ncruces/go-sqlite3/embed"
)

var AppDB *sql.DB

func InitDB() {
	var err error
	AppDB, err = sql.Open("sqlite3", "file:app.db")
	if err != nil {
		log.Fatal(err)
	}

}

func SetupDB() {

	tx, err := AppDB.Begin()
	if err != nil {
		log.Fatal(err)
	} else {
		_, err = tx.Exec(`CREATE TABLE IF NOT EXISTS users (
							id INTEGER PRIMARY KEY,
							username VARCHAR(100) NOT NULL,
							email VARCHAR(200) UNIQUE NOT NULL,
							pass VARCHAR(200),
							org_id INTEGER,
							key VARCHAR(200),
							FOREIGN KEY (org_id) REFERENCES orgs(id) ON DELETE SET NULL
						)`)
		if err != nil {
			_ = tx.Rollback()
			log.Fatal(err)
		}
		_, err = tx.Exec(`CREATE TABLE IF NOT EXISTS orgs (
								id INTEGER PRIMARY KEY,
								name VARCHAR(100) NOT NULL,
								ceo_id INTEGER NOT NULL,
								FOREIGN KEY (ceo_id) REFERENCES users(id) ON DELETE CASCADE
							)`)
		if err != nil {
			_ = tx.Rollback()
			log.Fatal(err)
		}
		_, err = tx.Exec(`CREATE TABLE IF NOT EXISTS teams (
								id INTEGER PRIMARY KEY,
								name VARCHAR(100) NOT NULL,
								org_id INTEGER NOT NULL,
								sup_id INTEGER NOT NULL,
								FOREIGN KEY (org_id) REFERENCES orgs(id) ON DELETE CASCADE
								FOREIGN KEY (sup_id) REFERENCES users(id) ON DELETE CASCADE

							)`)
		if err != nil {
			_ = tx.Rollback()
			log.Fatal(err)
		}
		_, err = tx.Exec(`CREATE TABLE IF NOT EXISTS user_team (
								id INTEGER PRIMARY KEY,
								user_id INTEGER NOT NULL,
								team_id INTEGER NOT NULL,
								user_role INTEGER NOT NULL,
								FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
								FOREIGN KEY (team_id) REFERENCES teams(id) ON DELETE CASCADE
							)`)
		if err != nil {
			_ = tx.Rollback()
			log.Fatal(err)
		}
		_, err = tx.Exec(`CREATE TABLE IF NOT EXISTS org_inv (
								id INTEGER PRIMARY KEY,
								user_id INTEGER NOT NULL,
								org_id INTEGER NOT NULL,
								key VARCHAR(200),
								FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
								FOREIGN KEY (org_id) REFERENCES orgs(id) ON DELETE CASCADE
							)`)
		if err != nil {
			_ = tx.Rollback()
			log.Fatal(err)
		}
		_, err = tx.Exec(`CREATE TABLE IF NOT EXISTS chats (
								id INTEGER PRIMARY KEY,
								team_id INTEGER NOT NULL,
								name VARCHAR(200),
								key VARCHAR(200),
								FOREIGN KEY (team_id) REFERENCES teams(id) ON DELETE CASCADE,
							)`)
		if err != nil {
			_ = tx.Rollback()
			log.Fatal(err)
		}
		_, err = tx.Exec(`CREATE TABLE IF NOT EXISTS msgs (
								id INTEGER PRIMARY KEY,
								user_id INTEGER NOT NULL,
								chat_id INTEGER NOT NULL,
								content TEXT,
								FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
								FOREIGN KEY (chat_id) REFERENCES chats(id) ON DELETE CASCADE
							)`)
		if err != nil {
			_ = tx.Rollback()
			log.Fatal(err)
		}

	}
	if err := tx.Commit(); err != nil {
		log.Fatal(err)
	}

}
