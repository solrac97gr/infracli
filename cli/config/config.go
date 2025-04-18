package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

const (
	// ConfigFileName es el nombre del archivo de configuración
	ConfigFileName = "infracli.json"
	// ConfigDirName es el nombre del directorio de configuración dentro de ~/.config
	ConfigDirName = "infracli"
)

// Config contiene la configuración para la herramienta InfraCLI
type Config struct {
	ServicesPath string   `json:"servicesPath"`
	ExcludedDirs []string `json:"excludedDirs"`
}

// GetDefaultConfig devuelve una configuración por defecto
func GetDefaultConfig() *Config {
	return &Config{
		ServicesPath: filepath.Join(os.Getenv("HOME"), "Development", "infrastructure", "services"),
		ExcludedDirs: []string{"config", "scripts", "cmd"},
	}
}

// GetConfigDir devuelve la ruta del directorio de configuración
func GetConfigDir() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("error getting home directory: %v", err)
	}
	
	// Ruta completa del directorio de configuración (~/.config/infracli)
	configDir := filepath.Join(homeDir, ".config", ConfigDirName)
	
	return configDir, nil
}

// GetConfigFilePath devuelve la ruta completa del archivo de configuración
func GetConfigFilePath() (string, error) {
	configDir, err := GetConfigDir()
	if err != nil {
		return "", err
	}
	
	return filepath.Join(configDir, ConfigFileName), nil
}

// SaveConfig guarda la configuración en el archivo
func SaveConfig(config *Config) error {
	configDir, err := GetConfigDir()
	if err != nil {
		return err
	}
	
	// Crear el directorio si no existe
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("error creating config directory: %v", err)
	}
	
	configPath, err := GetConfigFilePath()
	if err != nil {
		return err
	}
	
	// Serializar la configuración a JSON con formato legible
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("error serializing config: %v", err)
	}
	
	// Escribir en el archivo
	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("error writing config file: %v", err)
	}
	
	return nil
}

// LoadConfig carga la configuración desde el archivo de configuración
func LoadConfig() (*Config, error) {
	// Obtener la ruta del archivo de configuración
	configPath, err := GetConfigFilePath()
	if err != nil {
		return nil, err
	}
	
	// Comprobar si el archivo existe
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		// Si no existe, creamos una configuración por defecto
		config := GetDefaultConfig()
		
		// Guardamos la configuración en el archivo
		if err := SaveConfig(config); err != nil {
			return nil, fmt.Errorf("error creating default config: %v", err)
		}
		
		fmt.Println("Created default configuration at:", configPath)
		return config, nil
	}
	
	// Leer el archivo existente
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
	
	// Usar directamente la ruta de servicios configurada
	basePath := config.ServicesPath
	
	// Expandir la ruta si contiene ~/
	if len(basePath) >= 2 && basePath[:2] == "~/" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("error getting home directory: %v", err)
		}
		basePath = filepath.Join(homeDir, basePath[2:])
	}
	
	// Verificar que el directorio existe
	if _, err := os.Stat(basePath); os.IsNotExist(err) {
		return nil, fmt.Errorf("services directory not found: %s", basePath)
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
				if file.Name() == excl {
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