# Installation Guide for tconfigd

## Prerequisites

### Kubernetes Cluster:

Ensure your Kubernetes cluster is active and accessible.

### SPIRE:

tconfigd requires [SPIRE](https://github.com/spiffe/spire) for its operation. Ensure SPIRE is running in your Kubernetes cluster before proceeding to install tconfigd.

#### 1. Setup

tconfigd utilizes SPIRE's built-in [k8sbundle plugin](https://github.com/spiffe/spire/blob/main/doc/plugin_server_notifier_k8sbundle.md) to push CA certificates to its Kubernetes webhooks. Add the following configuration to your SPIRE server:

```
Notifier "k8sbundle" {
    plugin_data {
        webhook_label = "tratteria.io/webhook"
    }
}
```

For the above to operate, you need to add the following permissions to the SPIRE server ClusterRole:

```yaml
- apiGroups: ["admissionregistration.k8s.io"]
  resources: ["mutatingwebhookconfigurations", "validatingwebhookconfigurations"]
  verbs: ["get", "list", "patch", "watch"]
```

#### 2. Registering tconfigd

Register tconfigd to your running SPIRE. tconfigd operates within the `tratteria-system` namespace and uses the `tconfigd-service-account` service account. Below is a sample registration command:

```bash
kubectl exec -n spire spire-server-0 -- \
    /opt/spire/bin/spire-server entry create --dns tconfigd.tratteria-system.svc \
    -spiffeID spiffe://[your-trust-domain]/tconfigd \
    -parentID spiffe://[your-trust-domain]/ns/spire/sa/spire-agent \
    -selector k8s:ns:tratteria-system \
    -selector k8s:sa:tconfigd-service-account
```

Include the `--dns` option with the value `tconfigd.tratteria-system.svc` in the registration command; it is necessary for tconfigd operation.

#### 2. Registering tratteria service

Register [tratteria service](https://github.com/tratteria/tratteria), an open source Transaction Tokens (TraTs) Service, to your running SPIRE. tratteria operates within your application namespace and uses the `tratteria-service-account` service account. Below is a sample registration command:

```bash
kubectl exec -n spire spire-server-0 -- \
    /opt/spire/bin/spire-server entry create \
    -spiffeID spiffe://[your-trust-domain]/tratteria \
    -parentID spiffe://[your-trust-domain]/ns/spire/sa/spire-agent \
    -selector k8s:ns:[your-namespace] \
    -selector k8s:sa:tratteria-service-account
```

#### 3. Registering Microservices

Ensure all microservice that need to verify trats are registered to your running SPIRE.

<br>

For a reference implementation on setting up SPIRE for tconfigd, please check the example application's [spire installation](https://github.com/tratteria/example-application/tree/main/deploy/spire).


## Installation Steps

### 1. Clone the Repository
Clone the tconfigd repository and navigate to the installation directory:

```bash
git clone https://github.com/tratteria/tconfigd.git
```

```bash
cd tconfigd/installation
```

Alternatively, perform a shallow clone of just the installation subdirectory.


### 3. Configure tconfigd

Update the `config.yaml` file to match your specific deployment settings:

- **Configure the settings as described below:**
  - `enableTratInterception`: "`true`" 
    - **Description**: Set to "`true`" to enable interception of incoming requests for TraT verification. Set to "`false`" if using the [delegation method](https://github.com/tratteria/tratteria-agent?tab=readme-ov-file#operating-modes) for TraT verification.
  - `spireAgentHostDir`: `"/run/spire/sockets"`
    - **Description**: Host directory where the SPIRE agent's socket resides. Update this value if it is different in your SPIRE installation.
  - `tratteriaSpiffeId`: `"spiffe://[your-trust-domain]/tratteria"`
    - **Description**: SPIFFE ID used to register [tratteria service](https://github.com/tratteria/tratteria), an open source Transaction Tokens (TraTs) Service.
  - `agentApiPort`: "`9030`"
    - **Description**: Port number for the tratteria agent APIs. Do not change this unless you have some specific need.
  - `agentInterceptorPort`: "`9050`"
    - **Description**: The port number for the tratteria agent's incoming requests interceptor. Do not change this unless you have some specific need.


### 4. Run the Installation Script

Deploy tconfigd to your Kubernetes cluster by running the installation script:

```bash
./install.sh
```

### 5. Verification

Verify the installation by checking the status of the tconfigd pod in the `tratteria-system` namespace. Use the following command to view the pod:

```bash
kubectl get pods -n tratteria-system
```

<br>

For a practical example of installing tconfigd on a microservice application, refer to the example-application's [deployment resources](https://github.com/tratteria/example-application/tree/main/deploy).

