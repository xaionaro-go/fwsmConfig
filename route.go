package fwsmConfig

import (
	"bytes"
	"fmt"
	"github.com/xaionaro-go/networkControl"
	"io"
	"net"
)


type Route networkControl.Route
type Routes networkControl.Routes

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
		(Route)(*route).WriteTo(&buf)
	}
	return buf.String()
}
