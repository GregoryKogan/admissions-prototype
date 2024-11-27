package main

import (
	"github.com/L2SH-Dev/admissions/internal/config"
	"github.com/L2SH-Dev/admissions/internal/datastore"
	"github.com/L2SH-Dev/admissions/internal/logging"
	"github.com/L2SH-Dev/admissions/internal/ping"
	"github.com/L2SH-Dev/admissions/internal/server"
	"github.com/L2SH-Dev/admissions/internal/users"
)

func main() {
	config.Init()
	logging.Init()

	srv := server.NewServer()

	storage := datastore.InitStorage()

	srv.AddFrontend("ui/dist", "ui/dist/index.html")
	srv.AddHandlers(
		storage,
		ping.NewPingHandler,
		users.NewUsersHandler,
	)

	srv.Start()
}
