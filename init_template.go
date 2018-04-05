package main

import (
	"html/template"
	"io"
	"net/http"
	"path/filepath"

	"github.com/labstack/echo"
	"github.com/valyala/bytebufferpool"
)

type (
	// TemplateMap is a map containing html/template specs.
	TemplateMap map[string]*template.Template

	// Template interfaces the templates list that can be
	// access through mapping.
	Template struct {
		templates TemplateMap
	}
)

// Render is an interface that will connect the custom template
// function to echo.
func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	tmpl, ok := t.templates[name]
	if !ok {
		return echo.NewHTTPError(http.StatusInternalServerError, "template does not exist")
	}
	buf := bytebufferpool.Get()
	defer bytebufferpool.Put(buf)

	err := tmpl.ExecuteTemplate(w, "base.tmpl", data)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	buf.WriteTo(w)
	return nil
}

// LoadTemplates creates a map of all layout and include files.
func loadTemplates(e *echo.Echo) *Template {
	templates := make(TemplateMap)
	templatesDir := "view/templates/"

	layouts, err := filepath.Glob(templatesDir + "layouts/*tmpl")
	if err != nil {
		e.Logger.Fatal(err)
	}

	includes, err := filepath.Glob(templatesDir + "*.tmpl")
	if err != nil {
		e.Logger.Fatal(err)
	}

	for _, file := range includes {
		filename := filepath.Base(file)
		files := append(layouts, file)
		templates[filename] = template.Must(templates[filename].ParseFiles(files...))
	}

	t := &Template{
		templates,
	}
	return t
}
