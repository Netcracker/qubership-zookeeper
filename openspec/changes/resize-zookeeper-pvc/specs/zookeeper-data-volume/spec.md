# Spec Delta

## Purpose

Governs how the ZooKeeper operator sizes server data volumes and applies a larger size to an existing claim without replacing the volume.

## ADDED Requirements

### Requirement: Larger storage request is applied
The operator MUST update an existing managed volume's storage request when the desired size is greater than the current request.

#### Scenario: Server data volume grows
- **WHEN** a ZooKeeper server data claim already exists and `zooKeeper.storage.size` is greater than its storage request
- **THEN** the operator sets that claim's storage request to the desired size

#### Scenario: Snapshot or backup volume grows
- **WHEN** a snapshot or backup claim managed by the operator already exists and its configured size is greater than its storage request
- **THEN** the operator sets that claim's storage request to the desired size

### Requirement: Smaller storage request is rejected
The operator MUST reject a desired size that is smaller than the current storage request and MUST leave that claim unchanged.

#### Scenario: Shrink is requested
- **WHEN** the desired size of an existing managed volume is smaller than its current storage request
- **THEN** the reconcile fails with an error that shrinking is forbidden
- **AND** the claim's storage request remains the current size

### Requirement: Bound data volume expands with the deployment update
When a server data claim is bound and its capacity is below the desired size, the operator MUST publish that server's new Deployment spec with zero replicas, MUST wait until that server pod is gone, and MUST start the server again only while finishing the expansion.

#### Scenario: Size increase on a bound volume
- **WHEN** the server data claim is bound and its capacity is less than `zooKeeper.storage.size`
- **THEN** the Deployment update that carries the new pod spec also sets replicas to zero
- **AND** the operator does not start a replacement pod before the existing pod is gone

#### Scenario: Capacity reaches the desired size
- **WHEN** the server was stopped for expansion and the claim capacity is at least the desired size
- **THEN** the operator sets that server's replicas to one

#### Scenario: Filesystem is not yet expanded
- **WHEN** the server was started again and the claim capacity is still below the desired size
- **THEN** the operator sets that server's replicas back to zero and retries until the capacity is reached or the resize times out

### Requirement: Unbound data volume stays available for binding
The operator MUST update a server Deployment with one replica when its data claim is absent or not bound.

#### Scenario: Claim is not bound yet
- **WHEN** the server data claim exists, is not bound, and its capacity is below the desired size
- **THEN** the Deployment is updated with one replica

#### Scenario: Server has no data claim
- **WHEN** the server has no data PersistentVolumeClaim
- **THEN** the Deployment is updated with one replica

### Requirement: Sufficient capacity does not stop the server
The operator MUST keep one replica when a bound server data claim already has at least the desired capacity.

#### Scenario: Size is already satisfied
- **WHEN** the bound server data claim capacity is greater than or equal to `zooKeeper.storage.size`
- **THEN** the Deployment update keeps one replica

### Requirement: One server expands at a time
The operator MUST NOT expand the next server until the server being expanded is running again.

#### Scenario: Several servers need a larger volume
- **WHEN** more than one server data volume needs expansion
- **THEN** the operator expands one server
- **AND** starts the next server only after the previous server pod is running

### Requirement: Unfinished expansion fails the reconcile
The operator MUST fail the reconcile when a bound data volume does not reach the desired capacity within the resize limit.

#### Scenario: Expansion does not finish
- **WHEN** a bound data volume stays below the desired capacity until the resize limit
- **THEN** the reconcile returns an error

### Requirement: Claim update failure is logged and skipped
When the API rejects an update of an existing managed claim, the operator MUST log the error and MUST continue the reconcile.

#### Scenario: Storage class rejects the update
- **WHEN** updating an existing claim returns an error
- **THEN** the operator logs the error
- **AND** the reconcile continues
