# singlestore-helm

A Helm chart for deploying SingleStore (MemSQL) clusters via a Kubernetes operator.

## Features

* Installs a MemsqlCluster CRD  
* Creates service accounts and RBAC  
* Deploys the SingleStore Operator and a custom MemsqlCluster resource  
* All configurations are parameterized via a `values.yaml` file  
* Manages secrets for the license and a hash of the admin password

## Prerequisites

* [Helm](https://helm.sh/) 3 or later  
* A [SingleStore license](https://www.singlestore.com/pricing/?product=Self-Managed)

## Installation

```shell
helm install <release-name> ./singlestore
```

## Configuration

All values can be overridden in the `values.yaml` file or via the `--set` option on the command line.  

Anywhere `role` is used in the following tables, it can be substituted with master, aggregator, or leaf.

### Configuration Parameters

| Parameter | Description | Default |
| :---- | :---- | :---- |
| `clusterName` | Name of the SingleStore cluster | sdb-cluster |
| `operator.enabled` | Deploy the operator | true |
| `operator.image.repository` | Operator image repository | singlestore/operator |
| `operator.image.tag` | Operator image tag | 4.81.0-0d1ea9c8 |
| `spec.nodeImage.repository` | Node image repository | singlestore/node |
| `spec.nodeImage.tag` | Node image tag | alma-8.9.40-b828fb951c |
| `spec.rootPasswordSecret.name` | Name of the secret which contains the root password | empty (filled in by the operator) |
| `spec.rootPasswordSecret.key` | Key in the secret which contains the root password | empty (filled in by the operator) |
| `spec.license` | License key (base64 or plain) | empty (required) |
| `spec.adminPasswordSecret` | Hashed admin password | empty (required) |
| `spec.redundancyLevel` | Cluster redundancy level | 2 |
| `spec.serviceSpec.objectMetaOverrides` | Overrides for service metadata | empty |
| `spec.serviceSpec.type` | Service type for cluster services | LoadBalancer |
| `spec.schedulingDetails` | Defines scheduling details per role | empty (Refer to the [Scheduling Parameters](#scheduling-parameters) table) |
| `spec.master` | Specification for master aggregator pods | Refer to the [Role Parameters](#role-parameters) table |
| `spec.aggregator` | Specification for child aggregator pods | empty |
| `spec.leaf` | Specification for leaf node pods | Refer to the [Role Parameters](#role-parameters) table |
| `spec.defaultStorageConfig` | Remote storage configuration for unlimited storage | empty (Refer to the [Storage Parameters](#storage-parameters) table) |
| `spec.globalVariables` | Cluster-wide global variables | empty |
| `monitoringJob.enabled` | Enable the monitoring setup job | false |
| `sysreqDaemonSet.enabled` | Enable DaemonSet for Kubernetes nodes configuration | false |

Nodes of each role (master, aggregator, leaf) must be configured separately using a corresponding entry in the `values.yaml` file. For each node group, this Helm chart supports configuring the following parameters:

### Role Parameters

| Parameter | Description | Default |
| :---- | :---- | :---- |
| `spec.<role>.cores` | CPU request | 4 for master aggregator and leaf node |
| `spec.<role>.coresLimit` | CPU limit | empty (will be expanded to the number of host cores by the operator) |
| `spec.<role>.memoryMB` | Memory request in megabytes (MB) | 16384 |
| `spec.<role>.memoryLimitMB` | Memory limit in megabytes (MB) | empty (expanded to memoryMB by the operator) |
| `spec.<role>.globalVariables` | Global variables | empty |
| `spec.<role>.storageGB` | Amount of storage in gigabytes (GB) | 256 for master aggregator; 1024 for leaf node |
| `spec.<role>.storageClass` | Storage class | empty (required) |
| `spec.<role>.objectMetaOverrides` | Metadata overrides | empty |

To configure remote storage as unlimited storage, `spec.defaultStorageConfig` accepts following parameters:

### Storage Parameters

| Parameter | Description | Default |
| :---- | :---- | :---- |
| `spec.defaultStorageConfig.storageProvider` | Storage provider for remote storage | empty (required) |
| `spec.defaultStorageConfig.storageURI` | Storage URI | empty (required) |
| `spec.defaultStorageConfig.storageSecret` | Storage credentials secret details | empty (required) |
| `spec.defaultStorageConfig.storageSecret.name` | Name of the credentials secret | empty (required) |
| `spec.defaultStorageConfig.storageSecret.key` | Key used in credentials secret | empty (required) |

The objectMetaOverrides fields can be specified using the following format:

```yaml
objectMetaOverrides:
  labels:
    some-label: label-value
  annotations:
    some-annotations: annotation-value
```

Scheduling details are configured per role under the spec.schedulingDetails entry:

### Scheduling Parameters

| Parameter | Description | Default |
| :---- | :---- | :---- |
| `spec.schedulingDetails.<role>.nodeSelector` | Node selector for pods with a given role | empty |
| `spec.schedulingDetails.<role>.tolerations` | Tolerations for pods with given role | empty |

## Secret Management

The Helm chart creates a secret named `singlestore-helm-<clusterName>-secrets` containing the license key and hashed admin password from the `values.yaml` file.

## RBAC and Security

For RBAC, this Helm chart installs all required Roles, ClusterRoles, RoleBindings, and ServiceAccounts for the SingleStore Operator.

Operator permissions:

* Full access to pods, services, endpoints, PVCs, events, configmaps, secrets, deployments, daemonsets, replicasets, statefulsets, networkpolicies, and all memsql.com resources in the namespace.  
* Cluster-wide read access to storageclasses, persistentvolumes, and nodes.  
* Principle of least privilege: Review and restrict permissions for the target environment as needed.

## Uninstall

```shell
helm delete <release-name>
```

## Optional System Requirements DaemonSet

To ensure all Kubernetes nodes meet SingleStore's [system requirements](https://docs.singlestore.com/docs/system-requirements/), this Helm chart includes an optional DaemonSet that configures cluster nodes. The DaemonSet's pods run in a privileged security context.

To enable the DaemonSet, specify sysreqDaemonSet.enabled: true in the values override. If the DaemonSet should only run on a subset of Kubernetes cluster nodes, specify the sysreqDaemonSet.nodeSelector and sysreqDaemonSet.tolerations fields.

## Optional Monitoring Setup

This Helm chart can optionally deploy a one-time monitoring setup job using the SingleStore Toolbox image. When the monitoring job is enabled, the chart creates a secret to store the monitoring user’s credentials and dedicated service account.

To enable, set monitoringJob.enabled: true and provide the required values in the `values.yaml` file.

```yaml
monitoringJob:
  enabled: true
  image: singlestore/tools:alma-v1.11.6-1.17.2-cc87b449d97fd7cde78fdc4621c2aec45cc9a6cb
  username: "admin"
  password: "<admin-password>"
```

This job will use the dedicated ServiceAccount (`singlestore-helm-<clusterName>-monitoring`) which has all required permissions. Refer to the [SingleStore monitoring documentation](https://docs.singlestore.com/docs/monitoring/kubernetes/) for more information.  
