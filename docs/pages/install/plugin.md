---
icon: lucide/cable
title: Installation as proxy plugin
---

# :lucide-cable: Installation as proxy plugin

## :lucide-download: Installing

Add the `github.com/kasefuchs/lazygate` module to your project:

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
