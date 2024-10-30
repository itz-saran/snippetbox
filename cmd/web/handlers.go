package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"snippetbox.saran.net/internal/models"
	"snippetbox.saran.net/internal/validator"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	snippets, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, err)
	}
	data := app.newTemplateData(r)
	data.Snippets = snippets
	app.render(w, http.StatusOK, "home.templ", data)
}

func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	snippet, err := app.snippets.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	data := app.newTemplateData(r)
	data.Snippet = snippet
	app.render(w, http.StatusOK, "view.templ", data)
}

type snippetCreateFormData struct {
	Title               string `form:"title"`
	Content             string `form:"content"`
	Expires             int    `form:"expires"`
	validator.Validator `form:"-"`
}

func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	var form snippetCreateFormData
	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	// Validations
	form.CheckField(validator.NotBlank(form.Title), "title", "Title cannot be empty")
	form.CheckField(validator.MaxChars(form.Title, 100), "title", "Title cannot exceed 100 characters")
	form.CheckField(validator.NotBlank(form.Content), "content", "Content cannot be empty")
	form.CheckField(validator.PermittedValues(form.Expires, 1, 7, 365), "expires", "Expires must be one of the above values")
	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "create.templ", data)
		return
	}
	id, err := app.snippets.Insert(form.Title, form.Content, form.Expires)
	if err != nil {
		app.serverError(w, err)
	}
	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}

func (app *application) snippetCreateForm(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = snippetCreateFormData{
		Expires: 365,
	}
	app.render(w, http.StatusOK, "create.templ", data)
}
