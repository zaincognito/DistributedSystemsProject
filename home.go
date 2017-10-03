package main

import (
    "html/template"
    "net/http"
)

func handlerSignin(w http.ResponseWriter, r *http.Request) {
    temp, _ := template.ParseFiles("signin.html")
    temp.Execute(w,"signin.go")
}

func handlerSignup(w http.ResponseWriter, r *http.Request) {
    temp, _ := template.ParseFiles("signup.html")
    temp.Execute(w,"signin.go")
}

func main() {
    http.HandleFunc("/signin", handlerSignin)
    http.HandleFunc("/signup", handlerSignup)
    http.ListenAndServe(":8080", nil)
}