package auth

import (
	"encoding/json"
	"log"
	"net/http"
	mongodb "server/database/MongoDB"
	"server/utils"

	"github.com/dgrijalva/jwt-go"
)

type ValidateResponse struct {
	Username       string `json:"username"`
	Status         string `json:"status"`
	SitecraftLimit int    `json:"sitecraftlimit"`
}

func GateKeeper(w http.ResponseWriter, r *http.Request) {
	tkn := r.Header.Get("Authorization")
	if tkn == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	var claim JWTCredientials
	w.Header().Set("Content-Type", "application/json")
	token, err := jwt.ParseWithClaims(tkn, &claim, func(token *jwt.Token) (interface{}, error) {
		return []byte(utils.JWT_SECRET), nil
	})

	if err != nil && !token.Valid {
		log.Println(err)
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"status": "error"})
		return
	}
	limit, _ := mongodb.GetUserLimit(claim.Username)
	log.Println("User pass through gate: ", claim.Username)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ValidateResponse{
		Username:       claim.Username,
		Status:         "done",
		SitecraftLimit: limit,
	})
}
