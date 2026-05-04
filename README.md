# SRE Test – Go & Node Microservices

A complete microservices project demonstrating production-ready CI/CD practices, containerization, Kubernetes deployment, and comprehensive observability using industry-standard tools.

## 📋 Overview

This project showcases a full-stack SRE implementation featuring:

- **Two Microservices**: Go and Node.js applications with health checks and metrics
- **CI/CD Pipeline**: Automated testing, quality gates, and deployment via Jenkins
- **Container Orchestration**: Kubernetes-based deployment with configurable scaling
- **Monitoring Stack**: Prometheus metrics with Grafana visualizations
- **Logging Pipeline**: Centralized logging using ELK Stack (Elasticsearch, Filebeat, Kibana)
- **APM Integration**: New Relic for application performance monitoring

## 🏗️ Architecture

### Services

```
services/
├── go-service/         # Go microservice with Prometheus metrics (port 8080)
└── node-service/       # Node.js microservice with Prometheus metrics (port 8081)
```

### CI/CD Workflow

```
┌─────────────┐    ┌──────────────┐    ┌──────────────────┐
│   Checkout  │ -> │  Unit Tests  │ -> │ SonarQube Scan   │
└─────────────┘    └──────────────┘    └──────────────────┘
                                                │
       ┌────────────────────────────────────────┘
       │
       v
┌──────────────────┐    ┌──────────────┐    ┌─────────────────┐
│ Build Docker     │ -> │ Push Images  │ -> │ Deploy to K8s   │
└──────────────────┘    └──────────────┘    └─────────────────┘
```

### Monitoring & Logging

**Metrics Pipeline:**
```
Application → Prometheus → Grafana Dashboard
```

**Logging Pipeline:**
```
Application Logs → Filebeat → Elasticsearch → Kibana
```

## 🚀 Getting Started

### Prerequisites

- Go 1.19+
- Node.js 16+
- Docker
- Kubernetes cluster (minikube, kind, or cloud provider)
- kubectl configured
- Helm 3

### Running Services Locally

#### Go Service

```bash
cd services/go-service
go mod tidy
go run main.go
```

Test the service:
```bash
go test ./...
```

#### Node Service

```bash
cd services/node-service
npm install
npm start
```

Test the service:
```bash
npm test
```

### Building Docker Images

**Go Service:**
```bash
docker build -t sre-go:test services/go-service
```

**Node Service:**
```bash
docker build -t sre-node:test services/node-service
```

### Deploying to Kubernetes

Apply the namespace first, then set your image tag and deploy:
```bash
kubectl apply -f kubernetes-manifest/namespace.yaml

export IMAGE_TAG=test
envsubst < kubernetes-manifest/go-deployment.yaml | kubectl apply -f -
kubectl apply -f kubernetes-manifest/go-service.yaml

envsubst < kubernetes-manifest/node-deployment.yaml | kubectl apply -f -
kubectl apply -f kubernetes-manifest/node-service.yaml
```

To expose services via Ingress (replace `<PUBLIC_IP>` with your cluster's public IP):
```bash
# Example: export PUBLIC_IP=203.0.113.10
export PUBLIC_IP=<YOUR_CLUSTER_PUBLIC_IP>
envsubst < kubernetes-manifest/ingress.yaml | kubectl apply -f -
```

Services will be reachable at:
- `http://sre-go.<PUBLIC_IP>.nip.io`
- `http://sre-node.<PUBLIC_IP>.nip.io`

## 📊 Monitoring Setup

### Prometheus & Grafana

Install the monitoring stack using Helm:

```bash
# Add Prometheus community charts
helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
helm repo update

# Install Prometheus and Grafana
helm install monitoring prometheus-community/kube-prometheus-stack \
  -n monitoring --create-namespace
```

Access Grafana dashboard:
```bash
kubectl port-forward svc/monitoring-grafana -n monitoring 3000:80
```

Default credentials: `admin` / `prom-operator`

### New Relic APM

Configure New Relic by injecting your license key as an environment variable in your deployment manifests:

```yaml
env:
  - name: NEW_RELIC_LICENSE_KEY
    value: "your-license-key-here"
```

## 📝 Logging Setup (ELK Stack)

### 1. Install Elasticsearch

```bash
helm repo add elastic https://helm.elastic.co
helm repo update

helm install elasticsearch elastic/elasticsearch -n logging \
  --create-namespace \
  --set replicas=1 \
  --set minimumMasterNodes=1 \
  --set resources.requests.cpu=500m \
  --set resources.requests.memory=1Gi
```

### 2. Install Filebeat

```bash
helm install filebeat elastic/filebeat -n logging \
  --set daemonset.enabled=true \
  --set output.elasticsearch.hosts={http://elasticsearch-master.logging.svc.cluster.local:9200}
```

### 3. Install Kibana

```bash
helm install kibana elastic/kibana -n logging \
  --set elasticsearchHosts=http://elasticsearch-master.logging.svc.cluster.local:9200
```

Access Kibana dashboard:
```bash
kubectl port-forward svc/kibana-kibana 5601:5601 -n logging
```

Navigate to `http://localhost:5601` in your browser.

## 🔧 Configuration

### Environment Variables

Both services support the following environment variables:

- `PORT`: Service port (default: Go=8080, Node=8081)
- `NEW_RELIC_LICENSE_KEY`: New Relic APM license key
- `LOG_LEVEL`: Logging level (debug, info, warn, error)

### Kubernetes Resources

Manifests are located in `kubernetes-manifest/`:
- `namespace.yaml`: Creates the `sre-test` namespace
- `go-deployment.yaml`: Go service deployment
- `go-service.yaml`: Go ClusterIP service (port 80 → 8080)
- `node-deployment.yaml`: Node service deployment
- `node-service.yaml`: Node ClusterIP service (port 80 → 8081)
- `ingress.yaml`: NGINX Ingress with nip.io routing

## 🧪 Testing

Run tests for both services:

```bash
# Go service
cd services/go-service && go test -v ./...

# Node service
cd services/node-service && npm test
```

## 📦 CI/CD Pipeline

The Jenkins pipeline includes:

1. **Checkout**: Clone repository
2. **Unit Tests**: Run automated tests
3. **SonarQube Analysis**: Code quality and security scanning
4. **Docker Build**: Create container images
5. **Push Images**: Upload to container registry
6. **Deploy**: Roll out to Kubernetes cluster

## 📚 Additional Resources

- [Prometheus Documentation](https://prometheus.io/docs/)
- [Grafana Dashboards](https://grafana.com/grafana/dashboards/)
- [Kubernetes Best Practices](https://kubernetes.io/docs/concepts/configuration/overview/)
- [ELK Stack Guide](https://www.elastic.co/guide/index.html)

## 🔍 Service Endpoints

| Service | Port | Endpoints |
|---|---|---|
| go-service | 8080 | `GET /` greeting, `GET /healthz` health check, `GET /metrics` Prometheus |
| node-service | 8081 | `GET /` greeting, `GET /healthz` health check, `GET /metrics` Prometheus |

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📄 License

This project is licensed under the MIT License.

## 👥 Support

For issues and questions, please open a GitHub issue or contact the SRE team.

---

**Built with ❤️ by the SRE Team**