package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/svanhalla/ib-ui/internal/repository/filerepo"

	"github.com/svanhalla/ib-ui/internal/repository"

	"github.com/alexedwards/scs/v2"
)

const version = "1.0.0"
const ccsVersion = "1"

var session *scs.SessionManager

type config struct {
	port int
	env  string
	API  string
	db   struct {
		dsn string
	}
	stripe struct {
		secret string
		key    string
	}
}

type application struct {
	config        config
	infoLog       *log.Logger
	errorLog      *log.Logger
	templateCache map[string]*template.Template
	version       string
	Repo          repository.Repo
	Session       *scs.SessionManager
}

func (app *application) serve() error {
	srv := http.Server{
		Addr:              fmt.Sprintf(":%d", app.config.port),
		Handler:           app.routes(),
		IdleTimeout:       30 * time.Second,
		ReadTimeout:       10 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      5 * time.Second,
	}

	app.infoLog.Printf("Starting HTTP server in %q mode on port %d", app.config.env, app.config.port)

	return srv.ListenAndServe()
}

func main() {
	// gob.Register(TransactionData{})

	var cfg config

	flag.IntVar(&cfg.port, "port", 4000, "Server port to listen on")
	flag.StringVar(&cfg.env, "env", "development", "Application environment {development|production}")
	flag.StringVar(&cfg.db.dsn, "dsn", "hans:secret@tcp(localhost:3306)/widgets?parseTime=true&tls=false", "dsn")
	flag.StringVar(&cfg.API, "api", "http://localhost:4000", "Url to api")

	flag.Parse()

	cfg.stripe.key = os.Getenv("STRIPE_KEY")
	cfg.stripe.secret = os.Getenv("STRIPE_SECRET")

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	// conn, err := driver.OpenDB(cfg.db.dsn)
	// if err != nil {
	//	errorLog.Fatal(err)
	// }
	//
	// defer func() {
	//	_ = conn.Close()
	// }()

	// setup session
	session = scs.New()
	session.Lifetime = 24 * time.Hour

	tc := make(map[string]*template.Template)

	app := &application{
		config:        cfg,
		infoLog:       infoLog,
		errorLog:      errorLog,
		templateCache: tc,
		version:       version,
		Repo: &filerepo.FileModel{
			Dir: "~/.ib",
		},
		Session: session,
	}

	err := app.serve()
	if err != nil {
		app.errorLog.Println(err)
		log.Fatal(err)
	}
}
