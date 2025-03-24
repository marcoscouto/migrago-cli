package errors

import "fmt"

var (
	ErrCreateMigrationsDir = fmt.Errorf("failed to create migrations directory")
	ErrGetNextMigrationNum = fmt.Errorf("failed to get next migration number")
	ErrCreateMigrationFile = fmt.Errorf("failed to create migration file")
	ErrOpenDbConnection    = fmt.Errorf("failed to open database connection")
)
