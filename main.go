package main

import (
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	"github.com/marcoscouto/migrago-cli/internal/action"
	"github.com/marcoscouto/migrago-cli/internal/command"
	"github.com/spf13/cobra"
)

func main() {
		createAction := action.NewCreate()
		executeAction := action.NewExecute()

		startCommand := command.NewStart(
			createAction,
			executeAction,
		)

	var rootCmd = &cobra.Command{
		Use:   "migrago",
		Short: "Migrago is a database migration tool",
		Long:  `A simple database migration tool written in go that helps you manage your database schema changes.`,
	}

	var startCmd = &cobra.Command{
		Use:   "start",
		Short: "Start the migration tool",
		Run: startCommand.Start,
	}

	rootCmd.AddCommand(startCmd)
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
