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
echo "This script will install the infracli tool"

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

echo "Installing infracli from $PROJECT_DIR"

# Entrar al directorio del proyecto
cd "$PROJECT_DIR"

# Instalar el proyecto usando go install
echo "Installing infracli using go install..."
go install

echo -e "${GREEN}Installation completed successfully!${NC}"
echo "You can now use infracli by running: infracli"
echo "The configuration will be automatically created at ~/.config/infracli/infracli.json on first run"
echo "Try 'infracli --help' to see available commands"

# Notificar que go install coloca los binarios en GOPATH/bin
echo -e "${YELLOW}Note: infracli has been installed to your GOPATH/bin directory${NC}"
echo "Make sure GOPATH/bin is in your PATH environment variable"
