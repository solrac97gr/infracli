#!/bin/bash

# Script para desinstalar infracli
# Elimina la herramienta de línea de comandos infracli y sus archivos de configuración

set -e

# Colores para mensajes
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo -e "${YELLOW}InfraCLI Uninstallation Script${NC}"
echo "This script will remove the infracli tool and its configuration"

# Comprobar ubicaciones posibles del binario
POSSIBLE_LOCATIONS=(
    "/usr/local/bin/infracli"
    "$HOME/bin/infracli"
)

BINARY_FOUND=false
for location in "${POSSIBLE_LOCATIONS[@]}"; do
    if [ -f "$location" ]; then
        echo "Found infracli binary at $location"
        BINARY_FOUND=true
        
        # Solicitar confirmación
        read -p "Do you want to remove this binary? (y/n) " -n 1 -r
        echo
        if [[ $REPLY =~ ^[Yy]$ ]]; then
            echo "Removing binary from $location"
            rm "$location"
            echo -e "${GREEN}Binary removed successfully${NC}"
        else
            echo "Binary will not be removed"
        fi
    fi
done

if [ "$BINARY_FOUND" = false ]; then
    echo -e "${YELLOW}No infracli binary found in standard locations${NC}"
fi

# Comprobar ubicaciones posibles de la configuración
POSSIBLE_CONFIG_DIRS=(
    "/etc/infracli"
    "$HOME/.config/infracli"
)

CONFIG_FOUND=false
for config_dir in "${POSSIBLE_CONFIG_DIRS[@]}"; do
    if [ -d "$config_dir" ]; then
        echo "Found infracli configuration at $config_dir"
        CONFIG_FOUND=true
        
        # Solicitar confirmación
        read -p "Do you want to remove this configuration directory? (y/n) " -n 1 -r
        echo
        if [[ $REPLY =~ ^[Yy]$ ]]; then
            echo "Removing configuration from $config_dir"
            rm -rf "$config_dir"
            echo -e "${GREEN}Configuration removed successfully${NC}"
        else
            echo "Configuration will not be removed"
        fi
    fi
done

if [ "$CONFIG_FOUND" = false ]; then
    echo -e "${YELLOW}No infracli configuration found in standard locations${NC}"
fi

echo -e "${GREEN}Uninstallation process completed!${NC}"