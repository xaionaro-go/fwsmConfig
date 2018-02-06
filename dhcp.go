package fwsmConfig

import (
	"io"
	"net"
	"strings"
)

type DHCPCommon struct {
	NSs     NSs
	Options DHCPOptions
	Domain  Domain
}

type DHCP struct {
	DHCPCommon

	RangeStart net.IP
	RangeEnd   net.IP
}

type DHCPs []DHCP

type DHCPOptionValueType int

const (
	DHCPOPT_UNKNOWN = DHCPOptionValueType(0)
	DHCPOPT_ASCII = DHCPOptionValueType(1)
)

type DHCPOption struct {
	Id int
	ValueType DHCPOptionValueType
	Value []byte
}

type DHCPOptions []DHCPOption

func (dhcp DHCPCommon) WriteTo(writer io.Writer) error {
	return nil
}

func (dhcp DHCP) WriteTo(writer io.Writer) error {
	return nil
}
func (dhcp *DHCP) ParseRange(ipRangeString string) error {
	words := strings.Split(ipRangeString, "-")
	dhcp.RangeStart = net.ParseIP(words[0])
	dhcp.RangeEnd   = net.ParseIP(words[1])
	return nil
}

func parseDHCPOptionValueType(valueTypeString string) DHCPOptionValueType {
	switch valueTypeString {
	case "ascii":
		return DHCPOPT_ASCII
	}

	panic("Unknown DHCP option value type: <"+valueTypeString+">")
	return DHCPOPT_UNKNOWN
}

