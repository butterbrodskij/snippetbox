package main

import (
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"mysnippetbox.com/snippetbox/pkg/models"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		app.notFound(w)
		return
	}

	files := []string{
		"./ui/html/home.page.html",
		"./ui/html/base.layout.html",
		"./ui/html/footer.partial.html",
	}

	ts, err := template.ParseFiles(files...)

	if err != nil {
		app.serveError(w, err)
		return
	}

	s, err := app.snippets.Latest()

	if err != nil {
		app.serveError(w, err)
		return
	}

	err = ts.Execute(w, s)

	if err != nil {
		app.serveError(w, err)
	}
}

func (app *application) showSnippet(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))

	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	s, err := app.snippets.Get(id)

	if errors.Is(err, models.ErrNotFound) {
		app.notFound(w)
		return
	} else if err != nil {
		app.serveError(w, err)
		return
	}

	files := []string{
		"./ui/html/show.page.html",
		"./ui/html/base.layout.html",
		"./ui/html/footer.partial.html",
	}

	ts, err := template.ParseFiles(files...)

	if err != nil {
		app.serveError(w, err)
		return
	}

	err = ts.Execute(w, s)

	if err != nil {
		app.serveError(w, err)
	}
}

func (app *application) createSnippet(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		files := []string{
			"./ui/html/create.page.html",
			"./ui/html/base.layout.html",
			"./ui/html/footer.partial.html",
		}

		ts, err := template.ParseFiles(files...)

		if err != nil {
			app.serveError(w, err)
			return
		}

		err = ts.Execute(w, nil)

		if err != nil {
			app.serveError(w, err)
		}

		return
	}

	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}

	title := r.FormValue("title")
	content := r.FormValue("content")
	expires, err := strconv.Atoi(r.FormValue("days"))

	if err != nil {
		app.serveError(w, err)
		return
	}

	id, err := app.snippets.Insert(title, content, expires)

	if err != nil {
		app.serveError(w, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/snippet?id=%d", id), http.StatusSeeOther)
}
