package routes

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"harbor/main/db"
	"net/http"
)

type OrgInvites struct {
	Invites []db.OrgInvite `json:"invites"`
}

type org_st struct {
	OrgID int `json:"org_id"`
}

func GetUserData(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	u := ctx.Value("user").(*db.User)
	fmt.Println(u.Email)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Println(u)
	json.NewEncoder(w).Encode(u)
}

func GetOrgInvites(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	u := ctx.Value("user").(*db.User)
	fmt.Println(u)
	out := &OrgInvites{
		Invites: make([]db.OrgInvite, 0),
	}
	res, err := db.AppDB.Query("SELECT * FROM org_inv WHERE user_id=?", u.ID)
	if err != nil && err != sql.ErrNoRows {
		w.WriteHeader(401)
		fmt.Fprintf(w, "Error Is: %v", err)
		return
	}
	for res.Next() {
		o := &db.OrgInvite{}
		var key string
		res.Scan(&o.InvID, &o.UserID, &o.OrgID, &key)
		o, err = db.GetInvitebyKey(key)
		if err != nil {
			w.WriteHeader(401)
			fmt.Fprintf(w, "Error Is: %v\n", err.Error())
			return
		}
		out.Invites = append(out.Invites, *o)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(out)
}

func GetOrgData(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var v org_st
	json.NewDecoder(r.Body).Decode(&v)
	org_d, err := db.GetOrgData(v.OrgID)
	if err != nil {
		fmt.Println(err)
	}
	json.NewEncoder(w).Encode(org_d)
}

func HandleApiCalls() {
	http.HandleFunc("GET /api/org", AuthMiddleware(GetOrgData))
	http.HandleFunc("POST /api/user", AuthMiddleware(GetUserData))
	http.HandleFunc("POST /api/user/invites", AuthMiddleware(GetOrgInvites))
}
