package main

import (
    "html/template"
    "net/http"
    "strings"
    "time"
    "fmt"
)

var users = make(map[string][]string)

func handlerSignin(w http.ResponseWriter, r *http.Request) {
    if r.Method == "GET" {
	    temp, _ := template.ParseFiles("signin.html")
	    temp.Execute(w,"home.go")
	} else {
		r.ParseForm()
		username := strings.Join(r.Form["user"], " ")
		password := strings.Join(r.Form["pass"], " ")
		userArr, checkUser := users[username]
		if(checkUser == true && userArr[1] == password) {
			expiration := time.Now().Add(5*time.Minute)
			cookie := http.Cookie{Name: "username", Value: username, Expires: expiration}
			http.SetCookie(w, &cookie)
			http.Redirect(w,r,"/profile",http.StatusPermanentRedirect)
		} else {
			http.Redirect(w,r,"/signin",http.StatusPermanentRedirect)
		}
	}
}

func handlerLogout(w http.ResponseWriter, r *http.Request) {
	c, _ := r.Cookie("username")
	c.Value = ""
	c.Path = "/"
	c.MaxAge = -1
	http.SetCookie(w,c)
	http.Redirect(w,r,"/signin",http.StatusPermanentRedirect)
}

func handlerProfile(w http.ResponseWriter, r *http.Request) {
    temp, _ := template.ParseFiles("profile.html")
    temp.Execute(w,"home.go")
    userN, _ := r.Cookie("username")
    nameP := users[userN.Value][1]
	fmt.Fprintf(w,"Welcome ", nameP)
}

func handlerSignup(w http.ResponseWriter, r *http.Request) {
    if r.Method == "GET" {
	    temp, _ := template.ParseFiles("signup.html")
	    temp.Execute(w,"home.go")
	} else {
		r.ParseForm()
		username := strings.Join(r.Form["user"], " ")
		password := strings.Join(r.Form["pass"], " ")
		name := strings.Join(r.Form["name"], " ")
		users[username] = []string{password, name}
		expiration := time.Now().Add(5*time.Minute)
		cookie := http.Cookie{Name: "username", Value: username, Expires: expiration}
		http.SetCookie(w, &cookie)
		http.Redirect(w,r,"/profile",http.StatusPermanentRedirect)
	}
}

func main() {
    http.HandleFunc("/signin", handlerSignin)
    http.HandleFunc("/signup", handlerSignup)
    http.HandleFunc("/profile", handlerProfile)
    http.HandleFunc("/logout", handlerLogout)
    http.ListenAndServe(":8066", nil)
}