# AGENTS.md

Qubership ZooKeeper is an enterprise ZooKeeper distribution for Kubernetes/OpenShift.
It bundles a controller-runtime operator, one Helm chart, an Apache ZooKeeper
container image, a Python backup daemon, a Telegraf-based monitoring image,
and a Robot Framework integration test suite.

## Repository Layout

Monorepo. `go.mod` files live under `operator/` and `docker-zookeeper/`.
The `backup-daemon` scripts and the monitoring exec scripts are Python.

- `/operator` — controller-runtime operator (module
  `github.com/Netcracker/qubership-zookeeper/operator`). Single binary. It
  watches one namespace from `WATCH_NAMESPACE` and reconciles
  `ZooKeeperService` (`netcracker.com/v1`): the ensemble plus monitoring,
  backup daemon, and integration tests.
- `/operator/charts/helm/zookeeper-service` — Helm chart that installs the
  operator and a `ZooKeeperService` CR (`netcracker.com/v1`).
- `/backup-daemon` — Python backup/restore scripts on top of
  `qubership-backup-daemon-go`. Stores snapshots on a Persistent Volume
  (path configured in `backup-daemon.conf`).
- `/monitoring` — Telegraf image plus Python exec scripts for health, backup,
  and version metrics. Grafana dashboard and Prometheus alert rules ship in
  the Helm chart.
- `/docker-zookeeper` — Custom Apache ZooKeeper image (Temurin 21) with SASL,
  JMX exporter, and health checks. The Go binary `zookeeper-assistant`
  (`main.go`) serves the health check and the in-container backup endpoint.
- `/integration-tests` — Robot Framework suites at `robot/tests/`.
- `/demo` — Docker Compose stack (ZooKeeper + backup daemon + integration-tests
  web terminal on `localhost:8090`).
- `/docs/public` — User-facing documentation (architecture, installation,
  troubleshooting, security, backup, monitoring).
- `/docs/internal` — Developer docs (`developing.md`, `operator-guide.md`).
- `/operator/Makefile` — Drives operator codegen (`make generate`,
  `make manifests`).

## Operator Code Organization

- `operator/api/v1` — current `ZooKeeperService` type. `//+kubebuilder:storageversion`
  is on this version. Spec sections are `global`, `zooKeeper`, `monitoring`,
  `backupDaemon`, `vaultSecretManagement`, and `integrationTests`.
- `operator/api/v1alpha1` — previous API, still registered in the scheme.
  The reconciler does not watch it. Change `v1`.
- `operator/controllers/zookeeperservice_controller.go` — `ZooKeeperServiceReconciler`.
  `buildReconcilers` dispatches to the service reconcilers below.
- `operator/controllers/zookeeper_reconciler.go` — ZooKeeper ensemble.
- `operator/controllers/monitoring_reconciler.go` — Monitoring.
- `operator/controllers/backupdaemon_reconciler.go` — Backup daemon.
- `operator/controllers/integrationtests_reconciler.go` — Integration-test pod.
- `operator/controllers/condition.go` — CR conditions. Status writes live on
  the reconciler (`updateStatus`, `writeFailedStatus`).
- `operator/controllers/*_vault.go` and `vault_client.go` — Vault secret
  management.
- `operator/controllers/provider/` — Builders that construct Kubernetes
  objects (StatefulSet, Deployment, Service, Secret, ConfigMap, PVC). Keep
  manifest construction in this package.
- Flags and `WATCH_NAMESPACE` are read in `operator/main.go`.

## Compatibility Rules

All changes MUST be backward compatible.

### CRDs

- Add new functionality through new optional fields. Do not change the
  semantics or shape of existing fields.
- If a breaking change is unavoidable, stop and ask the user for explicit
  approval before adding a new API version. The procedure is in
  [docs/internal/operator-guide.md](docs/internal/operator-guide.md): new
  package under `operator/api/`, move `//+kubebuilder:storageversion`, point
  controllers and providers at the new package.
- After editing anything under `operator/api/`, run `make generate` and
  `make manifests` from `operator/`.
- Copy the generated CRD from
  `operator/config/crd/bases/netcracker.com_zookeeperservices.yaml` to
  `operator/charts/helm/zookeeper-service/crds/crd.yaml`.
  `make manifests` updates only the file under `config/crd/bases`.
- Keep `apiVersion: netcracker.com/v1` in
  `operator/charts/helm/zookeeper-service/templates/cr.yaml` unless the user
  approved a new storage version.

### Helm Chart

- Every new parameter MUST have a sensible default. Existing values files
  (`values.yaml`, `vals.yaml`, `resource-profiles/`) must keep working.
- Reuse shared templates in `templates/_helpers.tpl` for labels, annotations,
  images, and security contexts.
- Render the chart with `helm template` before submitting chart changes:

```sh
cd operator/charts/helm/zookeeper-service
helm template zookeeper-service .
```

## Dev-Kit (Operator Code Generation)

`operator/Makefile` downloads `controller-gen` and can run `make generate` and
`make manifests` on the host. Use the dev-kit when `operator-sdk` is not
installed locally. It runs `operator-sdk:v1.42` in Docker with `operator/`
mounted at `/zookeeper-service-operator`.

```sh
cd operator/dev-kit
./terminal.sh        # Builds and enters the operator-sdk container
make generate        # Regenerate DeepCopy code from operator/api/
make manifests       # Regenerate CRDs in config/crd/bases
```

## Integration Tests

Robot Framework suites live at `integration-tests/robot/tests/`. The `demo/`
directory provides a Docker Compose stack for local testing. The web terminal
is `http://localhost:8090`.

### Mandatory – every change must be tested

- Any change under a directory that contains a Dockerfile must be tested with
  a Docker build (`docker-zookeeper/docker`, `backup-daemon/docker`,
  `monitoring/docker`, `operator`, `integration-tests/docker`).
- Any code change must be covered with unit tests. Monitoring exec scripts
  already have `monitoring/exec-scripts/test_*.py`. Operator tests run with
  `make test` from `operator/`.
- Any functional feature that can be tested without a cluster redeploy must
  be covered with a Robot Framework test.
- Before submitting a change under `backup-daemon` or `integration-tests`,
  run at minimum the `zookeeper_crud` and `zookeeper_backup_daemon` suites.

```sh
cd demo
docker-compose up -d
docker-compose exec zookeeper-integration-tests robot -i zookeeper_crudORzookeeper_backup_daemon ./tests
docker-compose down
```

See [demo/README.md](demo/README.md) for the local environment
(services, environment variables, the web terminal).
