# 🛠️ InfraCLI

A powerful command-line tool to manage Docker Compose infrastructure services.

## ✨ Features

- Start and stop infrastructure services with a single command
- Get connection information and details for any service
- Manage multiple services at once
- Automatic service discovery
- Configurable via JSON configuration

## 📦 Installation

```bash
# Navigate to the infracli directory
cd infracli

# Run the installation script
./scripts/install.sh
```

The installation script will compile the CLI tool and install it to your system.

## 📚 Usage

### 📋 List Available Services

```bash
infracli list
```

### 📄 Get Service Information

```bash
# Get connection information for a service
infracli info mysql

# Get connection details for other services
infracli info postgres
infracli info mongo
infracli info elasticsearch-kibana
```

### 🚀 Start Services

```bash
# Start a single service
infracli run mysql

# Start multiple services
infracli run mysql mongo

# Start all available services
infracli run all
```

### 🛑 Stop Services

```bash
# Stop a single service
infracli down mysql

# Stop multiple services
infracli down mysql mongo

# Stop all services
infracli down all

# Stop services and remove volumes
infracli down mysql --volumes
```

### 🔍 Verbose Output

Add the `-v` or `--verbose` flag to get detailed output:

```bash
infracli run mysql -v
```

## 🗑️ Uninstallation

To remove the InfraCLI tool:

```bash
# Navigate to the infracli directory
cd infracli

# Run the uninstallation script
./scripts/uninstall.sh
```

## ⚙️ Configuration

The tool uses a `config.json` file located in the `config` directory. This file specifies:

- `servicesPath`: The relative path to the directory containing service directories
- `excludedDirs`: Directories to exclude from service discovery

Default configuration:

```json
{
  "servicesPath": "../",
  "excludedDirs": ["config", "scripts", "cmd"]
}
```

## 💻 Development

This tool is built using Go with the Cobra CLI framework. To contribute:

1. Make your changes
2. Test locally without installation: `go run main.go [command]`
3. Build for testing: `go build`