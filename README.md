# kube-resource

A simple CLI that provides an overview of the resource requests, limits in a Kubernetes cluster

## Installation

### Homebrew

```bash
brew tap ysicing/tap
brew install kr
```

### Krew

> not support

```bash
kubectl krew install kr
```

### Bash

```bash
curl -L --remote-name-all https://github.com/ysicing/kubectl-resource/releases/latest/download/kr_linux_amd64{,.sha256sum}
sha256sum --check kr_linux_amd64.sha256sum
mv kr_linux_amd64 /usr/local/bin/kubectl-kr
chmod +x /usr/local/bin/kubectl-kr
ln -s /usr/local/bin/kubectl-kr /usr/local/bin/kr
```

## Usage

```bash
kubectl kr [OPTIONS]
kr [OPTIONS]
```
