package auth

import (
	"encoding/json"
	"log"
	"net/http"
	mongodb "server/database/MongoDB"
	"server/utils"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type Credientials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type JWTCredientials struct {
	Username string
	jwt.StandardClaims
}

func Login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var data Credientials
	json.NewDecoder(r.Body).Decode(&data)
	log.Println("Request recived: ", data.Username)
	//fmt.Println(time.Now(), "User Login: ", data.Username)
	pass, err := mongodb.GetUserPassword(data.Username)
	if err != nil {
		log.Println("Request no user found with: ", data.Username)
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"status": "error1"})
		return
	}
	//database.TestPrintAllUser()

	if pass != data.Password {
		log.Println("Request failed password not matched: ", data)
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"status": "error2"})
		return
	}

	var claims = &JWTCredientials{
		Username: data.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
		},
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(utils.JWT_SECRET))
	if err != nil {
		log.Println("Request failed internal error: ", data)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"status": "error3"})
		return
	}
	log.Println("Request passed: ", data)
	w.Header().Set("Authorization", token)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"token": token, "status": "done"})
}
