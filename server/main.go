package main

import (
	"flag"
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/mysqlstore"
	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"
	_ "github.com/go-sql-driver/mysql"
	"github.com/mahmoud-shabban/snippetbox/internal/models"
)

type Application struct {
	logger         *slog.Logger
	snippets       *models.SnippetModel
	users          *models.UserModel
	templateCache  map[string]*template.Template
	formDecoder    *form.Decoder
	sessionManager *scs.SessionManager
}

func main() {

	// configs

	// database
	userName := "app"
	password := "pass"
	dbHost := "127.0.0.1" //"192.168.0.134"
	dbName := "snippetbox"

	// server data
	addr := flag.String("addr", ":8080", "http server address:port")
	dsn := flag.String("dsn", fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true", userName, password, dbHost, dbName), "db connection string (dsn)")
	flag.Parse()

	// logger
	var loggerOptions *slog.HandlerOptions = &slog.HandlerOptions{
		AddSource: true,
	}
	logger := slog.New(slog.NewJSONHandler(os.Stdout, loggerOptions))

	// database
	db, err := openDB(*dsn)

	if err != nil {
		logger.Error(err.Error(), slog.Any("source", "database"))
		os.Exit(1)
	}
	defer db.Close()

	// initialize template cache cache
	cache, err := newTemplateCache()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	// form decoder
	decoder := form.NewDecoder()

	// session manager
	sessionManager := scs.New()
	sessionManager.Store = mysqlstore.New(db)
	sessionManager.Lifetime = 12 * time.Hour

	// Initialize the APP and inject its dependencies
	app := Application{
		logger:         logger,
		snippets:       &models.SnippetModel{DB: db},
		users:          &models.UserModel{DB: db},
		templateCache:  cache,
		formDecoder:    decoder,
		sessionManager: sessionManager,
	}

	// http server
	serv := &http.Server{
		Addr:        *addr,
		Handler:     app.routes(),
		ErrorLog:    slog.NewLogLogger(logger.Handler(), slog.LevelError),
		IdleTimeout: 1 * time.Minute,
	}

	app.logger.Info("successfully connected to database")
	app.logger.Info("successfully initialized template cache")
	app.logger.Info("server starting at", slog.Any("address", *addr))
	// Start the Server
	// app.check(http.ListenAndServe(*addr, app.routes()))
	app.check(serv.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem"))
}
