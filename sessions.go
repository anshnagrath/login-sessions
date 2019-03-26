package main

import (
	"net/http"

	uuid "github.com/satori/go.uuid"
)

func getUser(w http.ResponseWriter, req *http.Request) User {
	c, err := req.Cookie("session")
	if err != nil {
		sid, _ := uuid.NewV4()
		c = &http.Cookie{
			Name:  "session",
			Value: sid.String(),
			//	secure:true for https
			HttpOnly: true, // this mean that thecookie would be aceesed by http/https request only not by js
		}
		http.SetCookie(w, c)
	}
	//if user already exist
	var u User
	if un, ok := dbSerssions[c.Value]; ok {
		u = dbUsers[un]
	}
	return u
}
func alreadyLoggedIn(req *http.Request) (bool, User) {
	c, _ := req.Cookie("session")
	// if err != nil {
	// 	return false,
	// }
	un := dbSerssions[c.Value]

	u, ok := dbUsers[un]
	return ok, u
}
