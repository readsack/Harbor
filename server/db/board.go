package db

type Item struct {
	Content      string
	ID           int
	ColumnID     int
	LastWriterID int
}

type Column struct {
	Items   []Item
	ID      int
	BoardID int
	Name    string
}

type Board struct {
	ID      int
	TeamID  int
	Columns []Column
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
	err := AppDB.QueryRow("SELECT id FROM columns WHERE team_id=?", team_id).Scan(&col_id)
	return col_id, err
}

func GetBoard(board_id int) {

}
