package db

import (
	"crypto/rand"
	"database/sql"
	_ "database/sql"
	"log"
	_ "log"
)

type OrgInvite struct {
	org_id   int
	user_id  int
	inv_key  int
	inv_id   int
	org_name string
	username string
}

func SendInvite(user_id int, org_id int) error {
	_, err := AppDB.Exec("INSERT INTO org_inv (user_id, org_id, key) VALUES (?, ?, ?)", user_id, org_id, rand.Text())
	return err
}

func AcceptOrDeclineInvite(invite_key string, accept bool) {
	var inv_id, org_id, user_id int
	err := AppDB.QueryRow("SELECT (id, user_id, org_id) FROM org_inv WHERE key=?", invite_key).Scan(&inv_id, &user_id, &org_id)
	if err != sql.ErrNoRows && err != nil {
		log.Fatal(err)
	}
	_, err = AppDB.Exec("DELETE FROM org_inv WHERE id=?", inv_id)
	if err != nil {
		log.Fatal(err)
	}
	if accept {
		_, er := AppDB.Exec("UPDATE users SET org_id=$1 WHERE id=$2", org_id, user_id)
		if er != nil {
			log.Fatal(er)
		}
	}

}
