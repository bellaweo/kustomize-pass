# kustomize-pass [![Build Status](https://cloud.drone.io/api/badges/bellaweo/kustomize-pass/status.svg)](https://cloud.drone.io/bellaweo/kustomize-pass)

kustomize plugin to use secrets stored in pass (passwordstore.org)

(this plugin modified from the kustomize secretsfromdatabase plugin example)

## Requirements

- Go 1.14
- `kustomize`

I use asdf (asdf-vm.com) for both

## Build & Install

set `$XDG_CONFIG_HOME` environment variable

build the plugin

```
go build -buildmode plugin -o $XDG_CONFIG_HOME/kustomize/plugin/someteam.example.com/v1/secretsfrompass/SecretsFromPass.so
```

create a configuration file

```
cat <<'EOF' >$XDG_CONFIG_HOME/secretFromPass.yaml
apiVersion: someteam.example.com/v1
kind: SecretsFromPass
metadata:
  name: mySecretGenerator
name: forbiddenValues
namespace: production
keys:
- ROCKET
- VEGETABLE
EOF
```

create the kustomization file referencing this plugin

```
cat <<'EOF' >$XDG_CONFIG_HOME/kustomization.yaml
generators:
- secretFromPass.yaml
EOF
```

finally, generate the yaml

```
kustomize build --enable_alpha_plugins $XDG_CONFIG_HOME
```

## Testing

to manually run the tests:

```
mkdir -p ${HOME}/plugin/someteam.example.com/v1/secretsfrompass
```
```
go build -buildmode=plugin -o ${HOME}/plugin/someteam.example.com/v1/secretsfrompass/SecretsFromPass.so
```
```
go test ./...
```