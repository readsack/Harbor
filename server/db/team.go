package db

type Team struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	OrgID int    `json:"org_id"`
	SupID int    `json:"sup_id"`
}

type TeamData struct {
	Team    `json:"team"`
	Chats   []Chat `json:"chats"`
	Members []User `json:"members"`
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
		return err
	}
	team_id, _ := row.LastInsertId()
	err = AddIntoTeam(int(team_id), user_id)
	err = CreateBoard(int(team_id))
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

func GetTeamData(team_id int) (*TeamData, error) {
	var teamData TeamData
	team, err := GetTeamByID(team_id)
	if err != nil {
		return &teamData, err
	}
	teamData.Team = *team
	chats, err := AppDB.Query("SELECT * FROM chats WHERE team_id=?", team_id)
	for chats.Next() {
		var c Chat
		err = chats.Scan(&c.ID, &c.TeamID, &c.Name, &c.Key)
		if err != nil {
			return &teamData, err
		}
		teamData.Chats = append(teamData.Chats, c)
	}
	users, err := AppDB.Query("SELECT user_id FROM user_team WHERE team_id=?", teamData.Team.ID)
	for users.Next() {
		var user_id int
		err := users.Scan(&user_id)
		if err != nil {
			return &teamData, err
		}
		u, err := FindUserByID(user_id)
		if err != nil {
			return &teamData, err
		}
		u.Password = ""
		//u.ID = -1
		teamData.Members = append(teamData.Members, *u)
	}
	return &teamData, nil
}

func DeleteTeam(team_id int) error {
	_, err := AppDB.Exec("DELETE FROM teams WHERE id=?", team_id)
	return err
}

func RemoveUserFromTeam(team_id int, user_id int) error {
	_, err := AppDB.Exec("DELETE FROM user_team WHERE user_id=? AND team_id=?", user_id, team_id)
	return err
}
