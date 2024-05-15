# Cluster Setup

Densify supports two use-cases for this repository.

## Single Kubernetes Cluster

Use this [config.yaml](config.yaml) file (with a single cluster with no identifiers) as a template.

## Multiple Kubernetes Clusters

Use this [config.yaml](config.yaml) file as a template.

> **_NOTE:_**  V4 of Densify Container Data Collection is backwards-compatible and has full support for the [deprecated **properties** format](config.properties) of the config of versions 1-3. However, new features introduced in V4 are configurable using [**yaml** format](config.yaml) only, and new configs should be created only using **yaml** format. The **properties** format will be removed in a feature release.
