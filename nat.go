package fwsmConfig

import (
	"bytes"
	"fmt"
	"github.com/xaionaro-go/networkControl"
	"io"
	"net"
)

const (
	EXTERNAL_NET = "outside"
)

type SNATSource networkControl.SNATSource
type SNATSources networkControl.SNATSources
type SNAT networkControl.SNAT
type DNAT networkControl.DNAT
type SNATs networkControl.SNATs
type DNATs networkControl.DNATs

func (snat SNAT) CiscoWriteTo(writer io.Writer) error {
	fmt.Fprintf(writer, "global ("+EXTERNAL_NET+") %v %v\n", snat.FWSMGlobalId, snat.NATTo.String())
	for _, source := range snat.Sources {
		fmt.Fprintf(writer, "nat (%v) %v %v %v\n", source.IfName, snat.FWSMGlobalId, source.IP.String(), net.IP(source.Mask).String())
	}
	return nil
}

func (snats SNATs) CiscoString() string {
	var buf bytes.Buffer
	for _, snat := range snats {
		SNAT(*snat).CiscoWriteTo(&buf)
	}
	return buf.String()
}

func (dnat DNAT) CiscoWriteTo(writer io.Writer) error {
	for _, dst := range dnat.Destinations {
		protocolPrefix := ""
		if dnat.NATTo.Protocol != nil {
			protocolPrefix = Protocol(*dnat.NATTo.Protocol).CiscoString()
			if protocolPrefix == "ip" {
				protocolPrefix = ""
			}
			if protocolPrefix != "" {
				protocolPrefix += " "
			}
		}
		fmt.Fprintf(writer, "static (%v,"+EXTERNAL_NET+") %v %v netmask 255.255.255.255\n", dnat.IfName, protocolPrefix+IPPort(dst).CiscoString(), IPPort(dnat.NATTo).CiscoString())
	}
	return nil
}

func (dnats DNATs) CiscoString() string {
	var buf bytes.Buffer
	for _, dnat := range dnats {
		DNAT(*dnat).CiscoWriteTo(&buf)
	}
	return buf.String()
}
