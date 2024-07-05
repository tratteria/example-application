if kubectl get namespace alpha-stocks &> /dev/null; then
    kubectl delete namespace alpha-stocks
else
    echo "Namespace 'alpha-stocks' does not exist, skipping deletion.."
fi
