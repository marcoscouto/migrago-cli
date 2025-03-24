package action

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/marcoscouto/migrago-cli/internal/errors"
)

const (
	migrationsDir = "migrations"
	permissions   = 0755
)

type Create interface {
	CreateMigration(name string) error
}

type create struct{}

func NewCreate() Create {
	return &create{}
}

func (c *create) CreateMigration(name string) error {
	if err := os.MkdirAll(migrationsDir, permissions); err != nil {
		return errors.ErrCreateMigrationsDir
	}

	nextNum, err := getNextMigrationNumber()
	if err != nil {
		return errors.ErrGetNextMigrationNum
	}

	sanitizedName := strings.ReplaceAll(name, " ", "_")
	filename := fmt.Sprintf("%d_%s.sql", nextNum, sanitizedName)
	filepath := filepath.Join(migrationsDir, filename)

	if err := os.WriteFile(filepath, []byte(""), 0644); err != nil {
		return errors.ErrCreateMigrationFile
	}

	return nil
}

func getNextMigrationNumber() (int, error) {
	if _, err := os.Stat(migrationsDir); os.IsNotExist(err) {
		return 1, nil
	}

	files, err := os.ReadDir(migrationsDir)
	if err != nil {
		return 0, err
	}

	maxNum := 0
	pattern := regexp.MustCompile(`^(\d+)_.*\.sql$`)

	for _, file := range files {
		matches := pattern.FindStringSubmatch(file.Name())
		if len(matches) > 1 {
			if num, err := strconv.Atoi(matches[1]); err == nil {
				if num > maxNum {
					maxNum = num
				}
			}
		}
	}

	return maxNum + 1, nil
}
