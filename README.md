# LazyGate

## Table of Contents

- [About](#about)
- [Getting Started](#getting_started)

## About <a name = "about"></a>

LazyGate is a [Gate proxy](https://github.com/minekube/gate) plugin that shuts down your Minecraft server when it's idle
and wakes it up when players connect.

## Getting Started <a name = "getting_started"></a>

### Installing

Add the `lazygate` module to your project:

```sh
go get github.com/kasefuchs/lazygate
```

Include the plugin in your list of proxy plugins:

```go
package main

import (
	lazygate "github.com/kasefuchs/lazygate/pkg/plugin"
	"go.minekube.com/gate/cmd/gate"
	"go.minekube.com/gate/pkg/edition/java/proxy"
)

func main() {
	proxy.Plugins = append(proxy.Plugins, lazygate.NewProxyPlugin())

	gate.Execute()
}
```

### Configuring

Configure the plugin using environment variables:

```sh
# Provider to use. Currently available nomad & docker.
LAZYGATE_PROVIDER="nomad"
```

### Usage

LazyGate matches registered Gate servers with provider's allocations using labels:

**Docker Compose:**

```yaml
services:
  minecraft-server-random:
    labels:
      lazygate.server: random_name
      lazygate.time.minimumOnline: 2m
      lazygate.time.inactivityThreshold: 5m
      lazygate.disconnectReasons.starting: random_name currently starting!
      lazygate.disconnectReasons.startFailed: Failed to start random_name :(
```

**Gate config:**

```yaml
---
config:
  servers:
    random_name: minecraft-server-random:25565
  try:
    - random_name
```

In this example, the `random_name` server will correspond to the `minecraft-server1` service.
