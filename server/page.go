package server

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"path"
)

var templateDir = "tmpl"

// Page is a abstraction for the Go Template handels
type Page struct {
	Language     string
	Title        string
	templateName string
	template     *template.Template
}

func AsJson(w http.ResponseWriter, v interface{}) error {
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}
	w.Header().Add("Content-Type", "application/json")
	w.Write(b)
	return nil
}

// OpenPage perpares a Page Instance to be used
func OpenPage(name string) (*Page, error) {
	err := checkPageName(name)
	if err != nil {
		return nil, err
	}
	ta, err := template.ParseGlob(fmt.Sprintf("%s/*.tmpl", templateDir))
	t := template.Must(ta, err)
	namedTemplate := t.Lookup(name)
	p := &Page{template: namedTemplate, templateName: name}
	return p, nil
}

func checkPageName(name string) error {
	fp := path.Join(templateDir, fmt.Sprintf("%s.tmpl", name))
	s, err := os.Stat(fp)
	if err != nil {
		return err
	}
	if s.IsDir() {
		return fmt.Errorf("%s.tmpl is a directory", name)
	}
	return nil
}

// Render writes a Template with given data to the Writer
func (p *Page) Render(wr io.Writer, data interface{}) error {
	return p.template.ExecuteTemplate(wr, p.templateName, data)
}
