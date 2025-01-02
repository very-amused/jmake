package main

import "net/netip"

type BridgeConfig struct {
	Name        string // Bridge interface name (i.e "bridge0")
	Description string // Bridge description for documentation purposes (i.e "DMZ")

	Network string // Combined bridge IP + netmask (i.e "192.168.1.2/24")
	IP      string // Bridge IP address (i.e "192.168.1.2")
	Netmask string // Bridge IP subnet (i.e "24" or "255.255.255.0")

	networkPrefix netip.Prefix // Parsed CIDR network prefix for the bridge
}

func (b *BridgeConfig) parseNetwork() {
	// TODO
}
