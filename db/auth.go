package db

import (
	_ "log"
	 "database/sql"
	
)

type User struct{
	ID int
	Email string
	Username string
	Password string
	OrgID sql.NullInt64
}

func FindUserByEmail(email string) (*User, error) {
	u := &User{}

	query := `SELECT id, username, email, pass, org_id FROM users WHERE email = ?`

	err := AppDB.QueryRow(query, email).Scan(
		&u.ID,
		&u.Username,
		&u.Email,
		&u.Password,
		&u.OrgID,
	)
	//log.Println(err)
	if err != nil {
		return nil, err
	}

	return u, nil
}

func FindUserByID(id int) (*User, error)  {
	u := &User{}
	err := AppDB.QueryRow("SELECT * FROM users WHERE id=?", id).Scan(&u.ID, &u.Username, &u.Email, &u.Password, &u.OrgID)
	if err != nil {
		//.Fatal(err)
		return &User{}, err
	}
	return u, nil
}

func CreateUser(u User) error {
	_, err := AppDB.Exec("INSERT INTO users (username, email, pass) VALUES (?, ?, ?)", u.Username, u.Email, u.Password)
	return err
}

