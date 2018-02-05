package fwsmConfig

import (
	"net"
	"io"
)

type SNAT struct {
	Sources IPNets
	NATTo   net.IP
}

type DNAT struct {
	Destinations IPPorts
	NATTo        IPPort
}

type SNATs []SNAT
type DNATs []DNAT

func (a SNATs) Len() int           { return len(a) }
func (a SNATs) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a SNATs) Less(i, j int) bool { return a[i].GetPos() < a[j].GetPos() }
func (a DNATs) Len() int           { return len(a) }
func (a DNATs) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a DNATs) Less(i, j int) bool { return a[i].GetPos() < a[j].GetPos() }

func (snat SNAT) GetPos() string {
	return snat.NATTo.String()
}

func (dnat DNAT) GetPos() string {
	return dnat.NATTo.String()
}

func (snat SNAT) WriteTo(writer io.Writer) error {
	return nil
}

func (dnat DNAT) WriteTo(writer io.Writer) error {
	return nil
}
