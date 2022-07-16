# Installing Monitoring Stack for NDM #
## Install Loki ##
```javascript
helm install loki -n monitoring grafana/loki-stack
```
## Install Prometheus ##
```javascript
helm install prometheus prometheus-community/prometheus -n monitoring -f prometheus/alerting_rules.yml -f prometheus/alertmanager.yml -f prometheus/prometheus-ingress.yml  -f prometheus/alertmanager-ingress.yml
```
This will configure the prometheus by reading the yml files. Alert Manager setting has pre-defined NDM slack group as receiver.
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
helm install grafana grafana/grafana -n monitoring -f grafana/values.yaml
```

After the above steps, the applications will load at following URLs
| Application   | URL                                   |
|---------------|---------------------------------------|
| Prometheus    | https://ndm.datapace.local/prometheus |
| Alert Manager | https://ndm.datapace.local/alerts     |
| Grafana       | https://ndm.datapace.local/grafana    |
