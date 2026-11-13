---
name: troubleshoot-zookeeper
description: Diagnose and resolve ZooKeeper failures in ensembles managed by the `zookeeper-service` operator (`operator/charts/helm/zookeeper-service` chart). Use when ZooKeeperService (netcracker.com/v1) is degraded or down, a node hits a CPU or memory limit, the data volume is full or a snapshot will not load, the last backup failed, clients report ConnectionLoss, a pod is refused with runAsNonRoot, or the backup daemon fails with PermissionError. Also use for ZooKeeper_Is_Degraded_Alarm, ZooKeeper_Is_Down_Alarm, ZooKeeper_CPU_Load_Alarm, ZooKeeper_Memory_Usage_Alarm, ZooKeeper_Last_Backup_Has_Failed_Alarm, heap size, autopurge, zoo.cfg, or backup storage.
---

# Troubleshoot ZooKeeper

Use this skill when investigating issues related to:

- **zookeeper-service operator** (reconciliation of the `ZooKeeperService` custom resource)
- **Ensemble health** (one node down, every node down, leader election, loss of majority)
- **Storage** (data directory `/var/opt/zookeeper/data`, autopurge, a snapshot that will not load, PVC versus emptyDir)
- **Resources** (heap, equal memory request and limit, CPU)
- **Clients** (`ConnectionLoss` while pods look healthy)
- **Pod startup** (`runAsNonRoot` on monitoring or the backup daemon)
- **Backup daemon** (last backup failed, backup volume or S3, `PermissionError`)
- **Monitoring** (Prometheus alerts and the monitoring deployment)

Failure-mode facts and resolution steps live in `references/`; this file provides routing and deployment-configuration guardrails. Use only the bundled references.

## Component Overview

### Operator architecture

This skill applies to a ZooKeeper ensemble managed by the `zookeeper-service` operator deployed from `operator/charts/helm/zookeeper-service`. The operator deployment `<name>-service-operator` reconciles one `ZooKeeperService` (`netcracker.com/v1`, plural `zookeeperservices`). `<name>` is `global.name`, default `zookeeper`. The operator owns the workloads it creates. A manual edit of a Deployment, of `HEAP_OPTS`, or of `zoo.cfg` inside a pod is replaced on the next reconciliation or container start.

Every ensemble member is the same kind of pod. The chart has no master, data, arbiter, or client role, and no `javaOpts` value. Server ids start at 1.

Key components:

- **ZooKeeper** — Deployments `<name>-1` … `<name>-N` (`zooKeeper.replicas`, default 3). One replica each, Recreate strategy. Container name `zookeeper`, working directory `/opt/zookeeper`. Client service `<name>` on port 2181 (2182 is the non-encrypted client port). Headless service `<name>-server` also exposes 2888 (followers) and 3888 (election). Jolokia defaults to `zooKeeper.jolokiaPort` 9087. Port 8080 on a ZooKeeper pod is the JMX exporter.
- **Monitoring** (`<name>-monitoring`) — created when `monitoring.install` is true (chart default). `monitoring.monitoringType` defaults to `prometheus`.
- **Backup daemon** (`<name>-backup-daemon`) — created only when `backupDaemon.install` is true (chart default is false). The container and Service listen on 8080, or on 8443 when backup-daemon TLS is on (`global.tls.enabled`, `backupDaemon.install`, and `backupDaemon.tls.enabled`).
- **Credentials** — `global.secrets.zooKeeper` for admin and client users. When `vaultSecretManagement.enabled` is true (default false), the operator injects those credentials from Vault.

### Critical configuration paths

The ZooKeeper solution consists of two parts:

- **ZooKeeper** (`zookeeper-service`) — Helm chart `operator/charts/helm/zookeeper-service/`. It deploys the operator, the ensemble, monitoring, and the backup daemon. The Helm values in the table below belong to this chart.
- **ZooKeeper CRD** (`zookeeper-crd`) — the `ZooKeeperService` custom resource definition (`netcracker.com/v1`). In this repository the definition is `operator/charts/helm/zookeeper-service/crds/crd.yaml`. Apply the CRD before the `zookeeper-service` chart.

Optional overlays in `resource-profiles/` (`small`, `medium`, `large`) replace heap, CPU, and memory together. Read the live `ZooKeeperService` spec before recommending a number. Chart defaults below apply when no profile is installed.

| Helm value | Purpose |
| --- | --- |
| `zooKeeper.replicas` | Ensemble size. Default 3. `INFRA_ZOOKEEPER_REPLICAS` replaces it when `global.cloudIntegrationEnabled` is true |
| `zooKeeper.heapSize` | Heap in MiB for every member. Default 256. The operator writes `HEAP_OPTS=-Xms<n>m -Xmx<n>m` |
| `zooKeeper.resources` | CPU request 50m and limit 300m. Memory request and limit default to 512Mi and stay equal |
| `zooKeeper.storage.size` | Data-volume request used when a PVC is created. Default 2Gi. A later change does not resize an existing PVC |
| `zooKeeper.storage.className`, `.volumes`, `.labels` | Choose a PVC `pvc-<name>-<id>`. With none of these set, the data volume is emptyDir and `storage.size` is unused |
| `zooKeeper.environmentVariables` | Extra `NAME=VALUE` strings. `CONF_ZOOKEEPER_*` is written into `zoo.cfg`. `SERVER_JVMFLAGS` is appended after the heap |
| `zooKeeper.snapshotStorage` | Optional snapshot volume on ZooKeeper pods at `/opt/zookeeper/backup-storage`. Default is emptyDir |
| `zooKeeper.rollingUpdate` | When true, the operator waits for each server before updating the next. Default false |
| `monitoring.install`, `monitoring.securityContext.runAsUser` | Monitoring deployment. Use `runAsUser: 1000` when the cluster forces non-root |
| `backupDaemon.install` | Backup-daemon deployment. Default false |
| `backupDaemon.backupStorage.volumeSize` | Backup volume. Default 1Gi |
| `backupDaemon.s3.*` | S3 backup storage. Default `enabled: false` |
| `backupDaemon.backupSchedule`, `backupDaemon.evictionPolicy` | Default `0 0 * * *` and `0/1d,7d/delete` |
| `backupDaemon.securityContext` | `runAsUser: 1000` and `fsGroup: 1000` for the documented startup failures. With cloud integration, `INFRA_ZOOKEEPER_FS_GROUP` sets `fsGroup` on the backup daemon and on ZooKeeper pods |
| `global.tls.enabled`, `zooKeeper.tls.enabled` | ZooKeeper TLS is on only when both are true. `global.tls.enabled` defaults to false |

#### `zoo.cfg`

The image ships `zoo.cfg` (`dataDir=/var/opt/zookeeper/data`, `autopurge.snapRetainCount=3`, `autopurge.purgeInterval=6`, `4lw.commands.whitelist=ruok,mntr,srvr,conf,dirs`). Override a property with `zooKeeper.environmentVariables`. The entrypoint turns each underscore in `CONF_ZOOKEEPER_<property>` into a dot:

```yaml
zooKeeper:
  environmentVariables:
    - CONF_ZOOKEEPER_autopurge_snapRetainCount=3
    - CONF_ZOOKEEPER_autopurge_purgeInterval=6
    - SERVER_JVMFLAGS=-Dzookeeper.globalOutstandingLimit=800
```

`SERVER_JVMFLAGS` is for extra JVM system properties. Put heap only in `zooKeeper.heapSize`. A `-Xmx` inside `SERVER_JVMFLAGS` is appended after `HEAP_OPTS` and becomes the heap the JVM uses.

#### JVM heap

`zooKeeper.heapSize` is an integer number of mebibytes. The operator sets `-Xms` and `-Xmx` to that same number. Keep `zooKeeper.resources.requests.memory` equal to `zooKeeper.resources.limits.memory`, and set `heapSize` to half of that memory in MiB (chart default 256 of 512Mi; the same half ratio is what `resource-profiles/` uses).

```yaml
zooKeeper:
  heapSize: 256
  resources:
    requests:
      cpu: 50m
      memory: 512Mi
    limits:
      cpu: 300m
      memory: 512Mi
```

### Forbidden actions

Never recommend these without explicit user approval and ticket evidence that the loaded section already requires the action:

- `rm -rf version-2` under `/var/opt/zookeeper/data`, unless the loaded section is Data Is Out of Space. That delete removes the node's snapshots and transaction logs.
- `kill -9 $(pidof java)`, unless the Connection Loss section's scale-down of every ZooKeeper deployment already failed `./bin/zkCli.sh ls /`.
- Treating a full `/opt/zookeeper/backup-storage/` as a full data volume, or the reverse. On the backup-daemon pod that path is backup storage. On a ZooKeeper pod it is snapshot storage. Ensemble data is `/var/opt/zookeeper/data`.
- Editing `zoo.cfg` or `zoo.cfg.dynamic` inside a pod, or setting `HEAP_OPTS` with `kubectl`. Reconciliation and the container entrypoint replace them. Change `zooKeeper.heapSize`, `zooKeeper.resources`, `zooKeeper.replicas`, or `zooKeeper.environmentVariables` and redeploy through the chart.
- Deleting a PVC, PV, Secret, or znode, or changing credentials, outside the documented procedure. The operator does not resize an existing data PVC; deleting it drops the ensemble data on that node.
- Raising `autopurge.snapRetainCount` to free disk space. A higher retain count keeps more snapshots.

## Common Failure Modes

### Data volume full

Symptom: `No space left on device` while writing `/var/opt/zookeeper/data/version-2/currentEpoch`, or `Unable to load database on disk` after the disk filled.

Check `df -h /var/opt/zookeeper/data` and whether the `data` volume is a PVC (`pvc-<name>-<id>`) or emptyDir before changing configuration. `autopurge.snapRetainCount` defaults to 3 snapshots and `autopurge.purgeInterval` defaults to 6 hours. Change them only through `CONF_ZOOKEEPER_autopurge_snapRetainCount` and `CONF_ZOOKEEPER_autopurge_purgeInterval` in `zooKeeper.environmentVariables`. Raising the retain count fills the disk sooner. `zooKeeper.storage.size` is applied when the PVC is created. The operator does not resize an existing data PVC. Follow the Disk Filled on All Nodes section and its scenario for cleanup. Follow Data Is Out of Space only when the node cannot load its snapshot.

### Heap pressure

Symptom: `java.lang.OutOfMemoryError: Java heap space`, `ZooKeeper_Memory_Usage_Alarm`.

Tune `zooKeeper.heapSize` on the ensemble. Keep `-Xms` and `-Xmx` equal by using that single integer, and keep the memory request equal to the memory limit at twice the heap. The Memory Limit scenario may lower outstanding requests with `SERVER_JVMFLAGS=-Dzookeeper.globalOutstandingLimit=800`. Raising `globalOutstandingLimit` does not relieve heap pressure. There is no per-role heap value.

### CPU at the pod limit

Symptom: `ZooKeeper_CPU_Load_Alarm`, `unable to create new native thread`.

The CPU Limit section owns this signal. Raise `zooKeeper.resources.requests.cpu` and `zooKeeper.resources.limits.cpu` together through the chart. The chart has no thread-pool queue setting. When the same log also contains `Java heap space`, read Memory Limit first.

### Backup storage filled

Symptom: `No space left on device` under `/opt/zookeeper/backup-storage/` on the backup daemon, `ZooKeeper_Last_Backup_Has_Failed_Alarm`.

That volume is `backupDaemon.backupStorage` (default 1Gi) or `backupDaemon.s3` when `backupDaemon.s3.enabled` is true. It is independent of `zooKeeper.storage`. The backup daemon exists only when `backupDaemon.install` is true. The documented health check, from a shell in the backup-daemon pod, is `curl -XGET http://localhost:8080/health`. When backup-daemon TLS is enabled the container listens on 8443 instead. ZooKeeper's own port 8080 is the JMX exporter and does not serve this health check.

## Using the Troubleshooting Guide and Scenarios

Never read `references/troubleshooting.md` in full; load only the section that matches the symptom.

[references/troubleshooting.md](./references/troubleshooting.md) is the fact source for known failure modes. It is a copy of `docs/public/troubleshooting.md`. [references/alerts.md](./references/alerts.md) and [references/troubleshooting-scenarios/](./references/troubleshooting-scenarios/) are the files that guide links to.

- Memory Limit Guide → [references/troubleshooting-scenarios/memory_limit.md](./references/troubleshooting-scenarios/memory_limit.md)
- Disk Filled Guide → [references/troubleshooting-scenarios/disk_filled_on_all_nodes.md](./references/troubleshooting-scenarios/disk_filled_on_all_nodes.md)

Local links resolve relative to the file containing them. GitHub links point to repository docs that are not bundled; mention them to the user instead of fetching them.

Steps:

1. Match the symptom in the table below. Tickets rarely quote alert names: match on pod state, error text, or user-visible behavior first. If the table looks stale, list headers with `grep -n '^## ' references/troubleshooting.md`.
2. Read from the matched header through the next header of the same or higher level. A `##` section includes its `###` children; stop at the next `##` or `#`.
3. Open the scenario linked from the matched section only. Do not read unrelated scenarios.

## Failure scenarios

Match one row, then read that section. No resolution steps live in this table.

| Symptom | Signal | Section header |
| --- | --- | --- |
| One ZooKeeper node failed, ensemble degraded | `ZooKeeper_Is_Degraded_Alarm`, `zookeeper_status_code == 5` | `## Single Node Failure` |
| Every ZooKeeper node failed | `ZooKeeper_Is_Down_Alarm`, `zookeeper_status_code == 10` | `## All Nodes Failure` |
| CPU at the pod limit | `ZooKeeper_CPU_Load_Alarm`, `unable to create new native thread` | `## CPU Limit` |
| Memory at the pod limit | `ZooKeeper_Memory_Usage_Alarm`, `Java heap space` | `## Memory Limit` |
| Data volume full | `No space left on device` writing `/var/opt/zookeeper/data/version-2/currentEpoch` | `## Disk Filled on All Nodes` |
| Last backup failed | `ZooKeeper_Last_Backup_Has_Failed_Alarm`, `zookeeper_backup_metric_last_backup_status == 4` | `## Last Backup Has Failed` |
| Node cannot load its snapshot | `Unable to load database on disk` | `## Data Is Out of Space` |
| Clients lose the session while pods look healthy | `ConnectionLoss`, `All connection tries failed` | `## Clients Cannot Connect To ZooKeeper with Connection Loss Error` |
| Monitoring or backup daemon refused at start | `container has runAsNonRoot and image will run as root` | `## Container Failed with Error: container has runAsNonRoot and image will run as root` |
| Backup daemon cannot create storage | `PermissionError` on `/opt/zookeeper/backup-storage/granular` | `## ZooKeeper Backup Daemon Failed with PermissionError` |

When several rows match, read the more specific section first: start failure (`runAsNonRoot`, `PermissionError`), then disk or memory, then backup storage, then `ConnectionLoss`, then degraded or down.

If no row matches (for example, `Pending` or `Evicted`), say the reference does not cover it and ask for `kubectl describe pod` events. Do not infer a capacity problem without evidence.

## Log Analysis Tips

Read logs from the ticket attachments before running `kubectl logs`. Use an attached operator, ZooKeeper, backup-daemon, or monitoring log when it covers the failing component and the time of the symptom. Run `kubectl logs` only when the ticket has no such attachment, or the attachment does not contain the error you need.

Use `<name>` from the release (`global.name`, default `zookeeper`). Run read-only commands only.

### Operator logs

```bash
# Reconciliation errors
kubectl logs -n <namespace> deployment/<name>-service-operator | grep -iE "error|failed"

# Custom resource, including the spec the operator last applied
kubectl get zookeeperservices -n <namespace> -o yaml
kubectl get pods,pvc -n <namespace> -o wide
```

### ZooKeeper pod logs

Members are Deployments `<name>-1` through `<name>-N`.

```bash
# Disk, heap, election, and client sessions
kubectl logs -n <namespace> deployment/<name>-1 | grep -iE "No space left on device|OutOfMemoryError|unable to create new native thread|Unable to load database|ConnectionLoss|Timeout while waiting for epoch"
```

From a shell in a ZooKeeper container (`/opt/zookeeper`). Four-letter words are limited to `ruok`, `mntr`, `srvr`, `conf`, and `dirs`. When ZooKeeper TLS is enabled, port 2181 is the secure client port and plaintext `nc` does not answer; use `zkCli.sh`.

```bash
df -h /var/opt/zookeeper/data
echo mntr | nc localhost 2181
echo srvr | nc localhost 2181
./bin/zkCli.sh ls /
```

### Backup daemon and monitoring

```bash
kubectl logs -n <namespace> deployment/<name>-backup-daemon | grep -iE "error|PermissionError|No space"
kubectl logs -n <namespace> deployment/<name>-monitoring | grep -iE "error|timeout"
```

Health, from a shell in the backup-daemon pod. The troubleshooting guide uses port 8080. When backup-daemon TLS is enabled, the container listens on 8443.

```bash
curl -XGET http://localhost:8080/health
```

## Safety

Commands in the references are recommendations for an authorized operator, not permission to run them. Before deleting `version-2`, killing the Java process, deleting a PVC, PV, or Secret, or changing credentials, get explicit user approval and confirm the namespace, the server id, and whether a restorable backup exists.

## Assessment

Report the section header, the evidence that selected it, the impact stated in that section, and the resolution steps from that section only. Say which check the ticket, pipeline log, or cluster could not answer. Do not invent a root cause.

### Recommendation discipline

- Phrase each `zoo.cfg` change as: “Set `CONF_ZOOKEEPER_<property_with_dots_as_underscores>=<value>` in `zooKeeper.environmentVariables` in the `zookeeper-service` chart values and redeploy through the normal delivery channel.”
- Phrase each heap change as: “Set `zooKeeper.heapSize: <MiB>` and set `zooKeeper.resources.requests.memory` and `zooKeeper.resources.limits.memory` to the same value, twice that heap, in the `zookeeper-service` chart values and redeploy.” Identify the live spec, including a `resource-profiles/` overlay when one is installed.
- Phrase each CPU change as: “Set `zooKeeper.resources.requests.cpu` and `zooKeeper.resources.limits.cpu` in the `zookeeper-service` chart values and redeploy.”
- Do not recommend `rm -rf version-2`, `kill -9`, znode deletion, or PVC deletion as a configuration fix. Those actions are operational, evidence-bound, and may destroy ensemble data.

## When to Escalate

Escalate to component maintainers when:

- The ensemble stays down after the All Nodes Failure section's reboot of every ZooKeeper deployment.
- A node cannot load its snapshot, the documented fix would delete `version-2`, and no usable backup exists.
- `ConnectionLoss` remains after the section's scale-down of every ZooKeeper deployment and the documented process kill.
- A majority of member pods are Running and no leader is elected.
- The backup daemon still fails after the permission and storage checks, and a restore is required.
- The operator crashes or loops on reconciliation (attach the stack trace).

Include in the escalation:

- Operator logs and logs of the affected `<name>-N` deployments.
- `kubectl get zookeeperservices -n <namespace> -o yaml` and `kubectl get pods,pvc -n <namespace> -o wide`.
- `df -h /var/opt/zookeeper/data`, plus `mntr` / `srvr` or `./bin/zkCli.sh ls /` when the process answers.
- Backup-daemon logs and the health output when the symptom is a backup.
