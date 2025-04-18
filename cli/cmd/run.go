package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/solrac97gr/infrastructure/infracli/config"
	"github.com/spf13/cobra"
)

var runCmd = &cobra.Command{
	Use:   "run [service1] [service2] ... or 'all'",
	Short: "Start one or more infrastructure services",
	Long: `Start one or more infrastructure services using docker-compose.
If 'all' is specified, it starts all available services.

Examples:
  infracli run mysql
  infracli run mongo elasticsearch-kibana
  infracli run all`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("Error: You must specify at least one service or 'all'")
			cmd.Help()
			return
		}

		verbose, _ := cmd.Flags().GetBool("verbose")

		// Obtener servicios disponibles
		availableServices, err := config.GetAvailableServices()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			return
		}

		// Obtener la configuración
		cfg, err := config.LoadConfig()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error loading configuration: %v\n", err)
			return
		}

		// Usar directamente la ruta de servicios configurada
		basePath := cfg.ServicesPath

		// Expandir la ruta si contiene ~/
		if len(basePath) >= 2 && basePath[:2] == "~/" {
			homeDir, err := os.UserHomeDir()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error getting home directory: %v\n", err)
				return
			}
			basePath = filepath.Join(homeDir, basePath[2:])
		}

		if verbose {
			fmt.Printf("Services path: %s\n", basePath)
			fmt.Printf("Available services: %s\n", strings.Join(availableServices, ", "))
		}

		// Comprobar si queremos iniciar todos los servicios
		if len(args) == 1 && args[0] == "all" {
			fmt.Println("Starting all available services...")
			runAllServices(availableServices, basePath, verbose)
			return
		}

		// Iniciar los servicios especificados
		for _, service := range args {
			serviceFound := false
			for _, availService := range availableServices {
				if service == availService {
					serviceFound = true
					break
				}
			}

			if !serviceFound {
				fmt.Printf("Warning: Service '%s' not found in available services\n", service)
				fmt.Printf("Available services: %s\n", strings.Join(availableServices, ", "))
				continue
			}

			runService(service, basePath, verbose)
		}
	},
}

func runService(service, basePath string, verbose bool) {
	servicePath := filepath.Join(basePath, service)
	fmt.Printf("Starting %s...\n", service)

	cmd := exec.Command("docker-compose", "up", "-d")
	cmd.Dir = servicePath

	if verbose {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	} else {
		output, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("Error starting %s: %v\n", service, err)
			fmt.Println(string(output))
			return
		}
	}

	fmt.Printf("%s started successfully\n", service)
}

func runAllServices(services []string, basePath string, verbose bool) {
	for _, service := range services {
		runService(service, basePath, verbose)
	}
	fmt.Println("All services have been started")
}

func init() {
	RootCmd.AddCommand(runCmd)
}
