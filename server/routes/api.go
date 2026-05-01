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
	//fmt.Println(u.Email)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	//fmt.Println(u)
	json.NewEncoder(w).Encode(u)
}

func GetOrgInvites(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	u := ctx.Value("user").(*db.User)
	//fmt.Println(u)
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
	//.Println(out)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(out)
}

func GetOrgData(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	u := ctx.Value("user").(*db.User)

	org_d, err := db.GetOrgData(int(u.OrgID.Int64), u.ID)
	if err != nil {
		fmt.Println(err)
	}
	json.NewEncoder(w).Encode(org_d)
}

func GetBoardData(w http.ResponseWriter, r *http.Request) {

	v := struct {
		TeamID int `json:"team_id"`
	}{
		-1,
	}
	err := json.NewDecoder(r.Body).Decode(&v)
	if err != nil {
		fmt.Println(err)
	}
	val, err := db.GetBoard(v.TeamID)
	if err != nil {
		fmt.Println(err)
	}
	json.NewEncoder(w).Encode(val)
}

func GetTeamData(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	u := ctx.Value("user").(*db.User)
	teamID := struct {
		TeamID int `json:"team_id"`
	}{
		-1,
	}
	json.NewDecoder(r.Body).Decode(&teamID)
	if teamID.TeamID == -1 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "Invalid Team ID")
		return
	}
	if db.CheckIfUserInTeam(teamID.TeamID, u.ID) != nil {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintln(w, "User Doesn't Belong To Team")
		return
	}
	teamData, err := db.GetTeamData(teamID.TeamID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "Server Error")
		return
	}
	json.NewEncoder(w).Encode(*teamData)
}

func HandleApiCalls() {
	http.HandleFunc("POST /api/org", AuthMiddleware(GetOrgData))
	http.HandleFunc("POST /api/user", AuthMiddleware(GetUserData))
	http.HandleFunc("POST /api/user/invites", AuthMiddleware(GetOrgInvites))
	http.HandleFunc("POST /api/team/board", AuthMiddleware(GetBoardData))
	http.HandleFunc("POST /api/team", AuthMiddleware(GetTeamData))
}
