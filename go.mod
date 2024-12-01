module snippetbox.saran.net

go 1.23.0

require (
	github.com/alexedwards/scs/postgresstore v0.0.0-20240316134038-7e11d57e8885
	github.com/alexedwards/scs/v2 v2.8.0
	github.com/go-playground/form/v4 v4.2.1
	github.com/lib/pq v1.10.9
)

require golang.org/x/crypto v0.29.0

require github.com/justinas/nosurf v1.1.1

replace snippetbox.saran.net => ./
