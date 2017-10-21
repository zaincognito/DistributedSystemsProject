Distributed and Parallel Systems Project Part 1
Contributors: Nabil Ahmed (na1489) and Zulkifl Arefin (za453)

ISSUES:
-For signin - if user has unknown username or incorrect password, redirects to signin page (correct), but does not show an error (so user doesn't know if incorrect password nor unknown username) -- we tried to display a message with the redirect, but it gives weird errors, and if we put it after the redirect, then that code is never reached.

-For signup - similar issue as above, where there is no error message, only a redirect when something is incorrect (if any fields for signup empty, is user already exists).

FILES:

home.go: Has all of the backend (with comments) which processes all of the html files. Creates a connection to the server, creates communication between the html files, creates data structures to store users, tweets, and other information.

logout.html: dummy file in order for the user to be able to click on a link and then the home.go file processes this and logs a user out.

navbar.html: is the navigation bar on top of the newsfeed and profile where it shows links such as "Newsfeed", "Profile", "Logout" and "Delete Account".

newsfeed.html: is the html file which shows the newsfeed, it has a form for the user to post tweets as well, see their own tweets and friends tweets.

profile.html: is the html file which shows the user profile, shows all of the friends the user has, it has a form for the user to post tweets, and it shows the tweets of the user only.

removeAcct: is a form to ask the user if they are sure they want to delete their account, if it is a yes home.go deletes the account, if no then they are redirected to their newsfeed.

signin.html: is a form for anyone to login, and a link for signing up

signup.html: is a form for anyone to signup