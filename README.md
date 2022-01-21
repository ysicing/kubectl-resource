# kube-resource

A simple CLI that provides an overview of the resource requests, limits in a Kubernetes cluster

## Installation

### Homebrew

```bash
brew tap ysicing/tap
brew install kr
```

### Krew

```bash
kubectl krew install kr
```

## Usage

```bash
kube-resource
NAMESPACE       NAME                                            TYPE            CPU REQUESTS    CPU LIMITS      MEMORY REQUESTS MEMORY LIMITS
kruise-system   kruise-controller-manager                       Deployment      200             200             512Mi           512Mi  
```