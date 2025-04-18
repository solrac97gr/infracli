#!/bin/bash

# Script para instalar infracli
# Compila e instala la herramienta de línea de comandos infracli

set -e

# Colores para mensajes
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo -e "${YELLOW}InfraCLI Installation Script${NC}"
echo "This script will build and install the infracli tool"

# Comprobar que Go está instalado
if ! command -v go &> /dev/null; then
    echo -e "${RED}Error: Go language is not installed${NC}"
    echo "Please install Go before continuing: https://golang.org/doc/install"
    exit 1
fi

# Verificar versión mínima de Go
GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
REQUIRED_VERSION="1.16"

echo "Detected Go version: $GO_VERSION"

# Determinar la ruta del script para encontrar el directorio del proyecto
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
PROJECT_DIR="$(dirname "$SCRIPT_DIR")"

echo "Building infracli from $PROJECT_DIR"

# Entrar al directorio del proyecto
cd "$PROJECT_DIR"

# Compilar el proyecto
echo "Compiling infracli..."
go build -o infracli

# Determinar el directorio de instalación
INSTALL_DIR="/usr/local/bin"
if [[ "$OSTYPE" == "darwin"* ]]; then
    # macOS
    if [ ! -d "$INSTALL_DIR" ] || [ ! -w "$INSTALL_DIR" ]; then
        echo "Installing to ~/bin instead of $INSTALL_DIR due to permissions"
        INSTALL_DIR="$HOME/bin"
        mkdir -p "$INSTALL_DIR"
    fi
elif [[ "$OSTYPE" == "linux-gnu"* ]]; then
    # Linux
    if [ ! -d "$INSTALL_DIR" ] || [ ! -w "$INSTALL_DIR" ]; then
        echo "Installing to ~/bin instead of $INSTALL_DIR due to permissions"
        INSTALL_DIR="$HOME/bin"
        mkdir -p "$INSTALL_DIR"
        
        # Verificar si ~/bin está en PATH y añadirlo si no lo está
        if [[ ":$PATH:" != *":$HOME/bin:"* ]]; then
            echo 'export PATH="$HOME/bin:$PATH"' >> ~/.bashrc
            echo "Added $HOME/bin to PATH in ~/.bashrc"
            
            # Verificar si también se usa zsh
            if [ -f "$HOME/.zshrc" ]; then
                echo 'export PATH="$HOME/bin:$PATH"' >> ~/.zshrc
                echo "Added $HOME/bin to PATH in ~/.zshrc"
            fi
        fi
    fi
fi

# Copiar el binario al directorio de instalación
echo "Installing infracli to $INSTALL_DIR"
cp infracli "$INSTALL_DIR/"

# Crear directorio de configuración si es necesario
CONFIG_DIR="/etc/infracli"
if [ ! -d "$CONFIG_DIR" ] || [ ! -w "$CONFIG_DIR" ]; then
    CONFIG_DIR="$HOME/.config/infracli"
    mkdir -p "$CONFIG_DIR"
fi

# Copiar la configuración
echo "Setting up configuration in $CONFIG_DIR"
cp -r config/config.json "$CONFIG_DIR/"

echo -e "${GREEN}Installation completed successfully!${NC}"
echo "You can now use infracli by running: infracli"
echo "Try 'infracli --help' to see available commands"

# Si estamos usando un directorio bin personalizado y es la primera instalación
if [ "$INSTALL_DIR" = "$HOME/bin" ]; then
    echo -e "${YELLOW}Note: You may need to open a new terminal or run 'source ~/.bashrc' to use infracli${NC}"
fi