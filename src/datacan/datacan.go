package datacan

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"io/ioutil"
	"net/http"

	"github.com/lahemi/conglo/scraps"
)

func DBInit(config map[string]string) error {
	db, err := sql.Open("sqlite3", config["DBfile"])
	if err != nil {
		return err
	}
	defer db.Close()
	_, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS musicks (
            id INTEGER NOT NULL PRIMARY KEY,
            url TEXT,
            artist TEXT,
            title TEXT
        )`)
	if err != nil {
		return err
	}
	return nil
}

func Save(config map[string]string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		DBInit(config) // In case the DB file is removed while server is running..
		db, err := sql.Open("sqlite3", config["DBfile"])
		if err != nil {
			scraps.PrintErr(err)
			return
		}
		defer db.Close()

		isBlank := func(s string) string {
			if s == "" {
				return "blank"
			}
			return s
		}
		lenLimit := func(s string) string {
			if len(s) > 100 {
				return "blank"
			}
			return s
		}
		url := lenLimit(isBlank(r.FormValue("url")))
		artist := lenLimit(isBlank(r.FormValue("artist")))
		title := lenLimit(isBlank(r.FormValue("title")))

		stmt, err := db.Prepare(`INSERT INTO musicks(url, artist, title) VALUES(?, ?, ?)`)
		if err != nil {
			scraps.PrintErr(err)
			return
		}
		defer stmt.Close()

		_, err = stmt.Exec(url, artist, title)
		if err != nil {
			scraps.PrintErr(err)
			return
		}
		http.Redirect(w, r, config["base"], http.StatusFound)
	}
}

func View(config map[string]string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		db, err := sql.Open("sqlite3", config["DBfile"])
		if err != nil {
			scraps.PrintErr(err)
			return
		}
		defer db.Close()

		rows, err := db.Query(`SELECT url, artist, title FROM musicks`)
		if err != nil {
			scraps.PrintErr(err)
			return
		}
		defer rows.Close()
		for rows.Next() {
			var url, artist, title string
			rows.Scan(&url, &artist, &title)
			fmt.Fprintf(w, "%s\n%s\n%s\n\n", url, artist, title)
		}
		rows.Close()
	}
}

func Index(config map[string]string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cnt, err := ioutil.ReadFile(config["datacanIndex"])
		if err != nil {
			scraps.PrintErr(err)
			return
		}
		fmt.Fprint(w, string(cnt))
	}
}
