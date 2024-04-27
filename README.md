# Rancher Information Agent
Rest tool for fetching information from a Rancher Management cluster about Cluster ID and Display Name of Cluster and Project Resources 

Requires a Service account with get list rights to apiGroup:  management.cattle.io resources: projects, clusters. See [deployment/authorization.yml](./deployment/authorization.yml) for examples.

responds with a json reply with format
```json
[{
  "name": "local", 
  "displayName": "local", 
  "projects": {
    "name": "p-42g36", 
    "displayName": "System"
  }
}]
```

## Command line options
| Option | Description |
| ------ | ----------- |
| -debug | Enable debugging output (developer focused) |
| -port=\[integer\] | Use a port different from 8080 |
| -onlyRootEndpoint=\[boolean\] | Enable/(Disable) Only reply json on / endpoint, disable if you want to have api  |
| -prometheus=\[boolean\] | (Enable)/Disable Prometheus endpoint |
| -prometheusEndpoint=\[string\] | custom prometheus endpoint (/metrics) |
| -healthEndpoint=\[string\] | custom health endpoint (/health) |

# Download
Docker image can be fetched from [github ghcr.io/simonstiil/rancher-info-agent](https://github.com/SimonStiil/rancher-info-agent/pkgs/container/rancher-info-agent)  
Can be build with go build .  
Will also be available as a release in releases in the future

