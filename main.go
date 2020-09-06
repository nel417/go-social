package main

import (
	"context"
	"html/template"
	"net/http"

	"github.com/go-redis/redis"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

//global redis client
var client *redis.Client

//global template
var templates *template.Template

//byte array for keys
var store = sessions.NewCookieStore([]byte("secret-lol"))

func main() {
	//redis client and host
	client = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	// render templates in templates folder
	templates = template.Must(template.ParseGlob("templates/*.html"))
	//initialzie mux router
	r := mux.NewRouter()
	//get route for home
	r.HandleFunc("/", indexGetHandler).Methods("GET")
	//post route for index
	r.HandleFunc("/", indexPostHandler).Methods("POST")
	//post route for login
	r.HandleFunc("/login", loginGetHandler).Methods("GET")
	//get route for login
	r.HandleFunc("/login", loginPostHandler).Methods("POST")
	//handles the home response

	//this was to test the cookie
	// r.HandleFunc("/test", testGetHandler).Methods("GET")

	http.Handle("/", r)

	//static file instantiation
	fs := http.FileServer(http.Dir("./static/"))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))
	//create server
	http.ListenAndServe(":8080", nil)

}

// INDEX AREA
func indexGetHandler(w http.ResponseWriter, r *http.Request) {
	// define context
	ctx := context.TODO()
	//comments and error handling
	comments, err := client.LRange(ctx, "comments", 0, 10).Result()
	//error handling
	if err != nil {
		return
	}
	//executes template, index.html
	templates.ExecuteTemplate(w, "index.html", comments)
}

func indexPostHandler(w http.ResponseWriter, r *http.Request) {
	//parses form from the request body
	r.ParseForm()
	//define context
	ctx := context.TODO()
	//hits comment name in html
	comment := r.PostForm.Get("comment")
	//pushes to redis db
	client.LPush(ctx, "comments", comment)
	//redirects to home page with new comment in array
	http.Redirect(w, r, "/", 302)
}

// LOGIN AREA
func loginGetHandler(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "login.html", nil)
}
func loginPostHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	username := r.PostForm.Get("username")
	session, _ := store.Get(r, "session")
	session.Values["username"] = username
	session.Save(r, w)
}

//COOKIE TEST
// func testGetHandler(w http.ResponseWriter, r *http.Request) {
// 	session, _ := store.Get(r, "session")
// 	untyped, ok := session.Values["username"]
// 	if !ok {
// 		return
// 	}
// 	username, ok := untyped.(string)
// 	if !ok {
// 		return
// 	}
// 	w.Write([]byte(username))
// }
