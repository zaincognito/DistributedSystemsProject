package main

import (
    "html/template"
    "net/http"
    "strings"
    "time"
    "fmt"
    "log"
)

//***If something is a log.Println(), that is for debugging and only shows in the terminal, not on web app.

//Person struct to store user information
type Person struct {
	username string
	password string
	name string
	friends []Person

}

//Post struct to have an easy way to store posts made
type Post struct {
	username string
	ID int
	post string
}

//Pseudo-database to store all the different posts. 
var posts = make([]Post,0)

//Pseudo-database to store all the different users and accounts
var users = make(map[string]*Person)

//Gives each post an individual ID
var postCount int

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

//takes in a Person array, which is a friends list and checks if a username is within the friends
//Basically checks if someone is a friend
func isFriend(friends []Person, printUser string) bool {
	for i:=0; i < len(friends); i++ {
		if friends[i].username == printUser {
			return true
		}
	}
	return false
}

//prints all the posts depending on newsfeed or profile
//prints the user's posts only if profile
//if in newsfeed, prints all user's posts and all friends posts as well
func printPosts(w http.ResponseWriter, username string, newsFeed bool, r *http.Request) {
	var i int
	for i=1; i < len(posts); i++ {
		printPost := posts[i]
		printUser := printPost.username
		theUser := users[username]
		myFriends := theUser.friends
		thePost := printPost.post

		if(len(thePost) >= 1) {
			if newsFeed == true {
				if isFriend(myFriends,printUser) || (printUser == username) {
					fmt.Fprintf(w, "<h2>%s: \n %s \n</h2>", printUser, thePost)
				}
			} else {
				if printUser == username {
					fmt.Fprintf(w, "<h2>%s: \n %s \n</h2>",
						printUser, thePost)
				}
			}
		}
	}
}

//processes the signin page, which has a way to signup as well
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
			http.Redirect(w,r,"/newsfeed",http.StatusPermanentRedirect)

		} else {

			//User doesn't have an account or incorrect password -> is redirected to the signin page
			if(checkUser == false) {
				http.Redirect(w,r,"/signin",302)
			} else {

				//password incorrect for the user
				http.Redirect(w,r,"/signin",302)
			}
		}
	}
}

//processes info from profile html file, displaying to screen
func newsfeed(w http.ResponseWriter, r *http.Request) {

	//Render the page
	userN, _ := r.Cookie("username")
    username := users[userN.Value].username
    renderHtml("navbar.html", "home.go", w)

    //gets user's name to display
    nameP := users[userN.Value].name
	fmt.Fprintf(w, "<h1>Welcome %s!</h1>", nameP)

	//form to search for a username to add to their friend's list (follow)
	fmt.Fprintf(w, "<form method = \"post\"> Search for Friends by Username:   "+ 
		"<input type=\"text\" name=\"search\"/>"+ 
		"<input type=\"submit\" value=\"Follow\"/></br></br>"+ 
		"</form>")

	//process the search of a username to add as friend
	if r.Method == "POST"{
		r.ParseForm()

		//Get what the user inputted as "friend's" username
		prospectFriendUsername := r.PostFormValue("search")

		//Get own username from cookie
		cookie, err := r.Cookie("username")
		if err != nil{
			log.Fatal(err)
		}
		username := cookie.Value

		//Check if person user is trying to follow exists
		if prosFriend := users[prospectFriendUsername]; prosFriend != nil{

			//Make sure user isnt trying to follow themselves
			if prospectFriendUsername == username{
				fmt.Fprintf(w, "<h3>You can't follow yourself!</h3>")
			} else{

				//Put other user into self's friends list
				self := users[username]

				//check if already your friend or not
				exists := false
				for _, friend := range self.friends{
					if friend.username == prospectFriendUsername{
						exists = true
					}
				}

				//if the person is not your friend, add them to your friend's list
				if !exists {
					self.friends = append(self.friends, *prosFriend)
					fmt.Fprintf(w, "<h3>You are now following %s</h3>", prospectFriendUsername)

				} else{

					//if the person is your friend, display a message to let them know
					fmt.Fprintf(w, "<h3>You are already following %s</h3>", prospectFriendUsername)
				}

			}
		} else{

			//User doesnt exist or hasnt signed up
			if prosFriend == nil && (len(prospectFriendUsername) >= 1) {
				fmt.Fprintf(w, "<h3>User doesn't exist</h3>")
			}
		}
	}

	//shows the user's newsfeed
	renderHtml("newsfeed.html", "home.go", w)

	//this processes' the newsfeed where you are allowed to post things as well
	if r.Method == "POST" {
		r.ParseForm()

		//create a post from the tweet
		curPost := strings.Join(r.Form["aPost"], " ")
		postCount++

		//create a post struct with the username, content, ID
		aPost := Post{
			username: username,
			ID: postCount,
			post: curPost,
		}

		//add to the posts pseudo database
		posts = append(posts,aPost)
	}

	//if there is more than one post (there will always be), print
	if(len(posts) > 1) {
		printPosts(w, username, true, r)
	}
}

//the user's profile, which has friends and their own posts
func profile(w http.ResponseWriter, r *http.Request) {

	//Render the page
	userN, _ := r.Cookie("username")
    username := users[userN.Value].username
    renderHtml("navbar.html", "home.go", w)	

    nameP := users[userN.Value].name
	fmt.Fprintf(w, "<h1>%s</h1>", nameP)

	fmt.Fprintf(w, "<h3>Friends:</h3>")

	//Get own username from cookie
	cookie, err := r.Cookie("username")
	if err != nil{
		log.Fatal(err)
	}
	username = cookie.Value
	self := *(users[username])
	log.Println(self.friends)

	//prints all friends on top of profile

	if len(self.friends) == 0 {
		fmt.Fprintf(w, "<h3>No friends</h3>")
	}

	for count:=1; count <= len(self.friends); count++{
		friend := self.friends[count-1]
		fmt.Fprintf(w, "%d. Username: %s  |  Name: %s </br>", count, friend.username, friend.name)
	}

	//Form in order to search for a friend and remove them (unfollow)
	fmt.Fprintf(w, "<form method = \"post\"> Remove friends by Username:   "+ 
	"<input type=\"text\" name=\"removed\"/>"+ 
	"<input type=\"submit\" value=\"Remove\"/></br></br>"+ 
	"</form>")

	if r.Method == "POST" {
		r.ParseForm()

		//Get username that user wants to unfollow
		removeFriend := r.PostFormValue("removed")

		//Get own username from cookie
		cookie, err := r.Cookie("username")
		if err != nil{
			log.Fatal(err)
		}
		username := cookie.Value


		//Check if person user is trying to unfollow exists
		if remFriend := users[removeFriend]; remFriend != nil{

			//Make sure user isnt trying to unfollow themself
			if removeFriend == username{
				fmt.Fprintf(w, "<h3>You can't unfollow yourself!</h3>")
			} else{

				self := users[username]
				log.Println(self.friends)

				//loops through friends list, if the username is there, gives the index, otherwise gives a -1
				idx := -1
				var removedName string
				for i, friend := range self.friends {
					if friend.username == removeFriend{
						idx = i
						removedName = friend.name
					}
				}

				//if -1 (doesn't exist in friends) then display this
				if idx == -1{
					fmt.Fprintf(w, "<h3>%s isn't in your friend's list.</h3>", removedName)
				} else{

					//unfollow them
					self.friends = append(self.friends[:idx], self.friends[idx+1:]...)

					fmt.Fprintf(w, "<h3>You have unfollowed %s</h3>", removedName)
				}

				log.Println(self.friends)

			}
		} else{
			//User doesnt exist or hasnt signed up
			if remFriend == nil && (len(removeFriend) >= 1) {
				fmt.Fprintf(w, "<h3>User doesn't exist</h3>")
			}
		}
	}

	//render profile page, which also has ability to tweet
	renderHtml("profile.html", "home.go", w)

	//process a tweet/post
	if r.Method == "POST" {
		r.ParseForm()

		//create a post from the tweet
		curPost := strings.Join(r.Form["aPost"], " ")
		postCount++

		aPost := Post{
			username: username,
			ID: postCount,
			post: curPost,
		}

		posts = append(posts,aPost)
	}

	if(len(posts) > 1) {
		printPosts(w, username, false, r)
	}
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

		if(username == "" || password == "" || name == "") {
			http.Redirect(w,r,"/signup",302)
		}

		_, checkUser := users[username]

		//if username is already taken
		if(checkUser == true) {
			http.Redirect(w,r,"/signup",302)
		}

		//stores user info into person struct
		pUser := Person{
			username: username,
			password: password,
			name: name,
			friends: nil,
		}
		//Store in our psuedo-database. Format: username : {password, name}
		users[username] = &pUser

		//Set a cookie for 5 minutes and redirect to their profile
		setCookie(w,username, 60)
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
	http.Redirect(w,r,"/signin",302)
}

//processes when the user clicks the delete account link
func removeAcct(w http.ResponseWriter, r *http.Request) {
	myUser, _ := r.Cookie("username")
	myUsername := users[myUser.Value].username
	if r.Method == "GET" {

    	//Render the page
    	renderHtml("removeAcct.html", "home.go", w)

	} else {
		r.ParseForm()

		//if they clicked yes to delete their account
		if r.PostFormValue("Yes") == "Yes" {

			//loop through all the users
			for _, theUser := range users {

				//loop through each users' friends
				for idx,friend := range theUser.friends {

					//if they have me as a friend, delete myself
					if friend.username == myUsername {
						theUser.friends = append(theUser.friends[:idx], theUser.friends[idx+1:]...)
					}
				}
			}

			//go through all posts, if I have a post, delete it
			for i:=0; i < len(posts); i++ {
				thePost := posts[i]
				if thePost.username == myUsername {
					posts = append(posts[:i],posts[i+1:]...)
				}
			}

			//delete myself from the users map
			delete(users,myUsername)
			
			//logout (delete cookie and redirect to signin page)
			logout(w,r)

		} else {

			//if I clicked no to delete my account, go to my newsfeed
			log.Println("No")
			http.Redirect(w,r,"/newsfeed",http.StatusPermanentRedirect)
		}
	}
}

//used so I can type localhost:(portnumber)/ and it will direct me to the signin page
//instead of writing signin
func start(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w,r,"/signin",http.StatusPermanentRedirect)
}

//call all of the handler functions
func main() {
	http.HandleFunc("/", start)
    http.HandleFunc("/signin", signin)
    http.HandleFunc("/signup", signup)
    http.HandleFunc("/newsfeed", newsfeed)
    http.HandleFunc("/logout", logout)
    http.HandleFunc("/profile", profile)
    http.HandleFunc("/removeAcct", removeAcct)
    http.ListenAndServe(":8070", nil)
}