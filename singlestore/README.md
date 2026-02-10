# singlestore-helm

A Helm chart for deploying SingleStore (MemSQL) clusters and operator on Kubernetes.

## Features
- Installs MemsqlCluster CRD, creates service accounts and RBAC
- Deploys the SingleStore Operator and a custom MemsqlCluster resource
- All configuration is parameterized via `values.yaml`
- Manages secret for license and hash of admin password

## Prerequisites

This chart requires Helm 3 or later.

## Installation

```sh
helm install <release-name> ./singlestore
```

## Configuration
All values can be overridden in your `values.yaml` or via `--set` on the command line.

In any place in the table where <role> is used, this can be substituted with one of: master, aggregator, leaf.

| Parameter                                   | Description                                         | Default                                  |
|---------------------------------------------|-----------------------------------------------------|------------------------------------------|
| `clusterName`                               | Name of the SingleStore cluster                     | sdb-cluster                              |
| `operator.enabled`                          | Deploy the operator                                 | true                                     |
| `operator.image.repository`                 | Operator image repository                           | singlestore/operator                     |
| `operator.image.tag`                        | Operator image tag                                  | 4.81.0-0d1ea9c8                          |
| `spec.nodeImage.repository`                 | Node image repository                               | singlestore/node                         |
| `spec.nodeImage.tag`                        | Node image tag                                      | alma-8.9.40-b828fb951c                   |
| `spec.rootPasswordSecret.name`              | Name of the secret containing root password         | empty(will be filled in by operator)     |
| `spec.rootPasswordSecret.key`               | Key in the secret containing root password          | empty(will be filled in by operator)     |
| `spec.license`                              | License key (base64 or plain)                       | empty(required)                          |
| `spec.adminPasswordSecret`                  | Hashed admin password                               | empty(required)                          |
| `spec.redundancyLevel`                      | Cluster redundancy level                            | 2                                        |
| `spec.serviceSpec.objectMetaOverrides`      | Overrides for service metadata                      | empty                                    |
| `spec.serviceSpec.type`                     | Service type for cluster services                   | LoadBalancer                             |
| `spec.schedulingDetails`                    | Define scheduling details per role                  | empty, see table below                   |
| `spec.master`                               | Specification for master pods                       | see table below                          |
| `spec.aggregator`                           | Specification for child aggregator pods             | empty                                    |
| `spec.leaf`                                 | Specification for leaf pods                         | see table below                          |
| `spec.defaultStorageConfig`                 | Remote storage configuration for unlimited storage  | empty, see table below for config        |
| `spec.globalVariables`                      | Cluster-wide global variables                       | empty                                    |
| `monitoringJob.enabled`                     | Enable the monitoring setup job                     | false                                    |
| `sysreqDaemonSet.enabled`                   | Enable DaemonSet for K8s nodes configuration        | false                                    |


Nodes of each role(master, aggregator, leaf) must be configured separately using corresponding entry in values file. For each node group chart supports configuring following options:

| Parameter                                   | Description            | Default                                                   |
|---------------------------------------------|------------------------|-----------------------------------------------------------|
| `spec.<role>.cores`                         | CPU request            | 4 for master and leaf                                     |
| `spec.<role>.coresLimit`                    | CPU limit              | empty, will be expanded to `cores` by operator            |
| `spec.<role>.memoryMB`                      | Memory request in MB   | 16384                                                     |
| `spec.<role>.memoryLimitMB`                 | Memory limit in MB     | empty, will be expanded to `memoryMB` by operator         |
| `spec.<role>.globalVariables`               | Global variables       | empty                                                     |
| `spec.<role>.storageGB`                     | Amount of storage (GB) | 256 for master, 1024 for leaf                             |
| `spec.<role>.storageClass`                  | Storage class          | empty(required)                                           |
| `spec.<role>.objectMetaOverrides`           | Metadata overrides     | empty                                                     |

To configure remote storage for unlimited storage, `spec.defaultStorageConfig` accepts following options:
| Parameter                                     | Description                        | Default                                                   |
|-----------------------------------------------|------------------------------------|-----------------------------------------------------------|
| `spec.defaultStorageConfig.storageProvider`   | Storage provider for remote storage| empty(required)                                           |
| `spec.defaultStorageConfig.storageURI`        | Storage URI                        | empty(required)                                           |
| `spec.defaultStorageConfig.storageSecret`     | Storage credentials secret details | empty(required)                                           |
| `spec.defaultStorageConfig.storageSecret.name`| Name of the credentials secret     | empty(required)                                           |
| `spec.defaultStorageConfig.storageSecret.key` | Key used in credentials secret     | empty(required)                                           |

`objectMetaOverrides` fields can be specified in a format:
```yaml
objectMetaOverrides:
  labels:
    some-label: label-value
  annotations:
    some-annotations: annotation-value
```

Scheduling details are configured per role under `spec.schedulingDetails` entry:

| Parameter                                   | Description                                         | Default                                  |
|---------------------------------------------|-----------------------------------------------------|------------------------------------------|
| `spec.schedulingDetails.<role>.nodeSelector`| Node selector for pods with given role              | empty                                    |
| `spec.schedulingDetails.<role>.tolerations` | Tolerations for pods with given role                | empty                                    |

## Secret Management
- The chart creates a secret named `singlestore-helm-<clusterName>-secrets` containing the license key and hashed admin password from `values.yaml`.

## RBAC and Security
- **RBAC:** The chart installs all required Roles, ClusterRoles, RoleBindings, and ServiceAccounts for SingleStore Operator.
- **Operator permissions:**
  - Full access to pods, services, endpoints, PVCs, events, configmaps, secrets, deployments, daemonsets, replicasets, statefulsets, networkpolicies, and all memsql.com resources in the namespace.
  - Cluster-wide read access to storageclasses, persistentvolumes, and nodes.
- **Principle of Least Privilege:** Review and restrict permissions for your environment as needed.

## Uninstall

```sh
helm delete <release-name>
```

## Example Commands

Upgrade the chart:

```sh
helm upgrade <release-name> <new-chart-version>
```

## Optional System Requirements DaemonSet

To ensure all Kubernetes nodes meet SingleStore's system requirements, this chart includes optional DaemonSet that configures cluster nodes.
DaemonSet's pods run in privileged security context.

To enable DaemonSet specify `sysreqDaemonSet.enabled: true` in the values override. If DaemonSet should only run on a subset of Kubernetes cluster nodes specify `sysreqDaemonSet.nodeSelector` and `sysreqDaemonSet.tolerations` fields.

## Optional Monitoring Setup Job

This chart can optionally deploy a one-time monitoring setup job using the official SingleStore Toolbox image. When monitoring job is enabled, chart will also create a secret to store monitoring user credentials and dedicated service account.

To enable, set `monitoringJob.enabled: true` and provide the required values in your `values.yaml`:

```yaml
monitoringJob:
  enabled: true
  image: singlestore/tools:alma-v1.11.6-1.17.2-cc87b449d97fd7cde78fdc4621c2aec45cc9a6cb
  username: "admin"
  password: "<admin-password>"
```

This job will use the dedicated ServiceAccount (`singlestore-helm-<clusterName>-monitoring`) which has all required permissions. 

See the [SingleStore monitoring documentation](https://docs.singlestore.com/db/v8.9/reference/singlestore-operator-reference/monitor-your-kubernetes-cluster/configure-cluster-monitoring-with-the-operator/#section-idm4604985730918433590373732789) for more details.
