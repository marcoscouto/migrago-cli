package main

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	"github.com/marcoscouto/migrago"
	"github.com/spf13/cobra"
)

const migrationsDir = "migrations"

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

func createMigration(name string) error {
	if err := os.MkdirAll(migrationsDir, 0755); err != nil {
		return fmt.Errorf("failed to create migrations directory: %v", err)
	}

	nextNum, err := getNextMigrationNumber()
	if err != nil {
		return fmt.Errorf("failed to get next migration number: %v", err)
	}

	sanitizedName := strings.ReplaceAll(name, " ", "_")
	filename := fmt.Sprintf("%d_%s.sql", nextNum, sanitizedName)
	filepath := filepath.Join(migrationsDir, filename)

	if err := os.WriteFile(filepath, []byte(""), 0644); err != nil {
		return fmt.Errorf("failed to create migration file: %v", err)
	}

	return nil
}

type DatabaseConfig struct {
	Driver   string
	Host     string
	Port     string
	Database string
	Username string
	Password string
}

func executeMigrations(config DatabaseConfig) error {
	var dsn string
	switch config.Driver {
	case "postgres":
		dsn = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			config.Host, config.Port, config.Username, config.Password, config.Database)
	case "mysql":
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
			config.Username, config.Password, config.Host, config.Port, config.Database)
	}

	db, err := sql.Open(config.Driver, dsn)
	if err != nil {
		return fmt.Errorf("failed to open database connection: %v", err)
	}

	defer db.Close()

	migrator := migrago.New(db, config.Driver)
	return migrator.ExecuteMigrations(migrationsDir)
}

func getDatabaseConfig() (DatabaseConfig, error) {
	config := DatabaseConfig{}

	questions := []*survey.Question{
		{
			Name: "driver",
			Prompt: &survey.Select{
				Message: "Choose database driver:",
				Options: []string{"postgres", "mysql"},
			},
		},
		{
			Name: "host",
			Prompt: &survey.Input{
				Message: "Enter database host:",
				Default: "localhost",
			},
		},
		{
			Name: "port",
			Prompt: &survey.Input{
				Message: "Enter database port:",
			},
		},
		{
			Name: "database",
			Prompt: &survey.Input{
				Message: "Enter database name:",
			},
		},
		{
			Name: "username",
			Prompt: &survey.Input{
				Message: "Enter database username:",
			},
		},
		{
			Name: "password",
			Prompt: &survey.Password{
				Message: "Enter database password:",
			},
		},
	}

	err := survey.Ask(questions, &config)
	return config, err
}

func main() {
	var rootCmd = &cobra.Command{
		Use:   "migrago",
		Short: "MigraGo is a database migration tool",
		Long:  `A simple database migration tool written in Go that helps you manage your database schema changes.`,
	}

	var startCmd = &cobra.Command{
		Use:   "start",
		Short: "Start the migration tool",
		Run: func(cmd *cobra.Command, args []string) {
			options := []string{"create", "execute"}
			var choice string

			prompt := &survey.Select{
				Message: "Choose an action:",
				Options: options,
			}

			if err := survey.AskOne(prompt, &choice); err != nil {
				fmt.Printf("Error: %v\n", err)
				os.Exit(1)
			}

			switch choice {
			case "create":
				var name string
				namePrompt := &survey.Input{
					Message: "Enter migration name:",
				}

				if err := survey.AskOne(namePrompt, &name); err != nil {
					fmt.Printf("Error: %v\n", err)
					os.Exit(1)
				}

				if err := createMigration(name); err != nil {
					fmt.Printf("Error: %v\n", err)
					os.Exit(1)
				}

				fmt.Println("Migration created successfully!")

			case "execute":
				config, err := getDatabaseConfig()
				if err != nil {
					fmt.Printf("Error getting database config: %v\n", err)
					os.Exit(1)
				}

				if err := executeMigrations(config); err != nil {
					fmt.Printf("Error executing migrations: %v\n", err)
					os.Exit(1)
				}

				fmt.Println("Migrations executed successfully!")
			}
		},
	}

	rootCmd.AddCommand(startCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
