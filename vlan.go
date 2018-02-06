package fwsmConfig

import (
	"bytes"
	"fmt"
	"io"
	"net"
)

type VLAN struct {
	net.Interface
	SecurityLevel    int
	IPs              IPs
	AttachedNetworks IPNets
}

type VLANs []*VLAN

func (a VLANs) Len() int           { return len(a) }
func (a VLANs) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a VLANs) Less(i, j int) bool { return a[i].Index < a[j].Index }

func (vlan VLAN) WriteTo(writer io.Writer) error {
	fmt.Fprintf(writer, "interface Vlan%v\n", vlan.Index)
	fmt.Fprintf(writer, " nameif %v\n", vlan.Name)
	fmt.Fprintf(writer, " security-level %v\n", vlan.SecurityLevel)
	if len(vlan.IPs) > 0 {
		ip := vlan.IPs[0]
		attachedNet := vlan.AttachedNetworks[0]
		if !attachedNet.Contains(ip) {
			panic("This shouldn't happend")
		}
		mask := net.IP(attachedNet.Mask)
		fmt.Fprintf(writer, " ip address %v %v\n", ip.String(), mask.String())
	} else {
		fmt.Fprintf(writer, " no ip address\n")
	}
	fmt.Fprintf(writer, "!\n")
	return nil
}

func (vlan VLAN) CiscoString() string {
	var buf bytes.Buffer
	vlan.WriteTo(&buf)
	return buf.String()
}
func (vlans VLANs) CiscoString() string {
	var buf bytes.Buffer
	for _, vlan := range vlans {
		vlan.WriteTo(&buf)
	}
	return buf.String()
}
