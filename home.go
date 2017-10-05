package main

import (
    "html/template"
    "net/http"
    "strings"
    "time"
    "fmt"
<<<<<<< HEAD
=======
    // "log"
>>>>>>> 64eab256242e828b3f70b0342a20274f85a3e0fb
)



//Pseudo-databse to store all the different users and accounts
var users = make(map[string][]string)

func renderHtml(htmlFile, goFile string, w http.ResponseWriter){
	temp, _ := template.ParseFiles(htmlFile)
	temp.Execute(w,goFile)
}

// func setCookie(username string, ) http.Cookie {

// }

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
		userArr, checkUser := users[username]
		if(checkUser == true && userArr[1] == password) {

			//Account credentials checked out. Set a cookie for 5 minutes and redirect to their profile
			expiration := time.Now().Add(5*time.Minute)
			cookie := http.Cookie{Name: "username", Value: username, Expires: expiration}
			http.SetCookie(w, &cookie)
			http.Redirect(w,r,"/profile",http.StatusPermanentRedirect)

		} else {

			//User doesn't have an account and is redirected to the home page
			http.Redirect(w,r,"/signin",http.StatusPermanentRedirect)

		}
	}
}

<<<<<<< HEAD
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
=======
func profile(w http.ResponseWriter, r *http.Request) {

	//Render the page
    renderHtml("profile.html", "home.go", w)

>>>>>>> 64eab256242e828b3f70b0342a20274f85a3e0fb
    userN, _ := r.Cookie("username")
    nameP := users[userN.Value][1]
	fmt.Fprintf(w,"Welcome ", nameP)
}

func signup(w http.ResponseWriter, r *http.Request) {
    if r.Method == "GET" {

    	//Render the page
	    renderHtml("signup.html", "home.go", w)

	} else {
		r.ParseForm()

		//Get the information from the form so we can store it
		username := strings.Join(r.Form["user"], " ")
		password := strings.Join(r.Form["pass"], " ")
		name := strings.Join(r.Form["name"], " ")

		//Store in our psuedo-database. Format: username : {password, name}
		users[username] = []string{password, name}

		//Set a cookie for 5 minutes and redirect to their profile
		expiration := time.Now().Add(5*time.Minute)
		cookie := http.Cookie{Name: "username", Value: username, Expires: expiration}
		http.SetCookie(w, &cookie)
		http.Redirect(w,r,"/profile",http.StatusPermanentRedirect)
	}
}

func logout(w http.ResponseWriter, r *http.Request) {
	c, _ := r.Cookie("username")
	c.Value = ""
	c.Path = "/"
	c.MaxAge = -1
	http.SetCookie(w,c)
	http.Redirect(w,r,"/signin",http.StatusPermanentRedirect)
}

func main() {
<<<<<<< HEAD
    http.HandleFunc("/signin", handlerSignin)
    http.HandleFunc("/signup", handlerSignup)
    http.HandleFunc("/profile", handlerProfile)
    http.HandleFunc("/logout", handlerLogout)
    http.ListenAndServe(":8066", nil)
=======
    http.HandleFunc("/signin", signin)
    http.HandleFunc("/signup", signup)
    http.HandleFunc("/profile", profile)
    http.HandleFunc("/logout", logout)
    http.ListenAndServe(":8080", nil)
>>>>>>> 64eab256242e828b3f70b0342a20274f85a3e0fb
}