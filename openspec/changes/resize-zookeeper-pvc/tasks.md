# Tasks

## 1. Claim size update

- [x] 1.1 In `updatePersistentVolumeClaim`, compare the desired storage request with the current one. When the desired size is smaller, return an error that contains `PVC shrinking is forbidden` and do not call Update. When it is larger, write that quantity onto the existing claim. Verify with `go build ./operator/controllers/` and by reading the compare before the Update call.
- [x] 1.2 Skip the claim Update when the size is unchanged and annotations match the stored PVC status. Still apply a larger size when annotations are unchanged. Keep the existing annotation copy and removal of annotations the operator used to manage. Verify by reading that the early return requires both size and annotations to be unchanged.
- [x] 1.3 When the claim Update returns an API error, log it and continue the reconcile instead of failing it. Verify the error branch returns nil after the log, and `go build ./operator/controllers/` succeeds.

## 2. Server rollout

- [x] 2.1 Replace the separate Deployment update plus later scale-down with `rolloutServer`. Always write the desired Deployment. Set `replicas` to 0 in that same object only when the server data claim is `Bound` and its capacity is below `zooKeeper.storage.size`. Leave `replicas` at 1 when the claim is missing or not bound. Verify the server loop calls `rolloutServer` once and does not call a second resize restart after it.
- [x] 2.2 After a zero-replica write, wait until that server's pod list is empty, then scale to 1 and allow 30 seconds for capacity to reach the desired size. If it does not, scale to 0 and retry with a delay that starts at 10 seconds and doubles, stopping with an error at 10 minutes. Verify the wait checks that the pod object is gone, not merely `Terminating`, and that the loop returns nil only after capacity is sufficient.
- [x] 2.3 Keep the existing per-server wait for a running pod after `rolloutServer`, so the next server is not updated until the expanded server is running. Verify the running-pod wait remains after `rolloutServer` and still runs for every server.

## 3. Documentation

- [x] 3.1 In `docs/public/installation.md`, state that increasing `zooKeeper.storage.size`, `zooKeeper.snapshotStorage.volumeSize`, and `backupDaemon.backupStorage.volumeSize` requires a StorageClass with `allowVolumeExpansion: true`, and that decreasing those sizes is not supported. Verify each of those three parameter rows contains both statements.
- [x] 3.2 Update the ZooKeeper troubleshooting skill so it no longer says an existing data PVC cannot be resized. State that a larger `zooKeeper.storage.size` expands a bound claim when the StorageClass allows it, and that a smaller size is rejected. Verify the skill copies under `.agents`, `.claude`, and `agent-packages` no longer contain "does not resize".
