# Patroni Exporter for Kubernetes

Connects to Patroni's metrics endpoint and exports metrics in Prometheus format.

Exports metrics on port 9394 at `/metrics`.

## Install

```sh
helm upgrade --install patroni-exporter charts/patroni-exporter
```

Check out the [chart readme](charts/patroni-exporter/README.md) for further details. The most basic config will require:

```yaml
env:
  - name: PATRONI_SERVER_URL
    value: "http://MY-DB.DB-NAMESPACE.svc.cluster.local:8009"

databaseService:
  create: true
  name: MY-DB-patroni
  port: 8009
  type: ClusterIP
  targetPort: 8009

  # how to locate the master DB pod. If using CrunchyData for example, this should work.
  # Run the command below to get more labels you can use as selectors
  #
  # kubectl get po -n DB-NAMESPACE --show-labels
  #
  selector:
    pg-cluster: MY-DB
    role: master
```
