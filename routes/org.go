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
		w.Write([]byte("You Don't Run Any Organization!"))
	} else {
		var inv invite
		//fmt.Println(r.Body)
		err := json.NewDecoder(r.Body).Decode(&inv)
		if err != nil {

		}
		//fmt.Printf("jsonContent: %v\n", inv)
		w.Write([]byte("Invite Sent"))
	}
}

func SetupOrgRoutes() {
	http.HandleFunc("POST /sendinvite", AuthMiddleware(sendInvite))
}
