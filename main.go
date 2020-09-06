package main

import (
	"context"
	"html/template"
	"net/http"

	"github.com/go-redis/redis"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
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
	r.HandleFunc("/", AuthRequired(indexGetHandler)).Methods("GET")
	//post route for index
	r.HandleFunc("/", AuthRequired(indexPostHandler)).Methods("POST")
	//get route for login
	r.HandleFunc("/login", loginGetHandler).Methods("GET")
	//post route for login
	r.HandleFunc("/login", loginPostHandler).Methods("POST")
	//get route for register
	r.HandleFunc("/register", registerGetHandler).Methods("GET")
	//post route for register
	r.HandleFunc("/register", registerPostHandler).Methods("POST")
	//handle index route
	http.Handle("/", r)

	//static file instantiation
	fs := http.FileServer(http.Dir("./static/"))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))
	//create server
	http.ListenAndServe(":8080", nil)

}

func AuthRequired(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, _ := store.Get(r, "session")
		_, ok := session.Values["username"]
		if !ok {
			http.Redirect(w, r, "/login", 302)
			return
		}
		handler.ServeHTTP(w, r)
	}
}

// INDEX AREA
func indexGetHandler(w http.ResponseWriter, r *http.Request) {

	// define context
	ctx := context.TODO()
	//comments and error handling
	comments, err := client.LRange(ctx, "comments", 0, 10).Result()
	//error handling
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("internal Server Error"))
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
	err := client.LPush(ctx, "comments", comment).Err()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("internal Server Error"))
		return
	}
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
	password := r.PostForm.Get("password")
	ctx := context.TODO()
	hash, err := client.Get(ctx, "user:"+username).Bytes()
	if err == redis.Nil {
		templates.ExecuteTemplate(w, "login.html", "Unknown User")
		return

	} else if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("internal Server Error"))
		return
	}

	err = bcrypt.CompareHashAndPassword(hash, []byte(password))
	if err != nil {
		templates.ExecuteTemplate(w, "login.html", "Invalid Login")
		return
	}
	//gets session
	session, _ := store.Get(r, "session")
	//sets session
	session.Values["username"] = username
	session.Save(r, w)
	//redirects to index
	http.Redirect(w, r, "/", 302)

}

func registerGetHandler(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "register.html", nil)

}

func registerPostHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	username := r.PostForm.Get("username")
	password := r.PostForm.Get("password")
	cost := bcrypt.DefaultCost
	ctx := context.TODO()
	hash, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("internal Server Error"))
		return
	}
	err = client.Set(ctx, "user:"+username, hash, 0).Err()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("internal Server Error"))
		return
	}
	//redirects to login
	http.Redirect(w, r, "/login", 302)
}
