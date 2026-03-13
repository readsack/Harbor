package db

import (
	"crypto/rand"
	"database/sql"
	_ "database/sql"
	"log"
)

type OrgInvite struct {
	OrgID    int
	UserID   int
	InvKey   int
	InvID    int
	OrgName  string
	Username string
}

type Org struct {
	ID    int
	Name  string
	CeoID int
}

func GetOrg(org_id int) (*Org, error) {
	u := &Org{}
	err := AppDB.QueryRow("SELECT * FROM orgs WHERE id=?", org_id).Scan(&u.ID, &u.Name, &u.CeoID)
	if err != nil {
		//.Fatal(err)
		return &Org{}, err
	}
	return u, nil
}

func CreateOrg(name string, ceo_id int) (sql.Result, error) {
	id, err := AppDB.Exec("INSERT INTO orgs (name, ceo_id) VALUES (?, ?)", name, ceo_id)
	id.LastInsertId()
	return id, err
}

func GetInvitebyKey(inv_key string) (*OrgInvite, error) {
	u := &OrgInvite{}
	err := AppDB.QueryRow(`SELECT org_inv.org_id, org_inv.user_id, org_inv.key, org_inv.id, orgs.name, users.username 
							FROM org_inv 
							INNER JOIN users ON users.id=org_inv.user_id 
							INNER JOIN orgs ON orgs.id=org_inv.org_id 
							WHERE org_inv.key=?`).Scan(
		&u.OrgID, &u.UserID, &u.InvKey, &u.InvID, &u.OrgName, &u.Username)
	if err != nil {
		return &OrgInvite{}, err
	}
	return u, nil
}

func SendInvite(user_id int, org_id int) error {
	_, err := AppDB.Exec("INSERT INTO org_inv (user_id, org_id, key) VALUES (?, ?, ?)", user_id, org_id, rand.Text())
	return err
}

func SetUserOrg(user_id int, org_id int) error {
	_, er := AppDB.Exec("UPDATE users SET org_id=$1 WHERE id=$2", org_id, user_id)
	if er != nil {
		return er
	}
	return nil
}

func AcceptOrDeclineInvite(invite_key string, accept bool) error {
	var inv_id, org_id, user_id int
	err := AppDB.QueryRow("SELECT (id, user_id, org_id) FROM org_inv WHERE key=?", invite_key).Scan(&inv_id, &user_id, &org_id)
	if err != sql.ErrNoRows && err != nil {
		log.Fatal(err)
	}
	if err == sql.ErrNoRows {
		return err
	}
	_, err = AppDB.Exec("DELETE FROM org_inv WHERE id=?", inv_id)
	if err != nil {
		log.Fatal(err)
	}
	if accept {
		er := SetUserOrg(user_id, org_id)
		if er != nil {
			log.Fatal(er)
		}
	}
	return nil
}
