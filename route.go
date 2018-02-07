package fwsmConfig

import (
	"bytes"
	"fmt"
	"io"
	"net"
)

type Route struct {
	Sources     IPNets
	Destination IPNet
	Gateway     net.IP
	Metric      int

// for FWSM config only:
	IfName string
}

type Routes []*Route

func (a Routes) Len() int           { return len(a) }
func (a Routes) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a Routes) Less(i, j int) bool { return a[i].GetPos() < a[j].GetPos() }

func (route Route) GetPos() string {
	return route.Gateway.String()
}

func (route Route) WriteTo(writer io.Writer) error {
	if len(route.Sources) != 1 {
		panic("This case is not implemented, yet")
	}
	source := route.Sources[0]

	if net.IP(source.Mask).String() != "0.0.0.0" {
		panic("This case is not implemented, yet")
	}

	_, err := fmt.Fprintf(writer, "route %v %v %v %v %v\n", route.IfName, route.Destination.IP.String(), net.IP(route.Destination.Mask).String(), route.Gateway.String(), route.Metric)
	return err
}

func (routes Routes) CiscoString() string {
	var buf bytes.Buffer
	for _, route := range routes {
		route.WriteTo(&buf)
	}
	return buf.String()
}

