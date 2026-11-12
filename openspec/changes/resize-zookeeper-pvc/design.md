# Design

## Context

ZooKeeper servers are Deployments with strategy `Recreate`, one replica each, and a data claim named `pvc-<service>-<id>`. The claim is created from `zooKeeper.storage.size` and, before this change, an existing claim was updated only when annotations changed. Snapshot and backup claims share that update function. See proposal.md for why the size must become mutable.

The ticket also asked for shrinking through a migration pod and for StatefulSet `volumeClaimTemplates`. Those are excluded: this operator does not use a StatefulSet, and the implemented path only grows a claim in place.

## Goals / Non-Goals

**Goals:**

- Apply a larger storage request on the existing claim.
- Take a bound server volume offline in the same Deployment write that publishes the new pod spec.
- Bring that server back only after the pod that held the volume is gone, and repeat until capacity catches up.
- Keep first-consumer binding working for a claim that is not bound yet.
- Expand one server at a time so the ensemble keeps quorum.

**Non-Goals:**

- Shrinking a claim, copying data to a replacement claim, or handling immutable `volumeClaimTemplates`.
- Restarting the backup daemon or snapshot consumer when those claims grow. Only their storage request is updated.
- Online filesystem growth while the server pod keeps the volume mounted.

## Decisions

### Grow the existing claim

The operator writes the larger size into `spec.resources.requests.storage` of the claim it already owns. A smaller quantity fails the reconcile before any write.

Alternative: create a second claim and copy data with a migration pod, which is what the ticket suggested for shrinking and for classes that cannot expand. That replaces the volume identity and is a separate procedure. In-place expansion matches classes with `allowVolumeExpansion: true` and keeps the same claim name the Deployment already mounts.

### Put replicas at zero in the Deployment update

`rolloutServer` always writes the desired Deployment. If the data claim is `Bound` and its capacity is still below the desired size, that write sets `replicas: 0`. The following wait does not run until the resize loop has started the server again.

Alternative: write `replicas: 1`, let `Recreate` start a pod, then scale to zero. The new pod is scheduled onto a volume that is still too small or still resizing, and the scale-down deletes it immediately. One write avoids that second start.

A claim that is missing or not `Bound` stays at one replica. Capacity is often empty until the first pod binds the volume (`WaitForFirstConsumer`). Scaling that server to zero would prevent the claim from ever binding.

### Wait until the pod object is gone

After `replicas: 0`, the loop waits until no pod with that server's selector remains, then waits a growing delay (10s, 20s, 40s, …) and scales to one. It allows 30 seconds for capacity to reach the desired size. If it does not, the server is scaled to zero and the attempt repeats. The whole loop stops after 10 minutes with an error.

`Terminating` is set as soon as deletion starts. The volume stays mounted until the container exits and kubelet unmounts it; the pod object is removed after that unmount. Waiting for an empty pod list is the signal that the volume is free. Waiting only for `Terminating` would start the next pod while the claim is still mounted.

The existing per-server wait for a running pod stays after this loop, so the next server is not rolled until the expanded server is up.

### Leave annotation handling in the claim update

Size is checked before the annotation early return. If neither size nor annotations changed, the claim is not written. A size change still updates the claim when annotations are unchanged. Snapshot and backup claims use this same function, so they grow and refuse to shrink, without the server scale-down.

If the API rejects the claim update, the error is logged and the reconcile continues. That preserves the previous behavior on clusters that cannot mutate claims.

## Risks / Trade-offs

- [One server is down while its volume expands] → The loop handles a single server, and the next server waits until this pod is running. A three-node ensemble keeps a majority if only one member is down.
- [StorageClass without `allowVolumeExpansion`, or an API that rejects the claim update] → The update error is only logged. The server is still scaled to zero because capacity is below the desired size, then the 10-minute loop fails. The next reconcile repeats that outage until the size is put back or the class allows expansion.
- [Rollback to an operator that does not resize] → An already expanded claim stays large. A Deployment left at zero replicas is restored to one only when a later reconcile writes the server Deployment again.
- [Filesystem growth that needs a mount] → The loop starts the pod, and if capacity is still short after 30 seconds it stops the pod and retries. Drivers that expand only while detached are covered by the delay before the first start.

## Migration Plan

1. Deploy the operator that contains the claim update and `rolloutServer`. No custom resource schema change is required.
2. Confirm the data StorageClass has `allowVolumeExpansion: true` before raising `zooKeeper.storage.size`.
3. Raise the size in the chart values and roll the release. The operator expands one server at a time.
4. Rollback of the operator binary does not shrink claims. Before rollback, confirm no server Deployment is stuck at zero replicas from an unfinished resize.

## Open Questions

None. Shrink-by-copy and StatefulSet claim templates stay out of this change.
