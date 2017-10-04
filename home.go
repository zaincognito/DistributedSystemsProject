package main

import (
    "html/template"
    "net/http"
)

var users = make(map[string]string)

func handlerSignin(w http.ResponseWriter, r *http.Request) {
    temp, _ := template.ParseFiles("signin.html")
    temp.Execute(w,"home.go")
}

func handlerSignup(w http.ResponseWriter, r *http.Request) {
    if r.Method == "get" {
	    temp, _ := template.ParseFiles("signup.html")
	    temp.Execute(w,"home.go")
	}
	else {
		r.ParseForm()
	}
}

func main() {
    http.HandleFunc("/signin", handlerSignin)
    http.HandleFunc("/signup", handlerSignup)
    http.ListenAndServe(":8080", nil)
}