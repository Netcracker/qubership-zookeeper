---
name: troubleshoot-zookeeper
description: Diagnose qubership-zookeeper failures in Jira tickets, CI/CD pipelines, and project support. Use when ZooKeeperService (netcracker.com/v1), the zookeeper-service chart, ZooKeeper monitoring, or the backup daemon shows ZooKeeper_Is_Degraded_Alarm, ZooKeeper_Is_Down_Alarm, ZooKeeper_CPU_Load_Alarm, ZooKeeper_Memory_Usage_Alarm, ZooKeeper_Last_Backup_Has_Failed_Alarm, ConnectionLoss, No space left on device, runAsNonRoot, or a backup-daemon PermissionError.
---

# Troubleshoot ZooKeeper

Fast entrypoint for ZooKeeper assessment. Failure steps stay in `references/troubleshooting.md`. This file does not repeat them.

`references/troubleshooting.md` is the fact source for known failure modes. It is `docs/public/troubleshooting.md` from this repository, with the memory and disk procedures that guide points at appended as their own sections. Do not read it in full for a single symptom.

1. Grep headers with line numbers first: `grep -n '^#' references/troubleshooting.md`
2. Match the symptom against the jump table below, or the raw headers if the table is stale.
3. Read from the matched header through the next header of the same or higher level. A `##` section includes its `###` children; stop at the next `##` or `#`. Never load the whole file for one lookup.

Read one more section the same way only when the section already loaded points at it:

- Memory Limit Guide → `## Memory limit procedure`
- Disk Filled Guide → `## Disk filled procedure`

## Deployment

- Chart `operator/charts/helm/zookeeper-service` renders one `ZooKeeperService` (`netcracker.com/v1`) from `templates/cr.yaml`.
- Spec groups: `global`, `zooKeeper`, `monitoring` (when `monitoring.install` is true), `backupDaemon` (only when `backupDaemon.install` is true; chart default is false). Backup storage is a persistent volume or S3 (`backupDaemon.s3`).
- Ensemble defaults: `zooKeeper.replicas` 3, data volume `zooKeeper.storage.size` 2Gi mounted at `/var/opt/zookeeper/data`. Client port 2181.
- `zooKeeper.heapSize` defaults to 256. Memory request and memory limit default to 512Mi and must stay equal. When the memory request changes, set `HEAP_OPTS` to half of the new request.
- Backup files live under `/opt/zookeeper/backup-storage`, not on the data volume. From a shell in the backup-daemon pod: `curl -XGET http://localhost:8080/health`.
- On a cluster that forces non-root: monitoring and backup daemon need `securityContext.runAsUser: 1000`. Backup-daemon volume writes need `securityContext.fsGroup: 1000` (`INFRA_ZOOKEEPER_FS_GROUP` when cloud integration is enabled).

## Forbidden actions

- Do not `rm -rf version-2` under `/var/opt/zookeeper/data` unless the loaded section is Data Is Out of Space. That delete removes the node's snapshots and transaction logs.
- Do not `kill -9 $(pidof java)` unless the Connection Loss section's scale-down of every ZooKeeper deployment already failed `zkCli.sh ls /`.
- Do not treat a full `/opt/zookeeper/backup-storage/` as a full data volume, or the reverse.

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

## Assessment

Report the section header, the evidence that selected it, the impact stated in that section, and the resolution steps from that section only. Say which check the ticket, pipeline log, or cluster could not answer.
