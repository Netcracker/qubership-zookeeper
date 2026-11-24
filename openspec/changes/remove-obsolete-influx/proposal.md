# Proposal

## Why

InfluxDB is no longer a ZooKeeper monitoring backend, but leftovers still describe it as one: an unused Grafana dashboard, the monitoring README, and an installation parameter. The CRD and the Helm values schema are shared across installs, so deleting those properties would reject values and custom resources that still send them.

## What Changes

- Remove the unused InfluxDB Grafana dashboard and the InfluxDB section of the monitoring README. The chart already ships a Prometheus dashboard.
- Stop describing `monitoring.zooKeeperVolumes`, `monitoring.needToCleanInfluxDb`, `monitoring.smDbHost`, and `monitoring.smDbName` as working settings. The operator already ignores them.
- Keep those four properties in the `netcracker.com/v1` and `netcracker.com/v1alpha1` CRD schemas and in `values.schema.json`, with the same types and the same `v1alpha1` required list. Existing values and custom resources must still validate.
- Leave Telegraf `data_format = "influx"` in place. That is the line-protocol parser for exec-script output, and the output plugin stays Prometheus.
- Zabbix removal is out of scope.

## Capabilities

### New Capabilities

- `zookeeper-monitoring`: Which monitoring output ZooKeeper supports, and how deprecated Influx fields stay valid on the shared schema while being ignored.

### Modified Capabilities

- None. This repository has no existing OpenSpec capabilities.

## Impact

- `monitoring/grafana/zookeeper-dashboard.json`: InfluxDB dashboard, not referenced by the Helm chart.
- `monitoring/README.md`: InfluxDB measurement section.
- `docs/public/installation.md`: `monitoring.zooKeeperVolumes` is documented as a live parameter.
- `operator/api/v1/zookeeperservice_types.go` and `operator/api/v1alpha1/zookeeperservice_types.go`: fields stay so generated CRDs keep the properties.
- `operator/charts/helm/zookeeper-service/values.schema.json`: `monitoring` has `additionalProperties: false`, so removing a property fails validation of shared values that still set it.
- `operator/config/crd/bases/netcracker.com_zookeeperservices.yaml` and `operator/charts/helm/zookeeper-service/crds/crd.yaml`: both served versions keep the properties. `v1alpha1` still requires `needToCleanInfluxDb`, `smDbHost`, and `smDbName`.
- `operator/charts/helm/zookeeper-service/templates/cr.yaml`: still emits `needToCleanInfluxDb: false`. `smDbHost`, `smDbName`, and `zooKeeperVolumes` are not copied onto the custom resource.
- No Influx-specific unit or Robot tests exist in this repository.
- Telegraf configs under the chart and `demo/monitoring/telegraf.conf` keep `data_format = "influx"`.
