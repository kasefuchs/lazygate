package main

import (
	lazygate "codeberg.org/kasefuchs/lazygate/pkg/plugin"
	"go.minekube.com/gate/cmd/gate"
	"go.minekube.com/gate/pkg/edition/java/proxy"
)

func main() {
	proxy.Plugins = append(proxy.Plugins, lazygate.NewProxyPlugin())

	gate.Execute()
}
