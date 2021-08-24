patroni-exporter
================
A Helm chart for Kubernetes

Current chart version is `0.1.0`





## Chart Values

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| affinity | object | `{}` |  |
| databaseService.create | bool | `true` |  |
| databaseService.name | string | `"MY-DB-patroni"` |  |
| databaseService.port | int | `8009` |  |
| databaseService.selector.pg-cluster | string | `"machine"` |  |
| databaseService.selector.role | string | `"master"` |  |
| databaseService.targetPort | int | `8009` |  |
| databaseService.type | string | `"ClusterIP"` |  |
| env[0].name | string | `"PATRONI_SERVER_URL"` |  |
| env[0].value | string | `"http://DB-SERVICE-NAME.pgo.svc.cluster.local:8009"` |  |
| env[1].name | string | `"PATRONI_EXPORTER_VERBOSE"` |  |
| env[1].value | string | `"true"` |  |
| fullnameOverride | string | `""` |  |
| image.pullPolicy | string | `"IfNotPresent"` |  |
| image.repository | string | `"digitalisdocker/patroni-exporter"` |  |
| image.tag | string | `""` |  |
| imagePullSecrets | list | `[]` |  |
| nameOverride | string | `""` |  |
| nodeSelector | object | `{}` |  |
| podAnnotations | object | `{}` |  |
| podSecurityContext | object | `{}` |  |
| resources | object | `{}` |  |
| securityContext | object | `{}` |  |
| service.port | int | `80` |  |
| service.targetPort | int | `9394` |  |
| service.type | string | `"ClusterIP"` |  |
| serviceAccount.annotations | object | `{}` |  |
| serviceAccount.create | bool | `true` |  |
| serviceAccount.name | string | `""` |  |
| serviceMonitor.annotations | object | `{}` |  |
| serviceMonitor.create | bool | `true` |  |
| serviceMonitor.jobLabel | string | `"patroni-exporter"` |  |
| serviceMonitor.labels | object | `{}` |  |
| serviceMonitor.namespace | string | `""` |  |
| tolerations | list | `[]` |  |
