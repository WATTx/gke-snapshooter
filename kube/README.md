# Deploying to k8s

## Credentials

### k8s

Since k8s mounts certificates into pods it's enough to provide `--in-cluster` flag.

### Compute Engine

At the moment compute engine doesn't have an integration with GKE. 

**If the key is not present** at k8s deployment yet, you need to do the following:
1. Create a new service account at the iam project page, for example: https://console.cloud.google.com/iam-admin/serviceaccounts/project?project=wattx-infra

1. Generate and save a new key for this account.

1. base64 the key, for example:

    ```
    cat ~/Downloads/wattx-infra-82ddbfadb1cb.json | base64 > .gce-key
    ```
    
1. upload the secret:

    ```
    kubectl create secret generic wattx-common --from-file=gceKey=./.gce-key
    ```

Now it's possible to reference the key in a pod spec.

### Slack

You can create a bot and get an API key in the slack settings.

## Deploying

```
kubectl create -f app.yaml
```
