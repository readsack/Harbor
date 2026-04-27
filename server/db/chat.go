package db

import "crypto/rand"

type Chat struct {
	Name   string `json:"name"`
	ID     int    `json:"id"`
	TeamID int    `json:"team_id"`
	Key    string `json:"key"`
}

type Message struct {
	Sender  string `json:"sender"`
	Content string `json:"content"`
	ID      int    `json:"id"`
}

type ChatData struct {
	Chat     `json:"chat"`
	Messages []Message `json:"messages"`
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

func GetChatHistory(c Chat) (ChatData, error) {
	var chatData ChatData
	chatData.Chat = c
	chatData.Messages = make([]Message, 0)
	msgs, err := AppDB.Query("SELECT msgs.id, msgs.content, users.username FROM msgs INNER JOIN users ON users.id=msgs.user_id WHERE msgs.chat_id=?", c.ID)
	if err != nil {
		return chatData, err
	}
	for msgs.Next() {
		var m Message
		err = msgs.Scan(m.ID, m.Content, m.Sender)
		if err != nil {
			return chatData, err
		}
		chatData.Messages = append(chatData.Messages, m)
	}
	return chatData, nil
}
