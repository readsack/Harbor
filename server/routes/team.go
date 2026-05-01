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
	TeamID int `json:"team_id"`
}

func createTeam(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(1048576)

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
	err = db.CreateTeam(teamName, u.ID, org.ID)

	fmt.Println(err)
	w.Write([]byte("Team Created Successfully"))
}

func addUserToTeam(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	u := ctx.Value("user").(*db.User)
	var req addUserReq
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil || req.TeamID == 0 || req.UserID == 0 {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Provided Content is not Valid"))
		return
	}
	org, err := db.GetOrg(int(u.OrgID.Int64))
	if err != nil || org.CeoID != u.ID {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("User Isn't CEO of The Organization"))
		return
	}
	team, err := db.GetTeamByID(req.TeamID)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("No Team Associated With ID"))
		return
	}

	err = db.AddIntoTeam(team.ID, req.UserID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error"))
		return
	}
	w.Write([]byte("Added Into Team Successfully"))
}

func createChat(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(1048576)

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

func getChatHistory(w http.ResponseWriter, r *http.Request) {
	chatKey := r.Header.Get("X-Chat-Key")
	if chatKey == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "No Chat Key")
		return
	}
	chatCh, err := db.GetChat(chatKey)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintln(w, "Chat Key Is Not Valid!")
		return
	}
	cData, err := db.GetChatHistory(chatCh)
	if err != nil {
		fmt.Println(err)
	}
	json.NewEncoder(w).Encode(cData)
}

func connectToChat(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	u := ctx.Value("user").(*db.User)
	chatKey := r.Header.Get("X-Chat-Key")
	if chatKey == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "No Chat Key")
		return
	}
	conn, err := chat.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		panic(err)
	}
	chatCh, err := db.GetChat(chatKey)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintln(w, "Chat Key Is Not Valid!")
		return
	}
	err = db.CheckIfUserInTeam(chatCh.TeamID, u.ID)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintln(w, "User Doesn't Belong To Team")
		return
	}
	c := chat.Client{}
	c.Conn = conn
	c.User = u
	c.Chat = chatCh
	c.Send = make(chan db.Message)
	c.HandleClientConnection()

}

func createColumnOrItem(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	u := ctx.Value("user").(*db.User)
	r.ParseMultipartForm(1048576)
	isCol := r.FormValue("isCol")
	//fmt.Println(isCol)
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

func deleteColumnOrItem(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(1048576)
	isCol := r.FormValue("isCol")
	switch isCol {
	case "1":
		valueStr := r.FormValue("col_id")
		colID, err := strconv.Atoi(valueStr)
		if err != nil {
			http.Error(w, "Invalid Column", http.StatusBadRequest)
			return
		}
		err = db.DeleteColumn(colID)
		if err != nil {
			fmt.Println(err)
			http.Error(w, "Invalid Column", http.StatusBadRequest)
			return
		}
	case "0":
		valueStr := r.FormValue("card_id")
		itemID, err := strconv.Atoi(valueStr)
		if err != nil {
			http.Error(w, "Invalid Card", http.StatusBadRequest)
			return
		}
		err = db.DeleteItem(itemID)
		if err != nil {
			http.Error(w, "Invalid Card", http.StatusBadRequest)
			return
		}
	default:
		http.Error(w, "Bad Request Sent", http.StatusBadRequest)
		return
	}
}

func removeUserFromTeam(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	u := ctx.Value("user").(*db.User)
	var data = struct {
		TeamID int `json:"team_id"`
		UserID int `json:"user_id"`
	}{}
	defer r.Body.Close()
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, "Bad Data", http.StatusBadRequest)
		return
	}
	team, _ := db.GetTeamByID(data.TeamID)
	if u.ID != team.SupID {
		http.Error(w, "Not Authenticated", http.StatusUnauthorized)
		return
	}
	err = db.RemoveUserFromTeam(data.TeamID, data.UserID)
	if err != nil {
		http.Error(w, "Internal Error", http.StatusInternalServerError)
	}
}
func deleteTeam(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	u := ctx.Value("user").(*db.User)
	var data = struct {
		TeamID int `json:"team_id"`
	}{}
	defer r.Body.Close()
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, "Bad Data", http.StatusBadRequest)
		return
	}
	team, _ := db.GetTeamByID(data.TeamID)
	if u.ID != team.SupID {
		http.Error(w, "Not Authenticated", http.StatusUnauthorized)
		return
	}
	err = db.DeleteTeam(data.TeamID)
	if err != nil {
		http.Error(w, "Internal Error", http.StatusInternalServerError)
	}
}

func SetupTeamRoutes() {
	http.HandleFunc("POST /createteam", AuthMiddleware(createTeam))
	http.HandleFunc("/msgs", AuthMiddleware(connectToChat))
	http.HandleFunc("POST /createchat", AuthMiddleware(createChat))
	http.HandleFunc("POST /addboard", AuthMiddleware(createColumnOrItem))
	http.HandleFunc("POST /chathistory", AuthMiddleware(getChatHistory))
	http.HandleFunc("POST /delboard", AuthMiddleware(deleteColumnOrItem))
	http.HandleFunc("POST /teamadd", AuthMiddleware(addUserToTeam))
	http.HandleFunc("POST /deleteteam", AuthMiddleware(deleteTeam))
	http.HandleFunc("POST /removeuser", AuthMiddleware(removeUserFromTeam))
}
