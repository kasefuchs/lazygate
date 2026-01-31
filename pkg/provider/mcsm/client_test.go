package mcsm_test

import (
	"context"
	"testing"

	"github.com/kasefuchs/lazygate/pkg/provider/mcsm"
)

func TestClient_ServerStart(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for receiver constructor.
		baseUrl   string
		api_key   string
		daemon_id string
		// Named input parameters for target function.
		id      string
		wantErr bool
	}{
		// TODO: Add test cases.
		struct {
			name      string
			baseUrl   string
			api_key   string
			daemon_id string
			id        string
			wantErr   bool
		}{name: "server", baseUrl: "http://192.168.50.102:23333", api_key: "edec08fcafd2410daa0a748dd8bf884c", daemon_id: "5b3e6cece5ff426283c2b5df0c735d00", id: "dc70e47aff0a4432b11454123265312b"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := mcsm.NewClient(context.Background(), tt.baseUrl, tt.api_key, tt.daemon_id)
			gotErr := c.ServerStart(tt.id)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("ServerStart() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("ServerStart() succeeded unexpectedly")
			}
		})
	}
}
