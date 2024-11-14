package main

import (
	"net/http"
	"path/filepath"
)

func (app *application) routes() *http.ServeMux {
	mux := http.NewServeMux()
	manageSession := func(handler func(http.ResponseWriter, *http.Request)) http.Handler {
		return app.sessionManager.LoadAndSave(http.HandlerFunc(handler))
	}
	fileServer := http.FileServer(&neuteredFileSystem{fs: http.Dir("./ui/static/")})

	mux.Handle("/static", http.NotFoundHandler())
	mux.Handle("/static/", http.StripPrefix("/static/", fileServer))
	mux.Handle("/{$}", manageSession(app.home))
	mux.Handle("GET /snippet/view/{id}", manageSession(app.snippetView))
	mux.Handle("GET /snippet/create", manageSession(app.snippetCreateForm))
	mux.Handle("POST /snippet/create", manageSession(app.snippetCreate))

	return mux
}

type neuteredFileSystem struct {
	fs http.FileSystem
}

// ? This will prevent the directory listing for directory paths
func (nfs *neuteredFileSystem) Open(path string) (http.File, error) {
	f, err := nfs.fs.Open(path)
	if err != nil {
		return nil, err
	}

	s, err := f.Stat()
	if err != nil {
		return nil, err
	}

	if s.IsDir() {
		index := filepath.Join(path, "index.html")
		if _, err := nfs.fs.Open(index); err != nil {
			closeErr := f.Close()

			if closeErr != nil {
				return nil, closeErr
			}

			return nil, err
		}
	}

	return f, nil
}
