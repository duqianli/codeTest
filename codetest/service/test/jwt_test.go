package test

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"testing"
)

type User struct {
	Identity string `json:"identity"`
	Name     string `json:"name"`
	jwt.RegisteredClaims
}

var key = []byte("gon_gorm_oj")

func TestClaim(t *testing.T) {

	userClaim := &User{
		Identity:         "user_1",
		Name:             "book",
		RegisteredClaims: jwt.RegisteredClaims{},
	}
	claim := jwt.NewWithClaims(jwt.SigningMethodHS256, userClaim)
	token, err := claim.SignedString(key)
	if err != nil {
		fmt.Println("error:", err)
	}
	fmt.Println("token:", token)
}

func TestPClaim(t *testing.T) {
	tokenString := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZGVudGl0eSI6InVzZXJfMSIsIm5hbWUiOiJib29rIn0.ruQU8Xo4WBZcR0YZC-qK-zLzru142KG86LAnhGRGuTU"
	claim := new(User)
	claimsss, err := jwt.ParseWithClaims(tokenString, claim, func(token *jwt.Token) (interface{}, error) {
		return key, nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if claimsss.Valid {
		fmt.Println(claim)
	}
	return
}
