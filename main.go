package main

import (
	"html/template"
	"net/http"

	"./models"
	"./routes"
	"./utils"
)

//global template
var templates *template.Template

func main() {
	models.Init()
	utils.LoadTemplates("templates/*.html")
	r := routes.NewRouter()
	//handle index route
	http.Handle("/", r)
	//create server
	http.ListenAndServe(":8080", nil)

}
