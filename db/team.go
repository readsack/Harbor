package db

type Team struct {
	ID    int
	Name  string
	OrgID int
	SupID int
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
	_, err := AppDB.Exec("INSERT INTO teams (name, org_id, sup_id) VALUES (?, ?, ?)", name, org_id, user_id)
	return err
}

func AddIntoTeam(team_id int, user_id int) error {
	_, err := AppDB.Exec("INSERT INTO user_team (user_id, team_id, user_role) VALUES (?, ?, ?)", user_id, team_id, 0)
	return err
}
