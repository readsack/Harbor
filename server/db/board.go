package db

import (
	"database/sql"
	"log"
)

type Card struct {
	ID       int    `json:"id"`
	ColumnID int    `json:"column_id"`
	Content  string `json:"content"`
	User     `json:"user"`
}

type Column struct {
	Cards   []Card `json:"cards"`
	ID      int    `json:"id"`
	BoardID int    `json:"board_id"`
	Name    string `json:"name"`
}

type Board struct {
	ID      int      `json:"id"`
	TeamID  int      `json:"team_id"`
	Columns []Column `json:"columns"`
}

func CreateBoard(team_id int) error {
	_, err := AppDB.Exec("INSERT INTO boards (team_id) VALUES (?)", team_id)
	return err
}

func CreateColumn(board_id int, name string) error {
	_, err := AppDB.Exec("INSERT INTO columns (board_id, name) VALUES (?, ?)", board_id, name)
	return err
}

func CreateItem(column_id int, user_id int, content string) error {
	_, err := AppDB.Exec("INSERT INTO items (column_id, content, last_writer) VALUES (?, ?, ?)", column_id, content, user_id)
	return err
}

func GetTeamBoardID(team_id int) (int, error) {
	var col_id int
	err := AppDB.QueryRow("SELECT id FROM boards WHERE team_id=?", team_id).Scan(&col_id)
	return col_id, err
}

func GetBoard(team_id int) (Board, error) {
	board_id, err := GetTeamBoardID(team_id)
	if err != nil {
		return Board{}, err
	}
	var b Board
	b.ID = board_id
	b.TeamID = team_id
	b.Columns = make([]Column, 0)
	cols, err := AppDB.Query("SELECT * FROM columns WHERE board_id=?", board_id)
	if err != nil {
		return b, err
	}
	for cols.Next() {
		var c Column
		err = cols.Scan(&c.ID, &c.BoardID, &c.Name)
		if err != nil && err != sql.ErrNoRows {
			return b, err
		}
		cards, err := AppDB.Query("SELECT * FROM items WHERE column_id=?", c.ID)
		c.Cards = make([]Card, 0)
		if err != nil && err != sql.ErrNoRows {
			log.Println("hiee")
			return b, err
		}
		for cards.Next() {
			var card Card
			var user_id int
			err = cards.Scan(&card.ID, &card.ColumnID, &card.Content, &user_id)
			if err != nil {
				log.Println("hiw")
				log.Fatal(err)
			}
			usr, err := FindUserByID(user_id)
			if err != nil {
				log.Println("hi")
				log.Fatal(err)
			}
			//log.Println(usr)
			card.User = *usr
			c.Cards = append(c.Cards, card)
		}
		b.Columns = append(b.Columns, c)
	}
	return b, nil
}
