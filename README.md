# Social Media Prototype

A non production social media with a twitter like timeline using reddis for database and caching and Golang for routing and back end

## Installation

This app uses a several libraries outside of the standard library   
Using the *go get* command, type in the following in your terminal
```
go get github.com/gorilla/mux  
go get github.com/go-redis/redis 
go get golang.org/x/crypto/bcrypt 

```

## External Library Breakdown
**Gorilla Mux** is used for our routing, session, and auth  
**Go Reddis** is our Redis client  
**Crypto Bcrypt** is for password hashing

## Usage

Make sure redis is installed and run *go run main.go* and *src/redis-server* in a split terminal
```go
go run main.go
src/redis-server
```

## Serving
The Redis server can be found in /models/db.go  
```go
package models

import (
	"github.com/go-redis/redis"
)

// global redis client
var client *redis.Client
// our server on port 6379
func Init() {
	// redis client and host
	client = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
}

```
and the Mux server can be found in our main.go file  
```go
package main

import (
  // Standard library imports
	"html/template"
	"net/http"

  // Packages
	"./models"
	"./routes"
	"./utils"
)

// Global template
var templates *template.Template

func main() {
	models.Init()
	utils.LoadTemplates("templates/*.html")
	r := routes.NewRouter()
	// handle index route
	http.Handle("/", r)
	// create server
	http.ListenAndServe(":8080", nil)

}

```
## Routing

```go
// Router function that takes a pointer to our Mux Router
func NewRouter() *mux.Router {
    // Initialize router
	r := mux.NewRouter()
    // Home get route
	r.HandleFunc("/", middleware.AuthRequired(indexGetHandler)).Methods("GET")
    // Home post route
	r.HandleFunc("/", middleware.AuthRequired(indexPostHandler)).Methods("POST")
    // Login get route
	r.HandleFunc("/login", loginGetHandler).Methods("GET")
    // Login post route
	r.HandleFunc("/login", loginPostHandler).Methods("POST")
    // Register get route
	r.HandleFunc("/register", registerGetHandler).Methods("GET")
    // Register post route
	r.HandleFunc("/register", registerPostHandler).Methods("POST")
    // Serves static files such as CSS and JS files
	fs := http.FileServer(http.Dir("./static/")) 
    // Strips prefix from file server for the static folder
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))
    // Path for going to someones page to show their updates
	r.HandleFunc("/{username}", middleware.AuthRequired(userGetHandler)).Methods("GET")
    // Return our server
	return r
}
```
### Documentation In Progress . . . 
## Contributing
Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

Please make sure to update tests as appropriate.

## License
[MIT](https://choosealicense.com/licenses/mit/)
