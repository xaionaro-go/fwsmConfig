package fwsmConfig

import (
	"fmt"
	"github.com/xaionaro-go/networkControl"
	"net"
	"strconv"
	"strings"
)

type IPPort networkControl.IPPort
type PortRange networkControl.PortRange
type IPs networkControl.IPs
type IPNet networkControl.IPNet
type IPNets networkControl.IPNets
type IPPorts networkControl.IPPorts
type NSs networkControl.NSs
type PortRanges networkControl.PortRanges
type Domain networkControl.Domain
type Protocol networkControl.Protocol

func (nss NSs) CiscoString() string {
	nssStr := []string{}
	for _, ns := range nss {
		nssStr = append(nssStr, ns.Host)
	}

	return strings.Join(nssStr, " ")
}

func (ipport IPPort) CiscoString() string {
	if ipport.Port != nil && ipport.Protocol == nil {
		panic(fmt.Errorf("This shouldn't happened: %v", ipport))
	}

	protocolPrefix := ""

	portSuffix := ""
	if ipport.Port != nil {
		portSuffix = " " + strconv.Itoa(int(*ipport.Port))
	}

	return protocolPrefix + ipport.IP.String() + portSuffix
}

func (ipport IPPort) String() string {
	protocolSuffix := ""
	if ipport.Protocol != nil {
		protocolSuffix = "/" + Protocol(*ipport.Protocol).CiscoString()
		if protocolSuffix == "/ip" {
			protocolSuffix = ""
		}
	}

	if ipport.Port == nil {
		return ipport.IP.String() + protocolSuffix
	}

	return ipport.IP.String() + ":" + strconv.Itoa(int(*ipport.Port)) + protocolSuffix
}

func parseIPNetUnmasked(ipStr string, maskStr string) (networkControl.IPNet, error) {
	return networkControl.IPNetUnmaskedFromStrings(ipStr, maskStr)
}
func parseIPNet(ipStr string, maskStr string) (networkControl.IPNet, error) {
	return networkControl.IPNetFromStrings(ipStr, maskStr)
}
func parseNS(nsStr string) net.NS {
	return networkControl.NSFromString(nsStr)
}
func parsePort(portStr string) uint16 {
	return networkControl.PortFromString(portStr)
}
func parseProtocol(protocolStr string) networkControl.Protocol {
	return networkControl.ProtocolFromString(protocolStr)
}
func (protocol Protocol) CiscoString() string {
	return networkControl.Protocol(protocol).String()
}

func (ipnet IPNet) CiscoString() string {
	maskString := net.IP(ipnet.Mask).String()
	switch maskString {
	case "255.255.255.255":
		return "host " + ipnet.IP.String()
	case "0.0.0.0":
		return "any"
	}

	return ipnet.IP.String() + " " + maskString

}

func (portRanges PortRanges) CiscoString() string {
	if len(portRanges) == 2 {
		if portRanges[0].Start != 0 || portRanges[1].End != 65535 || portRanges[0].End+2 != portRanges[1].Start {
			panic("This case is not implemented, yet")
		}
		return fmt.Sprintf("neg %v", portRanges[0].End+1)
	}

	portRange := portRanges[0]

	if portRange.Start == 0 && portRange.End == 65535 {
		return ""
	}

	if portRange.Start == portRange.End {
		return fmt.Sprintf("eq %v", portRange.Start)
	}

	if portRange.Start == 0 {
		return fmt.Sprintf("lt %v", portRange.End+1)
	}

	if portRange.End == 65535 {
		return fmt.Sprintf("gt %v", portRange.Start-1)
	}

	return fmt.Sprintf("range %v %v", portRange.Start, portRange.End)
}

func parseNetworkAndPortRanges(words []string) (ipnet networkControl.IPNet, portranges networkControl.PortRanges, unusedWords []string) {
	unusedWords = words

	switch unusedWords[0] { // TODO: add support of IPv6
	case "any":
		ipnet.IP = net.ParseIP("0.0.0.0")
		ipnet.Mask = net.IPMask(net.ParseIP("0.0.0.0"))
		unusedWords = unusedWords[1:]
	case "host":
		ipnet.IP = net.ParseIP(unusedWords[1])
		ipnet.Mask = net.IPMask(net.ParseIP("255.255.255.255"))
		unusedWords = unusedWords[2:]
	default:
		ipnet.IP = net.ParseIP(unusedWords[0])
		ipnet.Mask = net.IPMask(net.ParseIP(unusedWords[1]))
		unusedWords = unusedWords[2:]
	}

	if len(unusedWords) == 0 {
		portranges = append(portranges, networkControl.PortRange{Start: 0, End: 65535})
		return
	}

	switch unusedWords[0] {
	case "eq", "gt", "lt", "neg", "range":
		operator := unusedWords[0]
		port := parsePort(unusedWords[1])
		unusedWords = unusedWords[2:]

		switch operator {
		case "eq":
			portranges = append(portranges, networkControl.PortRange{Start: uint16(port), End: uint16(port)})
		case "gt":
			portranges = append(portranges, networkControl.PortRange{Start: uint16(port) + 1, End: 65535})
		case "lt":
			portranges = append(portranges, networkControl.PortRange{Start: 0, End: uint16(port) - 1})
		case "neg":
			portranges = append(portranges, networkControl.PortRange{Start: 0, End: uint16(port) - 1})
			portranges = append(portranges, networkControl.PortRange{Start: uint16(port) + 1, End: 65535})
		case "range":
			portEnd := parsePort(unusedWords[0])
			portranges = append(portranges, networkControl.PortRange{Start: uint16(port), End: uint16(portEnd)})
			unusedWords = unusedWords[1:]
		}
	}

	if len(portranges) == 0 {
		portranges = append(portranges, networkControl.PortRange{Start: 0, End: 65535})
	}
	return
}

func parseHostPort(words []string) (ip net.IP, port *uint16, unusedWords []string) {
	unusedWords = words

	ip = net.ParseIP(unusedWords[0])
	unusedWords = unusedWords[1:]

	if len(unusedWords) == 0 {
		return
	}

	if strings.Index(unusedWords[0], ".") == -1 && unusedWords[0] != "netmask" {
		parsedPort := parsePort(unusedWords[0])
		unusedWords = unusedWords[1:]
		port = &parsedPort
	}

	return
}

