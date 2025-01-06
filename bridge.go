package main

import (
	"encoding/binary"
	"errors"
	"fmt"
	"math/bits"
	"net/netip"
	"strconv"

	"github.com/very-amused/jmake/jtmp"
)

type BridgeConfig struct {
	Name        string `toml:"-"` // Bridge interface name (i.e "bridge0") [parsed from toml key]
	Description string // Bridge description for documentation purposes (i.e "DMZ")

	Network string // Combined bridge IP + netmask (i.e "192.168.1.2/24")
	IP      string // Bridge IP address (i.e "192.168.1.2")
	Netmask string // Bridge IP subnet (i.e "24" or "255.255.255.0")

	Interfaces []string // Host interfaces attached to the bridge using addm directives

	NetworkPrefix netip.Prefix `toml:"-"` // Parsed CIDR network prefix for the bridge
}

type BridgeConfigs map[string]*BridgeConfig

func (bc *BridgeConfigs) Generate(_ *Config) (errs []error) {
	for name, bridge := range *bc {
		bridge.Name = name
		if err := bridge.parsePrefix(); err != nil {
			errs = append(errs, err)
		}
	}

	errs = append(errs, jtmp.ExecTemplates(bc, jtmp.BridgeRC)...)

	return errs
}

func (b *BridgeConfig) parsePrefix() (err error) {
	if b.Network != "" {
		b.NetworkPrefix, err = netip.ParsePrefix(b.Network)
		return err
	}

	if b.IP != "" && b.Netmask != "" {
		var ip netip.Addr
		if ip, err = netip.ParseAddr(b.IP); err != nil {
			return err
		}
		// Try parsing netmask shorthand
		if netmask, err := strconv.ParseUint(b.Netmask, 10, 32); err == nil {
			b.NetworkPrefix, err = ip.Prefix(int(netmask))
			return err
		}
		// Try parsing netmask address form (only supported for ipv4 at the moment)
		if netmask, err := netip.ParseAddr(b.Netmask); err == nil && netmask.Is4() {
			// Pad 4 byte ip4 address to 8 bytes
			a4 := netmask.As4()
			a8 := make([]byte, 8)
			copy(a8[4:], a4[:])
			// Count 1s to get prefix len, then apply to bridge ip
			prefixLen := bits.OnesCount64(binary.LittleEndian.Uint64(a8))
			b.NetworkPrefix, err = ip.Prefix(prefixLen)
			return err
		}

		return fmt.Errorf("failed to parse netmask for bridge %s", b.Name)
	}

	return errors.New("missing both network and ip/netmask keys")
}
