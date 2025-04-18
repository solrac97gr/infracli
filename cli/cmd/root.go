package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// RootCmd representa el comando base cuando se llama sin subcomandos
var RootCmd = &cobra.Command{
	Use:   "infracli",
	Short: "CLI tool to manage infrastructure services",
	Long: `InfraCLI is a command line tool that simplifies the management of
infrastructure services defined in Docker Compose files.

It allows running and stopping multiple services at once from a centralized CLI.
The tool automatically detects available services based on the directory structure.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Si no se proporciona ningún subcomando, mostrar la ayuda
		cmd.Help()
	},
}

func init() {
	// Aquí se pueden agregar flags globales si se necesitan
	RootCmd.PersistentFlags().BoolP("verbose", "v", false, "Enable verbose output")

	fmt.Println(`
██╗███╗   ██╗███████╗██████╗  █████╗  ██████╗██╗     ██╗
██║████╗  ██║██╔════╝██╔══██╗██╔══██╗██╔════╝██║     ██║
██║██╔██╗ ██║█████╗  ██████╔╝███████║██║     ██║     ██║
██║██║╚██╗██║██╔══╝  ██╔══██╗██╔══██║██║     ██║     ██║
██║██║ ╚████║██║     ██║  ██║██║  ██║╚██████╗███████╗██║
╚═╝╚═╝  ╚═══╝╚═╝     ╚═╝  ╚═╝╚═╝  ╚═╝ ╚═════╝╚══════╝╚═╝`)
	fmt.Println()
}
