package main

import (
	"net/http"
	"runtime/debug"
)

func (app *application) serveError(w http.ResponseWriter, err error) {
	trace := err.Error() + "\n" + string(debug.Stack())
	app.errorLog.Output(2, trace)

	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (app *application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}
