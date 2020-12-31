package main

import (
  "bytes"
  "errors"
  "fmt"
  "log"
  "regexp"
  "net/http"
  "html/template"
  "io/ioutil"
)

// controllers

type Page struct {
  Title string
  Body []byte
}

func (p *Page) save() error {
  filename := p.Title + ".txt"
  return ioutil.WriteFile(filename, p.Body, 0600)
}

func loadPage(title string) (*Page, error) {
  body, err := ioutil.ReadFile(title + ".txt")
  if err != nil {
    return nil, err
  }
  return &Page{Title: title, Body: body}, nil
}

// server handlers

func defaultHandler(w http.ResponseWriter, r *http.Request) {
  // respond with the URL path without the leading /
  fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}

func viewHandler(w http.ResponseWriter, r *http.Request, title string) {
  p, _ := loadPage(title)
  renderTemplate(w, "view", p)
}

func editHandler(w http.ResponseWriter, r *http.Request, title string) {
  p, err := loadPage(title)
  if err != nil {
    http.Redirect(w, r, "/edit/"+title, http.StatusFound)
    return
  }
  renderTemplate(w, "edit", p)
}

func saveHandler(w http.ResponseWriter, r *http.Request, title string) {
  body := r.FormValue("body")
  p := &Page{Title: title, Body: []byte(body)}
  err2 := p.save()
  if err2 != nil {
    serverErr(w, err2)
    return
  }
  http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

// regexp
var validPath = regexp.MustCompile("^/(edit|save|view)/([a-zA-Z0-9]+)$")

// load templates
var templates = template.Must(template.ParseFiles("edit.html", "view.html"))

func main() {
  // load routes
  http.HandleFunc("/view/", makeHandler(viewHandler))
  http.HandleFunc("/edit/", makeHandler(editHandler))
  http.HandleFunc("/save/", makeHandler(saveHandler))
  http.HandleFunc("/", defaultHandler)
  log.Fatal(http.ListenAndServe(":8080", nil))
}

// misc helpers

func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
  return func(w http.ResponseWriter, r *http.Request) {
    title, err := getTitle(w, r)
    if err != nil {
      return
    }
    fn(w, r, title)
  }
}

func serverErr(w http.ResponseWriter, err error) {
  http.Error(w, err.Error(), http.StatusInternalServerError)
}

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
  var res bytes.Buffer
  err := templates.ExecuteTemplate(&res, tmpl+".html", p)
  if err != nil {
    serverErr(w, err)
  } else {
    res.WriteTo(w)
  }
}

func getTitle(w http.ResponseWriter, r *http.Request) (string, error) {
  m := validPath.FindStringSubmatch(r.URL.Path)
  if m == nil {
    http.NotFound(w, r)
    return "", errors.New("invalid Page Title")
  }
  return m[2], nil // The title is the second subexpression.
}
