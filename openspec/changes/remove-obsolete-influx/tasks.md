# Tasks

## 1. Remove the InfluxDB documents

- [x] 1.1 Delete `monitoring/grafana/zookeeper-dashboard.json`. Verify the file is gone and `operator/charts/helm/zookeeper-service/monitoring/zookeeper-dashboard.json` still contains `"type": "prometheus"`.
- [x] 1.2 In `monitoring/README.md`, remove the InfluxDB measurement section and stop the Grafana import example from loading that dashboard. Point the example at `operator/charts/helm/zookeeper-service/monitoring/zookeeper-dashboard.json`. Leave the Zabbix section in place. Verify the file has no `InfluxDB` heading and still has the `Zabbix` heading.
- [x] 1.3 In `docs/public/installation.md`, remove `monitoring.zooKeeperVolumes` from the working-parameter table. Add one note that `monitoring.zooKeeperVolumes`, `monitoring.needToCleanInfluxDb`, `monitoring.smDbHost`, and `monitoring.smDbName` are still accepted and ignored. Verify the parameter table no longer lists `zooKeeperVolumes` and the note names all four keys.

## 2. Confirm the shared schema and Prometheus path stay

- [x] 2.1 Do not edit the Go API types, either CRD file, or `values.schema.json`. Verify `values.schema.json` still has `additionalProperties: false` on monitoring and still defines `zooKeeperVolumes`, `needToCleanInfluxDb`, `smDbHost`, and `smDbName`. Verify both served versions in `operator/charts/helm/zookeeper-service/crds/crd.yaml` still define those properties, and the `v1alpha1` monitoring object still requires `needToCleanInfluxDb`, `smDbHost`, and `smDbName`.
- [x] 2.2 From `operator/charts/helm/zookeeper-service`, run `helm template zookeeper-service .`. Verify the rendered custom resource still contains `needToCleanInfluxDb: false` and does not contain `smDbHost`, `smDbName`, or `zooKeeperVolumes`. Verify the rendered Telegraf config contains `data_format = "influx"` and `[[outputs.prometheus_client]]`, and does not contain `outputs.influxdb`. Verify `demo/monitoring/telegraf.conf` still contains `data_format = "influx"`.
