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
GENERAL.IP-IFACE:RouterOS
IP4.ADDRESS[1]:10.10.20.5/32
`
	c := parseDetails("RouterOS", input)
	if c.Name != "RouterOS" {
		t.Errorf("Name = %q", c.Name)
	}
	if !c.Active {
		t.Error("Active should be true")
	}
	if c.Interface != "RouterOS" {
		t.Errorf("Interface = %q", c.Interface)
	}
	if c.AssignedIP != "10.10.20.5/32" {
		t.Errorf("AssignedIP = %q", c.AssignedIP)
	}
}

func TestParseWgShow(t *testing.T) {
	input := `interface: RouterOS
  public key: Oq7udzct5y+/KBauXAb3t+CLxSnM668004k6ymS/+wM=
  private key: (hidden)
  listening port: 51820

peer: EfVhhGXkkwHwMrfFSq2oZ+X1ITVlR/4+6YuaD+pNpz0=
  endpoint: 72.61.0.249:51822
  allowed ips: 0.0.0.0/0
  latest handshake: 1 minute, 41 seconds ago
  transfer: 81.45 MiB received, 6.00 MiB sent
`
	c := &Connection{}
	parseWgShow(c, input)

	if c.PublicKey != "Oq7udzct5y+/KBauXAb3t+CLxSnM668004k6ymS/+wM=" {
		t.Errorf("PublicKey = %q", c.PublicKey)
	}
	if c.PeerPublicKey != "EfVhhGXkkwHwMrfFSq2oZ+X1ITVlR/4+6YuaD+pNpz0=" {
		t.Errorf("PeerPublicKey = %q", c.PeerPublicKey)
	}
	if c.Endpoint != "72.61.0.249:51822" {
		t.Errorf("Endpoint = %q", c.Endpoint)
	}
	if c.RxBytes == 0 {
		t.Error("RxBytes should be non-zero")
	}
	if c.TxBytes == 0 {
		t.Error("TxBytes should be non-zero")
	}
	if c.LastHandshake.IsZero() {
		t.Error("LastHandshake should be set")
	}
}
