package routes

import (
	"encoding/json"
	"harbor/main/db"
	_ "log"
	"net/http"
)

type invite struct {
	Email string `json:"email"`
}

func sendInvite(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	ctx := r.Context()
	u := ctx.Value("user").(*db.User)
	if !u.OrgID.Valid {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("User Doesn't Have Any Associated Organization!"))
	} else {
		var inv invite
		//fmt.Println(r.Body)
		err := json.NewDecoder(r.Body).Decode(&inv)
		if err != nil {
		}
		usr, err := db.FindUserByEmail(inv.Email)
		if err != nil {
			w.Write([]byte("No user exists having the given email"))
		} else {
			//TODO: check if user is org owner.
			v := int(u.OrgID.Int64)
			db.SendInvite(usr.ID, v)
		}
		w.Write([]byte("Invite Sent"))
	}
}

func closeInvite(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	acc := params.Get("accept")
	inv_key := params.Get("invite")
	var accepted bool
	if acc == "0" {
		accepted = false
	} else {
		accepted = true
	}
	err := db.AcceptOrDeclineInvite(inv_key, accepted)
	if(err != nil){
		w.Write([]byte("Invite Doesn't Exist."))
	}
	w.Write([]byte("Invite Closed."))
}

func SetupOrgRoutes() {
	http.HandleFunc("POST /sendinvite", AuthMiddleware(sendInvite))
	http.HandleFunc("POST /closeinvite", AuthMiddleware(closeInvite))
}
