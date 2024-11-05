package registry

import (
	"context"
	"math/rand"
	"net"
	"time"

	"github.com/kasefuchs/lazygate/pkg/provider"
	"go.minekube.com/gate/pkg/edition/java/config"
	"go.minekube.com/gate/pkg/edition/java/netmc"
	"go.minekube.com/gate/pkg/edition/java/proto/packet"
	"go.minekube.com/gate/pkg/edition/java/proto/state"
	"go.minekube.com/gate/pkg/edition/java/proto/version"
	"go.minekube.com/gate/pkg/edition/java/proxy"
	"go.minekube.com/gate/pkg/gate/proto"
	"go.minekube.com/gate/pkg/util/netutil"
)

// Entry contains server metadata.
type Entry struct {
	Server     proxy.RegisteredServer // Server, obviously.
	Allocation provider.Allocation    // Allocation associated with server.

	LastActive      time.Time // Last time when someone reached the server.
	KeepOnlineUntil time.Time // Force server to be online until this.
}

// NewEntry creates new instance of Entry.
func NewEntry(srv proxy.RegisteredServer, alloc provider.Allocation) *Entry {
	return &Entry{
		Server:     srv,
		Allocation: alloc,
	}
}

// UpdateLastActive updates the last active time.
func (e *Entry) UpdateLastActive() {
	e.LastActive = time.Now()
}

// KeepOnlineFor force the server to be online for the given duration.
func (e *Entry) KeepOnlineFor(d time.Duration) {
	e.KeepOnlineUntil = time.Now().Add(d)
}

// ShouldStop decides whether server is going to sleep.
func (e *Entry) ShouldStop() bool {
	// Allocation must be started.
	if e.Allocation.State() != provider.AllocationStateStarted {
		return false
	}

	// Don't stop if players connected.
	if e.Server.Players().Len() > 0 {
		return false
	}

	// Don't stop until keep online ends.
	if e.KeepOnlineUntil.After(time.Now()) {
		return false
	}

	// Don't stop until we can get allocation configuration.
	cfg, err := e.Allocation.Config()
	if err != nil {
		return false
	}

	// Stop if last active is under activity threshold.
	return time.Since(e.LastActive) >= time.Duration(cfg.Time.InactivityThreshold)
}

// Ping pings underlying minecraft server.
func (e *Entry) Ping(ctx context.Context, cfg config.Config) bool {
	// Close everything after function.
	c, cancel := context.WithCancel(ctx)
	defer cancel()

	// Get server address.
	addr := e.Server.ServerInfo().Addr()

	// Dial to server.
	var dialer net.Dialer
	base, err := dialer.DialContext(c, addr.Network(), addr.String())
	if err != nil {
		return false
	}

	// Create client.
	conn, _ := netmc.NewMinecraftConn(
		c, base, proto.ClientBound,
		time.Duration(cfg.ReadTimeout), time.Duration(cfg.ConnectionTimeout), cfg.Compression.Level,
	)

	// Perform handshake.
	host, port := netutil.HostPort(addr)
	if err := conn.WritePacket(&packet.Handshake{
		ProtocolVersion: int(version.MinimumVersion.Protocol),
		NextStatus:      int(packet.StatusHandshakeIntent),

		ServerAddress: host,
		Port:          int(port),
	}); err != nil {
		return false
	}

	// Create ping packet.
	ping := &packet.StatusPing{
		RandomID: rand.Int63(),
	}

	// Ping server.
	conn.SetState(state.Status)
	if err := conn.WritePacket(ping); err != nil {
		return false
	}

	// Receive pong.
	pack, err := conn.Reader().ReadPacket()
	if err != nil {
		return false
	}

	// Verify pong.
	registry := state.FromDirection(proto.ServerBound, state.Status, pack.Protocol)
	if id, ok := registry.PacketID(ping); ok && id == pack.PacketID {
		return ping.RandomID == pack.Packet.(*packet.StatusPing).RandomID
	}

	return false
}
