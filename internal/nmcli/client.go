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
func GetDetails(name string) (Connection, error) {
	out, err := exec.Command("nmcli", "-t", "connection", "show", name).Output()
	if err != nil {
		return Connection{}, fmt.Errorf("nmcli show %q: %w", name, err)
	}
	return parseDetails(name, string(out)), nil
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
		case "wireguard.public-key":
			c.PublicKey = v
		case "wireguard.peers[1].endpoint":
			c.Endpoint = v
		case "wireguard.peers[1].public-key":
			c.PeerPublicKey = v
		}
	}
	return c
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
