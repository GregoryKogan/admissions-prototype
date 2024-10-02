package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/jackc/pgx"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// TODO: use configuration to store the database connection parameters
	// connect to the database
	conn, err := pgx.Connect(pgx.ConnConfig{Host: "database", Port: 5432, User: "l2shdev", Password: "l2sh", Database: "admissions"})
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
	api.GET("/hello", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	e.Logger.Fatal(e.Start(":8888"))
}

func getGreeting(conn *pgx.Conn) (string, error) {
	var greeting string
	if err := conn.QueryRow("select 'Hello, world!'").Scan(&greeting); err != nil {
		return "", err
	}
	return greeting, nil
}
