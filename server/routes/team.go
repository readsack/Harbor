package routes

import (
	"encoding/json"
	"fmt"
	"harbor/main/chat"
	"harbor/main/db"
	_ "log"
	"net/http"
	"strconv"
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

func createChat(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	valueStr := r.FormValue("team_id")
	name := r.FormValue("name")
	teamID, err := strconv.Atoi(valueStr)
	if err != nil {
		http.Error(w, "Invalid integer", http.StatusBadRequest)
		return
	}

	err = db.CreateChat(name, teamID)
	if err != nil {
		http.Error(w, "Internal Error In Creating Chat", http.StatusInternalServerError)
		return
	}
	w.Write([]byte("Created Chat Successfully"))
}

func connectToChat(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	u := ctx.Value("user").(*db.User)
	r.ParseForm()
	chatKey := r.Header.Get("X-Chat-Key")
	if chatKey == "" {
		return
	}
	conn, err := chat.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		panic(err)
	}
	chatCh, err := db.GetChat(chatKey)
	if err != nil {
		fmt.Fprintln(w, "Chat Key Is Not Valid!")
		return
	}
	err = db.CheckIfUserInTeam(chatCh.TeamID, u.ID)
	if err != nil {
		fmt.Fprintln(w, "User Doesn't Belong To Team")
		return
	}
	c := chat.Client{}
	c.Conn = conn
	c.User = u
	c.Send = make(chan chat.Message)
	c.HandleClientConnection()

}

func createColumnOrItem(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	u := ctx.Value("user").(*db.User)
	r.ParseForm()
	isCol := r.FormValue("isCol")

	switch isCol {
	case "1":
		name := r.FormValue("name")
		valueStr := r.FormValue("team_id")
		teamID, err := strconv.Atoi(valueStr)
		if err != nil {
			http.Error(w, "Invalid integer", http.StatusBadRequest)
			return
		}
		bid, err := db.GetTeamBoardID(teamID)
		if err != nil {
			http.Error(w, "Invalid Team", http.StatusBadRequest)
			return
		}
		db.CreateColumn(bid, name)
	case "0":
		valueStr := r.FormValue("col_id")
		colID, err := strconv.Atoi(valueStr)
		content := r.FormValue("content")
		if err != nil {
			http.Error(w, "Invalid Column", http.StatusBadRequest)
			return
		}
		db.CreateItem(colID, u.ID, content)
	default:
		http.Error(w, "Bad Request Sent", http.StatusBadRequest)
		return
	}
}

func SetupTeamRoutes() {
	http.HandleFunc("POST /createteam", AuthMiddleware(createTeam))
	http.HandleFunc("GET /msgs", AuthMiddleware(connectToChat))
	http.HandleFunc("POST /createchat", AuthMiddleware(createChat))
	http.HandleFunc("POST /addboard", AuthMiddleware(createColumnOrItem))
}
