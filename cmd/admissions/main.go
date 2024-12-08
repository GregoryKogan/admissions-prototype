package main

import (
	"github.com/L2SH-Dev/admissions/internal/admin"
	"github.com/L2SH-Dev/admissions/internal/config"
	"github.com/L2SH-Dev/admissions/internal/datastore"
	"github.com/L2SH-Dev/admissions/internal/logging"
	"github.com/L2SH-Dev/admissions/internal/ping"
	"github.com/L2SH-Dev/admissions/internal/regdata"
	"github.com/L2SH-Dev/admissions/internal/server"
	"github.com/L2SH-Dev/admissions/internal/users"
)

func main() {
	config.Init()
	logging.Init()

	storage := datastore.InitStorage()

	srv := server.NewServer(storage)

	srv.AddFrontend("ui/dist")
	srv.AddHandlers(
		ping.NewPingHandler,
		users.NewUsersHandler,
		regdata.NewRegistrationDataHandler,
	)

	admin.CreateDefaultAdmin(storage)

	srv.Start()
}
