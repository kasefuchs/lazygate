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
# Provider to use. Currently available nomad, docker and pufferpanel.
LAZYGATE_PROVIDER="nomad"
```

### Usage

#### DOCKER
LazyGate matches registered Gate servers with provider's allocations using labels:

**Docker Compose:**

```yaml
services:
  minecraft-server-random:
    labels:
      lazygate.server: random_name
      lazygate.time.minimumOnline: 2m
      lazygate.time.inactivityThreshold: 5m
      lazygate.queue.try: wait,kick
      lazygate.queue.wait.timeout: 10s
      lazygate.queue.wait.pingInterval: 2s
      lazygate.queue.kick.starting: random_name is currently starting!
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

#### PufferPanel

**Enviroment Variables for Gate**
```sh
LAZYGATE_PROVIDER="pufferpanel"
LAZYGATE_PUFFERPANEL_URL="<url>"
LAZYGATE_PUFFERPANEL_CLIENTID="<clientid>"
LAZYGATE_PUFFERPANEL_CLIENTSECRET="<clientsecret>"
```
***URL:*** The Url from pufferpanel e.g. `https://panel.example.com/`

***CLIENTID & CLIENTSECRET:*** The Client ID and Client Secret can you generate from Pufferpanel. Account (Top Right) -> OAuth2 Client -> Create New OAuth Client

**Gate config:**

```yaml
---
config:
  servers:
    random_name: minecraft-server-random:25565
  try:
    - random_name
```

**Server config**

Create a `lazygate.json` file inside pufferpanel files:
```json
{
  "lazygate.server": "random_name"
  "lazygate.time.minimumOnline": "2m"
  "lazygate.time.inactivityThreshold": "5m"
  "lazygate.queue.try": "wait,kick"
  "lazygate.queue.wait.timeout": "10s"
  "lazygate.queue.wait.pingInterval": "2s"
  "lazygate.queue.kick.starting": "random_name is currently starting!"
}
```
**Extras**

If you running your Gate Server also inside Docker on the same host as Pufferpanel you have to pass your pufferpanel domain as extrahost