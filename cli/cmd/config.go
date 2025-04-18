package cmd

import (
	"fmt"
	"os"

	"github.com/solrac97gr/infrastructure/infracli/config"
	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage InfraCLI configuration",
	Long: `View or update InfraCLI configuration settings.
This command allows you to see the current configuration and where it's stored.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Cargar la configuración actual
		cfg, err := config.LoadConfig()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error loading configuration: %v\n", err)
			return
		}

		// Obtener la ruta del archivo de configuración
		configPath, err := config.GetConfigFilePath()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting config path: %v\n", err)
			return
		}

		fmt.Println("Current InfraCLI Configuration:")
		fmt.Println("-------------------------------")
		fmt.Printf("Configuration file: %s\n\n", configPath)
		fmt.Printf("Services path: %s\n", cfg.ServicesPath)
		fmt.Printf("Excluded directories: %v\n", cfg.ExcludedDirs)

		fmt.Println("\nTo modify the configuration, edit the file directly or use:")
		fmt.Println("  infracli config set-path <new-services-path>")
	},
}

var configSetPathCmd = &cobra.Command{
	Use:   "set-path [path]",
	Short: "Update the services path in the configuration",
	Long: `Update the path where InfraCLI should look for services.
This path should point to the directory containing your infrastructure service directories.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		newPath := args[0]

		// Cargar la configuración actual
		cfg, err := config.LoadConfig()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error loading configuration: %v\n", err)
			return
		}

		// Actualizar la ruta de servicios
		oldPath := cfg.ServicesPath
		cfg.ServicesPath = newPath

		// Guardar la configuración
		if err := config.SaveConfig(cfg); err != nil {
			fmt.Fprintf(os.Stderr, "Error saving configuration: %v\n", err)
			return
		}

		fmt.Printf("Services path updated from '%s' to '%s'\n", oldPath, newPath)
	},
}

func init() {
	configCmd.AddCommand(configSetPathCmd)
	RootCmd.AddCommand(configCmd)
}
