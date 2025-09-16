package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

type Application struct {
	logger *slog.Logger
}

func main() {
	// server setup

	// database
	userName := "app"
	password := "pass"
	dbHost := "192.168.0.134"
	dbName := "snippetbox"

	// get configs
	addr := flag.String("addr", ":8080", "http server address:port")
	dsn := flag.String("dsn", fmt.Sprintf("%s:%s@%s/%s?parseTime=true", userName, password, dbHost, dbName), "db connection string (dsn)")
	flag.Parse()

	var loggerOptions *slog.HandlerOptions = &slog.HandlerOptions{
		AddSource: false,
	}
	logger := slog.New(slog.NewJSONHandler(os.Stdout, loggerOptions))

	app := Application{
		logger: logger,
	}

	db, err := openDB(*dsn)

	if err != nil {
		app.logger.Error(err.Error())
		os.Exit(1)
	}

	defer db.Close()

	app.logger.Info("server started", slog.Any("address", *addr))

	app.check(http.ListenAndServe(*addr, app.routes()))
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()

	if err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}
