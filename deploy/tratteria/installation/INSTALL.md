# Installation Guide for tconfigd

## Prerequisites
- **Kubernetes Cluster**: Ensure your Kubernetes cluster is active and accessible.
- **kubectl**: Confirm `kubectl` is installed and configured to interact with your Kubernetes cluster.

## Installation Steps

tconfigd requires SPIRE for its operation. If your infrastructure already has SPIRE installed, follow the [instructions for environments with existing SPIRE](#environment-with-existing-spire-installation) to install tconfigd. If SPIRE is not present, the steps below will include its installation as part of the tconfigd setup.


### 1. Clone the Repository
Clone the tconfigd repository and navigate to the installation directory:

```bash
git clone https://github.com/tratteria/tconfigd.git
```

```bash
cd tconfigd/installation
```

Alternatively, perform a shallow clone of just the installation subdirectory.


### 2. Configure tconfigd

Update the `config.yaml` file to match your specific deployment settings:

- **Configure the settings as described below:**
  - `enableTratInterception`: "`true`" 
    - **Description**: Set to "`true`" to enable interception of incoming requests for TraT verification. Set to "`false`" if using the [delegation method](https://github.com/tratteria/tratteria-agent) for TraT verification.
  - `agentApiPort`: "`9040`" 
    - **Description**: The port number for the tratteria agent API. Do not change this unless you have some specific need.
  - `agentInterceptorPort`: "`9050`" 
    - **Description**: The port number for the tratteria agent's incoming requests interceptor. Do not change this unless you have some specific need.
  - `spiffeEndpointSocket`: "`unix:///run/spire/sockets/agent.sock`" 
    - **Description**: Path to the SPIFFE Workload API Unix socket. Do not change this if you are using the tconfigd SPIRE installation.
  - `tconfigdSpiffeId`: "`spiffe://tratteria.io/tconfigd`" 
    - **Description**: SPIFFE ID for tconfigd. Do not change this if you are using the tconfigd SPIRE installation.


### 3. Run the Installation Script

Deploy the tconfigd to your Kubernetes cluster by running the installation script:

```bash
./install.sh
```

### 4. Verification

Verify the installation by checking the status of the tconfigd pod in the `tratteria-system` namespace. Use the following command to view the pod:

```bash
kubectl get pods -n tratteria-system
```

For a practical example of installing tconfigd on a microservice application, refer to the [example-application](https://github.com/tratteria/example-application).

## Environment with Existing SPIRE Installation
If you already have SPIRE installed, follow these steps to install tconfigd:

### 1. Register tconfigd with your running SPIRE:
Below is a sample registartion command:

```bash
kubectl exec -n spire spire-server-0 -- \
    /opt/spire/bin/spire-server entry create --dns tconfigd tratteria-system.svc \
    -spiffeID spiffe://<your-trust-domain>/tconfigd \
    -parentID spiffe://<your-trust-domain>/ns/spire/sa/spire-agent \
    -selector k8s:ns:tratteria-system \
    -selector k8s:sa:tconfigd-service-account
```

### 2. Clone the tconfig Repository:

Follow the [instructions to clone the repository](#1-clone-the-repository) from the general installation guide.

### 3. Configure tconfigd:

Adjust the following fields in your `config.yaml` file to align with your SPIRE settings.

```yaml
spiffeEndpointSocket: "unix:///run/spire/sockets/agent.sock" # Path to the SPIFFE Workload API Unix socket, replace this if your installation settings differ.
tconfigdSpiffeId: "spiffe://<trust-domain>/tconfigd" # SPIFFE ID used in the registration command above.
```

For configuring other fields, refer to the [general installation configuration guide above](#2-configure-tconfigd).

### 4. Run the Installation Script:

When installing tconfigd, add the `--no-spire` flag to the installation script to prevent reinstalling SPIRE in your environment. Similarly, use this flag when running the uninstallation script to keep SPIRE intact.

```bash
./install.sh --no-spire
```

### 5. Verify Installation:

Follow the [verification steps](#4-verification) provided in the general installation guide.

For a practical example of installing tconfigd on a microservice application with existing SPIRE setup, refer to the [example-application-with-existing-spire](https://github.com/tratteria/example-application-with-existing-spire).
