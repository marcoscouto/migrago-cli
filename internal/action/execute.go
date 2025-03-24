package action

import (
	"database/sql"
	"fmt"

	"github.com/marcoscouto/migrago"
	"github.com/marcoscouto/migrago-cli/internal/data"
	"github.com/marcoscouto/migrago-cli/internal/errors"
)

const (
	strConnPostgres = "host=%s port=%s user=%s password=%s dbname=%s sslmode=disable"
	strConnMysql    = "%s:%s@tcp(%s:%s)/%s"
)

type Execute interface {
	ExecuteMigrations(config data.DatabaseConfig) error
}

type execute struct{}

func NewExecute() Execute {
	return &execute{}
}

func (e *execute) ExecuteMigrations(config data.DatabaseConfig) error {
	var dsn string

	switch config.Driver {
	case data.Postgres:
		dsn = fmt.Sprintf(strConnPostgres,
			config.Host, config.Port, config.Username, config.Password, config.Database)
	case data.Mysql:
		dsn = fmt.Sprintf(strConnMysql,
			config.Username, config.Password, config.Host, config.Port, config.Database)
	}

	db, err := sql.Open(config.Driver, dsn)
	if err != nil {
		return errors.ErrOpenDbConnection
	}

	defer db.Close()

	migrator := migrago.New(db, config.Driver)
	return migrator.ExecuteMigrations(migrationsDir)
}
