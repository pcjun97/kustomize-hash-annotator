# Kustomize Hash Annotator

## Background

At the time of writing, for ConfigMap and Secret declared from ConfigMapGenerator
or SecretGenerator, [Kustomize](https://kustomize.io/) can append a hash suffix
to trigger a rolling update when there is a change in the ConfigMap or Secret.

This plugin allows you to annotate any resources with the hash values of any
other resources (including themselves), as long as they are all managed by Kustomize.

## Requirements

- [Go](https://golang.org/doc/install)
- [Kustomize](https://kubectl.docs.kubernetes.io/installation/kustomize)

```bash
# Verify Go is installed and your $GOPATH is set
go env

# Verify Kustomize is installed and your $XDG_CONFIG_HOME is set
kustomize version
echo $XDG_CONFIG_HOME # should be ~/.config or your custom path

# Set up your $XDG_CONFIG_HOME if it is unset
echo "export XDG_CONFIG_HOME=\$HOME/.config" >> $HOME/.bashrc
```

## Installation

### Remotely Download the Latest Release

```bash
mkdir -p $XDG_CONFIG_HOME/kustomize/plugin/pcjun97/v1/hashannotator/v1
wget -c https://github.com/pcjun97/kustomize-hash-annotator/releases/latest/download/hashannotator_latest_$(uname -s)_$(uname -m).tar.gz -O hashannotator_latest.tar.gz
tar -xz hashannotator_latest.tar.gz -C $XDG_CONFIG_HOME/kustomize/plugin/pcjun97/v1/hashannotator/
```

### Run `make install` with the cloned repo Locally

```bash
git clone https://github.com/pcjun97/kustomize-hash-annotator.git
cd kustomize-hash-annotator
make install
```

## Usage

### 1. Create your resources

#### `deployment.yaml`

```yaml
apiVersion: apps/v1beta1
kind: Deployment
metadata:
  name: example
spec:
  selector:
    matchLabels:
      app.kubernetes.io/name: example
  template:
    metadata:
      labels:
        app.kubernetes.io/name: example
    spec:
      containers:
        - name: example
          image: hello-world
          envFrom:
            - configMapRef:
                name: example
```

#### `configmap.yaml`

```yaml
apiVersion: v1beta1
kind: ConfigMap
metadata:
  name: example
data:
  FOO: bar
```

### 2. Add the HashAnnotator transformer

#### `hashannotator.yaml`

```yaml
apiVersion: pcjun97/v1
kind: HashAnnotator
metadata:
  name: not-important-to-example
targets:
  - kind: Deployment
resources:
  - kind: ConfigMap
fieldSpecs:
  - path: spec/template/metadata/annotations
    create: true
```

### 3. Create your `kustomization.yaml`

#### `kustomization.yaml`

```yaml
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization

resources:
  - deployment.yaml
  - configmap.yaml

transformers:
  - hashannotator.yaml
```

### 4. Build with kustomize

```bash
kustomize build --enable-alpha-plugins
```

#### Output

```yaml
apiVersion: v1beta1
data:
  FOO: bar
kind: ConfigMap
metadata:
  name: example
---
apiVersion: apps/v1beta1
kind: Deployment
metadata:
  name: example
spec:
  selector:
    matchLabels:
      app.kubernetes.io/name: example
  template:
    metadata:
      annotations:
        kustomize.config.k8s.io/hash-configmap-example: 59m54fbgh2
      labels:
        app.kubernetes.io/name: example
    spec:
      containers:
      - envFrom:
        - configMapRef:
            name: example
        image: hello-world
        name: example
```

## Transformer Options

### Targets

Select targets to add the annotations.

E.g. select resources with name matching `foo*`:

```yaml
targets:
  - name: foo*
```

Select all resources of kind `Deployment`:

```yaml
targets:
  - kind: Deployment
```

Using multiple fields just makes the target more specific.
The following selects only Deployments that also have the label
`app=hello` (full [label/annotation selector rules](https://kubernetes.io/docs/concepts/overview/working-with-objects/labels/#label-selectors)):

```yaml
targets:
  - kind: Deployment
    labelSelector: app=hello
```

### Resources

Select resources to hash. Uses the same selector rules as [targets](https://github.com/pcjun97/kustomize-hash-annotator#targets).

E.g.

```yaml
resources:
  - kind: ConfigMap
    labelSelector: app=hello
```

### FieldSpecs

Represents paths to fields in resources.
Determines which resource types and which fields within those types
the transformer can modify.

```yaml
fieldSpecs:
- kind: Deployment
  path: spec/template/metadata/annotations
- kind: Pod
  path: metadata/annotations
```
