# Proposal

## Why

ZooKeeper data volumes are created once from `zooKeeper.storage.size`. A later deployment with a larger size leaves the existing PersistentVolumeClaim unchanged, so a full disk cannot be fixed by changing the chart value. Deployment teams need to grow that storage as part of a normal rollout when the StorageClass allows expansion.

## What Changes

- Raise the storage request on an existing operator-managed PVC when the desired size is larger than the current request.
- Reject a smaller desired size. The reconcile stops with an error and does not change the claim.
- When a ZooKeeper server data PVC is already bound and its capacity is still below the desired size, publish the server Deployment with the new pod spec and `replicas: 0` in the same update, wait until that pod is gone, then start it again and repeat until the reported capacity reaches the desired size or the resize times out.
- Leave an unbound data PVC at `replicas: 1` so a volume that binds on first consumer can still be provisioned.
- Document that increasing `zooKeeper.storage.size`, snapshot volume size, and backup volume size requires `allowVolumeExpansion: true`, and that decreasing those sizes is not supported.
- Decreasing storage by copying data onto a new PVC with a migration pod is out of scope. ZooKeeper servers are Deployments, so StatefulSet `volumeClaimTemplates` immutability is also out of scope.

## Capabilities

### New Capabilities

- `zookeeper-data-volume`: How the operator creates and updates ZooKeeper server data volumes, including expansion of an existing claim and the server restart that lets the filesystem grow.

### Modified Capabilities

- None. This repository has no existing OpenSpec capabilities.

## Impact

- `operator/controllers/zookeeperservice_reconciler.go`: PVC update compares storage requests, forbids shrinking, and applies a larger request. The same update path also covers snapshot and backup claims.
- `operator/controllers/zookeeper_reconciler.go`: server rollout combines the Deployment update with an offline resize of `pvc-<name>-<id>`.
- `docs/public/installation.md`: size parameters describe expansion and the shrink rejection.
- Troubleshooting notes that said an existing data PVC is never resized.
- Requires a StorageClass with `allowVolumeExpansion: true`. Clusters that reject PVC updates still log the update error and continue; the resize then cannot finish.
