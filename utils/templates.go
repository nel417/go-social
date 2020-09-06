package utils

import (
	"net/http"

	"html/template"
)

//global template
var templates *template.Template

func LoadTemplates(pattern string) {
	// render templates in templates folder
	templates = template.Must(template.ParseGlob("templates/*.html"))
}

func ExecuteTemplate(w http.ResponseWriter, tmpl string, data interface{}) {
	//executes template, index.html
	templates.ExecuteTemplate(w, tmpl, data)
}
