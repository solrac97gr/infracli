package cmd

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/solrac97gr/infrastructure/infracli/config"
	"github.com/spf13/cobra"
)

// Simple helper functions to parse docker-compose.yml without using a YAML parser
func extractPorts(content string, servicePrefix string) []string {
	var ports []string
	inService := false
	inPorts := false
	
	scanner := bufio.NewScanner(strings.NewReader(content))
	re := regexp.MustCompile(`^\s*-\s*"?([^:]+):([^"]+)"?`)
	
	for scanner.Scan() {
		line := scanner.Text()
		
		// Check if we're entering a service section
		if strings.HasPrefix(line, servicePrefix) {
			inService = true
			continue
		}
		
		// If we've moved to a different section, stop processing
		if inService && len(line) > 0 && !strings.HasPrefix(line, " ") {
			inService = false
			inPorts = false
			continue
		}
		
		// Check if we're in the ports section of the current service
		if inService && strings.TrimSpace(line) == "ports:" {
			inPorts = true
			continue
		}
		
		// If we're in the ports section, extract port mappings
		if inService && inPorts {
			// If we've moved to a different subsection, stop processing ports
			if len(line) > 0 && strings.HasPrefix(line, "  ") && !strings.HasPrefix(line, "    -") && !strings.HasPrefix(line, "   -") {
				inPorts = false
				continue
			}
			
			// Extract port mappings
			matches := re.FindStringSubmatch(line)
			if len(matches) >= 3 {
				hostPort := strings.TrimSpace(matches[1])
				containerPort := strings.TrimSpace(matches[2])
				ports = append(ports, fmt.Sprintf("%s:%s", hostPort, containerPort))
			}
		}
	}
	
	return ports
}

func extractEnvironment(content string, servicePrefix string) map[string]string {
	env := make(map[string]string)
	inService := false
	inEnvironment := false
	
	scanner := bufio.NewScanner(strings.NewReader(content))
	reKeyValue := regexp.MustCompile(`^\s*([^:]+):\s*(.+)$`)
	reEnvVar := regexp.MustCompile(`^\s*-\s*([^=]+)=(.+)$`)
	
	for scanner.Scan() {
		line := scanner.Text()
		
		// Check if we're entering a service section
		if strings.HasPrefix(line, servicePrefix) {
			inService = true
			continue
		}
		
		// If we've moved to a different section, stop processing
		if inService && len(line) > 0 && !strings.HasPrefix(line, " ") {
			inService = false
			inEnvironment = false
			continue
		}
		
		// Check if we're in the environment section of the current service
		if inService && strings.TrimSpace(line) == "environment:" {
			inEnvironment = true
			continue
		}
		
		// If we're in the environment section, extract variables
		if inService && inEnvironment {
			// If we've moved to a different subsection, stop processing environment
			if len(line) > 0 && strings.HasPrefix(line, "  ") && !strings.HasPrefix(line, "    ") && !strings.Contains(line, ":") && !strings.Contains(line, "-") {
				inEnvironment = false
				continue
			}
			
			// Try to match key-value format
			matches := reKeyValue.FindStringSubmatch(line)
			if len(matches) >= 3 {
				key := strings.TrimSpace(matches[1])
				value := strings.TrimSpace(matches[2])
				env[key] = value
				continue
			}
			
			// Try to match array format with environment variables
			matches = reEnvVar.FindStringSubmatch(line)
			if len(matches) >= 3 {
				key := strings.TrimSpace(matches[1])
				value := strings.TrimSpace(matches[2])
				// Remove quotes if present
				value = strings.Trim(value, "\"'")
				parts := strings.SplitN(key, "=", 2)
				if len(parts) == 2 {
					env[parts[0]] = parts[1]
				} else {
					env[key] = value
				}
			}
		}
	}
	
	return env
}

func extractImageAndServices(content string) map[string]string {
	services := make(map[string]string)
	
	scanner := bufio.NewScanner(strings.NewReader(content))
	inServices := false
	currentService := ""
	
	reService := regexp.MustCompile(`^(\s*)([^:]+):$`)
	reImage := regexp.MustCompile(`^\s*image:\s*(.+)$`)
	
	for scanner.Scan() {
		line := scanner.Text()
		
		// Check if we're in the services section
		if strings.TrimSpace(line) == "services:" {
			inServices = true
			continue
		}
		
		// If we've exited the services section
		if inServices && len(line) > 0 && !strings.HasPrefix(line, " ") && line != "services:" {
			inServices = false
			continue
		}
		
		if inServices {
			// Try to match service name
			serviceMatches := reService.FindStringSubmatch(line)
			if len(serviceMatches) >= 3 && len(serviceMatches[1]) == 2 { // Check indent level to ensure it's a service
				currentService = strings.TrimSpace(serviceMatches[2])
				continue
			}
			
			// Try to match image definition
			if currentService != "" {
				imageMatches := reImage.FindStringSubmatch(line)
				if len(imageMatches) >= 2 {
					services[currentService] = strings.TrimSpace(imageMatches[1])
				}
			}
		}
	}
	
	return services
}

var infoCmd = &cobra.Command{
	Use:   "info [service]",
	Short: "Display information about a service",
	Long: `Display detailed information about a specific infrastructure service,
including connection strings, ports, and other relevant details.

Examples:
  infracli info mysql
  infracli info postgres
  infracli info mongo`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		serviceName := args[0]
		
		// Get verbose flag
		verbose, _ := cmd.Flags().GetBool("verbose")

		// Get available services
		availableServices, err := config.GetAvailableServices()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			return
		}

		// Check if service exists
		serviceFound := false
		for _, availService := range availableServices {
			if serviceName == availService {
				serviceFound = true
				break
			}
		}

		if !serviceFound {
			fmt.Printf("Error: Service '%s' not found in available services\n", serviceName)
			fmt.Printf("Available services: %s\n", strings.Join(availableServices, ", "))
			return
		}

		// Get configuration
		cfg, err := config.LoadConfig()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error loading configuration: %v\n", err)
			return
		}

		// Use the configured services path
		basePath := cfg.ServicesPath

		// Expand path if it contains ~/
		if len(basePath) >= 2 && basePath[:2] == "~/" {
			homeDir, err := os.UserHomeDir()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error getting home directory: %v\n", err)
				return
			}
			basePath = filepath.Join(homeDir, basePath[2:])
		}

		// Path to the docker-compose.yml file
		dockerComposePath := filepath.Join(basePath, serviceName, "docker-compose.yml")
		
		if verbose {
			fmt.Printf("Reading compose file: %s\n", dockerComposePath)
		}

		// Read docker-compose.yml
		composeData, err := os.ReadFile(dockerComposePath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading docker-compose.yml: %v\n", err)
			return
		}

		// Convert to string for our custom parsing
		composeContent := string(composeData)

		// Display service information
		fmt.Printf("Service: %s\n", serviceName)
		fmt.Println(strings.Repeat("=", 50))

		// Display connection information based on service type
		switch serviceName {
		case "mysql":
			displayMySQLInfo(composeContent)
		case "postgres":
			displayPostgresInfo(composeContent)
		case "mongo":
			displayMongoInfo(composeContent)
		case "redis":
			displayRedisInfo(composeContent)
		case "elasticsearch-kibana":
			displayElasticsearchKibanaInfo(composeContent)
		case "neo4j":
			displayNeo4jInfo(composeContent)
		default:
			// Generic display for other services
			displayGenericInfo(serviceName, composeContent)
		}
	},
}

func displayMySQLInfo(content string) {
	// Extract service information
	services := extractImageAndServices(content)
	
	// Find the MySQL service
	var mysqlServiceName string
	for name, image := range services {
		if strings.Contains(strings.ToLower(image), "mysql") {
			mysqlServiceName = name
			break
		}
	}
	
	if mysqlServiceName == "" {
		fmt.Println("MySQL service not found in docker-compose.yml")
		return
	}
	
	// Extract ports for MySQL service (looking for the service section with indentation)
	servicePrefix := "  " + mysqlServiceName + ":"
	ports := extractPorts(content, servicePrefix)
	
	// Find MySQL port
	var port string = "3306" // Default if not specified
	for _, portMapping := range ports {
		if strings.Contains(portMapping, "3306") {
			parts := strings.Split(portMapping, ":")
			if len(parts) == 2 {
				port = strings.TrimSpace(parts[0])
				break
			}
		}
	}
	
	// Extract environment variables
	env := extractEnvironment(content, servicePrefix)
	
	// Get database credentials
	database := env["MYSQL_DATABASE"]
	user := env["MYSQL_USER"]
	password := env["MYSQL_PASSWORD"]
	rootPassword := env["MYSQL_ROOT_PASSWORD"]
	
	// Display MySQL connection information
	fmt.Println("MySQL Connection Information:")
	fmt.Println(strings.Repeat("-", 40))
	fmt.Printf("Host: localhost\n")
	fmt.Printf("Port: %s\n", port)
	fmt.Printf("Database: %s\n", database)
	fmt.Printf("User: %s\n", user)
	fmt.Printf("Password: %s\n", password)
	fmt.Printf("Root Password: %s\n", rootPassword)
	
	// Display connection strings
	fmt.Println("\nConnection Strings:")
	fmt.Printf("JDBC: jdbc:mysql://localhost:%s/%s\n", port, database)
	fmt.Printf("URL: mysql://%s:%s@localhost:%s/%s\n", user, password, port, database)
	fmt.Printf("CLI: mysql -h localhost -P %s -u %s -p%s %s\n", port, user, password, database)
}

func displayPostgresInfo(content string) {
	// Extract service information
	services := extractImageAndServices(content)
	
	// Find the Postgres service
	var pgServiceName string
	for name, image := range services {
		if strings.Contains(strings.ToLower(image), "postgres") {
			pgServiceName = name
			break
		}
	}
	
	if pgServiceName == "" {
		fmt.Println("PostgreSQL service not found in docker-compose.yml")
		return
	}
	
	// Extract ports for Postgres service
	servicePrefix := "  " + pgServiceName + ":"
	ports := extractPorts(content, servicePrefix)
	
	// Find Postgres port
	var port string = "5432" // Default if not specified
	for _, portMapping := range ports {
		if strings.Contains(portMapping, "5432") {
			parts := strings.Split(portMapping, ":")
			if len(parts) == 2 {
				port = strings.TrimSpace(parts[0])
				break
			}
		}
	}
	
	// Extract environment variables
	env := extractEnvironment(content, servicePrefix)
	
	// Get database credentials
	user := env["POSTGRES_USER"]
	password := env["POSTGRES_PASSWORD"]
	database := env["POSTGRES_DB"]
	
	// If database name is not specified, it defaults to the username
	if database == "" {
		database = user
	}
	
	// Display Postgres connection information
	fmt.Println("PostgreSQL Connection Information:")
	fmt.Println(strings.Repeat("-", 40))
	fmt.Printf("Host: localhost\n")
	fmt.Printf("Port: %s\n", port)
	fmt.Printf("Database: %s\n", database)
	fmt.Printf("User: %s\n", user)
	fmt.Printf("Password: %s\n", password)
	
	// Display connection strings
	fmt.Println("\nConnection Strings:")
	fmt.Printf("JDBC: jdbc:postgresql://localhost:%s/%s\n", port, database)
	fmt.Printf("URL: postgresql://%s:%s@localhost:%s/%s\n", user, password, port, database)
	fmt.Printf("CLI: psql -h localhost -p %s -U %s -d %s\n", port, user, database)
}

func displayMongoInfo(content string) {
	// Extract service information
	services := extractImageAndServices(content)
	
	// Find the MongoDB service
	var mongoServiceName string
	for name, image := range services {
		if strings.Contains(strings.ToLower(image), "mongo") {
			mongoServiceName = name
			break
		}
	}
	
	if mongoServiceName == "" {
		fmt.Println("MongoDB service not found in docker-compose.yml")
		return
	}
	
	// Extract ports for MongoDB service
	servicePrefix := "  " + mongoServiceName + ":"
	ports := extractPorts(content, servicePrefix)
	
	// Find MongoDB port
	var port string = "27017" // Default if not specified
	for _, portMapping := range ports {
		if strings.Contains(portMapping, "27017") {
			parts := strings.Split(portMapping, ":")
			if len(parts) == 2 {
				port = strings.TrimSpace(parts[0])
				break
			}
		}
	}
	
	// Extract environment variables
	env := extractEnvironment(content, servicePrefix)
	
	// Get database credentials
	user := env["MONGO_INITDB_ROOT_USERNAME"]
	password := env["MONGO_INITDB_ROOT_PASSWORD"]
	
	// Display MongoDB connection information
	fmt.Println("MongoDB Connection Information:")
	fmt.Println(strings.Repeat("-", 40))
	fmt.Printf("Host: localhost\n")
	fmt.Printf("Port: %s\n", port)
	fmt.Printf("User: %s\n", user)
	fmt.Printf("Password: %s\n", password)
	fmt.Printf("Authentication Database: admin\n")
	
	// Display connection strings
	fmt.Println("\nConnection Strings:")
	fmt.Printf("URI: mongodb://%s:%s@localhost:%s/admin\n", user, password, port)
	fmt.Printf("CLI: mongosh mongodb://%s:%s@localhost:%s/admin\n", user, password, port)
}

func displayElasticsearchKibanaInfo(content string) {
	// Extract service information
	services := extractImageAndServices(content)
	
	// Find Elasticsearch and Kibana services
	var esServiceName, kibanaServiceName string
	for name, image := range services {
		if strings.Contains(strings.ToLower(name), "elastic") || 
		   strings.Contains(strings.ToLower(image), "elastic") {
			esServiceName = name
		}
		if strings.Contains(strings.ToLower(name), "kibana") || 
		   strings.Contains(strings.ToLower(image), "kibana") {
			kibanaServiceName = name
		}
	}
	
	// Extract ports for Elasticsearch service
	var esPort, kibanaPort string = "9200", "5601" // Defaults
	
	if esServiceName != "" {
		servicePrefix := "  " + esServiceName + ":"
		ports := extractPorts(content, servicePrefix)
		
		for _, portMapping := range ports {
			if strings.Contains(portMapping, "9200") {
				parts := strings.Split(portMapping, ":")
				if len(parts) == 2 {
					esPort = strings.TrimSpace(parts[0])
					break
				}
			}
		}
		
		// Check security settings
		securityEnabled := false
		
		// Check if security is enabled in environment - could be in various formats
		scanner := bufio.NewScanner(strings.NewReader(content))
		for scanner.Scan() {
			line := scanner.Text()
			if strings.Contains(line, "xpack.security.enabled=true") {
				securityEnabled = true
				break
			}
		}
		
		// Display Elasticsearch information
		fmt.Println("Elasticsearch Connection Information:")
		fmt.Println(strings.Repeat("-", 40))
		fmt.Printf("Elasticsearch URL: http://localhost:%s\n", esPort)
		if securityEnabled {
			fmt.Println("Security: Enabled (requires authentication)")
			fmt.Println("Default username: elastic")
		} else {
			fmt.Println("Security: Disabled (no authentication required)")
		}
	} else {
		fmt.Println("Elasticsearch service not found in docker-compose.yml")
	}
	
	// Extract ports for Kibana service
	if kibanaServiceName != "" {
		servicePrefix := "  " + kibanaServiceName + ":"
		ports := extractPorts(content, servicePrefix)
		
		for _, portMapping := range ports {
			if strings.Contains(portMapping, "5601") {
				parts := strings.Split(portMapping, ":")
				if len(parts) == 2 {
					kibanaPort = strings.TrimSpace(parts[0])
					break
				}
			}
		}
		
		// Display Kibana information
		fmt.Println("\nKibana Information:")
		fmt.Println(strings.Repeat("-", 40))
		fmt.Printf("Kibana URL: http://localhost:%s\n", kibanaPort)
	}
	
	// Display example commands
	fmt.Println("\nExample Commands:")
	fmt.Printf("Check Elasticsearch health: curl http://localhost:%s/_cluster/health?pretty\n", esPort)
	fmt.Printf("View indices: curl http://localhost:%s/_cat/indices\n", esPort)
}

func displayGenericInfo(serviceName string, content string) {
	fmt.Println("Service Configuration:")
	fmt.Println(strings.Repeat("-", 40))
	
	// Extract services and their images
	services := extractImageAndServices(content)
	
	for name, image := range services {
		fmt.Printf("Service: %s\n", name)
		fmt.Printf("Image: %s\n", image)
		
		// Extract ports for this service
		servicePrefix := "  " + name + ":"
		ports := extractPorts(content, servicePrefix)
		
		if len(ports) > 0 {
			fmt.Println("\nExposed Ports:")
			for _, port := range ports {
				fmt.Printf("- %s\n", port)
			}
		}
		
		// Extract environment variables
		env := extractEnvironment(content, servicePrefix)
		if len(env) > 0 {
			fmt.Println("\nEnvironment Variables:")
			for key, value := range env {
				fmt.Printf("- %s: %s\n", key, value)
			}
		}
		
		fmt.Println()
	}
	
	fmt.Println("To start this service:")
	fmt.Printf("  infracli run %s\n", serviceName)
	fmt.Println("\nTo stop this service:")
	fmt.Printf("  infracli down %s\n", serviceName)
}

func displayRedisInfo(content string) {
	// Extract service information
	services := extractImageAndServices(content)
	
	// Find the Redis service
	var redisServiceName string
	for name, image := range services {
		if strings.Contains(strings.ToLower(image), "redis") {
			redisServiceName = name
			break
		}
	}
	
	if redisServiceName == "" {
		fmt.Println("Redis service not found in docker-compose.yml")
		return
	}
	
	// Extract ports for Redis service
	servicePrefix := "  " + redisServiceName + ":"
	ports := extractPorts(content, servicePrefix)
	
	// Find Redis port
	var port string = "6379" // Default if not specified
	for _, portMapping := range ports {
		if strings.Contains(portMapping, "6379") {
			parts := strings.Split(portMapping, ":")
			if len(parts) == 2 {
				port = strings.TrimSpace(parts[0])
				break
			}
		}
	}
	
	// Check if password is configured
	password := ""
	requiresAuth := false
	
	// Look for password in command line arguments
	scanner := bufio.NewScanner(strings.NewReader(content))
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "requirepass") {
			requiresAuth = true
			parts := strings.Split(line, "requirepass")
			if len(parts) > 1 {
				password = strings.TrimSpace(parts[1])
			}
			break
		}
	}
	
	// Display Redis connection information
	fmt.Println("Redis Connection Information:")
	fmt.Println(strings.Repeat("-", 40))
	fmt.Printf("Host: localhost\n")
	fmt.Printf("Port: %s\n", port)
	
	if requiresAuth {
		fmt.Printf("Password: %s\n", password)
		fmt.Printf("Authentication: Enabled\n")
	} else {
		fmt.Printf("Authentication: Disabled (no password required)\n")
	}
	
	// Check if AOF persistence is enabled
	persistenceEnabled := false
	scanner = bufio.NewScanner(strings.NewReader(content))
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "appendonly yes") {
			persistenceEnabled = true
			break
		}
	}
	
	fmt.Printf("Persistence: %s\n", map[bool]string{true: "Enabled (appendonly)", false: "Standard RDB"}[persistenceEnabled])
	
	// Display connection strings
	fmt.Println("\nConnection Strings:")
	if requiresAuth {
		fmt.Printf("URI: redis://:%s@localhost:%s/0\n", password, port)
	} else {
		fmt.Printf("URI: redis://localhost:%s/0\n", port)
	}
	
	// Display example commands
	fmt.Println("\nExample Commands:")
	if requiresAuth {
		fmt.Printf("CLI: redis-cli -h localhost -p %s -a %s\n", port, password)
	} else {
		fmt.Printf("CLI: redis-cli -h localhost -p %s\n", port)
	}
	
	fmt.Printf("Ping test: redis-cli -h localhost -p %s ping\n", port)
}

func displayNeo4jInfo(content string) {
	// Extract service information
	services := extractImageAndServices(content)
	
	// Find the Neo4j service
	var neo4jServiceName string
	for name, image := range services {
		if strings.Contains(strings.ToLower(image), "neo4j") {
			neo4jServiceName = name
			break
		}
	}
	
	if neo4jServiceName == "" {
		fmt.Println("Neo4j service not found in docker-compose.yml")
		return
	}
	
	// Extract ports for Neo4j service
	servicePrefix := "  " + neo4jServiceName + ":"
	ports := extractPorts(content, servicePrefix)
	
	// Find Neo4j ports
	var httpPort string = "7474" // Default HTTP port
	var boltPort string = "7687" // Default Bolt port
	var httpsPort string = "7473" // Default HTTPS port
	
	for _, portMapping := range ports {
		if strings.Contains(portMapping, "7474") {
			parts := strings.Split(portMapping, ":")
			if len(parts) == 2 {
				httpPort = strings.TrimSpace(parts[0])
			}
		} else if strings.Contains(portMapping, "7687") {
			parts := strings.Split(portMapping, ":")
			if len(parts) == 2 {
				boltPort = strings.TrimSpace(parts[0])
			}
		} else if strings.Contains(portMapping, "7473") {
			parts := strings.Split(portMapping, ":")
			if len(parts) == 2 {
				httpsPort = strings.TrimSpace(parts[0])
			}
		}
	}
	
	// Extract environment variables
	env := extractEnvironment(content, servicePrefix)
	
	// Get database credentials
	authInfo := env["NEO4J_AUTH"]
	var user, password string
	
	if authInfo != "" {
		// Default format is neo4j/password
		parts := strings.Split(authInfo, "/")
		if len(parts) == 2 {
			user = parts[0]
			password = parts[1]
		} else {
			user = "neo4j" // Default username
			password = "neo4j" // Default password when authentication is enabled
		}
	} else {
		user = "neo4j" // Default username
		password = "(disabled)" // If NEO4J_AUTH not set or empty
	}
	
	// Display Neo4j connection information
	fmt.Println("Neo4j Connection Information:")
	fmt.Println(strings.Repeat("-", 40))
	fmt.Printf("Host: localhost\n")
	fmt.Printf("HTTP Port: %s\n", httpPort)
	fmt.Printf("Bolt Port: %s\n", boltPort)
	fmt.Printf("HTTPS Port: %s\n", httpsPort)
	fmt.Printf("Username: %s\n", user)
	fmt.Printf("Password: %s\n", password)
	
	// Display connection strings
	fmt.Println("\nConnection Information:")
	fmt.Printf("Browser UI: http://localhost:%s\n", httpPort)
	fmt.Printf("Bolt URI: bolt://localhost:%s\n", boltPort)
	fmt.Printf("HTTPS UI: https://localhost:%s\n", httpsPort)
	
	// Display cypher-shell command
	fmt.Println("\nConnect with Cypher Shell:")
	fmt.Printf("cypher-shell -a bolt://localhost:%s -u %s -p %s\n", boltPort, user, password)
	
	// Display Docker connection commands
	fmt.Println("\nDocker Commands:")
	fmt.Printf("Cypher Shell: docker exec -it neo4j cypher-shell -u %s -p %s\n", user, password)
	fmt.Printf("Interactive Shell: docker exec -it neo4j bash\n")
}

func init() {
	RootCmd.AddCommand(infoCmd)
}