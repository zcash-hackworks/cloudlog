# cloudlog
CO AUTHOR SANITY CHECK
A small application to forward received messages.

Right now, it will take in a [CloudEvent](https://github.com/cloudevents/spec) message from a [Tekton](https://github.com/tektoncd/pipeline) task and send Slack notifications.

## Usage

### Local usage for testing

```
$ slackURL='https://hooks.slack.com/services/T...' ./cloudlog
```

### Kubernetes deployment

```
# Create Slack Webhook secret
kubectl create secret generic cloudlog-secret --from-literal slackURL='https://hooks.slack.com/services/T...'

# Create the cloudlog service
kubectl apply -f kubernetes/service.yml

# Create the cloudlog deployment (build your own image?)
kubectl apply -f kubernetes/deployment.yml

# Check that the pod is running
kubectl get pods  -l app=cloudlog

# Follow the logs
kubectl logs -f -l app=cloudlog
```

### Tekton usage

Create a resource for the service

```
---
apiVersion: tekton.dev/v1alpha1
kind: PipelineResource
metadata:
  name: event-to-cloudlog
spec:
  type: cloudEvent
  params:
    - name: targetURI
      value: http://cloudlog.default.svc.cluster.local:3000/inbox
```

Utilize the resource as the output for a task
```
---
apiVersion: tekton.dev/v1alpha1
kind: Task
metadata:
  name: send-a-message
spec:
  inputs:
...
  outputs:
    resources:
      - name: notification
        type: cloudEvent
```

Provide TaskRun parameters for the resource
```
---
apiVersion: tekton.dev/v1alpha1
kind: TaskRun
metadata:
  generateName: send-message-
  label: send-message
spec:
...
  outputs:
    resources:
      - name: notification
        resourceRef:
          name: event-to-cloudlog
```
