package db

type Team struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	OrgID int    `json:"org_id"`
	SupID int    `json:"sup_id"`
}

type Card struct {
	ID       int    `json:"id"`
	ColumnID int    `json:"column_id"`
	Content  string `json:"content"`
	User     `json:"user"`
}

type TeamData struct {
	Team  `json:"team"`
	Chats []Chat `json:"chats"`
}

func GetTeamByID(team_id int) (*Team, error) {
	u := &Team{}
	err := AppDB.QueryRow("SELECT * FROM teams WHERE id=?", team_id).Scan(&u.ID, &u.Name, &u.OrgID, &u.SupID)
	if err != nil {
		//.Fatal(err)
		return &Team{}, err
	}
	return u, nil
}

func handleErr(err error) {
	if err != nil {
		panic(err)
	}
}

func CreateTeam(name string, user_id int, org_id int) error {
	row, err := AppDB.Exec("INSERT INTO teams (name, org_id, sup_id) VALUES (?, ?, ?)", name, org_id, user_id)
	if err != nil {
		team_id, _ := row.LastInsertId()
		err := CreateBoard(int(team_id))
		return err
	}
	return err
}

func AddIntoTeam(team_id int, user_id int) error {
	_, err := AppDB.Exec("INSERT INTO user_team (user_id, team_id, user_role) VALUES (?, ?, ?)", user_id, team_id, 0)
	return err
}

func CheckIfUserInTeam(team_id int, user_id int) error {
	row := AppDB.QueryRow("SELECT * FROM user_team WHERE user_id=$1 AND team_id=$2", user_id, team_id)
	return row.Err()
}
