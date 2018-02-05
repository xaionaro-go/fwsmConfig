package fwsmConfig

import (
	"io"
	"net"
)

type VLAN struct {
	net.Interface
	SecurityLevel    int
	IPs              IPs
	AttachedNetworks IPNets
}

type VLANs []VLAN

func (a VLANs) Len() int           { return len(a) }
func (a VLANs) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a VLANs) Less(i, j int) bool { return a[i].Index < a[j].Index }

func (vlan VLAN) WriteTo(writer io.Writer) error {
	return nil
}
