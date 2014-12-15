package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/lahemi/conglo/datacan"
	"github.com/lahemi/conglo/filecan"
	"github.com/lahemi/conglo/pastecan"

	"github.com/lahemi/conglo/scraps"
)

var (
	port string

	basePath  = os.Getenv("HOME") + "/.local/share/conglo/"
	htmlPath  = basePath + "htmls/"
	cssPath   = basePath + "styles/"
	imagePath = basePath + "images/"
	jsPath    = basePath + "javascript/"

	indexPage     = htmlPath + "index.html"
	filmsPage     = htmlPath + "films.html"
	theartistPage = htmlPath + "theartist.html"

	config map[string]map[string]string
)

func makeBasicHandler(page string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cnt, err := ioutil.ReadFile(page)
		if err != nil {
			panic(err)
		}
		fmt.Fprint(w, string(cnt))
	}
}

func init() {
	flag.StringVar(&port, "port", "15003", "The port to use.")
	flag.Parse()

	config = getConfig()

	if err := datacan.DBInit(config["datacan"]); err != nil {
		scraps.PrintErr(err)
		os.Exit(1)
	}
}

func main() {
	http.HandleFunc("/", makeBasicHandler(indexPage))
	http.HandleFunc("/films/", makeBasicHandler(filmsPage))
	http.HandleFunc("/films/theartist/", makeBasicHandler(theartistPage))

	http.HandleFunc("/filecan/", filecan.Index(config["filecan"]))
	http.HandleFunc("/filecan/save", filecan.Upload(config["filecan"]))
	http.HandleFunc("/filecan/v/", filecan.View(config["filecan"]))

	http.HandleFunc("/pastecan/", pastecan.Index(config["pastecan"]))
	http.HandleFunc("/pastecan/save", pastecan.Upload(config["pastecan"]))
	http.HandleFunc("/pastecan/v/", pastecan.View(config["pastecan"]))

	http.HandleFunc("/datacan/", datacan.Index(config["datacan"]))
	http.HandleFunc("/datacan/save", datacan.Save(config["datacan"]))
	http.HandleFunc("/datacan/v/", datacan.View(config["datacan"]))

	http.Handle(
		"/fserve/",
		http.StripPrefix(
			"/fserve",
			http.FileServer(http.Dir(config["fserve"]["dir"])),
		),
	)

	http.Handle(
		"/styles/",
		http.StripPrefix(
			"/styles/",
			http.FileServer(http.Dir(cssPath)),
		),
	)
	http.Handle(
		"/images/",
		http.StripPrefix(
			"/images/",
			http.FileServer(http.Dir(imagePath)),
		),
	)
	http.Handle(
		"/javascript/",
		http.StripPrefix(
			"/javascript/",
			http.FileServer(http.Dir(jsPath)),
		),
	)

	http.ListenAndServe(":"+port, nil)
}
