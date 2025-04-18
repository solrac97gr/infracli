#!/bin/bash

# Function to display script usage
show_usage() {
    echo "Usage: ./run.sh [service1] [service2] ... [serviceN]"
    echo "  Where services can be: mongo, postgres, elasticsearch-kibana, mysql"
    echo ""
    echo "Examples:"
    echo "  ./run.sh mongo                  # Start only MongoDB"
    echo "  ./run.sh mysql postgres         # Start MySQL and PostgreSQL"
    echo "  ./run.sh all                    # Start all available services"
}

# Check if no arguments were provided
if [ $# -eq 0 ]; then
    show_usage
    exit 1
fi

# Array of available services
available_services=("mongo" "postgres" "elasticsearch-kibana" "mysql")

# Function to start a service
start_service() {
    local service=$1
    
    # Check if the service directory exists
    if [ ! -d "$service" ]; then
        echo "Error: Service '$service' not found in the current directory"
        return 1
    fi
    
    echo "Starting $service..."
    (cd "$service" && docker-compose up -d)
    
    if [ $? -eq 0 ]; then
        echo "$service started successfully"
    else
        echo "Failed to start $service"
        return 1
    fi
    
    return 0
}

# Function to start all services
start_all_services() {
    local failed=0
    
    echo "Starting all services..."
    
    for service in "${available_services[@]}"; do
        start_service "$service" || failed=1
    done
    
    if [ $failed -eq 0 ]; then
        echo "All services started successfully"
    else
        echo "Some services failed to start"
        return 1
    fi
    
    return 0
}

# Process arguments
if [ "$1" = "all" ]; then
    start_all_services
else
    # Start each specified service
    for service in "$@"; do
        # Check if the service is in the available services list
        if [[ ! " ${available_services[*]} " =~ " ${service} " ]]; then
            echo "Warning: '$service' is not in the list of known services"
            echo "Known services: ${available_services[*]}"
            echo "Attempting to start anyway..."
        fi
        
        start_service "$service"
    done
fi

echo "Done"