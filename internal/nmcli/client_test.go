package nmcli

import (
	"testing"
)

func TestParseConnections(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		want   []Connection
	}{
		{
			name:  "active wireguard",
			input: "RouterOS:wireguard:activated\neth0:ethernet:activated\n",
			want:  []Connection{{Name: "RouterOS", Active: true}},
		},
		{
			name:  "inactive wireguard",
			input: "Internacional:wireguard:--\nNacional:wireguard:--\n",
			want: []Connection{
				{Name: "Internacional", Active: false},
				{Name: "Nacional", Active: false},
			},
		},
		{
			name:  "no wireguard connections",
			input: "eth0:ethernet:activated\n",
			want:  nil,
		},
		{
			name:  "empty output",
			input: "",
			want:  nil,
		},
		{
			name:  "mixed connections",
			input: "RouterOS:wireguard:activated\nInternacional:wireguard:--\neth0:ethernet:activated\n",
			want: []Connection{
				{Name: "RouterOS", Active: true},
				{Name: "Internacional", Active: false},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parseConnections(tt.input)
			if len(got) != len(tt.want) {
				t.Fatalf("got %d connections, want %d", len(got), len(tt.want))
			}
			for i, c := range got {
				if c.Name != tt.want[i].Name || c.Active != tt.want[i].Active {
					t.Errorf("conn[%d] = {%q, %v}, want {%q, %v}", i, c.Name, c.Active, tt.want[i].Name, tt.want[i].Active)
				}
			}
		})
	}
}

func TestParseDetails(t *testing.T) {
	input := `GENERAL.NAME:RouterOS
GENERAL.STATE:100 (activated)
GENERAL.IP-IFACE:wg0
IP4.ADDRESS[1]:10.0.0.2/24
wireguard.public-key:abc123==
wireguard.peers[1].endpoint:vpn.example.com:51820
wireguard.peers[1].public-key:srv456==
`
	c := parseDetails("RouterOS", input)
	if c.Name != "RouterOS" {
		t.Errorf("Name = %q", c.Name)
	}
	if !c.Active {
		t.Error("Active should be true")
	}
	if c.Interface != "wg0" {
		t.Errorf("Interface = %q", c.Interface)
	}
	if c.AssignedIP != "10.0.0.2/24" {
		t.Errorf("AssignedIP = %q", c.AssignedIP)
	}
	if c.Endpoint != "vpn.example.com:51820" {
		t.Errorf("Endpoint = %q", c.Endpoint)
	}
}
