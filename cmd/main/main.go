package main

import (
	"context"

	"github.com/kasefuchs/lazygate/internal/app/lazygate"
	"go.minekube.com/gate/cmd/gate"
	"go.minekube.com/gate/pkg/edition/java/proxy"
)

func main() {
	proxy.Plugins = append(proxy.Plugins, proxy.Plugin{
		Name: "LazyGate",
		Init: func(ctx context.Context, proxy *proxy.Proxy) error {
			return lazygate.NewPlugin(ctx, proxy).Init()
		},
	})

	gate.Execute()
}
