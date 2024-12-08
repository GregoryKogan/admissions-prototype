package admin

import (
	"errors"
	"time"

	"github.com/L2SH-Dev/admissions/internal/datastore"
	"github.com/L2SH-Dev/admissions/internal/regdata"
	"github.com/L2SH-Dev/admissions/internal/users"
	"github.com/L2SH-Dev/admissions/internal/users/auth"
	"github.com/L2SH-Dev/admissions/internal/users/auth/passwords"
	"github.com/L2SH-Dev/admissions/internal/users/roles"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

func CreateDefaultAdmin(storage datastore.Storage) {
	usersService := users.NewUsersService(
		users.NewUsersRepo(storage),
		roles.NewRolesService(roles.NewRolesRepo(storage)),
	)
	passwordsService := passwords.NewPasswordsService(passwords.NewPasswordsRepo(storage))
	authService := auth.NewAuthService(auth.NewAuthRepo(storage), passwordsService)
	regdataService := regdata.NewRegistrationDataService(
		regdata.NewRegistrationDataRepo(storage),
		usersService,
		authService,
		passwordsService,
	)

	// Check if admin already exists
	admin, err := usersService.GetByLogin("admin")
	if err == nil && admin != nil {
		return
	} else if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		panic(err)
	}

	registrationData := regdata.RegistrationData{
		Email:           viper.GetString("admin.email"),
		EmailVerified:   true,
		FirstName:       "Админ",
		LastName:        "Админов",
		Gender:          "N",
		BirthDate:       time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
		Grade:           6,
		OldSchool:       "Лицей \"Первая школа\"",
		ParentFirstName: "Разработчик",
		ParentLastName:  "Админов",
		ParentPhone:     "+70000000000",
	}

	if err := regdataService.Create(&registrationData); err != nil {
		panic(err)
	}

	admin, err = usersService.CreateDefaultAdmin(registrationData.ID)
	if err != nil {
		panic(err)
	}

	if err := authService.Register(admin.ID, viper.GetString("secrets.admin_password")); err != nil {
		panic(err)
	}
}
