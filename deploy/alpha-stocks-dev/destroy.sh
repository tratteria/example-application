if kubectl get namespace alpha-stocks-dev &> /dev/null; then
    kubectl delete namespace alpha-stocks-dev
else
    echo "Namespace 'alpha-stocks-dev' does not exist, skipping deletion.."
fi
