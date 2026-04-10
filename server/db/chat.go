package db

import "crypto/rand"

type Chat struct {
	Name   string `json:"name"`
	ID     int    `json:"id"`
	TeamID int    `json:"team_id"`
	Key    string `json:"key"`
}

func CreateChat(name string, team_id int) error {
	_, err := AppDB.Exec("INSERT INTO chats (name, team_id, key) VALUES (?, ?, ?)", name, team_id, rand.Text())
	return err
}

func AddMsg(content string, chat_id int, user_id int) error {
	_, err := AppDB.Exec("INSERT INTO msgs (user_id, chat_id, content) VALUES (?, ?, ?)", user_id, chat_id, content)
	return err
}

func GetChat(key string) (Chat, error) {
	var c Chat
	err := AppDB.QueryRow("SELECT * FROM chats WHERE key=?", key).Scan(&c.ID, &c.TeamID, &c.Name, &c.Key)
	return c, err
}
