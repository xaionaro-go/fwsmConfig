package fwsmConfig

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"github.com/xaionaro-go/networkControl"
	"io"
	"net"
	"sort"
	"strings"
)

type DHCP networkControl.DHCP
type DHCPRange networkControl.DHCPRange
type DHCPSubnet struct {
	networkControl.DHCPSubnet
}

var isAsciiString = networkControl.IsDHCPAsciiString

func NewDHCP() *DHCP {
	return (*DHCP)(networkControl.NewDHCP())
}

func (r *DHCPRange) Parse(ipRangeString string) error {
	words := strings.Split(ipRangeString, "-")
	r.Start = net.ParseIP(words[0])
	r.End   = net.ParseIP(words[1])
	return nil
}

func parseDHCPOptionValue(valueTypeString string, valueString string) ([]byte, networkControl.DHCPValueType) {
	switch valueTypeString {
	case "hex":
		v, err := hex.DecodeString(valueString)
		if err != nil {
			panic(err)
		}
		return v, networkControl.DHCPValueType_BYTEARRAY
	case "ascii":
		return []byte(valueString), networkControl.DHCPValueType_ASCIISTRING
	}

	panic("Unknown DHCP option value type: <" + valueTypeString + ">")
	return []byte{}, networkControl.DHCPValueType_UNKNOWN
}

func (dhcp DHCP) CiscoWriteTo(writer io.Writer, vlans VLANs) (err error) {
	if len(dhcp.Options.DomainNameServers) == 0 || dhcp.Options.DomainName == "" {
		panic("This case is not implemented, yet")
	}

	_, err = fmt.Fprintf(writer, "dhcpd dns %v\ndhcpd domain %v\n", NSs(dhcp.Options.DomainNameServers.ToNetNSs()).CiscoString(), dhcp.Options.DomainName)
	if err != nil {
		return
	}

	var codes []int
	for code := range dhcp.Options.Custom {
		codes = append(codes, code)
	}
	sort.IntSlice(codes).Sort()

	for _, code := range codes {
		value := dhcp.Options.Custom[code]
		if isAsciiString(string(value)) {
			_, err = fmt.Fprintf(writer, "dhcpd option %v ascii %v\n", code, string(value))
		} else {
			_, err = fmt.Fprintf(writer, "dhcpd option %v hex %v\n", code, hex.EncodeToString(value))
		}
		if err != nil {
			return
		}
	}

	subnetToVlanMap := map[string]*VLAN{}
	for _, vlan := range vlans {
		for _, ip := range vlan.IPs {
			subnetToVlanMap[ip.IP.Mask(ip.Mask).String()] = vlan
		}
	}

	var subnetKeys []string
	for k := range dhcp.Subnets {
		subnetKeys = append(subnetKeys, k)
	}
	sort.StringSlice(subnetKeys).Sort()

	for _, k := range subnetKeys {
		subnet := dhcp.Subnets[k]
		vlan := subnetToVlanMap[subnet.Network.IP.String()]
		if vlan == nil {
			panic(fmt.Errorf("Cannot find a VLAN for subnet %v: %v", k, subnet))
			continue
		}
		fmt.Fprintf(writer, "dhcpd address %v-%v %v\n", subnet.Options.Range.Start.String(), subnet.Options.Range.End.String(), vlan.Name)
	}

	return nil
}
func (dhcp DHCP) CiscoString(vlans VLANs) string {
	var buf bytes.Buffer
	dhcp.CiscoWriteTo(&buf, vlans)
	return buf.String()
}

