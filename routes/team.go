package routes

import (
	"encoding/json"
	"harbor/main/db"
	_ "log"
	"net/http"

	"github.com/gorilla/websocket"
)

type createTeamReq struct {
	Name string `json:"name"`
}

type addUserReq struct {
	UserID int `json:"user_id"`
	TeamID int `json:"org_id"`
}

func createTeam(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	u := ctx.Value("user").(*db.User)
	r.ParseForm()
	teamName := r.FormValue("name")
	if teamName == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Request is not Valid"))
		return
	}
	org, err := db.GetOrg(int(u.OrgID.Int64))
	if err != nil || org.CeoID != u.ID {
		w.Write([]byte("User Isn't CEO of The Organization"))
		return
	}
	db.CreateTeam(teamName, u.ID, org.ID)
	w.Write([]byte("Team Created Successfully"))
}

func addUserToTeam(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	u := ctx.Value("user").(*db.User)
	var req addUserReq
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil || req.TeamID == 0 || req.UserID == 0 {
		w.Write([]byte("Provided Content is not Valid"))
		return
	}
	org, err := db.GetOrg(int(u.OrgID.Int64))
	if err != nil || org.CeoID != u.ID {
		w.Write([]byte("User Isn't CEO of The Organization"))
		return
	}
	team, err := db.GetTeamByID(req.TeamID)
	if err != nil {
		w.Write([]byte("No Team Associated With ID"))
		return
	}
	err = db.AddIntoTeam(team.ID, req.UserID)
	if err != nil {
		w.Write([]byte("Internal Server Error"))
		return
	}
	w.Write([]byte("Added Into Team Successfully"))
}

func connectToChat(w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		panic(err)
	}
	conn.WriteMessage(websocket.TextMessage, []byte("Hi From Server!"))
	_, msg, err := conn.ReadMessage()
	println(string(msg))
}

func SetupTeamRoutes() {
	http.HandleFunc("POST /createteam", AuthMiddleware(createTeam))
	http.HandleFunc("GET /msgs", connectToChat)
}
