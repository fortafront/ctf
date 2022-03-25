package main

import (
	"github.com/golang-jwt/jwt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"time"
)

type CustomClaims struct {
	Authorized     bool
	Signature      string
	UserID         string
	Expires        int64
	StandardClaims jwt.StandardClaims
}

func (c CustomClaims) Valid() error {
	//TODO implement me
	return nil
}

func CreateToken() (string, error) {
	var err error
	//Creating Access Token
	os.Setenv("ACCESS_SECRET", "Iamhwd8gtA3zwj4G78pT") //this should be in an env file
	claims := CustomClaims{
		Authorized: false,
		Signature:  "Iamhwd8gtA3zwj4G78pT",
		UserID:     "ritsecuser",
		Expires:    time.Now().Add(time.Minute * 15).Unix(),
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute * 15).Unix(),
			Issuer:    "RITSEC",
		},
	}

	at := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := at.SignedString([]byte(os.Getenv("ACCESS_SECRET")))
	if err != nil {
		return "", err
	}
	return token, nil
}
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	token, err := CreateToken()
	if r.Header.Get("token") != "" {
		a := CustomClaims{}
		_, err := jwt.ParseWithClaims(
			r.Header.Get("token"),
			&a,
			func(token *jwt.Token) (interface{}, error) {
				return []byte("Iamhwd8gtA3zwj4G78pT"), nil
			},
		)
		if err != nil {
			w.Write([]byte("Error parsing JWT"))
			w.Write([]byte(err.Error()))
		} else {
			if a.Authorized && a.UserID == "ritsecuser" && a.Expires > time.Now().Unix() {
				w.Write([]byte("RS{DREAM_IN_JWT}"))
			} else {
				w.Write([]byte("Invalid token"))
			}
		}

	} else {
		w.Write([]byte("missing header token  \n"))
		if err != nil {
			w.Write([]byte("Error"))
		} else {
			w.Write([]byte(token))
		}
	}

}
func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", HomeHandler)

	srv := &http.Server{
		Handler: r,
		Addr:    "127.0.0.1:8000",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
