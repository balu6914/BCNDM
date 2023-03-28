# Installing Monitoring Stack for BCNDM #
## Configure Loki ##
### Configure AWS S3 as storage ###
1. Create an IAM user. 2. Create a S3 bucket and provide access to the created IAM user. 3. Create a secret using IAM 
credentials. ```javascript kubectl create secret generic loki-s3 --from-literal=AWS_ACCESS_KEY_ID='<access_key>' 
--from-literal=AWS_SECRET_ACCESS_KEY='<secret_key>' -n monitoring ``` 4. Install Loki with the S3 details. 
```javascript helm repo add grafana https://grafana.github.io/helm-charts helm repo update helm install loki -n 
monitoring grafana/loki-stack -f loki/loki_aws_s3.yaml -f loki/promtail.yaml --set 
loki.config.storage_config.aws.s3=s3://<region_name>/<bucket_name> ```
### Configure Azure Blob as storage ###
1. Create a storage account and a container in the Azure cloud. 2. Pass the details in the command as below 
```javascript helm repo add grafana https://grafana.github.io/helm-charts helm repo update helm install loki -n 
monitoring grafana/loki-stack -f loki/loki_azure_blob.yaml -f loki/promtail.yaml --set 
loki.config.storage_config.azure.container_name=<container_name> --set 
loki.config.storage_config.azure.account_name=<account_name> --set 
loki.config.storage_config.azure.account_key=<account_key> ```
### Configure Local Filesystem as storage ###
```javascript helm repo add grafana https://grafana.github.io/helm-charts helm repo update helm install loki -n 
monitoring grafana/loki-stack -f loki/loki_local.yaml -f loki/promtail.yaml ```
## Configure Prometheus ##
### Create secret to login ###
```javascript htpasswd -c auth admin New password: Re-type new password: Adding password for user admin ``` Add the 
credential as secret ```javascript kubectl create secret generic prometheuspwd --from-file=auth -n monitoring ```
### Install Prometheus ###
```javascript helm repo add prometheus-community https://prometheus-community.github.io/helm-charts helm repo update 
helm install prometheus prometheus-community/prometheus -n monitoring -f prometheus/alerting_rules.yml -f 
prometheus/alertmanager.yml -f prometheus/prometheus-ingress.yml -f prometheus/alertmanager-ingress.yml ``` This will 
configure the prometheus by reading the yml files. Alert Manager setting has pre-defined BCNDM slack group as receiver. 
To login into prometheus, use the created secret.
## Configure Grafana ##
### Configure datasources ###
```javascript kubectl create cm grafana-datasources --from-file=grafana/datasources.yaml -n monitoring kubectl label 
cm grafana-datasources grafana_datasource=1 -n monitoring ```
### Configure dashboards ###
```javascript kubectl create cm grafana-charts --from-file=grafana/charts -n monitoring kubectl label cm 
grafana-charts grafana_dashboard=1 -n monitoring ```
### Install Grafana ###
```javascript helm repo add grafana https://grafana.github.io/helm-charts helm repo update helm install grafana 
grafana/grafana -n monitoring -f grafana/values.yaml ``` To get password to login to grafana, ```javascript kubectl 
get secret --namespace monitoring grafana -o jsonpath="{.data.admin-password}" | base64 --decode ; echo ``` After the 
above steps, the applications will load at following URLs
| Application | URL | ---------------|---------------------------------------| Prometheus | 
| https://ndm.datapace.local/prometheus | Alert Manager | https://ndm.datapace.local/alerts | Grafana | 
| https://ndm.datapace.local/grafana |
