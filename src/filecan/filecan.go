package filecan

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/lahemi/conglo/scraps"
)

func Upload(config map[string]string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		uploadPath := config["uploadPath"]

		files, err := ioutil.ReadDir(uploadPath)
		if err != nil {
			scraps.PrintErr(err)
			return
		}

		// Don't keep too many files around, just in case.
		if len(files) > 30 {
			for _, f := range files {
				if err := os.Remove(uploadPath + f.Name()); err != nil {
					scraps.PrintErr("Failed to remove redundant files.")
					return
				}
			}
		}

		// Naive, though works if there are no nasty complications nor,
		// more likely, a malicious party giving false ContentLength.
		// Should read to a buffer and check against the limit.
		flimit, err := strconv.Atoi(config["flimit"])
		if err != nil {
			scraps.PrintErr(err)
			return
		}
		if r.ContentLength > int64(flimit) {
			fmt.Fprintf(w, "Too large file...")
			return
		}
		file, _, err := r.FormFile("file")
		if err != nil {
			scraps.PrintErr(err)
			return
		}
		defer file.Close()

		var title = scraps.GenRandTitle()
		out, err := os.Create(uploadPath + title) // possibility of overwrite
		if err != nil {
			scraps.PrintErr("Unable to create the file for writing. Check your privilege.")
			return
		}
		defer out.Close()

		// write the content from POST to the file
		_, err = io.Copy(out, file)
		if err != nil {
			scraps.PrintErr(err)
			return
		}
		fmt.Println("File uploaded - " + title)

		http.Redirect(w, r, "/filecan/v/"+title, http.StatusFound)
	}
}

func View(config map[string]string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		up := r.URL.Path
		if up == "" {
			http.NotFound(w, r)
			return
		}
		// Remove the /filecan/v? part from the url.
		ipath := config["uploadPath"] + up[strings.LastIndex(up, "/")+1:]
		http.ServeFile(w, r, ipath)
	}
}

func Index(config map[string]string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cnt, err := ioutil.ReadFile(config["uploadHTML"])
		if err != nil {
			scraps.PrintErr(err)
			return
		}
		fmt.Fprint(w, string(cnt))
	}
}
