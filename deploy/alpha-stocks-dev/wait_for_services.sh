echo "\n\nWaiting for the readiness of the application...\n"

check_service_readiness() {
    local service=$1
    local namespace=${2:-alpha-stocks-dev}
    local attempts=0
    local max_attempts=30

    while ! kubectl get pods -n $namespace | grep "$service" | grep -q 'Running'; do
        if [ $attempts -ge $max_attempts ]; then
            echo "Failed to verify the readiness of $service."
            return 1
        fi
        attempts=$((attempts + 1))
        echo "Waiting for $service to be ready..."
        sleep 10
    done
    echo "$service is ready!"
    return 0
}

main() {
    local services=("gateway" "order" "stocks")
    
    for service in "${services[@]}"; do
        if ! check_service_readiness $service; then
            exit 1
        fi
    done
}

main