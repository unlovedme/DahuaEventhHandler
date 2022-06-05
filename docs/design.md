# Kubewatch

Kubewatch contains three components: controller, config, handler

![Kubewatch Diagram](kubewatch.png?raw=true "Kubewatch Overview")

## Config

The config object contains `kubewatch` configuration, like handlers, filters.

A config object is used to creating new client.

## Controller

The controller initializes using the config object by reading the `.kubewatch.yaml` or c