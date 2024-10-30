package main

import (
	"database/sql"
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/go-playground/form/v4"
	_ "github.com/lib/pq"
	"snippetbox.saran.net/internal/models"
)

func main() {
	addr := flag.String("addr", ":3000", "HTTP port the server runs on")
	dsn := flag.String("dsn", "postgresql://snb_web_user:snbwebuser@localhost:5432/snippetbox?sslmode=disable", "Postgres data source/connection string  name")
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	db, err := openDB(*dsn)
	if err != nil {
		errorLog.Fatal(err)
	}
	defer db.Close()

	templateCache, err := newTemplateCache()
	if err != nil {
		errorLog.Fatal(err)
	}
	formDecoder := form.NewDecoder()
	app := &application{
		infoLog:       infoLog,
		errorLog:      errorLog,
		snippets:      &models.SnippetModel{DB: db},
		templateCache: templateCache,
		formDecoder:   formDecoder,
	}

	middlewares := CreateStack(app.recoverPanic, app.logRequest, secureHeaders)
	server := &http.Server{
		Addr:     *addr,
		Handler:  middlewares(app.routes()),
		ErrorLog: errorLog,
	}
	infoLog.Printf("Listening on port %s\n", *addr)
	err = server.ListenAndServe()
	errorLog.Fatal(err)
}

type application struct {
	errorLog      *log.Logger
	infoLog       *log.Logger
	snippets      *models.SnippetModel
	templateCache map[string]*template.Template
	formDecoder   *form.Decoder
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
