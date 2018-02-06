package fwsmConfig

import (
	"bytes"
	"io"
	"net"
)

type SNAT struct {
	Sources      IPNets
	NATTo        net.IP
	FWSMGlobalId int
}

type DNAT struct {
	Destinations IPPorts
	NATTo        IPPort
}

type SNATs []*SNAT
type DNATs []*DNAT

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

func (snats SNATs) CiscoString() string {
	var buf bytes.Buffer
	for _, snat := range snats {
		snat.WriteTo(&buf)
	}
	return buf.String()
}

func (dnat DNAT) WriteTo(writer io.Writer) error {
	return nil
}

func (dnats DNATs) CiscoString() string {
	var buf bytes.Buffer
	for _, dnat := range dnats {
		dnat.WriteTo(&buf)
	}
	return buf.String()
}
