#!/bin/bash

info() {
    printf "\e[34m%s\n\e[0m" "$1"
}

error() {
    printf "\e[31m%s\n\e[0m" "$1"
}

success() {
    printf "\e[32m%s\n\e[0m" "$1"
}

usage() {
    echo "Usage: $0 [OPTIONS]"
    echo "Options:"
    echo "  --no-spire    Skips the uninstallation of SPIRE."
    echo "  -h, --help    Displays this help message."
}

resource_exists() {
    resource_type=$1
    resource_name=$2
    kubectl get ${resource_type} ${resource_name} > /dev/null 2>&1
}

delete_k8s_resource() {
    resource_type=$1
    resource_name=$2
    if resource_exists ${resource_type} ${resource_name}; then
        kubectl delete ${resource_type} ${resource_name} || { error "Failed to delete ${resource_type} ${resource_name}"; exit 1; }
    else
        info "${resource_type} ${resource_name} does not exist, skipping deletion."
    fi
}

uninstall_spire=true

while [[ "$#" -gt 0 ]]; do
    case $1 in
        --no-spire) uninstall_spire=false ;;
        -h|--help) usage; exit 0 ;;
        *) echo "Unknown option: $1"; usage; exit 1 ;;
    esac
    shift
done

if [ "$uninstall_spire" = true ]; then
    info "Uninstalling spire..."
    cd resources/spire
    chmod +x uninstall.sh
    ./uninstall.sh
    cd ../../
fi

info "Uninstalling tconfigd..."

delete_k8s_resource namespace tratteria-system
delete_k8s_resource clusterrole tconfigd-service-account-role
delete_k8s_resource clusterrolebinding tconfigd-service-account-binding
delete_k8s_resource mutatingwebhookconfiguration tratteria-agent-injector
delete_k8s_resource crd trats.tratteria.io
delete_k8s_resource crd tratteriaconfigs.tratteria.io

success "tconfigd uninstalled successfully."