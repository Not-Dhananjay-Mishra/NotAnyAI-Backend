package router

import (
	"server/utils"

	"github.com/dgrijalva/jwt-go"
)

type JWTCredientials struct {
	Username string
	jwt.StandardClaims
}

func WSSGateKeeper(tkn string) string {

	if tkn == "" {
		return "error"
	}

	var claim JWTCredientials
	token, err := jwt.ParseWithClaims(tkn, &claim, func(token *jwt.Token) (interface{}, error) {
		return []byte(utils.JWT_SECRET), nil
	})

	if err != nil && !token.Valid {
		return "error"
	}
	//log.Println("User pass through gate: ", claim.Username)
	return claim.Username
}
