# keda-actions-webhook

Keda has a github runner scaler. The scaler that has been implemented in Keda uses a polling mechanism which can be troublesome for organizatinos and/or users 


## How it works

* Github sends a webhook of type `workflow_job` with the field `action` set to `queued` to the endpoint of keda-actions-webhook (:1234/webhook)
* keda-actions-webhook adds it to a list that is maintained in Redis
* Keda watches the redis list and is configured to scale a job or deployment
* Github sends a webhook of type `workflow_job` with the field `action` set to `in_progress` to the endpoint of keda-actions-webhook (:1234/webhook)
* keda-actions-webhook removes it a list that is maintained in Redis
* The redis scaler that needs to be used has documentation [here](https://keda.sh/docs/2.11/scalers/redis-lists/) 

Notes: 
* The Github webhook payload must be configured to send JSON
* The keda-actions-webhook will check the validity of the webhook. If this is not validated then it will not be picked up and a bad request is returned.


## Configuration

|Variable|Default|Description|
|---|---|---|
|`SERVER_ADRESS`|":1234"|Server address and port on which to expose the webhook server |
|`SECRET_TOKEN`|"CHANGE_ME!"|Github token for webhook validation. [Github docs](https://docs.github.com/en/webhooks-and-events/webhooks/securing-your-webhooks)|
|`REDIS_ADDRESS`|"localhost:6379"|Redis server address|
|`REDIS_PASSWORD`|""|Redis password|
|`REDIS_DB`|0|Redis database index|


## Helm deployment

The following Helm values can be used in combination with [this chart](https://github.com/orgs/lab42/packages/container/package/charts%2Fone). Inspect the package values.yaml if you want to see what other settings can be used.

```yaml
image:
  repository: ghcr.io/lab42/keda-actions-webhook
  pullPolicy: IfNotPresent
  tag: latest

securityContext:
  capabilities:
    drop:
    - ALL
  readOnlyRootFilesystem: true
  runAsNonRoot: true
  runAsUser: 1234

containerPort: 1234

ingress:
  enabled: false
  className: ""
  annotations: {}
  kubernetes.io/ingress.class: nginx
  # kubernetes.io/tls-acme: "true"
  hosts:
    - host: chart-example.local
      paths:
        - path: /
          pathType: ImplementationSpecific
  tls: []
  #  - secretName: chart-example-tls
  #    hosts:
  #      - chart-example.local

resources:
  limits:
    cpu: 100m
    memory: 128Mi
  requests:
    cpu: 50m
    memory: 64Mi

env: {}

envFrom: {}

livenessProbe:
  httpGet:
    path: /healthz
    port: 1234

readinessProbe:
  httpGet:
    path: /healthz
    port: 1234

startuProbe:
  httpGet:
    path: /healthz
    port: 1234

```

## Keda redis scaler configuration

It is advised to use a `ScaledJob` object to configure Github runners and scale them with the `EPHEMERAL` setting.

* Use the same redis server
* ListName should be `runner_count`
* See this page for more information: https://keda.sh/docs/2.11/scalers/redis-lists/
