package main

import (
	"net/http"
	"path/filepath"

	"snippetbox.saran.net/ui"
)

func (app *application) routes() *http.ServeMux {
	mux := http.NewServeMux()
	fileServer := http.FileServer(&neuteredFileSystem{fs: http.FS(ui.Files)})

	mux.Handle("/static", http.NotFoundHandler())
	mux.Handle("/static/*filepath", fileServer)

	defaultMiddleware := CreateStack(app.manageSession, noSurf, app.authenticate)

	mux.Handle("/{$}", defaultMiddleware(http.HandlerFunc(app.home)))
	mux.Handle("GET /snippet/view/{id}", defaultMiddleware(http.HandlerFunc(app.snippetView)))
	mux.Handle("GET /user/signup", defaultMiddleware(http.HandlerFunc(app.UserSignupForm)))
	mux.Handle("POST /user/signup", defaultMiddleware(http.HandlerFunc(app.UserSignup)))
	mux.Handle("GET /user/login", defaultMiddleware(http.HandlerFunc(app.UserLoginForm)))
	mux.Handle("POST /user/login", defaultMiddleware(http.HandlerFunc(app.UserLogin)))

	// ProtectedRoutes
	protected := CreateStack(app.manageSession, app.requireAuth)
	mux.Handle("POST /user/logout", protected(http.HandlerFunc(app.UserLogout)))
	mux.Handle("GET /snippet/create", protected(http.HandlerFunc(app.snippetCreateForm)))
	mux.Handle("POST /snippet/create", protected(http.HandlerFunc(app.snippetCreate)))

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
