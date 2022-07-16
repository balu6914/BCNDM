# Installing Monitoring Stack for NDM #
## Install Loki ##
```javascript
helm repo add grafana https://grafana.github.io/helm-charts
helm repo update
helm install loki -n monitoring grafana/loki-stack
```
## Configure Prometheus ##
### Create secret to login ###
```javascript
htpasswd -c auth admin
New password:
Re-type new password:
Adding password for user admin
```
Add the credential as secret
```javascript
kubectl create secret generic prometheuspwd --from-file=auth -n monitoring
```
### Install Prometheus ###
```javascript
helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
helm repo update
helm install prometheus prometheus-community/prometheus -n monitoring -f prometheus/alerting_rules.yml -f prometheus/alertmanager.yml -f prometheus/prometheus-ingress.yml  -f prometheus/alertmanager-ingress.yml
```
This will configure the prometheus by reading the yml files. Alert Manager setting has pre-defined NDM slack group as receiver. To login into prometheus, use the created secret.
## Configure Grafana ##
### Configure datasources ###
```javascript
kubectl create cm grafana-datasources --from-file=grafana/datasources.yaml -n monitoring
kubectl label cm grafana-datasources grafana_datasource=1 -n monitoring
```
### Configure dashboards ###
```javascript
kubectl create cm grafana-charts --from-file=grafana/charts -n monitoring
kubectl label cm grafana-charts grafana_dashboard=1 -n monitoring
```
### Install Grafana ###
```javascript
helm repo add grafana https://grafana.github.io/helm-charts
helm repo update
helm install grafana grafana/grafana -n monitoring -f grafana/values.yaml
```
To get password to login to grafana,
```javascript
kubectl get secret --namespace monitoring grafana -o jsonpath="{.data.admin-password}" | base64 --decode ; echo
```

After the above steps, the applications will load at following URLs
| Application   | URL                                   |
|---------------|---------------------------------------|
| Prometheus    | https://ndm.datapace.local/prometheus |
| Alert Manager | https://ndm.datapace.local/alerts     |
| Grafana       | https://ndm.datapace.local/grafana    |
