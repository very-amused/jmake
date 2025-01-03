package main

import (
	"bufio"
	"encoding/binary"
	"errors"
	"fmt"
	"math/bits"
	"net/netip"
	"os"
	"path"
	"regexp"
	"strconv"

	jtmp "github.com/very-amused/jmake/template"
)

type BridgeConfig struct {
	Name        string // Bridge interface name (i.e "bridge0")
	Description string // Bridge description for documentation purposes (i.e "DMZ")

	Network string // Combined bridge IP + netmask (i.e "192.168.1.2/24")
	IP      string // Bridge IP address (i.e "192.168.1.2")
	Netmask string // Bridge IP subnet (i.e "24" or "255.255.255.0")

	bridgeNo      int          // Parsed bridge interface number
	networkPrefix netip.Prefix // Parsed CIDR network prefix for the bridge
}

// Bridge template execution context obtained from (BridgeConfig *).context()
type BridgeCtx struct {
	Name          string
	Description   string
	BridgeNo      int
	NetworkPrefix netip.Prefix
}

func (b *BridgeConfig) context() (ctx BridgeCtx) {
	return BridgeCtx{
		Name:          b.Name,
		Description:   b.Description,
		BridgeNo:      b.bridgeNo,
		NetworkPrefix: b.networkPrefix}
}

func (b *BridgeConfig) makeTemplates(_ *Config) (err error) {
	if b.Name == "" {
		return errors.New("missing bridge interface name")
	}
	if err = b.parsePrefix(); err != nil {
		return err
	}
	if err = b.parseBridgeNo(); err != nil {
		fmt.Println(err)
		return err
	}

	// Exit early if rc template has already been written
	rcPath := path.Join(jtmp.TemplateDir, jtmp.BridgeRC)
	if _, err = os.Stat(rcPath); err == nil {
		return nil // template file exists
	}

	var (
		rcFile *bufio.Writer
	)

	if file, err := os.Create(path.Join(jtmp.TemplateDir, jtmp.BridgeRC)); err != nil {
		return err
	} else {
		defer file.Close()
		rcFile = bufio.NewWriter(file)
		defer rcFile.Flush()
	}

	rcFile.WriteString("# Bridge {{.Name}} ({{.Description}}) config\n")
	vtnet := "vtnet{{.BridgeNo}}"
	jtmp.WriteRc(rcFile, "ifconfig_{{.Name}}", fmt.Sprintf("inet {{.NetworkPrefix.String}} addm %s", vtnet))
	jtmp.WriteRc(rcFile, fmt.Sprintf("ifconfig_%s", vtnet), "up")

	return nil
}

func (b *BridgeConfig) execTemplates(c *Config) {
	jtmp.ExecAppend(b.context(), jtmp.BridgeRC)
}

func (b *BridgeConfig) parseBridgeNo() (err error) {
	regex, err := regexp.Compile("([A-Za-z]+)(\\d+)")
	if err != nil {
		return err
	}
	if matches := regex.FindStringSubmatch(b.Name); len(matches) >= 3 {
		b.bridgeNo, err = strconv.Atoi(matches[2])
	}
	return err
}

func (b *BridgeConfig) parsePrefix() (err error) {
	if b.Network != "" {
		b.networkPrefix, err = netip.ParsePrefix(b.Network)
		return err
	}

	if b.IP != "" && b.Netmask != "" {
		var ip netip.Addr
		if ip, err = netip.ParseAddr(b.IP); err != nil {
			return err
		}
		// Try parsing netmask shorthand
		if netmask, err := strconv.ParseUint(b.Netmask, 10, 32); err == nil {
			b.networkPrefix, err = ip.Prefix(int(netmask))
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
			b.networkPrefix, err = ip.Prefix(prefixLen)
			return err
		}

		return fmt.Errorf("failed to parse netmask for bridge %s", b.Name)
	}

	return errors.New("missing both network and ip/netmask keys")
}
