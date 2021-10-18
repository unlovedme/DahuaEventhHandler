
## WARNING: Kubewatch is no longer actively maintained by VMware.

VMware has made the difficult decision to stop driving this project and therefore we will no longer actively respond to issues or pull requests. The project will be externally maintained in the following fork: https://github.com/robusta-dev/kubewatch

Thank You.

<p align="center">
  <img src="./docs/kubewatch-logo.jpeg">
</p>


[![Build Status](https://travis-ci.org/bitnami-labs/kubewatch.svg?branch=master)](https://travis-ci.org/bitnami-labs/kubewatch) [![Go Report Card](https://goreportcard.com/badge/github.com/bitnami-labs/kubewatch)](https://goreportcard.com/report/github.com/bitnami-labs/kubewatch) [![GoDoc](https://godoc.org/github.com/bitnami-labs/kubewatch?status.svg)](https://godoc.org/github.com/bitnami-labs/kubewatch) [![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://github.com/bitnami-labs/kubewatch/blob/master/LICENSE)

**kubewatch** is a Kubernetes watcher that currently publishes notification to available collaboration hubs/notification channels. Run it in your k8s cluster, and you will get event notifications through webhooks.

# Usage
```
$ kubewatch -h

Kubewatch: A watcher for Kubernetes

kubewatch is a Kubernetes watcher that publishes notifications
to Slack/hipchat/mattermost/flock channels. It watches the cluster
for resource changes and notifies them through webhooks.

supported webhooks:
 - slack
 - hipchat
 - mattermost
 - flock
 - webhook
 - smtp

Usage:
  kubewatch [flags]
  kubewatch [command]

Available Commands:
  config      modify kubewatch configuration
  resource    manage resources to be watched
  version     print version

Flags:
  -h, --help   help for kubewatch

Use "kubewatch [command] --help" for more information about a command.

```

# Install

### Cluster Installation
#### Using helm:

When you have helm installed in your cluster, use the following setup:

```console
helm install --name kubewatch bitnami/kubewatch --set='rbac.create=true,slack.channel=#YOUR_CHANNEL,slack.token=xoxb-YOUR_TOKEN,resourcesToWatch.pod=true,resourcesToWatch.daemonset=true'
```

You may also provide a values file instead:

```yaml
rbac:
  create: true
resourcesToWatch:
  deployment: false
  replicationcontroller: false
  replicaset: false
  daemonset: false
  services: true
  pod: true
  job: false
  node: false
  clusterrole: true
  serviceaccount: true
  persistentvolume: false
  namespace: false
  secret: false
  configmap: false
  ingress: false
slack:
  channel: '#YOUR_CHANNEL'
  token: 'xoxb-YOUR_TOKEN'
```

And use that:

```console