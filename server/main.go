package main

import (
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/mahmoud-shabban/snippetbox/internal/models"
)

type Application struct {
	logger   *slog.Logger
	snippets *models.SnippetModel
}

func main() {

	// configs

	// database
	userName := "app"
	password := "pass"
	dbHost := "127.0.0.1" //"192.168.0.134"
	dbName := "snippetbox"

	// server
	addr := flag.String("addr", ":8080", "http server address:port")
	dsn := flag.String("dsn", fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true", userName, password, dbHost, dbName), "db connection string (dsn)")
	flag.Parse()

	// logger
	var loggerOptions *slog.HandlerOptions = &slog.HandlerOptions{
		AddSource: false,
	}
	logger := slog.New(slog.NewJSONHandler(os.Stdout, loggerOptions))

	// database
	db, err := openDB(*dsn)

	if err != nil {
		logger.Error(err.Error(), slog.Any("source", "database"))
		os.Exit(1)
	}
	defer db.Close()

	// Initialize the APP and inject its dependencies
	app := Application{
		logger:   logger,
		snippets: &models.SnippetModel{DB: db},
	}

	app.logger.Info("successfuly connected to database")
	app.logger.Info("server started", slog.Any("address", *addr))

	// Start the Server
	app.check(http.ListenAndServe(*addr, app.routes()))
}
