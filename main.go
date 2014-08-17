package main

import (
	"bytes"
	"crypto/sha1"
	"database/sql"
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/elazarl/go-bindata-assetfs"
	_ "github.com/mattn/go-sqlite3"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"os"
	"runtime"
	"strings"
	"time"
)

type MemoDecorator struct {
	AccessKey string
	Body      string
	CreatedAt string
	UpdatedAt string
	Error     string
}

var formTmpl = template.Must(template.New("form").Parse(AssetFiles("views/base.html", "views/form.html")))
var memoTmpl = template.Must(template.New("memo").Parse(AssetFiles("views/base.html", "views/memo.html")))
var db *sql.DB

func AssetFiles(filenames ...string) string {
	var buffer bytes.Buffer
	for _, filename := range filenames {
		src, err := Asset(filename)
		if err != nil {
			log.Fatal(err)
		}
		buffer.WriteString(string(src))
	}
	return buffer.String()
}

func dbInit(filename string) *sql.DB {
	var db *sql.DB
	var err error
	db, _ = sql.Open("sqlite3", filename)
	if _, err = os.Stat(filename); os.IsNotExist(err) {
		sqlStmt := `
    CREATE TABLE IF NOT EXISTS memos (
    	id INTEGER PRIMARY KEY AUTOINCREMENT,
    	access_key TEXT NOT NULL,
    	body TEXT NOT NULL,
    	created_at INTEGER NOT NULL,
    	updated_at INTEGER NOT NULL
    );
    CREATE UNIQUE INDEX i1 ON memos (access_key);
		`
		_, err = db.Exec(sqlStmt)
		if err != nil {
			log.Printf("%q: %s\n", err, sqlStmt)
		}
	}
	return db
}

func formHandler(w http.ResponseWriter, r *http.Request) {
	formTmpl.ExecuteTemplate(w, "base", nil)
}

func createHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.Header().Set("Allow", "POST")
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	body := r.FormValue("body")
	if body == "" {
		http.Error(w, "Body is required", http.StatusBadRequest)
		return
	}

	now := time.Now().Unix()
	rand.Seed(now)
	data := []byte(fmt.Sprintf("%d -- %s -- %f", now, body, rand.Intn(10000)))
	key := fmt.Sprintf("%x", sha1.Sum(data))

	query := "INSERT INTO memos (access_key, body, created_at, updated_at) VALUES (?, ?, ?, ?)"
	_, err := db.Exec(query, key, body, now, now)
	if err != nil {
		var memo MemoDecorator
		memo.Error = err.Error()
		formTmpl.ExecuteTemplate(w, "base", memo)
		return
	}

	http.Redirect(w, r, "/memos/"+key, http.StatusFound)
}

func memoHandler(w http.ResponseWriter, r *http.Request) {
	key := strings.TrimLeft(r.URL.Path, "/memos/")
	if r.Method == "DELETE" {
		var id int64
		err := db.QueryRow("SELECT id FROM memos WHERE access_key = ?", key).Scan(&id)
		switch {
		case err == sql.ErrNoRows:
			http.NotFound(w, r)
			return
		case err != nil:
			log.Fatal(err)
		}

		_, err = db.Exec("DELETE FROM memos WHERE access_key = ?", key)
		if err != nil {
			log.Fatal(err)
		}
		http.Redirect(w, r, "/", http.StatusFound)
	} else if r.Method == "GET" {
		var memo MemoDecorator
		var createdAt int64
		query := "SELECT body, created_at FROM memos WHERE access_key = ?"
		err := db.QueryRow(query, key).Scan(&memo.Body, &createdAt)
		switch {
		case err == sql.ErrNoRows:
			http.NotFound(w, r)
			return
		case err != nil:
			log.Fatal(err)
		}
		memo.CreatedAt = time.Unix(createdAt, 0).Format("2006/01/02 15:04")

		memoTmpl.ExecuteTemplate(w, "base", memo)
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

var (
	flHelp        bool
	flVersion     bool
	flHost        string
	flPort        string
	flDatabaseUrl string
)

func main() {
	app := cli.NewApp()
	app.Name = "go-pastel"
	app.Version = Version
	app.Usage = "A copy and paste sharing web application like git"
	app.Author = "Naotoshi Seo"
	app.Email = "sonots@gmail.com"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "host",
			Value: "0.0.0.0",
			Usage: "Address to serve this service",
		},
		cli.StringFlag{
			Name:  "port",
			Value: "5050",
			Usage: "Port number to serve this service",
		},
		cli.StringFlag{
			Name:  "database_url",
			Value: "pastel.db",
			Usage: "Path to sqlite storage file",
		},
	}
	app.Action = start
	app.Commands = []cli.Command{
		{
			Name:        "start",
			ShortName:   "s",
			Usage:       "Start up",
			Description: "Start the pastel server",
			Action:      start,
			Flags:       app.Flags,
		},
	}
	app.Run(os.Args)
}

func start(c *cli.Context) {
	cpus := runtime.NumCPU()
	runtime.GOMAXPROCS(cpus)

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(&assetfs.AssetFS{Asset, AssetDir, "/static/"})))
	// http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/", formHandler)
	http.HandleFunc("/create", createHandler)
	http.HandleFunc("/memos/", memoHandler)
	db = dbInit(c.String("database_url"))
	address := c.String("host") + ":" + c.String("port")
	fmt.Println("Start pastel at " + address)
	http.ListenAndServe(address, nil)
}
