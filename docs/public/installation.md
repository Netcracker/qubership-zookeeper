The following topics are covered in this chapter:

# Prerequisites

## Common

Before you start the installation and configuration of a ZooKeeper cluster, ensure the following requirements are met:

* Kubernetes 1.19+ or OpenShift 4.9+
* `kubectl` 1.19+ or `oc` 4.9+ CLI
* Helm 3.0+
* All required CRDs are installed

### Custom Resource Definitions

The following Custom Resource Definitions should be installed to the cloud before the installation of ZooKeeper:

* `ZooKeeperService` 
* `GrafanaDashboard`, `PrometheusRule`, and `ServiceMonitor` - It should be installed when you install ZooKeeper monitoring with `monitoring.monitoringType=prometheus`.
You need to install the Monitoring Operator service before the ZooKeeper installation.

**Important**: To create CRDs, you must have cloud rights for `CustomResourceDefinitions`.

### Deployment Permissions

To avoid using `cluster-wide` rights during the deployment, the following conditions are required:

* The cloud administrator creates the namespace/project in advance.
* The following grants should be provided for the `Role` of deployment user:

    <details>
    <summary>Click to expand YAML</summary>

    ```yaml
    rules:
      - apiGroups:
          - qubership.org
        resources:
          - "*"
        verbs:
          - create
          - get
          - list
          - patch
          - update
          - watch
          - delete
      - apiGroups:
          - ""
        resources:
          - pods
          - services
          - endpoints
          - persistentvolumeclaims
          - configmaps
          - secrets
          - pods/exec
          - pods/portforward
          - pods/attach
          - serviceaccounts
        verbs:
          - create
          - get
          - list
          - patch
          - update
          - watch
          - delete
      - apiGroups:
          - apps
        resources:
          - deployments
          - deployments/scale
          - deployments/status
        verbs:
          - create
          - get
          - list
          - patch
          - update
          - watch
          - delete
          - deletecollection
      - apiGroups:
          - batch
        resources:
          - jobs
        verbs:
          - create
          - get
          - list
          - patch
          - update
          - watch
          - delete
      - apiGroups:
          - ""
        resources:
          - events
        verbs:
          - create
      - apiGroups:
          - apps
        resources:
          - statefulsets
        verbs:
          - create
          - delete
          - get
          - list
          - patch
          - update
      - apiGroups:
          - networking.k8s.io
        resources:
          - ingresses
        verbs:
          - create
          - delete
          - get
          - list
          - patch
          - update
      - apiGroups:
          - rbac.authorization.k8s.io
        resources:
          - roles
          - rolebindings
        verbs:
          - create
          - delete
          - get
          - list
          - patch
          - update
      - apiGroups:
          - integreatly.org
        resources:
          - grafanadashboards
        verbs:
          - create
          - delete
          - get
          - list
          - patch
          - update
      - apiGroups:
          - monitoring.coreos.com
        resources:
          - servicemonitors
          - prometheusrules
        verbs:
          - create
          - delete
          - get
          - list
          - patch
          - update
    ```

    </details>

The following rules require `cluster-wide` permissions. If it is not possible to provide them to the deployment user, you have to apply the resources manually.

* To avoid applying manual CRD, the following grants should be provided for `ClusterRole` of the deployment user:

  ```yaml
    rules:
      - apiGroups: ["apiextensions.k8s.io"]
        resources: ["customresourcedefinitions"]
        verbs: ["get", "create", "patch"]
    ```

* Custom resource definition `ZooKeeperService` should be created/applied before the installation if the corresponding
  rights cannot be provided to the deployment user.
<!-- #GFCFilterMarkerStart# -->  
The CRD for this version is stored in [crd.yaml](../../operator/charts/helm/zookeeper-service/crds/crd.yaml) and can be
applied with the following command:

  ```sh
  kubectl replace -f crd.yaml
  ```
<!-- #GFCFilterMarkerEnd# -->

## Kubernetes

* It is required to upgrade the component before upgrading Kubernetes. Follow the information in tags regarding Kubernetes
  certified versions.

## OpenShift

* It is required to upgrade the component before upgrading OpenShift. Follow the information in tags regarding Openshift
  certified versions.
* The following annotations should be specified for the project:

  ```sh
  oc annotate --overwrite ns ${OS_PROJECT} openshift.io/sa.scc.supplemental-groups='1000/1000'
  oc annotate --overwrite ns ${OS_PROJECT} openshift.io/sa.scc.uid-range='1000/1000'
  ```

## Google Cloud

The `Google Storage` bucket is created if a backup is needed.

## AWS

The `AWS S3` bucket is created if a backup is needed.

# Best Practices and Recommendations

## HWE

The provided values do not guarantee that these values are correct for all cases. It is a general recommendation.
Resources should be calculated and estimated for each project case with test load on the SVT stand, especially the HDD size.

### Small

Recommended for development purposes, PoC, and demos.

| Module                  | CPU   | RAM, Gi | Storage, Gb |
|-------------------------|-------|---------|-------------|
| ZooKeeper (x3)          | 0.5   | 1       | 10          |
| ZooKeeper Monitoring    | 0.1   | 0.2     | 0           |
| ZooKeeper Backup Daemon | 0.1   | 0.2     | 10          |
| ZooKeeper Operator      | 0.1   | 0.2     | 0           |
| **Total (Rounded)**     | **2** | **4**   | **40**      |

<details>
<summary>Click to expand YAML</summary>

```yaml
operator:
  resources:
    requests:
      cpu: 50m
      memory: 128Mi
    limits:
      cpu: 100m
      memory: 256Mi
zooKeeper:
  heapSize: 512
  resources:
    requests:
      cpu: 200m
      memory: 1Gi
    limits:
      cpu: 500m
      memory: 1Gi
monitoring:
  resources:
    requests:
      cpu: 50m
      memory: 128Mi
    limits:
      cpu: 100m
      memory: 256Mi
backupDaemon:
resources:
  requests:
    cpu: 50m
    memory: 128Mi
  limits:
    cpu: 100m
    memory: 256Mi
```

</details>

### Medium

Recommended for deployments with average load.

| Module                  | CPU   | RAM, Gi | Storage, Gb |
|-------------------------|-------|---------|-------------|
| ZooKeeper (x3)          | 1     | 2       | 30          |
| ZooKeeper Monitoring    | 0.2   | 0.2     | 0           |
| ZooKeeper Backup Daemon | 0.1   | 0.2     | 30          |
| ZooKeeper Operator      | 0.1   | 0.2     | 0           |
| **Total (Rounded)**     | **4** | **7**   | **120**     |

<details>
<summary>Click to expand YAML</summary>

```yaml
operator:
  resources:
    requests:
      cpu: 50m
      memory: 128Mi
    limits:
      cpu: 100m
      memory: 256Mi
zooKeeper:
  heapSize: 1024
  resources:
    requests:
      cpu: 300m
      memory: 2Gi
    limits:
      cpu: 1
      memory: 2Gi
monitoring:
  resources:
    requests:
      cpu: 50m
      memory: 128Mi
    limits:
      cpu: 200m
      memory: 256Mi
backupDaemon:
resources:
  requests:
    cpu: 50m
    memory: 128Mi
  limits:
    cpu: 100m
    memory: 256Mi
```

</details>

### Large

Recommended for deployments with high workload and large amount of data.

| Module                  | CPU     | RAM, Gi | Storage, Gb |
|-------------------------|---------|---------|-------------|
| ZooKeeper (x3)          | 2       | 4       | 50          |
| ZooKeeper Monitoring    | 0.2     | 0.2     | 0           |
| ZooKeeper Backup Daemon | 0.1     | 0.2     | 50          |
| ZooKeeper Operator      | 0.1     | 0.2     | 0           |
| **Total (Rounded)**     | **6.5** | **13**  | **200**     |

<details>
<summary>Click to expand YAML</summary>

```yaml
operator:
  resources:
    requests:
      cpu: 50m
      memory: 128Mi
    limits:
      cpu: 100m
      memory: 256Mi
zooKeeper:
  heapSize: 2048
  resources:
    requests:
      cpu: 500m
      memory: 4Gi
    limits:
      cpu: 2
      memory: 4Gi
monitoring:
  resources:
    requests:
      cpu: 50m
      memory: 128Mi
    limits:
      cpu: 200m
      memory: 256Mi
backupDaemon:
  resources:
    requests:
      cpu: 50m
      memory: 128Mi
    limits:
      cpu: 100m
      memory: 256Mi
```

</details>

# Parameters

## Cloud Integration Parameters

| Parameter                   | Type    | Mandatory | Default value         | Description                                                                                                                                                                                                                                                                                                                                         |
|-----------------------------|---------|-----------|-----------------------|-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| SERVICE_NAME                | string  | yes       | `"zookeeper-service"` | The name of the Service exposed for the database.                                                                                                                                                                                                                                                                                                   |
| PART_OF                     | string  | yes       | `"zookeeper-service"` | Microservice belonging to a group.                                                                                                                                                                                                                                                                                                                  |
| DELIMITER                   | string  | yes       | `"-"`                 | Delimiter for labels.                                                                                                                                                                                                                                                                                                                               |
| ARTIFACT_DESCRIPTOR_VERSION | string  | yes       | `""`                  | Artifact descriptor version which is installed.                                                                                                                                                                                                                                                                                                     |
| ZOOKEEPER_ADMIN_PASSWORD    | string  | yes       | `""`                  | The password of the ZooKeeper administrator user. ZooKeeper nodes use these credentials to communicate. This parameter enables ZooKeeper authentication. The username for this user is `ZOOKEEPER_ADMIN_USERNAME`.                                                                                                                                  |
| ZOOKEEPER_CLIENT_USERNAME   | string  | yes       | `""`                  | The username of the ZooKeeper client user. These credentials are used by the ZooKeeper client to establish a connection with the ZooKeeper server if SASL authentication is enabled, where the `ZOOKEEPER_ADMIN_USERNAME` and `ZOOKEEPER_ADMIN_PASSWORD` parameter values are specified. The password for this user is `ZOOKEEPER_CLIENT_PASSWORD`. |
| ZOOKEEPER_CLIENT_PASSWORD   | string  | yes       | `""`                  | The password of the ZooKeeper client user. ZooKeeper client uses these credentials to establish a connection with the ZooKeeper server.                                                                                                                                                                                                             |
| INFRA_ZOOKEEPER_REPLICAS    | integer | no        | `3`                   | The number of ZooKeeper servers.                                                                                                                                                                                                                                                                                                                    |
| MONITORING_ENABLED          | boolean | no        | `false`               | Specifies whether ZooKeeper Monitoring component is to be deployed or not.                                                                                                                                                                                                                                                                          |
| STORAGE_RWO_CLASS           | string  | yes       | `""`                  | This parameter specifies the ZooKeeper storage class.                                                                                                                                                                                                                                                                                               |

## Global

| Parameter                                  | Type    | Mandatory | Default value | Description                                                                                                                                                                                                                                                                                                                                                                                            |
|--------------------------------------------|---------|-----------|---------------|--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| global.name                                | string  | no        | zookeeper     | The custom resource name that is used to form service names for ZooKeeper, ZooKeeper monitoring, and ZooKeeper backup daemon. <br/>**Important**: If you modify this parameter, you always need to add the `CUSTOM_RESOURCE_NAME` parameter with the same value when deploying.                                                                                                                        |
| global.waitForPodsReady                    | boolean | no        | true          | Whether the operator should wait for the pods to be ready in order to publish the status to the Custom Resource.                                                                                                                                                                                                                                                                                       |
| global.podReadinessTimeout                 | integer | no        | 300           | The timeout in seconds for how long the operator should wait for the pods to be ready for each service.                                                                                                                                                                                                                                                                                                |
| global.customLabels                        | object  | no        | {}            | The custom labels for all pods that are related to the ZooKeeper Service. These labels can be overridden by the component `customLabel` parameter.                                                                                                                                                                                                                                                     |
| global.securityContext                     | object  | no        | {}            | The pod security context for all pods which are related to the ZooKeeper Service. The security context can be overridden by the component `securityContext` parameter.                                                                                                                                                                                                                                 |
| global.tls.enabled                         | boolean | no        | false         | Whether to use TLS to connect to the appropriate components, such as ZooKeeper and ZooKeeper Backup Daemon.                                                                                                                                                                                                                                                                                            |
| global.tls.cipherSuites                    | list    | no        | []            | The list of cipher suites that are used to negotiate the security settings for a network connection using TLS or SSL network protocol. By default, all the available cipher suites are supported.                                                                                                                                                                                                      |
| global.tls.allowNonencryptedAccess         | boolean | no        | false         | Whether to allow non-encrypted access to ZooKeeper by port `2182` or not.                                                                                                                                                                                                                                                                                                                              |
| global.tls.generateCerts.enabled           | boolean | no        | true          | Whether to generate TLS certificates by Helm or not.                                                                                                                                                                                                                                                                                                                                                   |
| global.tls.generateCerts.certProvider      | string  | no        | cert-manager  | The provider used to generate TLS certificates. The possible values are `helm` and `cert-manager`.                                                                                                                                                                                                                                                                                                     |
| global.tls.generateCerts.durationDays      | integer | no        | 365           | The TLS certificate validity duration in days.                                                                                                                                                                                                                                                                                                                                                         |
| global.tls.generateCerts.clusterIssuerName | string  | no        | ""            | The name of the `ClusterIssuer` resource. If the parameter is not set or empty, the `Issuer` resource in the current Kubernetes namespace is used. It is used when the `global.tls.generateCerts.certProvider` parameter is set to `cert-manager`.                                                                                                                                                     |
| global.cloudIntegrationEnabled             | boolean | no        | true          | The parameter specifies whether to apply global cloud parameters instead of parameters described in ZooKeeper service (`ZOOKEEPER_ADMIN_USERNAME`, `ZOOKEEPER_ADMIN_PASSWORD`, `ZOOKEEPER_CLIENT_USERNAME`, `ZOOKEEPER_CLIENT_PASSWORD`, `MONITORING_ENABLED`, `STORAGE_RWO_CLASS`). If it is set to `false` or global parameter is absent, corresponding parameter from ZooKeeper service is applied. |

### Secrets

| Parameter                                           | Type   | Mandatory | Default value | Description                                                                                                                                                                                                                                                                                                                                                                                                                                                                                 |
|-----------------------------------------------------|--------|-----------|---------------|---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| global.secrets.zooKeeper.adminUsername              | string | no        | ""            | The username of the ZooKeeper administrator. ZooKeeper nodes use these credentials to communicate. This parameter enables ZooKeeper authentication. The password for this user is `global.secrets.zooKeeper.adminPassword`. If the parameter value is empty, ZooKeeper deploys with a disabled SASL.                                                                                                                                                                                        |
| global.secrets.zooKeeper.adminPassword              | string | no        | ""            | The password of the ZooKeeper administrator user. ZooKeeper nodes use these credentials to communicate. This parameter enables ZooKeeper authentication. The username for this user is `global.secrets.zooKeeper.adminUsername`. If the parameter value is empty, ZooKeeper deploys with a disabled SASL. <br/>**Note**: The value of this parameter is not used if the `vaultSecretManagement.enabled` parameter value is `true`. In this case, the passwords are generated automatically. |
| global.secrets.zooKeeper.clientUsername             | string | no        | ""            | The username of the ZooKeeper client user. These credentials are used by the ZooKeeper client to establish a connection with the ZooKeeper server if SASL authentication is enabled, where the `global.secrets.zooKeeper.adminUsername` and `global.secrets.zooKeeper.adminPassword` parameter values are specified. The password for this user is `global.secrets.zooKeeper.clientPassword`. The parameter can be empty if authentication is disabled.                                     |
| global.secrets.zooKeeper.clientPassword             | string | no        | ""            | The password of the ZooKeeper client user. ZooKeeper client uses these credentials to establish a connection with the ZooKeeper server. The parameter can be empty if authentication is disabled. <br/>**Note**: The value of this parameter is not used if the `vaultSecretManagement.enabled` parameter value is `true`. In this case, the passwords are generated automatically.                                                                                                         |
| global.secrets.zooKeeper.additionalUsers            | string | no        | ""            | The comma-separated pairs, `username:password`, of additional users that are used by clients for authentication in ZooKeeper. The parameter can be empty. For example, `user1_name:user1_password,user2_name:user2_password`. **Note**: If the `vaultSecretManagement.enabled` parameter value is `true`, the value of this parameter should include only usernames. In this case, the passwords are generated automatically. For example, `user1_name,user2_name`.                         |
| global.secrets.backupDaemon.username                | string | no        | ""            | The username of the ZooKeeper Backup Daemon API user. This parameter enables the ZooKeeper Backup Daemon authentication. The password for this user is `global.secrets.backupDaemon.password`. If the parameter value is empty, ZooKeeper Backup Daemon deploys with disabled authentication.                                                                                                                                                                                               |
| global.secrets.backupDaemon.password                | string | no        | ""            | The password of the ZooKeeper Backup Daemon API user. This parameter enables the ZooKeeper Backup Daemon authentication. If the parameter value is empty, ZooKeeper Backup Daemon deploys with disabled authentication. <br/>**Note**: The value of this parameter is not used if the `vaultSecretManagement.enabled` parameter value is `true`. In this case, the passwords are generated automatically.                                                                                   |
| global.secrets.integrationTests.prometheus.user     | string | no        | ""            | The username for authentication on Prometheus/VictoriaMetrics secured endpoints.                                                                                                                                                                                                                                                                                                                                                                                                            |
| global.secrets.integrationTests.prometheus.password | string | no        | ""            | The password for authentication on Prometheus/VictoriaMetrics secured endpoints.                                                                                                                                                                                                                                                                                                                                                                                                            |

**Note**: The username and passwords can contain only the following symbols:

* Alphabets: `a-zA-Z`
* Numerals: `0-9`
* Punctuation marks: `.`, `;`, `!`, `?`
* Mathematical symbols: `-`, `+`, `*`, `/`, `%`
* Brackets: `(`, `)`, `{`, `}`, `<`, `>`
* Additional symbols: `_`, `|`, `&`, `@`, `$`, `^`, `#`, `~`

## Operator

| Parameter                          | Type     | Mandatory | Default value            | Description                                                                                                                                                                                                                                                                                                                       |
|------------------------------------|----------|-----------|--------------------------|-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| operator.dockerImage               | string   | no        | Calculates automatically | The image of ZooKeeper Service Operator.                                                                                                                                                                                                                                                                                          |
| operator.affinity                  | object   | no        | {}                       | The affinity scheduling rules in `json` format.                                                                                                                                                                                                                                                                                   |
| operator.tolerations               | list     | no        | []                       | The list of toleration policies for ZooKeeper service operator pod in `json` format.                                                                                                                                                                                                                                              |
| operator.priorityClassName         | string   | no        | ""                       | The priority class to be used by the ZooKeeper Service Operator pod. You should create the priority class beforehand. For more information about this feature, refer to [https://kubernetes.io/docs/concepts/configuration/pod-priority-preemption/](https://kubernetes.io/docs/concepts/configuration/pod-priority-preemption/). |
| operator.customLabels              | object   | no        | {}                       | The custom labels for the ZooKeeper Service operator pod in `json` format.                                                                                                                                                                                                                                                        |
| operator.securityContext           | object   | no        | {}                       | The pod security context for the ZooKeeper Service operator pod.                                                                                                                                                                                                                                                                  |
| operator.resources.limits.cpu      | string   | no        | 100m                     | This parameter specifies the ZooKeeper operator CPU limits.                                                                                                                                                                                                                                                                       |
| operator.resources.limits.memory   | string   | no        | 256Mi                    | This parameter specifies the ZooKeeper operator CPU limits.                                                                                                                                                                                                                                                                       |
| operator.resources.requests.cpu    | string   | no        | 50m                      | This parameter specifies the ZooKeeper operator CPU requests.                                                                                                                                                                                                                                                                     |
| operator.resources.requests.memory | string   | no        | 128Mi                    | This parameter specifies the ZooKeeper operator memory requests.                                                                                                                                                                                                                                                                  |

## ZooKeeper

| Parameter                                                  | Type    | Mandatory | Default value                                                                       | Description                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                           |
|------------------------------------------------------------|---------|-----------|-------------------------------------------------------------------------------------|-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| zooKeeper.dockerImage                                      | string  | no        | Calculates automatically                                                            | The Docker image of ZooKeeper.                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                        |
| zooKeeper.affinity                                         | object  | no        | `{}`                                                                                | The affinity scheduling rules for ZooKeeper pods.                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                     |
| zooKeeper.tolerations                                      | object  | no        | `{}`                                                                                | The list of toleration policies for ZooKeeper pods.                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                   |
| zooKeeper.priorityClassName                                | string  | no        | `""`                                                                                | The priority class to be used by ZooKeeper pods.                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                      |
| zooKeeper.disruptionBudget.enabled                         | boolean | no        | false                                                                               | Whether to create PodDisruptionBudget to prevent voluntary degradation of the ZooKeeper server cluster.                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                               |
| zooKeeper.disruptionBudget.minAvailable                    | integer | no        | `2`                                                                                 | The minimal number of pods that must still be available after the eviction. Calculated as `(n/2)+1`, where `n` is the number of replicas (or 0, if the number of replicas is 1).                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                      |
| zooKeeper.replicas                                         | integer | no        | `3`                                                                                 | The number of ZooKeeper servers.                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                      |
| zooKeeper.storage.volumes                                  | list    | no        | `[]`                                                                                | The list of persistent volume names that are used to bind with the persistent volume claims. The number of persistent volume names must be equal to the value of the `zooKeeper.replicas` parameter.                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                  |
| zooKeeper.storage.nodes                                    | list    | no        | `[]`                                                                                | The list of node names that is used to schedule on which nodes the pods run. This parameter is mandatory if ZooKeeper uses storage.                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                   |
| zooKeeper.storage.labels                                   | list    | no        | `[]`                                                                                | The list of labels that is used to bind suitable persistent volumes with the persistent volume claims. The number of labels must be equal to the value of the `zooKeeper.replicas` parameter, one label per persistent volume in `key=value` format. You must specify this parameter only for the label selector volume binding.                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                      |
| zooKeeper.storage.className                                | list    | yes       | `[]`                                                                                | The list of storage class names used to dynamically provide volumes. The number of storage classes should be equal to `1` if one storage class is used for all persistent volumes, or the value of the `zooKeeper.replicas` parameter if persistent volumes use different storage classes. If this parameter is empty (set to `""`), the persistent volumes without storage class are bound with the persistent volume claims. You must specify this parameter only for the dynamic volume provisioning and for the label selector volume binding.                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                    |
| zooKeeper.storage.size                                     | string  | yes       | `2Gi`                                                                               | The size of the persistent volume in Gi.                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                              |
| zooKeeper.snapshotStorage.persistentVolumeType             | string  | no        | `""`                                                                                | The type of persistent volume for snapshots. If this parameter is empty, the persistent volume and the persistent volume claim for snapshots are not created or updated. There are three possible values available: <br><br>* `predefined` uses the already prepared shared persistent volume for snapshots. You can specify the name of the prepared persistent volume in the `zooKeeper.snapshotStorage.persistentVolumeName` parameter. You can specify the name of the persistent volume claim that is created at the time of the installation in the `zooKeeper.snapshotStorage.persistentVolumeClaimName` parameter. If the prepared persistent volume is created by dynamic volume provisioning, you can specify the storage class in the `zooKeeper.snapshotStorage.storageClass` parameter.<br>* `predefined_claim` uses the already prepared shared persistent volume claim for snapshots. You can specify the name of the prepared persistent volume claim in the `zooKeeper.snapshotStorage.persistentVolumeClaimName` parameter.<br>* `storage_class` uses dynamically provided shared volumes. You can specify the name of the storage class in the `zooKeeper.snapshotStorage.storageClass` parameter. |
| zooKeeper.snapshotStorage.persistentVolumeName             | string  | no        | `""`                                                                                | Specifies the snapshots' persistent volume name that is used to bind with the snapshots' persistent volume claim. You must specify this parameter for the `predefined` persistent volume type.                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                        |
| zooKeeper.snapshotStorage.persistentVolumeClaimName        | string  | no        | `pvc-<name>-snapshots`, where `<name>` is the value of the `global.name` parameter. | Specifies the name of the snapshots' persistent volume claim. If the parameter is empty, `pvc-<name>-snapshots`, the default name of the persistent volume claim is used, where `<name>` is the value of the `global.name` parameter.                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                 |
| zooKeeper.snapshotStorage.volumeSize                       | string  | yes       | `1Gi`                                                                               | Specifies the size of the persistent volume for snapshots in Gi.                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                      |
| zooKeeper.snapshotStorage.storageClass                     | string  | yes       | `""`                                                                                | Specifies the name of the storage class used to dynamically provide volumes. If the `zooKeeper.snapshotStorage.storageClass` value is empty (set to `""`), the persistent volumes without storage class are bound with the persistent volume claims. If the `zooKeeper.snapshotStorage.storageClass` value is not specified, the persistent volumes with default storage class are bound with the persistent volume claims. This parameter is used for the `predefined` and `storage_class` persistent volume types.                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                  |
| zooKeeper.heapSize                                         | integer | no        | `256`                                                                               | Specifies the heap size of JVM in Mi.                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                 |
| zooKeeper.jolokiaPort                                      | integer | no        | `9087`                                                                              | Specifies the jolokia agent port. This agent defines the JMX-HTTP bridge. Set this parameter when you want to use monitoring for ZooKeeper. The default value is "9087".                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                              |
| zooKeeper.resources.requests.cpu                           | string  | no        | `50m`                                                                               | Specifies the minimum number of CPUs the container should use.                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                        |
| zooKeeper.resources.requests.memory                        | string  | no        | `512Mi`                                                                             | Specifies the minimum amount of memory the container should use. The value can be specified with SI suffixes (E, P, T, G, M, K, m) or their power-of-two-equivalents (Ei, Pi, Ti, Gi, Mi, Ki).                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                        |
| zooKeeper.resources.limits.cpu                             | string  | no        | `300m`                                                                              | Specifies the maximum number of CPUs the container can use.                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                           |
| zooKeeper.resources.limits.memory                          | string  | no        | `512Mi`                                                                             | Specifies the maximum amount of memory the container can use. The value can be specified with SI suffixes (E, P, T, G, M, K, m) or their power-of-two-equivalents (Ei, Pi, Ti, Gi, Mi, Ki).                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                           |
| zooKeeper.quorumAuthEnabled                                | boolean | no        | true                                                                                | Enables internal authentication between ZooKeeper nodes if SASL authentication is enabled, where the `global.secrets.zooKeeper.adminUsername` and `global.secrets.zooKeeper.adminPassword` parameter values are specified. The parameter can be empty if authentication is disabled.                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                  |
| zooKeeper.tls.enabled                                      | boolean | no        | true                                                                                | Specifies whether to use TLS to connect ZooKeeper. If the `zooKeeper.tls.enabled` value is "true", it requires the `global.tls.enabled` value to be set to "true" to enable TLS for ZooKeeper connection.                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                             |
| zooKeeper.tls.certificates.crt                             | string  | no        | ""                                                                                  | The certificate in BASE64 format. It is required if `global.tls.enabled` parameter is set to `true`, `global.tls.generateCerts.certProvider` parameter is set to `helm` and `global.tls.generateCerts.enabled` parameter is set to `false`.                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                           |
| zooKeeper.tls.certificates.key                             | string  | no        | ""                                                                                  | The private key in BASE64 format. It is required if `global.tls.enabled` parameter is set to `true`, `global.tls.generateCerts.certProvider` parameter is set to `helm` and `global.tls.generateCerts.enabled` parameter is set to `false`.                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                           |
| zooKeeper.tls.certificates.ca                              | string  | no        | ""                                                                                  | The root CA certificate in BASE64 format. It is required if `global.tls.enabled` parameter is set to `true`, `global.tls.generateCerts.certProvider` parameter is set to `helm` and `global.tls.generateCerts.enabled` parameter is set to `false`.                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                   |
| zooKeeper.tls.secretName                                   | string  | no        | `{name}-tls-secret`                                                                 | Specifies the secret that contains TLS certificates. If the `global.tls.generateCerts.enabled` parameter is set to "true", the default value is `{name}-tls-secret`, where `{name}` is the value of the `global.name` parameter.                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                      |
| zooKeeper.tls.cipherSuites                                 | list    | no        | `[]`                                                                                | Specifies the list of cipher suites that are used to negotiate the security settings for a network connection using TLS or SSL network protocol. By default, all the available cipher suites are supported.                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                           |
| zooKeeper.tls.mTLS                                         | boolean | no        | false                                                                               | Specifies whether to enable two-way TLS authentication or not.                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                        |
| zooKeeper.tls.subjectAlternativeName.additionalDnsNames    | list    | no        | `[]`                                                                                | Specifies the list of additional DNS names to be added to the **Subject Alternative Name** field of a TLS certificate.                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                |
| zooKeeper.tls.subjectAlternativeName.additionalIpAddresses | list    | no        | `[]`                                                                                | Specifies the list of additional IP addresses to be added to the **Subject Alternative Name** field of a TLS certificate.                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                             |
| zooKeeper.securityContext                                  | object  | no        | `{}`                                                                                | Specifies the pod-level security attributes and common container settings. The parameter value can be empty and should be specified in the `json` format. For example, you can add `{"fsGroup": 1000}`.                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                               |
| zooKeeper.auditEnabled                                     | boolean | no        | false                                                                               | Specifies whether to enable audit logging for ZooKeeper.                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                              |
| zooKeeper.environmentVariables                             | list    | no        | `[]`                                                                                | Specifies the list of additional environment variables for ZooKeeper deployments in `key=value` format. The parameter value can be empty.                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                             |
| zooKeeper.rollingUpdate                                    | boolean | no        | false                                                                               | Specifies either to redeploy ZooKeeper pods during an update one by one or all in the same time. If "true" is specified after every ZooKeeper server update, the status of all servers is checked.                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                    |
| zooKeeper.customLabels                                     | object  | no        | `{}`                                                                                | The custom labels for all ZooKeeper pods.                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                             |
| zooKeeper.diagnostics.mode                                 | string  | no        | `disable`                                                                           | The parameter specifies mode of Cloud Diagnostic Toolset. Allowed values are `disable`/`dev`/`prod`:<br>* `disable` - to disable CDT integration.<br>* `dev`/`prod` - to enable CDT integration. **Note**: The production mode does not store to disk java calls that lasted less than 1ms.                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                           |
| zooKeeper.diagnostics.agentService                         | string  | no        | `nc-diagnostic-agent`                                                               | The parameter specifies the location to Cloud Diagnostic Toolset (host to which will send data).                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                      |

## Monitoring

| Parameter                            | Type    | Mandatory | Default value            | Description                                                                                                                                                                                     |
|--------------------------------------|---------|-----------|--------------------------|-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| monitoring.install                   | boolean | no        | true                     | Specifies whether ZooKeeper Monitoring component is to be deployed or not.                                                                                                                      |
| monitoring.dockerImage               | string  | no        | Calculates automatically | The Docker image of ZooKeeper Monitoring.                                                                                                                                                       |
| monitoring.affinity                  | object  | no        | {}                       | The affinity scheduling rules.                                                                                                                                                                  |
| monitoring.tolerations               | object  | no        | {}                       | The list of toleration policies for ZooKeeper Monitoring pod.                                                                                                                                   |
| monitoring.priorityClassName         | string  | no        | ""                       | The priority class to be used by a ZooKeeper Monitoring pod. You should create the priority class beforehand.                                                                                   |
| monitoring.resources.requests.cpu    | string  | no        | 25m                      | The minimum number of CPUs the container should use.                                                                                                                                            |
| monitoring.resources.requests.memory | string  | no        | 128Mi                    | The minimum amount of memory the container should use.                                                                                                                                          |
| monitoring.resources.limits.cpu      | string  | no        | 200m                     | The maximum number of CPUs the container can use.                                                                                                                                               |
| monitoring.resources.limits.memory   | string  | no        | 256Mi                    | The maximum amount of memory the container can use.                                                                                                                                             |
| monitoring.customLabels              | object  | no        | {}                       | The custom labels for a ZooKeeper Service monitoring pod.                                                                                                                                       |
| monitoring.monitoringType            | string  | no        | prometheus               | The type of output plugin that is used for service monitoring. Currently only Prometheus type is supported. The Prometheus plugin does not require additional parameters for the configuration. |
| monitoring.installGrafanaDashboard   | boolean | no        | true                     | Whether the ZooKeeper Grafana dashboard is to be applied or not.                                                                                                                                |
| monitoring.zooKeeperVolumes          | string  | no        | ""                       | The persistent volume name or a comma-separated list of names that are used by ZooKeeper Service.                                                                                               |
| monitoring.securityContext           | object  | yes       | {}                       | The pod-level security attributes and common container settings. The parameter value can be empty and should be specified in the `json` format. For example, you can add `{"runAsUser": 1000}`. |

## Backup Daemon

| Parameter                                                     | Type    | Mandatory | Default value            | Description                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                          |
|---------------------------------------------------------------|---------|-----------|--------------------------|--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| backupDaemon.install                                          | boolean | no        | false                    | Whether the ZooKeeper Backup Daemon component is to be deployed or not. The value should be equal to "true" to install ZooKeeper Backup Daemon.                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                      |
| backupDaemon.dockerImage                                      | string  | no        | Calculates automatically | The Docker image of ZooKeeper Backup Daemon.                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                         |
| backupDaemon.affinity                                         | object  | no        | {}                       | The affinity scheduling rules.                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                       |
| backupDaemon.tolerations                                      | object  | no        | {}                       | The list of toleration policies for ZooKeeper Backup Daemon pod.                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                     |
| backupDaemon.priorityClassName                                | string  | no        | ""                       | The priority class to be used by a ZooKeeper Backup Daemon pod. You should create the priority class beforehand. For more information about this feature, refer to [https://kubernetes.io/docs/concepts/configuration/pod-priority-preemption/](https://kubernetes.io/docs/concepts/configuration/pod-priority-preemption/)                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                          |
| backupDaemon.backupStorage.persistentVolumeType               | string  | no        | ""                       | The type of persistent volume for snapshots. If this parameter is empty, the persistent volume and the persistent volume claim for snapshots are not created or updated. There are four possible values available: <br> <br> * `standalone` uses the non-shared persistent volume. You can specify the name of the prepared persistent volume in the `backupDaemon.backupStorage.persistentVolumeName` parameter. You can specify the name of the persistent volume claim that is created at the time of installation in the `backupDaemon.backupStorage.persistentVolumeClaimName` parameter. For dynamic volume provisioning, you can specify the storage class in the `backupDaemon.backupStorage.storageClass` parameter. For label selector volume binding, you can specify the label in the `backupDaemon.backupStorage.persistentVolumeLabel` parameter and storage class in the `backupDaemon.backupStorage.storageClass` parameter. <br> <br> * `predefined` uses the already prepared shared persistent volume. You can specify the name of the prepared persistent volume in the `backupDaemon.backupStorage.persistentVolumeName` parameter. You can specify the name of the persistent volume claim that creates at the time of installation in the `backupDaemon.backupStorage.persistentVolumeClaimName` parameter. If prepared persistent volume is created by dynamic volume provisioning, you can specify the storage class in the `backupDaemon.backupStorage.storageClass` parameter. It must be the same persistent volume that is used for the ZooKeeper service to store snapshots. <br> <br> * `predefined_claim` uses the already prepared shared persistent volume claim. You can specify the name of the prepared persistent volume claim in the `backupDaemon.backupStorage.persistentVolumeClaimName` parameter. It must be the same persistent volume claim that is used for the ZooKeeper service to store snapshots. <br> <br> * `storage_class` uses dynamically provided shared volumes. You can specify the name of the storage class in the `backupDaemon.backupStorage.storageClass` parameter. |
| backupDaemon.backupStorage.persistentVolumeName               | string  | no        | ""                       | The persistent volume name that is used to bind with the persistent volume claim. You must specify this parameter for the standalone and predefined persistent volume types.                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                         |
| backupDaemon.backupStorage.persistentVolumeClaimName          | string  | no        | ""                       | The name of the persistent volume claim. If the parameter is empty, the default name of the persistent volume claim (`pvc-<name>-snapshots`) is used, where `<name>` is the value of the `global.name` parameter.                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                    |
| backupDaemon.backupStorage.persistentVolumeLabel              | string  | no        | ""                       | The label that is used to bind a suitable persistent volume with a persistent volume claim. You must specify this parameter only for the label selector volume binding. For example, `disk-id=3befc1c1`.                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                             |
| backupDaemon.backupStorage.storageClass                       | string  | yes       | ""                       | The name of the storage class used to dynamically provide volumes. If this parameter is empty (set to `""`), the persistent volume without storage class is bound with the persistent volume claim. This parameter can be specified for `standalone`, `predefined`, and `storage_class` `backupDaemon.backupStorage.persistentVolumeType`.                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                           |
| backupDaemon.backupStorage.nodeName                           | string  | no        | ""                       | The node name that is used to schedule on which node the pod runs.                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                   |
| backupDaemon.backupStorage.volumeSize                         | string  | yes       | 1Gi                      | The size of the persistent volume for snapshots in Gi.                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                               |
| backupDaemon.resources.requests.cpu                           | string  | no        | 25m                      | The minimum number of CPUs the container should use.                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                 |
| backupDaemon.resources.requests.memory                        | string  | no        | 512Mi                    | The minimum amount of memory the container should use. The value can be specified with SI suffixes (E, P, T, G, M, K, m) or their power-of-two-equivalents (Ei, Pi, Ti, Gi, Mi, Ki).                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                 |
| backupDaemon.resources.limits.cpu                             | string  | no        | 300m                     | The maximum number of CPUs the container can use.                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                    |
| backupDaemon.resources.limits.memory                          | string  | no        | 512Mi                    | The maximum amount of memory the container can use. The value can be specified with SI suffixes (E, P, T, G, M, K, m) or their power-of-two-equivalents (Ei, Pi, Ti, Gi, Mi, Ki).                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                    |
| backupDaemon.backupSchedule                                   | string  | no        | 0 0 * * *                | The cron-like backup schedule. If this parameter is empty, the default schedule (`"0 * * * *"`), defined in the ZooKeeper Backup Daemon configuration is used. The value `0 * * * *` means that snapshots are created every hour at the beginning of the hour.                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                       |
| backupDaemon.evictionPolicy                                   | string  | no        | 0/1d,7d/delete           | The backup eviction policy. It is a comma-separated string of policies written as `$start_time/$interval`. This policy splits all backups older than `$start_time` to numerous time intervals `$interval` time long. Then it deletes all backups in every interval, except the newest one. For example, `1d/7d` policy means "take all backups older then one day, split them in groups by a 7-day interval, and leave only the newest." If this parameter is empty, the default eviction policy (`"0/1d,7d/delete"`) defined in the ZooKeeper Backup Daemon configuration is used.                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                  |
| backupDaemon.ipv6                                             | boolean | no        | false                    | If ZooKeeper Backup Daemon REST API should be started on an IPv6 interface. If the service is deployed in an environment with IPv6 network interfaces, set this parameter value to "true".                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                           |
| backupDaemon.securityContext                                  | object  | no        | {}                       | The pod-level security attributes and common container settings. The parameter value can be empty and should be specified in the `json` format. For example, you can add `{"fsGroup": 1000}`.                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                        |
| backupDaemon.customLabels                                     | object  | no        | {}                       | The parameter allows specifying custom labels for the ZooKeeper Service Backup daemon pod.                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                           |
| backupDaemon.tls.enabled                                      | boolean | no        | true                     | Whether to use TLS to connect ZooKeeper Backup Daemon. If `backupDaemon.tls.enabled` value is "true", it requires `global.tls.enabled` to be set to "true" to enable TLS for the ZooKeeper connection.                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                               |
| backupDaemon.tls.certificates.crt                             | string  | no        | ""                       | The certificate in BASE64 format. It is required if `global.tls.enabled` parameter is set to `true`, `global.tls.generateCerts.certProvider` parameter is set to `helm` and `global.tls.generateCerts.enabled` parameter is set to `false`.                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                          |
| backupDaemon.tls.certificates.key                             | string  | no        | ""                       | The private key in BASE64 format. It is required if `global.tls.enabled` parameter is set to `true`, `global.tls.generateCerts.certProvider` parameter is set to `helm` and `global.tls.generateCerts.enabled` parameter is set to `false`.                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                          |
| backupDaemon.tls.certificates.ca                              | string  | no        | ""                       | The root CA certificate in BASE64 format. It is required if `global.tls.enabled` parameter is set to `true`, `global.tls.generateCerts.certProvider` parameter is set to `helm` and `global.tls.generateCerts.enabled` parameter is set to `false`.                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                  |
| backupDaemon.tls.secretName                                   | string  | no        | ""                       | The secret that contains TLS certificates. If the `global.tls.generateCerts.enabled` parameter is set to "true", the default value is `{name}-backup-daemon-tls-secret`, where `{name}` is the value of the `global.name` parameter.                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                 |
| backupDaemon.tls.cipherSuites                                 | list    | no        | []                       | The list of cipher suites that are used to negotiate the security settings for a network connection using TLS or SSL network protocol. By default, all the available cipher suites are supported.                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                    |
| backupDaemon.tls.subjectAlternativeName.additionalDnsNames    | list    | no        | []                       | The list of additional DNS names to be added to the **Subject Alternative Name** field of a TLS certificate.                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                         |
| backupDaemon.tls.subjectAlternativeName.additionalIpAddresses | list    | no        | []                       | The list of additional IP addresses to be added to the **Subject Alternative Name** field of a TLS certificate.                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                      |
| backupDaemon.s3.enabled                                       | boolean | no        | false                    | Whether to store backups to S3 storage (AWS, Google, MinIO, and so on). A clipboard storage is needed to be mounted to Backup Daemon and it can be `emptyDir` volume. As soon as backup is uploaded to S3, it is removed from the clipboard storage.<br>A restore procedure works in the same way - a backup is downloaded from S3 to the clipboard and restored from it, then it is removed from the clipboard but stays on S3. Eviction removes backups directly from S3. <br> **Note**: Only the hierarchical backup mode is supported for S3 backup storage.                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                     |
| backupDaemon.s3.url                                           | string  | no        | ""                       | The URL of the S3 storage. For example, `https://s3.amazonaws.com`.                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                  |
| backupDaemon.s3.sslVerify                                     | boolean | no        | `true`                   | This parameter specifies whether or not to verify SSL certificates for S3 connections.                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                               |
| backupDaemon.s3.sslSecretName                                 | string  | no        | `""`                     | This parameter specifies name of the secret with CA certificate for S3 connections. If secret not exists and parameter `backupDaemon.s3.sslCert` is specified secret will be created, else boto3 certificates will be used.                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                          |
| backupDaemon.s3.sslCert                                       | string  | no        | `""`                     | The root CA certificate in BASE64 format. It is required if pre-created secret with certificates not exists and default boto3 certificates will not be used.                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                         |
| backupDaemon.s3.keyId                                         | string  | no        | ""                       | The key ID for the S3 storage. A user must have access to the bucket.                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                |
| backupDaemon.s3.keySecret                                     | string  | no        | ""                       | The key secret for the S3 storage. A user must have access to the bucket.                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                            |
| backupDaemon.s3.bucket                                        | string  | no        | ""                       | The bucket in the S3 storage that is used to store backups.                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                          |
| backupDaemon.zooKeeperHost                                    | string  | no        | "zookeeper"              | The host name of ZooKeeper.                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                          |
| backupDaemon.zooKeeperPort                                    | integer | no        | 2181                     | The port of ZooKeeper.                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                               |

## Integration Tests

| Parameter                                     | Type    | Mandatory | Default value            | Description                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                |
|-----------------------------------------------|---------|-----------|--------------------------|--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| integrationTests.install                      | boolean | no        | false                    | Whether the ZooKeeper Service Integration Tests component is to be deployed or not. The value should be equal to "true" to perform ZooKeeper Service Integration Tests after the deployment. <br> **Important**: If you perform integration tests with a Deployer job, you are not able to run integration tests as a separate job in the future. It is also not possible to run integration tests with Vault Credentials Management enabled.                                                                                              |
| integrationTests.service.name                 | string  | no        | Calculates automatically | The name of the ZooKeeper integration tests service.                                                                                                                                                                                                                                                                                                                                                                                                                                                                                       |
| integrationTests.waitForResult                | boolean | no        | true                     | Whether the operator should wait for the integration tests to be completed successfully in order to publish the status to the Custom Resource. You can enable it only if the `global.waitForPodsReady` parameter value is "true". <br> **Important**: If this property is enabled, the operator waits for integration tests pod to be ready with all tests passed. If the tests are complete with unsuccessful result, the operator does not stop the checking until timeout is reached. It allows manually restarting the pod with tests. |
| integrationTests.timeout                      | integer | no        | 420                      | The timeout in seconds for how long the operator should wait for the integration tests' result.                                                                                                                                                                                                                                                                                                                                                                                                                                            |
| integrationTests.affinity                     | object  | no        | {}                       | The affinity scheduling rules for the ZooKeeper integration tests pod. The value should be specified in the `json` format. The parameter can be empty.                                                                                                                                                                                                                                                                                                                                                                                     |
| integrationTests.image                        | string  | no        | Calculates automatically | The Docker image of ZooKeeper Service.                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                     |
| integrationTests.tags                         | string  | no        | "zookeeper_crud"         | The tags in combination with AND, OR, and NOT operators that select the test cases to run. Do not specify any value for this parameter if you want to run all the tests. Information about available tags can be found in the [Integration test tags description](#integration-test-tags-description) article.                                                                                                                                                                                                                             |
| integrationTests.zookeeperVolumeSize          | integer | no        | 2                        | The size of the ZooKeeper persistent volume in Gi.                                                                                                                                                                                                                                                                                                                                                                                                                                                                                         |
| integrationTests.pvType                       | string  | no        | "nfs"                    | The type of persistent volume for ZooKeeper snapshots. The possible values are "standalone", "predefined", "predefined_claim", "nfs", or can be empty. If the parameter is empty, all integration tests regarding transactional ZooKeeper backups are skipped. The parameter should match ZooKeeper Backup Daemon `PV_TYPE` and `backupDaemon.backupStorage.persistentVolumeType` parameter.                                                                                                                                               |
| integrationTests.prometheusUrl                | string  | no        | ""                       | The URL (with schema and port) to Prometheus. For example, `http://prometheus.cloud.openshift.sdntest.example.com:80`. This parameter must be specified if you want to run integration tests with the "prometheus" tag. **Note:** This parameter could be used as VictoriaMetrics URL instead of Prometheus. For example, `http://vmauth-k8s.monitoring:8427`.                                                                                                                                                                             |
| integrationTests.resources.requests.memory    | string  | no        | 256Mi                    | The minimum amount of memory the container should use. The value can be specified with SI suffixes (E, P, T, G, M, K, m) or their power-of-two-equivalents (Ei, Pi, Ti, Gi, Mi, Ki).                                                                                                                                                                                                                                                                                                                                                       |
| integrationTests.resources.requests.cpu       | string  | no        | 200m                     | The minimum number of CPUs the container should use.                                                                                                                                                                                                                                                                                                                                                                                                                                                                                       |
| integrationTests.resources.limits.memory      | string  | no        | 256Mi                    | The maximum amount of memory the container can use. The value can be specified with SI suffixes (E, P, T, G, M, K, m) or their power-of-two-equivalents (Ei, Pi, Ti, Gi, Mi, Ki).                                                                                                                                                                                                                                                                                                                                                          |
| integrationTests.resources.limits.cpu         | string  | no        | 400m                     | The maximum number of CPUs the container can use.                                                                                                                                                                                                                                                                                                                                                                                                                                                                                          |
| integrationTests.customLabels                 | object  | no        | {}                       | The custom labels for the ZooKeeper Service integration tests pod.                                                                                                                                                                                                                                                                                                                                                                                                                                                                         |
| integrationTests.securityContext              | object  | no        | {}                       | The pod security context for the ZooKeeper Service integration tests pod.                                                                                                                                                                                                                                                                                                                                                                                                                                                                  |
| integrationTests.zookeeperHost                | string  | no        | "zookeeper"              | The host name of ZooKeeper.                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                |
| integrationTests.zookeeperPort                | integer | no        | 2181                     | The port of ZooKeeper.                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                     |
| integrationTests.url                          | string  | yes       |                          | The parameter specifies the URL of Kubernetes/OpenShift server.                                                                                                                                                                                                                                                                                                                                                                                                                                                                            |
| integrationTests.zookeeperIsManagedByOperator | boolean | no        | `true`                   | The parameter specifies whether ZooKeeper is managed by Kubernetes operator                                                                                                                                                                                                                                                                                                                                                                                                                                                                |

## Vault Credentials Management

There is an ability to store ZooKeeper service credentials in Vault secrets instead of Kubernetes secrets.
<!-- #GFCFilterMarkerStart# -->
For more information, see [ZooKeeper Vault Credentials Management](vault.md).
<!-- #GFCFilterMarkerEnd# -->

| Parameter                                         | Type    | Mandatory | Default value            | Description                                                                                                                                                                                                                                                                                                                                                                                                                  |
|---------------------------------------------------|---------|-----------|--------------------------|------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| vaultSecretManagement.enabled                     | boolean | no        | false                    | Whether to store credentials in Vault. If enabled, no passwords from parameters are used, and auto generated passwords are used instead.                                                                                                                                                                                                                                                                                     |
| vaultSecretManagement.url                         | string  | no        | ""                       | The address of Vault service. It can be internal Kubernetes address or external URL. For example, ```http://vault-service.vault:8200```.                                                                                                                                                                                                                                                                                     |
| vaultSecretManagement.method                      | string  | no        | kubernetes               | The name of the Vault authentication method for operator login.                                                                                                                                                                                                                                                                                                                                                              |
| vaultSecretManagement.role                        | string  | no        | kubernetes-operator-role | The name of Vault role for operator login.                                                                                                                                                                                                                                                                                                                                                                                   |
| vaultSecretManagement.path                        | string  | no        | secret                   | The path to Vault secret to store credentials.                                                                                                                                                                                                                                                                                                                                                                               |
| vaultSecretManagement.writePolicies               | string  | no        | true                     | The operator to create policies and roles for services. If the operator role does not allow creating policies, this parameter should be set to "false" and the corresponding policies should be created manually before the installation. <!-- #GFCFilterMarkerStart# --><br> For more information, see [Vault Prerequisites](vault.md#vault-prerequisites)<!-- #GFCFilterMarkerEnd# -->.                                    |
| vaultSecretManagement.passwordGenerationMechanism | string  | no        | operator                 | The mechanism that should be used to generate passwords. There are two options: <br> `operator` - The passwords are generated internally by the operator.<br> `vault` - The passwords are generated by Vault with the corresponding password policies. This option is available with Vault 1.5+.                                                                                                                             |
| vaultSecretManagement.refreshCredentials          | string  | no        | false                    | Whether to refresh credentials if they exist in the Vault. If set to "true", the operator generates new passwords even if Vault already has corresponding secrets. If set to "false", new passwords are generated only if Vault does not have the corresponding secrets. <!-- #GFCFilterMarkerStart# -->The parameter is used as part of [Credentials Rotation](vault.md#credentials-rotation).<!-- #GFCFilterMarkerEnd# --> |

## Integration test tags description

This section contains information about integration test tags that can be used in order to test ZooKeeper service. You can use the following tags:

* `zookeeper` tag runs all presented tests except `Test Disk Is Filled On One Node` and `Full Eviction Test` tests:
    * `zookeeper_acl` tag runs all tests connected to ACL scenario.
    * `zookeeper_sasl` tag runs all tests connected to SASL scenarios.
    * `prometheus` tag runs all tests connected to Prometheus scenarios:
        * `zookeeper_prometheus_alert` tag runs all tests for Prometheus alert cases:
            * `zookeeper_is_degraded_alert` tag runs `ZooKeeper Is Degraded Alert` test.
            * `zookeeper_is_down_alert` tag runs `ZooKeeper Is Down Alert` test.
        * `jmx_metrics` tag runs `Check JMX Metrics` test.
    * `zookeeper_backup_daemon` tag runs all tests connected to the backup scenarios:
        * `hierarchical_backup` tag runs tests for hierarchical backup cases:
            * `restore_hierarchical_backup` tag runs `Restore Hierarchical Backup` test.
            * `restore_hierarchical_backup_advanced` tag runs `Restore Hierarchical Backup Advanced` test.
            * `restore_hierarchical_backup_high_load` tag runs `Restore Hierarchical Backup High Load` test.
            * `create_and_delete_hierarchical_backup` tag runs `Create And Delete Hierarchical Backup` test.
        * `transactional_backup` tag runs tests for transactional backup case:
            * `create_and_delete_transactional_backup` tag runs `Create And Delete Transactional Backup` test.
        * `Full Eviction Test` test is performed when `full_eviction_test` tag is specified explicitly.
        * `unauthorized_access` tag runs `Unauthorized Access` test.
    * `zookeeper_crud` tag runs all tests for creating, reading, updating and removing ZooKeeper data.
    * `zookeeper_ha` tag runs all tests connected to HA scenarios:
        * `zookeeper_ha_without_leader` tag runs `Test Zookeeper Without Leader Node` test.
        * `zookeeper_ha_disk_is_filled` tag runs `Test Disk Is Filled On One Node` test.

# Installation

## Before You Begin

* Make sure the environment corresponds the requirements in the [Prerequisites](#prerequisites) section.
* Before doing major upgrade, it is recommended to make a backup.
* Check if the application is already installed and find its previous deployments' parameters to make changes.

### Ops Portal Preparation

Make sure all YAML values are escaped in accordance with the Ops portal syntax.

### Helm

To deploy via Helm you need to prepare yaml file with custom deploy parameters and run the following
command in [ZooKeeper Chart](/operator/charts/helm/zookeeper-service):

```bash
helm install [release-name] ./ -f [parameters-yaml] -n [namespace]
```

If you need to use resource profile then you can use the following command:

```bash
helm install [release-name] ./ -f ./resource-profiles/[profile-name-yaml] -f [parameters-yaml] -n [namespace]
```

**Warning**: pure Helm deployment does not support the automatic CRD upgrade procedure, so you need to perform it manually.

```bash
kubectl replace -f ./crds/crd.yaml
```

## On-Prem Examples

### HA Scheme

The minimal template for HA scheme is as follows.

```yaml
global:
  name: zookeeper
  secrets:
    zooKeeper:
      adminUsername: admin
      adminPassword: admin
      clientUsername: client
      clientPassword: client
    backupDaemon:
      username: admin
      password: admin
zooKeeper:
  replicas: 3
  storage:
    className:
      - {applicable_to_env_storage_class}
    size: 10Gi
  heapSize: 512
  resources:
    requests:
      cpu: 200m
      memory: 1Gi
    limits:
      cpu: 500m
      memory: 1Gi
  quorumAuthEnabled: true
  securityContext:
    fsGroup: 1000
monitoring:
  install: true
  securityContext:
    runAsUser: 1000
backupDaemon:
  install: true
  backupStorage:
    storageClass: {applicable_to_env_storage_class}
    size: 10Gi
  securityContext:
    runAsUser: 1000
integrationTests:
  install: false
DEPLOY_W_HELM: true
ESCAPE_SEQUENCE: true
```

### DR Scheme

Not applicable

## Google Cloud Examples

### HA Scheme

<details>
<summary>Click to expand YAML</summary>

```yaml
global:
  name: zookeeper
  secrets:
    zooKeeper:
      adminUsername: admin
      adminPassword: admin
      clientUsername: client
      clientPassword: client
    backupDaemon:
      username: admin
      password: admin
zooKeeper:
  replicas: 3
  storage:
    className:
      - {applicable_to_env_storage_class}
    size: 10Gi
  heapSize: 512
  resources:
    requests:
      cpu: 200m
      memory: 1Gi
    limits:
      cpu: 500m
      memory: 1Gi
  quorumAuthEnabled: true
  securityContext:
    fsGroup: 1000
monitoring:
  install: true
  securityContext:
    runAsUser: 1000
backupDaemon:
  install: true
  backupStorage:
  s3:
    enabled: true
    url: "https://storage.googleapis.com"
    bucket: {google_cloud_storage_bucket}
    keyId: {google_cloud_storage_key_id}
    keySecret: {google_cloud_storage_secret}
  securityContext:
    runAsUser: 1000
integrationTests:
  install: false
DEPLOY_W_HELM: true
ESCAPE_SEQUENCE: true
```

</details>

### DR Scheme

## AWS Examples

### HA Scheme

<details>
<summary>Click to expand YAML</summary>

```yaml
global:
  name: zookeeper
  secrets:
    zooKeeper:
      adminUsername: admin
      adminPassword: admin
      clientUsername: client
      clientPassword: client
    backupDaemon:
      username: admin
      password: admin
zooKeeper:
  replicas: 3
  storage:
    className:
      - {applicable_to_env_storage_class}
    size: 10Gi
  heapSize: 512
  resources:
    requests:
      cpu: 200m
      memory: 1Gi
    limits:
      cpu: 500m
      memory: 1Gi
  quorumAuthEnabled: true
  securityContext:
    fsGroup: 1000
monitoring:
  install: true
  securityContext:
    runAsUser: 1000
backupDaemon:
  install: true
  backupStorage:
  s3:
    enabled: true
    url: "https://s3.amazonaws.com"
    bucket: {amazon_s3_bucket}
    keyId: {amazon_s3_account_key_id}
    keySecret: {amazon_s3_account_key_secret}
  securityContext:
    runAsUser: 1000
integrationTests:
  install: false
DEPLOY_W_HELM: true
ESCAPE_SEQUENCE: true
```

</details>

### DR Scheme

Not applicable

## Azure Examples

### HA Scheme

The same as [On-Prem Examples HA Scheme](#on-prem-examples).

### DR Scheme

Not applicable

# Upgrade

## Common

In the common way, the upgrade procedure is the same as the initial deployment.
You need to follow `Release Notes` and `Breaking Changes` in the version you install to find details.
If you upgrade to a version which has several major diff changes from the installed version (e.g. 0.3.1 over 0.1.1), you need to
check `Release Notes` and `Breaking Changes` sections for `0.2.0` and `0.3.0` versions.

## Rolling Upgrade

ZooKeeper supports rolling upgrade feature with near-zero downtime.
It can be enabled with `zooKeeper.rollingUpdate: true`, by default it is disabled.

## CRD Upgrade

Custom resource definition `ZooKeeperService` should be upgraded before the installation if the new version has major
changes.
<!-- #GFCFilterMarkerStart# -->
The CRD for this version is stored in [crd.yaml](../../operator/charts/helm/zookeeper-service/crds/crd.yaml) and can be
applied with the following command:

  ```sh
  kubectl replace -f crd.yaml
  ```

<!-- #GFCFilterMarkerEnd# -->

## HA to DR Scheme

Not applicable

# Rollback

ZooKeeper does not support rollback with downgrade of a major version.

# Additional Features

## Multiple Availability Zone Deployment

When deploying to a cluster with several availability zones, it is important that ZooKeeper pods start in different
availability zones.

### Affinity

You can manage pods' distribution using `affinity` rules to prevent Kubernetes from running ZooKeeper pods on nodes of
the same availability zone.

#### Affinity cross Kubernetes Nodes

Default affinity cross kubernetes Nodes (1 pod per 1 node). 

<details>
<summary>Click to expand YAML</summary>

```yaml
zooKeeper:
  affinity: {
    "podAntiAffinity": {
      "requiredDuringSchedulingIgnoredDuringExecution": [
        {
          "labelSelector": {
            "matchExpressions": [
              {
                "key": "component",
                "operator": "In",
                "values": [
                  "zookeeper"
                ]
              }
            ]
          },
          "topologyKey": "kubernetes.io/hostname"
        }
      ]
    },
    "nodeAffinity": {
      "requiredDuringSchedulingIgnoredDuringExecution": {
        "nodeSelectorTerms": [
          {
            "matchExpressions": [
              {
                "key": "role",
                "operator": "In",
                "values": [
                  "compute"
                ]
              }
            ]
          }
        ]
      }
    }
  }
```

</details>

Where:

* `kubernetes.io/hostname` is the name of the label that defines the node name. This is the default name for Kubernetes.
* `role` and `compute` are the sample name and value of label that defines the region to run ZooKeeper pods.

**Note**: This section describes deployment only for `storage class` Persistent Volumes (PV) type because with
Predefined PV, the ZooKeeper pods are started on the nodes that are specified explicitly with Persistent Volumes. In
that way, it is necessary to take care of creating PVs on nodes belonging to different AZ in advance.

#### Replicas Fewer Than Availability Zones

For cases when the number of ZooKeeper pods (value of the `zooKeeper.replicas` parameter) is equal or less than the number of
availability zones, you need to restrict the start of pods to one pod per availability zone.
You can also specify additional node affinity rule to start pods on allowed Kubernetes nodes.

For this, you can use the following affinity rules:

<details>
<summary>Click to expand YAML</summary>

```yaml
zooKeeper:
  affinity: {
    "podAntiAffinity": {
      "requiredDuringSchedulingIgnoredDuringExecution": [
        {
          "labelSelector": {
            "matchExpressions": [
              {
                "key": "component",
                "operator": "In",
                "values": [
                  "zookeeper"
                ]
              }
            ]
          },
          "topologyKey": "topology.kubernetes.io/zone"
        }
      ]
    },
    "nodeAffinity": {
      "requiredDuringSchedulingIgnoredDuringExecution": {
        "nodeSelectorTerms": [
          {
            "matchExpressions": [
              {
                "key": "role",
                "operator": "In",
                "values": [
                  "compute"
                ]
              }
            ]
          }
        ]
      }
    }
  }
```

</details>

Where:

* `topology.kubernetes.io/zone` is the name of the label that defines the availability zone. This is the default name for Kubernetes
  1.17+. Earlier, `failure-domain.beta.kubernetes.io/zone` was used.
* `role` and `compute` are the sample name and value of label that defines the region to run ZooKeeper pods.

#### Replicas More Than Availability Zones

For cases when the number of ZooKeeper pods (value of the `zooKeeper.replicas` parameter) is greater than the number of availability
zones, you need to restrict the start of pods to one pod per node and specify the preferred rule to start on different
availability zones.
You can also specify an additional node affinity rule to start the pods on allowed Kubernetes nodes.

For this, you can use the following affinity rules:

<details>
<summary>Click to expand YAML</summary>

```yaml
zooKeeper:
  affinity: {
    "podAntiAffinity": {
      "requiredDuringSchedulingIgnoredDuringExecution": [
        {
          "labelSelector": {
            "matchExpressions": [
              {
                "key": "component",
                "operator": "In",
                "values": [
                  "zookeeper"
                ]
              }
            ]
          },
          "topologyKey": "kubernetes.io/hostname"
        }
      ],
      "preferredDuringSchedulingIgnoredDuringExecution": [
        {
          "weight": 100,
          "podAffinityTerm": {
            "labelSelector": {
              "matchExpressions": [
                {
                  "key": "component",
                  "operator": "In",
                  "values": [
                    "zookeeper"
                  ]
                }
              ]
            },
            "topologyKey": "topology.kubernetes.io/zone"
          }
        }
      ]
    },
    "nodeAffinity": {
      "requiredDuringSchedulingIgnoredDuringExecution": {
        "nodeSelectorTerms": [
          {
            "matchExpressions": [
              {
                "key": "role",
                "operator": "In",
                "values": [
                  "compute"
                ]
              }
            ]
          }
        ]
      }
    }
  }
```

</details>

Where:

* `kubernetes.io/hostname` is the name of the label that defines the Kubernetes node. This is a standard name for Kubernetes.
* `topology.kubernetes.io/zone` is the name of the label that defines the availability zone. This is a standard name for
  Kubernetes 1.17+. Earlier, `failure-domain.beta.kubernetes.io/zone` was used.
* `role` and `compute` are the sample name and value of the label that defines the region to run ZooKeeper pods.

# Frequently Asked Questions

## Deploy job failed with some error in templates

Make sure you performed the necessary [Prerequisites](#prerequisites). Fill the [Parameters](#parameters) correctly and compare
with [Examples](#on-prem-examples).
