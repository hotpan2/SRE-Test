----------------------------------------------------------------------
SRE Test – Go & Node Microservices + CI/CD + Kubernetes Deployment
----------------------------------------------------------------------
This project contains two microservices (Go & Node.js) deployed through a complete CI/CD process using Jenkins, Docker, SonarQube, and Kubernetes.

It includes:
Go service (services/go-service)
Node service (services/node-service)
Jenkins CI/CD pipeline
Docker image builds
Image tagging
Kubernetes deployments
Prometheus metrics 
Grafana dashboards

---------------------------------------------------------
 CI/CD Workflow
-------------------------------------------------------
Checkout -> Run Unit Test -> Sonarqube Analysis -> Build Docker Images -> Push Images ->  Deploy to Kubernetes

---------------------------------------------------------
 How to Run Locally
---------------------------------------------------------
- Run Go Service
cd services/go-service
go mod tidy
go run main.go

Test:
go test ./...

- Run Node Service
cd services/node-service
npm install
npm start

Test:
npm test

- Build Docker Images
Go:
docker build -t sre-go:test services/go-service
Node:
docker build -t sre-node:test services/node-service

- Deploy to Kubernetes
export IMAGE_TAG=test
envsubst < kubernetes-manifest/services/go-deployment.yaml | kubectl apply -f -
envsubst < kubernetes-manifest/services/node-deployment.yaml | kubectl apply -f -
--------------------------------------------------------
Monitoring/Logging
--------------------------------------------------------

- New Relic Monitoring
Inject license key via environment variable in deployment.yaml   

- Install Prometheus&Grafana
helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
helm repo update
helm install monitoring prometheus-community/kube-prometheus-stack -n monitoring --create-namespace

- Running Grafana
kubectl port-forward svc/monitoring-grafana -n monitoring 3000:80

- Running ELK (elasticsearch, logstash, kibana)
Install Elasticsearch : 
helm repo add elastic https://helm.elastic.co
helm repo update
helm install elasticsearch elastic/elasticsearch -n logging \
  --set replicas=1 \
  --set minimumMasterNodes=1 \
  --set resources.requests.cpu=500m \
  --set resources.requests.memory=1Gi

- Install Filebeat / Fluent Bit:
helm repo add elastic https://helm.elastic.co
helm install filebeat elastic/filebeat -n logging \
  --set daemonset.enabled=true \
  --set output.elasticsearch.hosts={http://elasticsearch-master.logging.svc.cluster.local:9200}

- Install Kibana:
helm install kibana elastic/kibana -n logging \
  --set elasticsearchHosts=http://{elasticservice}:9200

- Kibana Dashboard:
kubectl port-forward svc/kibana-kibana 5601:5601 -n logging




