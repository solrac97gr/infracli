package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Config contiene la configuración para la herramienta InfraCLI
type Config struct {
	ServicesPath string   `json:"servicesPath"`
	ExcludedDirs []string `json:"excludedDirs"`
}

// LoadConfig carga la configuración desde el archivo config.json
func LoadConfig() (*Config, error) {
	// Obtener la ruta del ejecutable para localizar config.json relativo al binario
	execPath, err := os.Executable()
	if err != nil {
		return nil, fmt.Errorf("error getting executable path: %v", err)
	}
	
	// Durante el desarrollo, utilizamos una ruta relativa
	configPath := filepath.Join(filepath.Dir(execPath), "../config/config.json")
	
	// Verificar si estamos en desarrollo (el archivo no existirá en la ubicación relativa)
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		// En desarrollo, usamos una ruta relativa al directorio de trabajo
		configPath = "config/config.json"
	}
	
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("error reading config file: %v", err)
	}
	
	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("error parsing config file: %v", err)
	}
	
	return &config, nil
}

// GetAvailableServices devuelve una lista de servicios disponibles
func GetAvailableServices() ([]string, error) {
	config, err := LoadConfig()
	if err != nil {
		return nil, err
	}
	
	// Determinar la ruta base para buscar servicios
	execPath, err := os.Executable()
	if err != nil {
		return nil, fmt.Errorf("error getting executable path: %v", err)
	}
	
	basePath := filepath.Join(filepath.Dir(execPath), config.ServicesPath)
	
	// Verificar si estamos en desarrollo (si no existe la ruta con el ejecutable)
	if _, err := os.Stat(basePath); os.IsNotExist(err) {
		// En desarrollo, usar ruta relativa al directorio de trabajo
		basePath = config.ServicesPath
	}
	
	// Leer los directorios en la ubicación configurada
	files, err := os.ReadDir(basePath)
	if err != nil {
		return nil, fmt.Errorf("error reading services directory: %v", err)
	}
	
	var services []string
	for _, file := range files {
		if file.IsDir() {
			// Comprobar si el directorio no está excluido
			excluded := false
			for _, excl := range config.ExcludedDirs {
				if file.Name() == excl || file.Name() == "infracli" {
					excluded = true
					break
				}
			}
			
			// Si no está excluido y contiene un docker-compose.yml, agregarlo a la lista
			if !excluded {
				dockerComposePath := filepath.Join(basePath, file.Name(), "docker-compose.yml")
				if _, err := os.Stat(dockerComposePath); err == nil {
					services = append(services, file.Name())
				}
			}
		}
	}
	
	return services, nil
}