# External IP Init

A configuration helper to provide external IP discovery.

Replaces a placeholder in a configuration file with the external IP reported
by kubernetes.

### Docker image

Can be released to your private [GCP registry](https://cloud.google.com/container-registry/) with:

    REGISTRY=gcr.io/yourproject make release

### Example k8s config

Deployment spec
```yaml
apiVersion: extensions/v1beta1
kind: Deployment
metadata:
  name: app
spec:
  initContainers:
  - name: app-config
    image: foo:bar
    args:
      - -configmap=config
      - -filename=/app/config/appsettings.xml
    volumeMounts:
      - name: conf
        mountPath: /app/config
  containers:
  - name: app
    image: app:latest
    volumeMounts:
      - name: conf
        mountPath: /app/config
  volumes:
    - name: conf
      emptyDir: {}
```

Config map

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: config
data:
  appsettings.xml: |
    <?xml version="1.0" encoding="utf-8"?>
    <Settings>
      <MatchPoolName>MatchPool</MatchPoolName>
      <ExternalAddress>K8S_EXTERNALADDRESS</ExternalAddress>
    </Settings>
```

This will replace appsettings.xml values before it is read by the container.
