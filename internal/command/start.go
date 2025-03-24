package command

import (
	"fmt"
	"os"

	"github.com/AlecAivazis/survey/v2"
	"github.com/marcoscouto/migrago-cli/internal/action"
	"github.com/marcoscouto/migrago-cli/internal/data"
	"github.com/spf13/cobra"
)

type Start interface {
	Start(cmd *cobra.Command, args []string)
}

type start struct {
	create  action.Create
	execute action.Execute
}

func NewStart(create action.Create, execute action.Execute) Start {
	return &start{
		create:  create,
		execute: execute,
	}
}

func (s *start) Start(cmd *cobra.Command, args []string) {

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

		if err := s.create.CreateMigration(name); err != nil {
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

		if err := s.execute.ExecuteMigrations(config); err != nil {
			fmt.Printf("Error executing migrations: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("Migrations executed successfully!")
	}
}

func getDatabaseConfig() (data.DatabaseConfig, error) {
	questions := []*survey.Question{
		{
			Name: "driver",
			Prompt: &survey.Select{
				Message: "Choose database driver:",
				Options: []string{data.Mysql, data.Postgres},
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
				Default: "3306",
			},
		},
		{
			Name: "database",
			Prompt: &survey.Input{
				Message: "Enter database name:",
				Default: "migragocli",
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

	var config data.DatabaseConfig
	err := survey.Ask(questions, &config)
	return config, err
}
