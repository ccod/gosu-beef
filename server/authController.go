package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	bnet "github.com/ccod/go-bnet"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gorilla/context"
	"golang.org/x/oauth2"
)

func (s *Server) jwtMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// for possible token errors, I should return an http.Error handler | later
		reqToken := r.Header.Get("Authorization")
		reqToken = strings.Split(reqToken, "Bearer ")[1]

		token, err := jwt.ParseWithClaims(reqToken, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(s.env.jwtSecret), nil
		})

		claims, ok := token.Claims.(*jwt.StandardClaims)
		if !(ok && token.Valid) {
			fmt.Printf("error with Parse with Claims: %s", err)
			w.Write([]byte("{\"failure\":true}"))
			return
		}

		accountID, err := strconv.Atoi(claims.Id)
		if err != nil {
			fmt.Printf("Atoi call failed: %s", err)
			w.Write([]byte("{\"failure\":true}"))
			return
		}

		context.Set(r, JAccID, accountID)

		next(w, r)
	})
}

func (s *Server) handlBnetLogin(w http.ResponseWriter, r *http.Request) {
	url := s.oauthCfg.AuthCodeURL(s.oauthStateString, oauth2.AccessTypeOnline)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

// TODO: first print a jwt with user info, later make it a redirect to webapp with jwt in the url or something
func (s *Server) handleBnetCallback(w http.ResponseWriter, r *http.Request) {
	state := r.FormValue("state")
	if state != s.oauthStateString {
		fmt.Printf("invalid oauth state expected '%s', go '%s'\n", s.oauthStateString, state)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	code := r.FormValue("code")
	token, err := s.oauthCfg.Exchange(oauth2.NoContext, code)
	if err != nil {
		fmt.Printf("oauthConf.Exchange() failed with '%s'\n", err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	oauthClient := s.oauthCfg.Client(oauth2.NoContext, token)
	client := bnet.NewClient("us", oauthClient)
	user, _, err := client.UserInfo()
	if err != nil {
		fmt.Printf("client.Profile().SC2() failed with '%s'\n", err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	claims := &jwt.StandardClaims{
		Id:     strconv.Itoa(user.ID),
		Issuer: "gosu-beef",
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := jwtToken.SignedString([]byte(s.env.jwtSecret))
	if err != nil {
		fmt.Printf("jwt signing failed: %s", err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	fmt.Printf("UserInfo is: %v\n", user)
	fmt.Printf("Jwt token: %s", tokenString)
	http.Redirect(w, r, s.env.clientDomain+"/callback#"+tokenString, http.StatusTemporaryRedirect)
}
