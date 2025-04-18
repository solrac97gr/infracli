package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/solrac97gr/infrastructure/infracli/config"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all available infrastructure services",
	Long: `List all available infrastructure services that can be managed by this CLI.
These are the services that can be used with the run and down commands.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Obtener servicios disponibles
		availableServices, err := config.GetAvailableServices()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			return
		}

		if len(availableServices) == 0 {
			fmt.Println("No services found. Check your configuration.")
			return
		}

		fmt.Println("Available services:")
		fmt.Println(strings.Repeat("-", 20))
		for _, service := range availableServices {
			fmt.Printf("- %s\n", service)
		}
		fmt.Println(strings.Repeat("-", 20))
		fmt.Printf("Total: %d services\n", len(availableServices))
		fmt.Println("\nYou can run any of these services with: infracli run <service-name>")
		fmt.Println("You can stop any of these services with: infracli down <service-name>")
		fmt.Println("You can manage all services at once with: infracli run all or infracli down all")
	},
}

func init() {
	RootCmd.AddCommand(listCmd)
}
