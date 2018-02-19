package fwsmConfig

import (
	"bytes"
	"fmt"
	"io"
	"net"
)

const (
	EXTERNAL_NET = "outside"
)

type SNATSource struct {
	IPNet

	// for FWSM config only:
	IfName string
}

type SNATSources []SNATSource

type SNAT struct {
	Sources SNATSources
	NATTo   net.IP

	// for FWSM config only:
	FWSMGlobalId int
}

type DNAT struct {
	Destinations IPPorts
	NATTo        IPPort

	// for FWSM config only:
	IfName string
}

type SNATs []*SNAT
type DNATs []*DNAT

func (a SNATs) Len() int           { return len(a) }
func (a SNATs) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a SNATs) Less(i, j int) bool { return a[i].FWSMGlobalId < a[j].FWSMGlobalId }

func (a DNATs) Len() int           { return len(a) }
func (a DNATs) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a DNATs) Less(i, j int) bool { return a[i].GetPos() < a[j].GetPos() }

/*func (snat SNAT) GetPos() string {
	return snat.NATTo.String()
}*/

func (dnat DNAT) GetPos() string {
	return dnat.NATTo.String()
}

func (snat SNAT) WriteTo(writer io.Writer) error {
	fmt.Fprintf(writer, "global ("+EXTERNAL_NET+") %v %v\n", snat.FWSMGlobalId, snat.NATTo.String())
	for _, source := range snat.Sources {
		fmt.Fprintf(writer, "nat (%v) %v %v %v\n", source.IfName, snat.FWSMGlobalId, source.IP.String(), net.IP(source.Mask).String())
	}
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
	for _, dst := range dnat.Destinations {
		protocolPrefix := ""
		if dnat.NATTo.Protocol != nil {
			protocolPrefix = dnat.NATTo.Protocol.CiscoString()
			if protocolPrefix == "ip" {
				protocolPrefix = ""
			}
			if protocolPrefix != "" {
				protocolPrefix += " "
			}
		}
		fmt.Fprintf(writer, "static (%v,"+EXTERNAL_NET+") %v %v netmask 255.255.255.255\n", dnat.IfName, protocolPrefix+dst.CiscoString(), dnat.NATTo.CiscoString())
	}
	return nil
}

func (dnats DNATs) CiscoString() string {
	var buf bytes.Buffer
	for _, dnat := range dnats {
		dnat.WriteTo(&buf)
	}
	return buf.String()
}
