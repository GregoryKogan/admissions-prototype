package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/L2SH-Dev/admissions/internal/config"
	"github.com/L2SH-Dev/admissions/internal/ping"
	"github.com/L2SH-Dev/admissions/internal/secrets"
	"github.com/jackc/pgx"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/spf13/viper"
)

func main() {
	config.Init()

	// connect to the database
	dbConfig := viper.Sub("database")
	db_password, err := secrets.ReadSecret("db_password")
	if err != nil {
		log.Fatalf("Failed to read database password: %v", err)
	}
	conn, err := pgx.Connect(pgx.ConnConfig{
		Host:     dbConfig.GetString("host"),
		Port:     uint16(dbConfig.GetInt("port")),
		User:     dbConfig.GetString("user"),
		Password: db_password,
		Database: dbConfig.GetString("name"),
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close()

	// ping the database to check if it's alive
	if err := conn.Ping(context.Background()); err != nil {
		fmt.Fprintf(os.Stderr, "Unable to ping database: %v\n", err)
		os.Exit(1)
	} else {
		fmt.Println("Connected to database")
	}

	// get greeting from the database
	greeting, err := getGreeting(conn)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to get greeting: %v\n", err)
		os.Exit(1)
	}
	fmt.Println(greeting)

	e := echo.New()

	e.Static("/", "ui/dist")
	e.File("/", "ui/dist/index.html")

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:8888"},
		AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete},
	}))

	api := e.Group("/api")
	api.GET("/ping", ping.PingHandler)

	e.Logger.Fatal(e.Start(":8888"))
}

func getGreeting(conn *pgx.Conn) (string, error) {
	var greeting string
	if err := conn.QueryRow("select 'Hello, world!'").Scan(&greeting); err != nil {
		return "", err
	}
	return greeting, nil
}
