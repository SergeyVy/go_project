package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"golang.org/x/crypto/bcrypt"
	"url-shorter/internal/auth"
	"url-shorter/internal/http-server/middleware"
)

type LoginReq struct {
	Email string `json:"email"`
	Pass  string `json:"password"`
}
type LoginResp struct {
	Access  string `json:"access_token"`
	Refresh string `json:"refresh_token"`
}

func Login(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req LoginReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "bad json", http.StatusBadRequest)
			return
		}
		var id int64
		var hash string
		err := db.QueryRowContext(r.Context(),
			`SELECT id, pass_hash FROM users WHERE email=$1`, req.Email).
			Scan(&id, &hash)
		if err != nil || bcrypt.CompareHashAndPassword([]byte(hash), []byte(req.Pass)) != nil {
			http.Error(w, "invalid credentials", http.StatusUnauthorized)
			return
		}
		access, err := auth.NewAccess(id)
		if err != nil {
			http.Error(w, "token", 500)
			return
		}
		refresh, err := auth.NewRefresh(id)
		if err != nil {
			http.Error(w, "token", 500)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(LoginResp{Access: access, Refresh: refresh})
	}
}

func Me(w http.ResponseWriter, r *http.Request) {
	uid, ok := middleware.UserID(r)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	w.Write([]byte(fmt.Sprintf("your user_id: %d", uid)))
}
