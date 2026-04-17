package nmcli

import (
	"errors"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

// Connection represents a WireGuard profile managed by NetworkManager.
type Connection struct {
	Name          string
	Active        bool
	Interface     string
	AssignedIP    string
	Endpoint      string
	PublicKey     string
	PeerPublicKey string
	LastHandshake time.Time
	RxBytes       int64
	TxBytes       int64
	AllPeers      []Peer
}

// Peer represents a single WireGuard peer from `wg show`.
type Peer struct {
	PublicKey     string
	Endpoint      string
	LastHandshake time.Time
	RxBytes       int64
	TxBytes       int64
}

// CheckPrerequisites returns an error if nmcli is unavailable or NetworkManager is not running.
func CheckPrerequisites() error {
	if _, err := exec.LookPath("nmcli"); err != nil {
		return errors.New("nmcli not found in PATH — install NetworkManager")
	}
	out, err := exec.Command("nmcli", "-t", "general", "status").Output()
	if err != nil {
		return fmt.Errorf("NetworkManager not responding: %w", err)
	}
	if !strings.Contains(string(out), "connected") && !strings.Contains(string(out), "disconnected") {
		return errors.New("NetworkManager is not running")
	}
	return nil
}

// ListConnections returns all WireGuard connections with their active state.
func ListConnections() ([]Connection, error) {
	out, err := exec.Command("nmcli", "-t", "-f", "NAME,TYPE,STATE", "connection", "show").Output()
	if err != nil {
		return nil, fmt.Errorf("nmcli list: %w", err)
	}
	return parseConnections(string(out)), nil
}

func parseConnections(output string) []Connection {
	var conns []Connection
	for _, line := range strings.Split(strings.TrimSpace(output), "\n") {
		if line == "" {
			continue
		}
		parts := strings.SplitN(line, ":", 3)
		if len(parts) != 3 {
			continue
		}
		if parts[1] != "wireguard" {
			continue
		}
		conns = append(conns, Connection{
			Name:   parts[0],
			Active: parts[2] == "activated",
		})
	}
	return conns
}

// GetDetails returns extended details for a named connection.
// For active connections it enriches data from `wg show`.
func GetDetails(name string) (Connection, error) {
	out, err := exec.Command("nmcli", "-t", "connection", "show", name).Output()
	if err != nil {
		return Connection{}, fmt.Errorf("nmcli show %q: %w", name, err)
	}
	c := parseDetails(name, string(out))

	if c.Active {
		enrichFromWg(&c)
	}
	return c, nil
}

func parseDetails(name, output string) Connection {
	c := Connection{Name: name}
	for _, line := range strings.Split(output, "\n") {
		k, v, ok := strings.Cut(line, ":")
		if !ok {
			continue
		}
		v = strings.TrimSpace(v)
		switch k {
		case "GENERAL.STATE":
			c.Active = strings.Contains(v, "activated")
		case "GENERAL.IP-IFACE":
			c.Interface = v
		case "IP4.ADDRESS[1]":
			c.AssignedIP = v
		}
	}
	return c
}

// enrichFromWg populates peer/traffic fields from `wg show <iface>`.
func enrichFromWg(c *Connection) {
	iface := c.Interface
	if iface == "" {
		iface = c.Name
	}
	out, err := exec.Command("wg", "show", iface).Output()
	if err != nil {
		// wg not available or not root — skip silently
		return
	}
	parseWgShow(c, string(out))
}

func parseWgShow(c *Connection, output string) {
	var current *Peer
	for _, line := range strings.Split(output, "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "interface:") {
			current = nil
			continue
		}
		if strings.HasPrefix(line, "peer:") {
			p := Peer{PublicKey: strings.TrimSpace(strings.TrimPrefix(line, "peer:"))}
			c.AllPeers = append(c.AllPeers, p)
			current = &c.AllPeers[len(c.AllPeers)-1]
			continue
		}
		k, v, ok := strings.Cut(line, ":")
		if !ok {
			continue
		}
		k = strings.TrimSpace(k)
		v = strings.TrimSpace(v)
		if current == nil {
			if k == "public key" {
				c.PublicKey = v
			}
			continue
		}
		switch k {
		case "endpoint":
			current.Endpoint = v
		case "latest handshake":
			current.LastHandshake = parseHandshake(v)
		case "transfer":
			parseTransferInto(&current.RxBytes, &current.TxBytes, v)
		}
	}
	// populate top-level fields from first peer for backward compat
	if len(c.AllPeers) > 0 {
		c.PeerPublicKey = c.AllPeers[0].PublicKey
		c.Endpoint = c.AllPeers[0].Endpoint
		c.LastHandshake = c.AllPeers[0].LastHandshake
		c.RxBytes = c.AllPeers[0].RxBytes
		c.TxBytes = c.AllPeers[0].TxBytes
	}
}

func parseHandshake(s string) time.Time {
	// `wg show` returns relative strings like "1 minute, 41 seconds ago"
	// We store a synthetic time by subtracting the duration from now.
	// Simple parse: look for known units.
	s = strings.TrimSuffix(s, " ago")
	var total time.Duration
	parts := strings.Split(s, ", ")
	for _, p := range parts {
		p = strings.TrimSpace(p)
		var n int
		var unit string
		fmt.Sscanf(p, "%d %s", &n, &unit)
		unit = strings.TrimSuffix(unit, "s") // plurals
		switch unit {
		case "second":
			total += time.Duration(n) * time.Second
		case "minute":
			total += time.Duration(n) * time.Minute
		case "hour":
			total += time.Duration(n) * time.Hour
		case "day":
			total += time.Duration(n) * 24 * time.Hour
		}
	}
	if total == 0 {
		return time.Time{}
	}
	return time.Now().Add(-total)
}

func parseTransferInto(rx, tx *int64, s string) {
	var rxVal, txVal float64
	var rxUnit, txUnit string
	fmt.Sscanf(s, "%f %s received, %f %s sent", &rxVal, &rxUnit, &txVal, &txUnit)
	*rx = toBytes(rxVal, rxUnit)
	*tx = toBytes(txVal, txUnit)
}

func parseTransfer(c *Connection, v string) {
	parseTransferInto(&c.RxBytes, &c.TxBytes, v)
}

func toBytes(val float64, unit string) int64 {
	switch strings.ToLower(strings.TrimSuffix(unit, ",")) {
	case "kib":
		return int64(val * 1024)
	case "mib":
		return int64(val * 1024 * 1024)
	case "gib":
		return int64(val * 1024 * 1024 * 1024)
	case "b":
		return int64(val)
	}
	return int64(val)
}

// Connect activates a WireGuard connection by name.
func Connect(name string) error {
	out, err := exec.Command("nmcli", "connection", "up", name).CombinedOutput()
	if err != nil {
		return fmt.Errorf("%s", strings.TrimSpace(string(out)))
	}
	return nil
}

// Disconnect deactivates a WireGuard connection by name.
func Disconnect(name string) error {
	out, err := exec.Command("nmcli", "connection", "down", name).CombinedOutput()
	if err != nil {
		return fmt.Errorf("%s", strings.TrimSpace(string(out)))
	}
	return nil
}

// AddConnection creates a new WireGuard profile in NetworkManager.
func AddConnection(name, endpoint, serverPubKey, assignedIP, privKey, presharedKey, dns string) error {
	args := []string{
		"connection", "add",
		"type", "wireguard",
		"con-name", name,
		"wireguard.private-key", privKey,
		"wireguard.peers", fmt.Sprintf("public-key=%s,endpoint=%s,allowed-ips=0.0.0.0/0", serverPubKey, endpoint),
		"ipv4.method", "manual",
		"ipv4.addresses", assignedIP,
	}
	if presharedKey != "" {
		args = append(args, "wireguard.peers.preshared-key", presharedKey)
	}
	if dns != "" {
		args = append(args, "ipv4.dns", dns)
	}
	out, err := exec.Command("nmcli", args...).CombinedOutput()
	if err != nil {
		return fmt.Errorf("%s", strings.TrimSpace(string(out)))
	}
	return nil
}

// DeleteConnection removes a WireGuard profile from NetworkManager.
func DeleteConnection(name string) error {
	out, err := exec.Command("nmcli", "connection", "delete", name).CombinedOutput()
	if err != nil {
		return fmt.Errorf("%s", strings.TrimSpace(string(out)))
	}
	return nil
}
