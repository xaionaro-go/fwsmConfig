package fwsmConfig

import (
	"bytes"
	"fmt"
	"github.com/xaionaro-go/networkControl"
	"io"
	"net"
	"strings"
)

type DHCPCommon networkControl.DHCPCommon

func (dhcp DHCPCommon) CiscoString() (result string) {
	result = fmt.Sprintf("dhcpd dns %v\ndhcpd domain %v\n", NSs(dhcp.NSs).CiscoString(), dhcp.Domain)
	for _, option := range dhcp.Options {
		switch option.ValueType {
		case networkControl.DHCPOPT_ASCII:
			result += fmt.Sprintf("dhcpd option %v ascii %v\n", option.Id, string(option.Value))
		default:
			panic(fmt.Errorf("Unknown DHCP option value type: %v", option.ValueType))
		}
	}
	return
}

type DHCP networkControl.DHCP
type DHCPs networkControl.DHCPs
type DHCPOptionValueType networkControl.DHCPOptionValueType
type DHCPOption networkControl.DHCPOption
type DHCPOptions networkControl.DHCPOptions

func (dhcp DHCPCommon) WriteTo(writer io.Writer) error {
	fmt.Fprintf(writer, "%v", dhcp.CiscoString())
	return nil
}

func (dhcp DHCP) WriteTo(writer io.Writer) error {
	if len(dhcp.NSs) != 0 || len(dhcp.Options) != 0 || dhcp.Domain != "" {
		panic(fmt.Errorf("This case is not implemented, yet: %v", dhcp))
	}

	fmt.Fprintf(writer, "dhcpd address %v-%v %v\n", dhcp.RangeStart.String(), dhcp.RangeEnd.String(), dhcp.IfName)

	return nil
}
func (dhcp *DHCP) ParseRange(ipRangeString string) error {
	words := strings.Split(ipRangeString, "-")
	dhcp.RangeStart = net.ParseIP(words[0])
	dhcp.RangeEnd = net.ParseIP(words[1])
	return nil
}

func parseDHCPOptionValueType(valueTypeString string) networkControl.DHCPOptionValueType {
	switch valueTypeString {
	case "ascii":
		return networkControl.DHCPOPT_ASCII
	}

	panic("Unknown DHCP option value type: <" + valueTypeString + ">")
	return networkControl.DHCPOPT_UNKNOWN
}

func (dhcps DHCPs) CiscoString() string {
	var buf bytes.Buffer
	for _, dhcp := range dhcps {
		DHCP(dhcp).WriteTo(&buf)
	}
	return buf.String()
}
