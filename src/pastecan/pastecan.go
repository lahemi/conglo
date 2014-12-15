package pastecan

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/lahemi/conglo/pbnf"
	"github.com/lahemi/conglo/scraps"
)

type Page struct {
	Title, Body string
}

func (p *Page) save(config map[string]string) (*Page, error) {
	pastePath := config["pastePath"]
	filename := pastePath + p.Title
	_, e := os.Stat(filename)
	if e == nil {
		filename = pastePath + scraps.GenRandTitle()
		for {
			_, err := os.Stat(filename)
			if err != nil {
				break
			}
			filename = pastePath + scraps.GenRandTitle()
		}
		p.Title = strings.Split(filename, "/")[0]
	}
	err := ioutil.WriteFile(filename, []byte(p.Body), 0600)
	if err != nil {
		return p, err
	}
	return p, nil
}

func loadPage(title string, config map[string]string) (*Page, error) {
	filename := config["pastePath"] + title
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: string(body)}, nil
}

func Upload(config map[string]string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body := r.FormValue("pastebody")
		lang := r.FormValue("lang")
		title := scraps.GenRandTitle()
		switch lang {
		case "go", "lua":
			body = pbnf.Colourify(lang, body)
		}
		p := &Page{Title: title, Body: body}
		p, err := p.save(config)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/pastecan/v/"+p.Title, http.StatusFound)
	}
}

func View(config map[string]string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		validPath := regexp.MustCompile("^" + config["base"] + "/v/([A-Za-z]+)$")
		m := validPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			http.NotFound(w, r)
			return
		}
		p, err := loadPage(m[1], config)
		if err != nil {
			return
		}
		cnt, err := ioutil.ReadFile(config["viewHTML"])
		if err != nil {
			return
		}
		html := string(cnt)
		re := regexp.MustCompile("{{.WAT}}")
		html = re.ReplaceAllString(html, p.Body)
		fmt.Fprint(w, html)
	}
}

func Index(config map[string]string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cnt, err := ioutil.ReadFile(config["pasteIndex"])
		if err != nil {
			scraps.PrintErr(err)
			return
		}
		fmt.Fprint(w, string(cnt))
	}
}
