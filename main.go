package main

import (
	"context"
	"html/template"
	"net/http"

	"github.com/go-redis/redis"
	"github.com/gorilla/mux"
)

//global redis client
var client *redis.Client

//global template
var templates *template.Template

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
	//handles the home response
	http.Handle("/", r)
	//makes server
	http.ListenAndServe(":8080", nil)

}

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
