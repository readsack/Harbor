package routes

import (
	"encoding/json"
	_ "fmt"
	"harbor/main/db"
	"net/http"
)

type invite struct {
	Email string `json:"email"`
}

type org struct {
	Name string `json:"name"`
}

func sendInvite(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(1048576)

	defer r.Body.Close()
	ctx := r.Context()
	u := ctx.Value("user").(*db.User)
	if !u.OrgID.Valid {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("User Doesn't Have Any Associated Organization!"))
		return
	} else {
		email := r.FormValue("email")
		if email == "" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Not Enough Data"))
			return
		}
		usr, err := db.FindUserByEmail(email)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("No user exists having the given email"))
			return
		} else {
			org, err := db.GetOrg(int(u.OrgID.Int64))
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("User Doesn't Have Any Associated Organization!"))
				return
			} else if org.CeoID != u.ID {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("User Isn't the CEO of his Organization"))
				return
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
	//fmt.Println(inv)
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
	r.ParseMultipartForm(1048576)

	ctx := r.Context()
	u := ctx.Value("user").(*db.User)
	orgName := r.FormValue("name")
	if orgName == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("No Name Provided"))
		return
	}

	id, err := db.CreateOrg(orgName, u.ID)
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

func removeUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	u := ctx.Value("user").(*db.User)
	var data = struct {
		UserID int `json:"user_id"`
	}{}
	defer r.Body.Close()
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, "Bad Data", http.StatusBadRequest)
		return
	}
	org, _ := db.GetOrg(int(u.OrgID.Int64))
	if u.ID != org.CeoID {
		http.Error(w, "Not Authenticated", http.StatusUnauthorized)
		return
	}
	err = db.DeleteUserFromOrg(data.UserID)
	if err != nil {
		http.Error(w, "Internal Error", http.StatusInternalServerError)
	}
}
func deleteOrg(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	u := ctx.Value("user").(*db.User)
	var data = struct {
		OrgID int `json:"org_id"`
	}{}
	defer r.Body.Close()
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, "Bad Data", http.StatusBadRequest)
		return
	}
	org, _ := db.GetOrg(data.OrgID)
	if u.ID != org.CeoID {
		http.Error(w, "Not Authenticated", http.StatusUnauthorized)
		return
	}
	err = db.DeleteOrg(data.OrgID)
	if err != nil {
		http.Error(w, "Internal Error", http.StatusInternalServerError)
	}
}

func SetupOrgRoutes() {
	http.HandleFunc("POST /deleteorg", AuthMiddleware(deleteOrg))
	http.HandleFunc("POST /deleteuser", AuthMiddleware(removeUser))
	http.HandleFunc("POST /sendinvite", AuthMiddleware(sendInvite))
	http.HandleFunc("POST /closeinvite", AuthMiddleware(closeInvite))
	http.HandleFunc("POST /createorg", AuthMiddleware(createOrg))
}
