package main

import (
    "html/template"
    "net/http"
    "strings"
    "time"
    "fmt"
    "log"
)

//Person struct to store user information
type Person struct {
	username string
	password string
	name string
	friends []Person
}

//NewsFeed struct to have an easy way to store posts made
type NewsFeed struct {
	username string
	post string
}

//Pseudo-database to store all the different posts. 
var posts = make([]NewsFeed,0)

//Pseudo-database to store all the different users and accounts
var users = make(map[string]Person)

//takes an html file and renders it to the front end
func renderHtml(htmlFile, goFile string, w http.ResponseWriter){
	temp, _ := template.ParseFiles(htmlFile)
	temp.Execute(w,goFile)
}

//sets a cookie in our case using the string for a username and the time in minutes
func setCookie(w http.ResponseWriter, username string, minutes time.Duration) {
		expiration := time.Now().Add(minutes*time.Minute)
		cookie := http.Cookie{Name: "username", Value: username, Expires: expiration}
		http.SetCookie(w, &cookie)
}

// func deleteCookie(cookie )

func signin(w http.ResponseWriter, r *http.Request) {
    if r.Method == "GET" {

    	//Render the page
    	renderHtml("signin.html", "home.go", w)

	} else {
		r.ParseForm()

		//get the username and password from the Form
		username := strings.Join(r.Form["user"], " ")
		password := strings.Join(r.Form["pass"], " ")
		//check if user exists by checking our map "database"
		userPerson, checkUser := users[username]
		if(checkUser == true && userPerson.password == password) {

			//Account credentials checked out. Set a cookie for 5 minutes and redirect to their profile
			setCookie(w,username,5)
			http.Redirect(w,r,"/profile",http.StatusPermanentRedirect)

		} else {

			//User doesn't have an account and is redirected to the home page
			if(checkUser == false) {
				http.Redirect(w,r,"/signin",http.StatusPermanentRedirect)
			} else {

				//password incorrect for the user
				http.Redirect(w,r,"/signin",http.StatusPermanentRedirect)
			}
		}
	}
}

//processes info from profile html file, displaying to screen
func newsfeed(w http.ResponseWriter, r *http.Request) {

	//Render the page
    renderHtml("newsfeed.html", "home.go", w)

    userN, _ := r.Cookie("username")
    nameP := users[userN.Value].name
	fmt.Fprintf(w, "<h1>Welcome %s</h1>", nameP)
}

//processes when the user clicks sign up
func signup(w http.ResponseWriter, r *http.Request) {
	//parses the sign-up html file
    if r.Method == "GET" {

    	//Render the page
	    renderHtml("signup.html", "home.go", w)

	} else {

		//this else statement is processing the sign-up form input

		r.ParseForm()

		//Get the information from the form so we can store it
		username := strings.Join(r.Form["user"], " ")
		password := strings.Join(r.Form["pass"], " ")
		name := strings.Join(r.Form["name"], " ")
		log.Printf(name)
		//stores user info into person struct
		pUser := Person{
			username: username,
			password: password,
			name: name,
			friends: nil,
		}
		//Store in our psuedo-database. Format: username : {password, name}
		users[username] = pUser

		//Set a cookie for 5 minutes and redirect to their profile
		setCookie(w,username, 5)
		http.Redirect(w,r,"/newsfeed",http.StatusPermanentRedirect)
	}
}

//processes when the user clicks the logout link
func logout(w http.ResponseWriter, r *http.Request) {
	//deletes cookie
	c, _ := r.Cookie("username")
	c.Value = ""
	c.Path = "/"
	c.MaxAge = -1
	http.SetCookie(w,c)
	http.Redirect(w,r,"/signin",http.StatusPermanentRedirect)
}

func main() {
    http.HandleFunc("/signin", signin)
    http.HandleFunc("/signup", signup)
    http.HandleFunc("/newsfeed", newsfeed)
    http.HandleFunc("/logout", logout)
    http.ListenAndServe(":8053", nil)
}