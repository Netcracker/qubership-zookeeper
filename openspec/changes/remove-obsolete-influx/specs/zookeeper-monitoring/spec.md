# Spec Delta

## Purpose

Defines ZooKeeper monitoring output and how deprecated Influx fields stay valid on the shared schema without configuring InfluxDB.

## ADDED Requirements

### Requirement: Prometheus is the monitoring output
Installed monitoring MUST publish metrics through Prometheus and MUST NOT configure an InfluxDB output.

#### Scenario: Default monitoring install
- **WHEN** monitoring is installed with the default monitoring type
- **THEN** the Telegraf configuration contains a Prometheus client output
- **AND** it does not contain an InfluxDB output plugin

### Requirement: Deprecated Influx fields stay on the shared schema
The Helm values schema and both served CRD versions MUST still accept `monitoring.zooKeeperVolumes`, `monitoring.needToCleanInfluxDb`, `monitoring.smDbHost`, and `monitoring.smDbName` with their current types. The `v1alpha1` monitoring object MUST still require `needToCleanInfluxDb`, `smDbHost`, and `smDbName`.

#### Scenario: Existing values pass schema validation
- **WHEN** a values file sets any of those four monitoring keys
- **THEN** Helm values schema validation accepts the file

#### Scenario: Existing custom resource stays valid
- **WHEN** a `ZooKeeperService` of `netcracker.com/v1` or `netcracker.com/v1alpha1` includes those monitoring fields
- **THEN** the CRD schema still accepts the object

### Requirement: Deprecated Influx fields are ignored
The operator MUST NOT use those four fields to configure InfluxDB, clean a database, or select volumes.

#### Scenario: Fields are present on the custom resource
- **WHEN** a custom resource sets `needToCleanInfluxDb`, `smDbHost`, `smDbName`, or `zooKeeperVolumes`
- **THEN** reconciliation does not create an InfluxDB client, job, or output
- **AND** monitoring still uses the Prometheus configuration

#### Scenario: Chart renders a stable custom resource
- **WHEN** the Helm chart renders a `ZooKeeperService`
- **THEN** it still sets `needToCleanInfluxDb` to false
- **AND** it does not copy `smDbHost`, `smDbName`, or `zooKeeperVolumes` onto the custom resource

### Requirement: InfluxDB dashboard is not part of the product
The repository MUST NOT ship an InfluxDB Grafana dashboard or document InfluxDB as a monitoring destination. The dashboard installed by the chart MUST remain the Prometheus dashboard.

#### Scenario: Chart dashboard
- **WHEN** the chart installs the Grafana dashboard
- **THEN** that dashboard queries Prometheus

#### Scenario: Obsolete dashboard is gone
- **WHEN** a reader looks for the InfluxDB dashboard and the InfluxDB measurement section
- **THEN** those documents are not in the repository

### Requirement: Telegraf line protocol parser stays
Exec-script metrics MUST still be parsed with Telegraf's influx line-protocol data format. That setting is the parser for script output, not an InfluxDB server.

#### Scenario: Exec input format
- **WHEN** monitoring configuration is rendered
- **THEN** the exec input data format remains `influx`
- **AND** the output plugin remains Prometheus
