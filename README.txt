I need a helm chart from multiple helm repos using a single url

helm repo add http://localhost:8080/helm-group
helm repo update

A helm group
    stable  https://charts.helm.sh/stable
    grafana https://grafana.github.io/helm-charts
    bitnami https://charts.bitnami.com/bitnami
    minio   https://helm.min.io

Group URL: 
    http://localhost:8080/helm-group/index.yaml

1. Client fetches index.yaml file
2. Server serves a merged index.yaml file from the group.
   The merged index chart urls are rewritten.
3. Helm fetches a chart
4. Server 


