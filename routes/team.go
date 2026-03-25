package routes

import (
	"harbor/main/db"
	"net/http"
)

func createTeam(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	u := ctx.Value("user").(*db.User)
}
