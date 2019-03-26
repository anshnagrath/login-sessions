package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
)

var tpl *template.Template

func init() {
	tpl = template.Must(template.ParseGlob("./Templates/*.html"))
}
func main() {
	http.HandleFunc("/", index)
	http.HandleFunc("/bar", bar)
	http.HandleFunc("/signup", signup)
	http.Handle("/favicon.ico", http.NotFoundHandler())
	log.Fatalln(http.ListenAndServe(":8080", nil))

}

type User struct {
	UserName  string
	FirstName string
	LastName  string
	Password  []byte
}

var dbUsers = make(map[string]User)
var dbSerssions = make(map[string]string)

func index(w http.ResponseWriter, req *http.Request) {
	u := getUser(w, req)
	fmt.Println(u)
	tpl.ExecuteTemplate(w, "index.html", u)
}
func signup(w http.ResponseWriter, req *http.Request) {
	ok, u := alreadyLoggedIn(req)
	if ok {
		tpl.ExecuteTemplate(w, "bar.html", u)
		return
	} else {
		http.Redirect(w, req, "/", http.StatusSeeOther)
		return
	}
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
	if req.Method == http.MethodPost {
		um := req.FormValue("username")
		p := req.FormValue("password")
		fn := req.FormValue("firstname")
		ln := req.FormValue("lastname")
		bs, err := bcrypt.GenerateFromPassword([]byte(p), bcrypt.MinCost)
		if err != nil {
			log.Fatalln(err)
		}
		u = User{
			um, fn, ln, bs,
		}
		if _, ok := dbUsers[um]; ok {
			http.Error(w, "Username already taken", http.StatusForbidden)
		}
		dbSerssions[c.Value] = um
		dbUsers[um] = u
		http.Redirect(w, req, "/", http.StatusSeeOther)
		tpl.ExecuteTemplate(w, "signupForm.html", u)
	}
	tpl.ExecuteTemplate(w, "signupForm.html", u)

}
func bar(w http.ResponseWriter, req *http.Request) {

	u := getUser(w, req)
	if ok, _ := alreadyLoggedIn(req); !ok {
		http.Redirect(w, req, "/", http.StatusSeeOther)
		return
	}
	tpl.ExecuteTemplate(w, "bar.html", u)

}
