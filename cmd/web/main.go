package main

import (
	"database/sql"
	"flag"
	"log"
	"net/http"
	"os"

	"mysnippetbox.com/snippetbox/pkg/models/postgres"

	_ "github.com/lib/pq"
)

type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
	snippets *postgres.SnippetModel
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn)

	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

func main() {
	addr := flag.String("addr", ":4000", "network address http")
	dsn := flag.String("dsn", "user=webuser password=password dbname=snippetbox sslmode=disable", "data source")
	flag.Parse()

	infos := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errs := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Llongfile)

	db, err := openDB(*dsn)

	if err != nil {
		errs.Fatal(err)
	}

	defer db.Close()

	app := &application{
		errorLog: errs,
		infoLog:  infos,
		snippets: &postgres.SnippetModel{DB: db},
	}

	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errs,
		Handler:  app.routes(),
	}

	app.infoLog.Printf("starting the server on localhost%s", *addr)
	err = srv.ListenAndServe()
	errs.Fatal(err)
}
