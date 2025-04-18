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

var downCmd = &cobra.Command{
	Use:   "down [service1] [service2] ... or 'all'",
	Short: "Stop one or more infrastructure services",
	Long: `Stop one or more infrastructure services using docker-compose.
If 'all' is specified, it stops all available services.

Examples:
  infracli down mysql
  infracli down mongo elasticsearch-kibana
  infracli down all`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("Error: You must specify at least one service or 'all'")
			cmd.Help()
			return
		}

		// Verificar si se debe eliminar volúmenes
		removeVolumes, _ := cmd.Flags().GetBool("volumes")
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

		// Determinar la ruta base para los servicios
		execPath, err := os.Executable()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting executable path: %v\n", err)
			return
		}

		basePath := filepath.Join(filepath.Dir(execPath), cfg.ServicesPath)
		
		// Verificar si estamos en desarrollo
		if _, err := os.Stat(basePath); os.IsNotExist(err) {
			// En desarrollo, usar ruta relativa al directorio de trabajo
			basePath = cfg.ServicesPath
		}

		if verbose {
			fmt.Printf("Services path: %s\n", basePath)
			fmt.Printf("Available services: %s\n", strings.Join(availableServices, ", "))
			if removeVolumes {
				fmt.Println("Volumes will be removed")
			}
		}

		// Comprobar si queremos detener todos los servicios
		if len(args) == 1 && args[0] == "all" {
			fmt.Println("Stopping all available services...")
			stopAllServices(availableServices, basePath, removeVolumes, verbose)
			return
		}

		// Detener los servicios especificados
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

			stopService(service, basePath, removeVolumes, verbose)
		}
	},
}

func stopService(service, basePath string, removeVolumes, verbose bool) {
	servicePath := filepath.Join(basePath, service)
	fmt.Printf("Stopping %s...\n", service)
	
	args := []string{"down"}
	if removeVolumes {
		args = append(args, "-v")
	}
	
	cmd := exec.Command("docker-compose", args...)
	cmd.Dir = servicePath
	
	if verbose {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	} else {
		output, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("Error stopping %s: %v\n", service, err)
			fmt.Println(string(output))
			return
		}
	}
	
	fmt.Printf("%s stopped successfully\n", service)
}

func stopAllServices(services []string, basePath string, removeVolumes, verbose bool) {
	for _, service := range services {
		stopService(service, basePath, removeVolumes, verbose)
	}
	fmt.Println("All services have been stopped")
}

func init() {
	downCmd.Flags().BoolP("volumes", "v", false, "Remove volumes when stopping services")
	RootCmd.AddCommand(downCmd)
}