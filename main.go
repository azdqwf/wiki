package main

import (
	"errors"
	"html/template"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"regexp"
)

var templates = template.Must(template.ParseFiles("tmpl/edit.html", "tmpl/view.html"))
var validPath = regexp.MustCompile("^/(edit|save|view)/([a-zA-Z0-9]+)$")

//Page ...
type Page struct {
	Title string
	Body  []byte
}

func (p *Page) save() error {
	filename := filepath.Join("data", p.Title+".txt")
	return ioutil.WriteFile(filename, p.Body, 0600)
}

func loadPage(title string) (*Page, error) {
	filename := filepath.Join("data", title+".txt")
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	t, err := template.ParseFiles(filepath.Join("tmpl", tmpl+".html"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = t.Execute(w, p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func getTitle(w http.ResponseWriter, r *http.Request) (string, error) {
	m := validPath.FindStringSubmatch(r.URL.Path)
	if m == nil {
		http.NotFound(w, r)
		return "", errors.New("Invalid Page Title")
	}
	return m[2], nil // The title is the second subexpression.
}
func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/view/", makeHandler(viewHandler))
	mux.HandleFunc("/edit/", makeHandler(editHandler))
	mux.HandleFunc("/save/", makeHandler(saveHandler))
	http.Handle("/files/", http.StripPrefix("/files/",
		http.FileServer(http.Dir("."))))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if m := validPath.Find([]byte(r.URL.Path)); m == nil {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		mux.ServeHTTP(w, r)
	})

	http.ListenAndServe(":8080", nil)
}
