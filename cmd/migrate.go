package cmd

import (
	"fmt"
	"github.com/rideziro/go-storage/storage"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var MigrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Relational database migration",
	RunE: func(cmd *cobra.Command, args []string) error {
		migrationPath := viper.GetString("MIGRATION_PATH")
		migrator, err := storage.NewMigrator(migrationPath)
		if err != nil {
			return err
		}
		err = migrator.Up()
		if err != nil {
			return err
		}

		fmt.Println("Successfully migrated files")

		return nil
	},
}
