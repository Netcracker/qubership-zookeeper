# Design

## Context

InfluxDB is already unused at runtime. `Monitoring` in `netcracker.com/v1` still carries `zooKeeperVolumes`, `needToCleanInfluxDb`, `smDbHost`, and `smDbName`, marked deprecated, and the reconciler never reads them. `v1alpha1` still serves the same fields, and its monitoring object requires `needToCleanInfluxDb`, `smDbHost`, and `smDbName`. The chart writes `needToCleanInfluxDb: false` and does not copy the other three onto the custom resource. See proposal.md for why the leftovers should go and why the schema stays.

What still looks like an Influx product is documentation and one dashboard. `monitoring/grafana/zookeeper-dashboard.json` queries `influxdb` and is not the file the chart installs. The chart dashboard under `operator/charts/helm/zookeeper-service/monitoring/zookeeper-dashboard.json` queries Prometheus. Telegraf's `data_format = "influx"` is the line-protocol parser in front of `outputs.prometheus_client`, in both the chart template and `demo/monitoring/telegraf.conf`.

## Goals / Non-Goals

**Goals:**

- Delete the InfluxDB dashboard and the docs that present InfluxDB, or those four fields, as working monitoring configuration.
- Leave the four properties, their types, and the `v1alpha1` required list unchanged in the CRD and in `values.schema.json`.
- Leave Prometheus output, the chart dashboard, and the line-protocol parser unchanged.
- Leave the rendered custom resource field `needToCleanInfluxDb: false` unchanged.

**Non-Goals:**

- Removing Zabbix documentation or templates.
- Deleting the four Go fields or regenerating the CRD. Regeneration would drop the properties from the shared schema.
- Changing `additionalProperties: false` on the monitoring values object.
- Treating `data_format = "influx"` as the InfluxDB feature.
- Adding operator behavior that reads or rejects the deprecated fields.

## Decisions

### Keep the properties, stop describing them as settings

The shared Helm schema rejects unknown monitoring keys because `additionalProperties` is false. The shared CRD still serves both versions, and `v1alpha1` requires three of the fields. Removing a property, changing its type, or dropping it from that required list makes existing values or custom resources fail validation.

The Go json tags stay, so a later `make manifests` does not delete the properties. This change does not run `make generate` or `make manifests`.

Alternative: delete the fields from the API and accept a new storage version. That is a breaking CRD change and needs a separate approval. It is not this change.

### Remove the unused dashboard and the Influx docs only

Delete `monitoring/grafana/zookeeper-dashboard.json`. Remove the InfluxDB measurement section from `monitoring/README.md`, and remove or retarget the Grafana import example in that file so it does not tell readers to load the deleted dashboard. Point readers at the Prometheus dashboard the chart already installs.

In `docs/public/installation.md`, remove `monitoring.zooKeeperVolumes` from the table of working parameters. Add one compatibility note that the four keys are still accepted by the schema and are ignored. Do not document `smDbHost` or `smDbName` as new parameters.

Alternative: leave the installation row and only delete the dashboard. The row still tells operators that `zooKeeperVolumes` selects volumes, which the reconciler does not do.

### Do not touch the Telegraf parser or the chart custom resource

`data_format = "influx"` stays in the chart template and in the demo config. The output plugin stays `outputs.prometheus_client`.

`templates/cr.yaml` keeps `needToCleanInfluxDb: false`. Dropping that line would change every rendered custom resource without turning Influx off, because nothing reads the field. The other three values stay schema-legal and stay off the custom resource, which is already the case.

## Risks / Trade-offs

- [A later API cleanup deletes the four fields and regenerates the CRD] → This change does not regenerate manifests. The spec requires the properties and the `v1alpha1` required list to remain.
- [Readers treat `data_format = "influx"` as leftover InfluxDB configuration and delete it] → Tasks call out both Telegraf files and require the Prometheus output to remain.
- [The monitoring README Grafana section is the only import guide for the deleted file] → Retarget or remove that example in the same edit that deletes the dashboard.
- [Shared values still contain the four keys] → Schema validation keeps accepting them. The chart does not gain an Influx code path.

## Migration Plan

No cluster migration. Existing custom resources and values files stay valid. Deploying this change removes repository documents only; it does not rewrite live objects.

Rollback is restoring the deleted dashboard and the documentation sections. The CRD and the Helm schema do not change, so rollback does not require a CRD downgrade.

## Open Questions

None. Zabbix stays for a separate change. No Influx tests were found to delete.
