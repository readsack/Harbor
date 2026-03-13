package routes

import (
	"encoding/json"
	"harbor/main/db"
	"log"
	"net/http"
)

type invite struct {
	Email string `json:"email"`
}

type org struct {
	Name string `json: "name"`
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
		err := json.NewDecoder(r.Body).Decode(&inv)
		if err != nil {
			w.WriteHeader(http.StatusUnsupportedMediaType)
			w.Write([]byte("Provided Content is not JSON"))
			log.Fatal("Can't Decode JSON")
			return
		}
		usr, err := db.FindUserByEmail(inv.Email)
		if err != nil {
			w.Write([]byte("No user exists having the given email"))
		} else {
			org, err := db.GetOrg(int(u.OrgID.Int64))
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("User Doesn't Have Any Associated Organization!"))
			} else if org.CeoID != u.ID {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("User Isn't the CEO of his Organization"))
			} else {
				v := int(u.OrgID.Int64)
				db.SendInvite(usr.ID, v)
			}
		}
		w.Write([]byte("Invite Sent"))
	}
}

func closeInvite(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	acc := params.Get("accept")
	inv_key := params.Get("invite")
	var accepted bool
	ctx := r.Context()
	u := ctx.Value("user").(*db.User)
	switch acc {
	case "0":
		accepted = false
	case "1":
		accepted = true
	}
	inv, err := db.GetInvitebyKey(inv_key)
	if inv.UserID != u.ID {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invite Doesn't Belong To User"))
		return
	}
	err = db.AcceptOrDeclineInvite(inv_key, accepted)
	if err != nil {
		w.Write([]byte("Invite Doesn't Exist."))
	}
	w.Write([]byte("Invite Closed."))
}

func createOrg(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	u := ctx.Value("user").(*db.User)
	var Org org
	err := json.NewDecoder(r.Body).Decode(&Org)
	if err != nil {
		w.WriteHeader(http.StatusUnsupportedMediaType)
		w.Write([]byte("Provided Content is not JSON"))
		log.Fatal("Can't Decode JSON")
		return
	}
	id, err := db.CreateOrg(Org.Name, u.ID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Could not Create Org"))
		return
	}
	org_id, err := id.LastInsertId()
	db.SetUserOrg(u.ID, int(org_id))
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("New Organization Was Created"))

}

func SetupOrgRoutes() {
	http.HandleFunc("POST /sendinvite", AuthMiddleware(sendInvite))
	http.HandleFunc("POST /closeinvite", AuthMiddleware(closeInvite))
	http.HandleFunc("POST /createorg", AuthMiddleware(createOrg))
}
