package fwsmConfig

import (
	"net"
	"io"
)

type DHCPCommon struct {
	NSs []IPs
}

type DHCP struct {
	DHCPCommon

	RangeStart net.IP
	RangeEnd   net.IP
}

func (dhcp DHCPCommon) WriteTo(writer io.Writer) error {
	return nil
}

func (dhcp DHCP) WriteTo(writer io.Writer) error {
	return nil
}
