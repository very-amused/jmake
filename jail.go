package main

import (
	"fmt"
	"log"
	"net/netip"
	"strconv"
	"strings"
)

type JailConfig struct {
	name string

	// Jail IP addresses by bridge
	IP map[string]string

	// Parsed IP field
	addrs map[string]netip.Addr
}

// Parse and verify jail IP addresses
func (j *JailConfig) parseIPs(bridges map[string]*BridgeConfig) (err error) {
	j.addrs = make(map[string]netip.Addr)
	for bridgeName, ip := range j.IP {
		configStr := fmt.Sprintf("%s.ip.%s = %s", j.name, bridgeName, ip)
		// Ignore and warn about IPs that don't attach to a defined bridge
		bridge, ok := bridges[bridgeName]
		if !ok {
			log.Println("Warning - ignoring IP attached to an undefined bridge:", configStr)
			delete(j.IP, ip)
			continue
		}

		// Parse the address and validate that it falls under the bridge's subnet
		var addr netip.Addr
		if addr, err = netip.ParseAddr(ip); err != nil {
			return err
		}
		if !bridge.NetworkPrefix.Contains(addr) {
			return fmt.Errorf("jail IP is outside bridge subnet: %s", configStr)
		}
		// Save parsed + verified addr
		j.addrs[bridgeName] = addr
	}

	return nil
}

// Get an IP's host identifier in a concatenated, base 10 form with no separator characters.
// This form is suitable for use in generating epair names: epair{hostID}{bridgeNo}.
// NOTE: this function assumes that {ip} is contained within {prefix}
func hostID(ip netip.Addr, prefix netip.Prefix) string {
	mask := prefix.Masked().Addr().AsSlice()
	addr := ip.AsSlice()
	if len(mask) != len(addr) {
		return "" // ip protocol mismatch
	}

	var stem strings.Builder
	var i int
	for i = range addr {
		// Remove subnet prefix
		addr[i] &= ^mask[i]
		// Start the host ID at the first non-zero segment
		if addr[i] != 0 {
			break
		}
	}
	for _, a := range addr[i:] {
		stem.WriteString(strconv.FormatUint(uint64(a), 10))
	}

	return stem.String()
}
