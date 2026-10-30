---
description: >-
  Diagnostic playbook for ZooKeeper ensemble, backup, disk, and client
  connection failures. Never invoke automatically — only when explicitly
  requested by the user.
applyTo: "**/*"
---

When diagnosing a ZooKeeper problem — a degraded or down ensemble, CPU or
memory pressure, a full data volume, a failed backup, ConnectionLoss,
runAsNonRoot, or a backup-daemon PermissionError — apply the
`troubleshoot-zookeeper` skill.
