package routes

import (
	"net/http"

	"../middleware"
	"../models"
	"../sessions"
	"../utils"
	"github.com/gorilla/mux"
)

func NewRouter() *mux.Router {
	r := mux.NewRouter()
	//get route for home
	r.HandleFunc("/", middleware.AuthRequired(indexGetHandler)).Methods("GET")
	//post route for index
	r.HandleFunc("/", middleware.AuthRequired(indexPostHandler)).Methods("POST")
	//get route for login
	r.HandleFunc("/login", loginGetHandler).Methods("GET")
	//post route for login
	r.HandleFunc("/login", loginPostHandler).Methods("POST")
	//get route for register
	r.HandleFunc("/register", registerGetHandler).Methods("GET")
	//post route for register
	r.HandleFunc("/register", registerPostHandler).Methods("POST")

	//static file instantiation
	fs := http.FileServer(http.Dir("./static/"))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))
	return r
}

// INDEX AREA
func indexGetHandler(w http.ResponseWriter, r *http.Request) {

	//comments and error handling
	comments, err := models.GetComments()
	//error handling
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("internal Server Error"))
		return
	}
	//executes template, index.html
	utils.ExecuteTemplate(w, "index.html", comments)
}

func indexPostHandler(w http.ResponseWriter, r *http.Request) {
	//parses form from the request body
	r.ParseForm()
	//hits comment name in html
	comment := r.PostForm.Get("comment")
	//pushes to redis db
	err := models.PostComment(comment)
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
	utils.ExecuteTemplate(w, "login.html", nil)
}

func loginPostHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	username := r.PostForm.Get("username")
	password := r.PostForm.Get("password")
	err := models.AuthenticateUser(username, password)

	if err != nil {
		switch err {
		case models.ErrUserNotFound:
			utils.ExecuteTemplate(w, "login.html", "Unknown User")
		case models.ErrInvalidLogin:
			utils.ExecuteTemplate(w, "login.html", "Invalid Login")
		default:
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("internal Server Error"))

		}
		return
	}

	//gets session
	session, _ := sessions.Store.Get(r, "session")
	//sets session
	session.Values["username"] = username
	session.Save(r, w)
	//redirects to index
	http.Redirect(w, r, "/", 302)

}

func registerGetHandler(w http.ResponseWriter, r *http.Request) {
	utils.ExecuteTemplate(w, "register.html", nil)

}

func registerPostHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	username := r.PostForm.Get("username")
	password := r.PostForm.Get("password")
	err := models.RegisterUser(username, password)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("internal Server Error"))
		return
	}
	//redirects to login
	http.Redirect(w, r, "/login", 302)
}
